/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

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
	Use:   "init-img",
	Short: "Initialize qemu image from local or remote ISO image file",
	Long: `
	This command will search for an existing ISO image,
	or else it will download the latest from our repository.
	Next, it will kick off a base qcow2 image creation,
	which in turn can be used by the snapshot-image command.`,
	Run: func(cmd *cobra.Command, args []string) {
		distro, _ := cmd.Flags().GetString("distro")
		if distro == "ubuntu-focal" {
			ubuntuFocalInit()
		} else {
			fmt.Println(distro, " distro is not supported")
		}
	},
}

func ubuntuFocalInit() {
	ubuntuFocalIsoInit()
	ubuntuFocalImgCreateAndInstall()
}

func ubuntuFocalIsoInit() {
	home, _ := os.UserHomeDir()
	isocfgdir := viper.Get("ubuntu-iso-dir")
	dirstr := fmt.Sprintf("%v", isocfgdir)
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
	fmt.Println("Next run itest snapshot-qcow2. See itest snapshot-qcow2 --help for additional info ...")
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

func ubuntuFocalImgCreateAndInstall() {
	home, _ := os.UserHomeDir()
	t := time.Now()
	timeFormatted := fmt.Sprintf("%d-%02d-%02dT%02d-%02d-%02d-",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	imgcfgdir := viper.Get("ubuntu-images-dir")
	imgdirstr := fmt.Sprintf("%v", imgcfgdir)
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
	hdsize := viper.GetString("ubuntu-hd-size")
	hdsizestr := fmt.Sprintf("%v", hdsize)
	imgfullpath := imgdir + "/" + timeFormatted + id.String() + "-itest-base-ubuntu.img.qcow2"
	fmt.Println("*** Creating Ubuntu Focal qcow2 with .itest.yaml ubuntu-hd-size setting ***")
	imgcreatecmd := exec.Command("qemu-img", "create", "-f", "qcow2", imgfullpath, hdsizestr)
	err := imgcreatecmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(imgfullpath, " successfully created ....")
	fmt.Println("*** Ubuntu Focal qcow2 image creation complete. Installing Ubuntu on image now ... ***")

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
}
