/*

 */
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// snapshotImageCmd represents the snapshotImage command
var snapshotImageCmd = &cobra.Command{
	Use:   "snapshot-image",
	Short: "Createa SNAPSHOT of a BASE image",
	Long:  `Run floxbox list-images --distro=ubunt-focal to list images. Grab the BASE image name you want to snapshot from the list, and run this command to create a snapshot.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("snapshot-image called")
		distro, _ := cmd.Flags().GetString("distro")
		baseimagename, _ := cmd.Flags().GetString("base-image-name")
		snapshotname, _ := cmd.Flags().GetString("snapshot-name")
		if distro == "ubuntu-focal" {
			//ubuntuFocalInit(memory, cores, hdsize)
			fmt.Println("snapshot-image called")
			ubuntuSnapshotImage(baseimagename, snapshotname)
		} else {
			fmt.Println(distro, " distro is not supported")
		}
	},
}

func ubuntuSnapshotImage(baseimagename string, snapshotname string) {
	ssimgdir := viper.Get("ubuntu-snapshot-images-dir")
	ssdirstr := fmt.Sprintf("%v", ssimgdir)
	if _, err := os.Stat(ssdirstr); os.IsNotExist(err) {
		err := os.MkdirAll(ssdirstr, 0750)
		if err != nil && !os.IsExist(err) {
			log.Print(err)
		} else {
			fmt.Println("Created ", ssdirstr)
		}
	} else {
		fmt.Println(ssdirstr, " already exists, moving on")
	}
	//	t := time.Now()
	//	timeFormatted := fmt.Sprintf("%d-%02d-%02dT%02d-%02d-%02d-",
	//		t.Year(), t.Month(), t.Day(),
	//		t.Hour(), t.Minute(), t.Second())
	imgcfgdir := viper.Get("ubuntu-base-images-dir")
	baseimgdirstr := fmt.Sprintf("%v", imgcfgdir)
	basefullpath := baseimgdirstr + "/" + baseimagename
	snapshotfullpath := ssdirstr + "/" + snapshotname
	fmt.Println("*** Creating SNAPSHOT " + snapshotfullpath + " of " + baseimagename + " ***")
	snapshotcreatecmd := exec.Command("qemu-img", "create", "-b", basefullpath, "-f", "qcow2", "-F", "qcow2", snapshotfullpath)
	err := snapshotcreatecmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(snapshotfullpath, " successfully created ....")
}

func init() {

	id := uuid.New()
	rootCmd.AddCommand(snapshotImageCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// snapshotImageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// snapshotImageCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	snapshotImageCmd.Flags().String("distro", "ubuntu-focal", "Desired distro to snapshot. Currently supported distros include: ubuntu-focal")
	snapshotImageCmd.Flags().String("base-image-name", "", "Grab your desired BASE image name from floxbox list-images and put it here")
	snapshotImageCmd.Flags().String("snapshot-name", id.String(), "Create a unique name for your SNAPSHOT image. Otherwise will default to a UUID if no value is given")
}
