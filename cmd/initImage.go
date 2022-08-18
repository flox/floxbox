/*

 */
package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initImageCmd represents the initImage command
var initImageCmd = &cobra.Command{
	Use:   "init-image",
	Short: "Initialize qemu image from local or remote ISO image file",
	Long: `
	This command will search for an existing ISO image,
	or else it will download the latest from our repository.
	Next, it will kick off a base qcow2 image creation,
	which in turn can be used by the snapshot-image command.`,
	Run: func(cmd *cobra.Command, args []string) {
		distro, _ := cmd.Flags().GetString("distro")
		memory, _ := cmd.Flags().GetString("memory")
		cores, _ := cmd.Flags().GetString("cores")
		hdsize, _ := cmd.Flags().GetString("hd-size")
		if distro == "ubuntu-focal" {
			ubuntuFocalInit(memory, cores, hdsize)
		} else {
			fmt.Println(distro, " distro is not supported")
		}
	},
}

func ubuntuFocalInit(memory string, cores string, hdsize string) {
	ubuntuFocalIsoInit()
	ubuntuFocalImgCreateAndInstall(memory, cores, hdsize)
}

func ubuntuIsoDirStr() string {
	//var isodirstr string
	isocfgdir := viper.Get("ubuntu-iso-dir")
	isodirstr := fmt.Sprintf("%v", isocfgdir)
	return isodirstr
}

func ubuntuImgDirStr() string {
	imgcfgdir := viper.Get("ubuntu-images-dir")
	imgdirstr := fmt.Sprintf("%v", imgcfgdir)
	return imgdirstr
}

func ubuntuFocalIsoInit() {
	home, _ := os.UserHomeDir()
	dirstr := ubuntuIsoDirStr()
	isodir := home + "/" + dirstr
	if _, err := os.Stat(isodir); os.IsNotExist(err) {
		err := os.MkdirAll(isodir, 0750)
		if err != nil && !os.IsExist(err) {
			log.Print(err)
		} else {
			fmt.Println("Created ", isodir)
		}
	} else {
		fmt.Println(isodir, " already exists, moving on")
	}
	isopath := isodir + "/" + "itest-ubuntu.iso"
	fileUrl := "https://itestiso.s3.amazonaws.com/itest-ubuntu.iso"
	if _, err := os.Stat(isopath); os.IsNotExist(err) {
		err := isoDownload(isopath, fileUrl)
		if err != nil {
			panic(err)
		}
		fmt.Println("Downloaded: " + fileUrl)
	} else {
		fmt.Println(isopath, " already exists, moving on")
	}
}

func isoDownload(filepath string, url string) error {
	fmt.Println("*** DOWNLOADING ISO ***")
	// Get the iso
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func ubuntuFocalImgCreateAndInstall(memory string, cores string, hdsize string) {
	home, _ := os.UserHomeDir()
	t := time.Now()
	timeFormatted := fmt.Sprintf("%d-%02d-%02dT%02d-%02d-%02d-",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	imgdirstr := ubuntuImgDirStr()
	isodir := ubuntuIsoDirStr()
	imgdir := home + "/" + imgdirstr
	if _, err := os.Stat(imgdir); os.IsNotExist(err) {
		err := os.MkdirAll(imgdir, 0750)
		if err != nil && !os.IsExist(err) {
			log.Print(err)
		} else {
			fmt.Println("Created ", imgdir)
		}
	} else {
		fmt.Println(imgdir, " already exists, moving on")
	}
	id := uuid.New()

	imgfullpath := imgdir + "/" + timeFormatted + id.String() + "-itest-BASE-ubuntu.img.qcow2"
	isofullpath := home + "/" + isodir + "/" + "itest-ubuntu.iso"
	//TODO programmatically recognize iso name or create command flag
	fmt.Println("*** Creating Ubuntu Focal qcow2 with .itest.yaml ubuntu-hd-size setting ***")
	imgcreatecmd := exec.Command("qemu-img", "create", "-f", "qcow2", imgfullpath, hdsize)
	err := imgcreatecmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(imgfullpath, " successfully created ....")
	fmt.Println("*** Ubuntu Focal qcow2 image creation complete. Installing Ubuntu on image now. This may take several minutes. If the install process crashes, you will need to run itest init-image again. Once the install succeeds, you'll need to hit <enter> in the QEMU window that popped up, and login with user: itest password: itest enter poweroff to complete this step. Now, you can run test images-list --distro=ubuntu-focal, and then run itest snapshot-image --distro=ubuntu-focal --base-image-name=<imgname> --snapshot-name=<mysnapshotname> ... ***")

	imginstallcmd := exec.Command("qemu-system-x86_64", "-cdrom", isofullpath, "-drive", "file="+imgfullpath+",format=qcow2", "-enable-kvm", "-m", memory, "-smp", cores)
	fmt.Println("Installing with ", imginstallcmd)
	installerr := imginstallcmd.Run()
	if installerr != nil {
		log.Fatal(installerr)
	}
}

func init() {
	rootCmd.AddCommand(initImageCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initImageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	initImageCmd.Flags().String("distro", "ubuntu-focal", "Desired distro to init. Currently supported distros include: ubuntu-focal")
	initImageCmd.Flags().String("memory", "10G", "The amount of memory you want to assign to the machine you are creating")
	initImageCmd.Flags().String("cores", "2", "The number of host machine cores you want to dedicate to this machine")
	initImageCmd.Flags().String("hd-size", "20G", "The size of the hard drive for virtual machine")
}
