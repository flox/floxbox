/*

 */
package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initImageCmd represents the initImage command
var initImageCmd = &cobra.Command{
	Use:   "init-image",
	Short: "Initialize qemu image from local or remote ISO image file",
	Long: `
	This command will search for an existing base image,
	or else it will download the latest from the ditro creator.
	Next, it will kick off a base qcow2 image creation,
	which in turn can be used by the snapshot-image command.`,
	Run: func(cmd *cobra.Command, args []string) {
		dirstro, _ := cmd.Flags().GetString("distro")
		memory, _ := cmd.Flags().GetString("memory")
		cores, _ := cmd.Flags().GetString("cores")
		hdsize, _ := cmd.Flags().GetString("hd-size")
		metadata, _ := cmd.Flags().GetString("metadata-file")
		userdata, _ := cmd.Flags().GetString("user-data-file")
		ubuntu_img_url, _ := cmd.Flags().GetString("ubuntu-img-url")
		if dirstro == "ubuntu-focal" {
			ubuntuFocalInit(memory, cores, hdsize, metadata, userdata, ubuntu_img_url)
		} else {
			fmt.Println(dirstro, " dirstro is not supported")
		}
	},
}

func ubuntuFocalInit(memory string, cores string, hdsize string, metadata string, userdata string, ubuntu_img_url string) {
	//ubuntuFocalIsoInit()
	//ubuntuCloudImageDownload(ubuntu_img_url)
	ubuntuImgCreateAndInstall(memory, cores, hdsize, metadata, userdata, ubuntu_img_url)
}

func ubuntuBaseImgDirStr() string {
	imgcfgdir := viper.Get("ubuntu-base-images-dir")
	imgdirstr := fmt.Sprintf("%v", imgcfgdir)
	return imgdirstr
}

func ubuntuImgCreateAndInstall(memory string, cores string, hdsize string, metadata string, userdata string, ubuntu_img_url string) {
	fmt.Println("*** DOWNLOADING: " + ubuntu_img_url + " ***")

	// Get the cloud image
	dirstr := ubuntuBaseImgDirStr()
	if _, err := os.Stat(dirstr); os.IsNotExist(err) {
		err := os.MkdirAll(dirstr, 0750)
		if err != nil && !os.IsExist(err) {
			log.Print(err)
		} else {
			fmt.Println("Created ", dirstr)
		}
	} else {
		fmt.Println(dirstr, " already exists, moving on")
	}
	os.Chdir(dirstr)
	parsed_url, err := url.Parse(ubuntu_img_url)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.Get(ubuntu_img_url)
	if err != nil {
		log.Fatal(err)
	}
	path := parsed_url.Path
	segments := strings.Split(path, "/")
	fileName := segments[len(segments)-1]
	defer resp.Body.Close()

	out, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	//TODO SHA256SUM based on https://ubuntu.com/tutorials/how-to-verify-ubuntu#6-check-the-iso
	//create seed image
	seedimgcreate := exec.Command("cloud-localds", "seed.img", userdata, metadata)
	err = seedimgcreate.Run()
	//TODO seed handling
	if err != nil {
		log.Fatal(err)
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
	initImageCmd.Flags().String("distro", "ubuntu-focal", "Desired dirstro to init. Currently supported dirstros include: ubuntu-focal")
	initImageCmd.Flags().String("memory", "10G", "The amount of memory you want to assign to the machine you are creating")
	initImageCmd.Flags().String("cores", "2", "The number of host machine cores you want to dedicate to this machine")
	initImageCmd.Flags().String("hd-size", "20G", "The size of the hard drive for virtual machine")
	initImageCmd.Flags().String("metadata-file", "", "The path to the metadata.yml file for use by cloud-init")
	initImageCmd.Flags().String("user-data-file", "", "The path to the user-dta.yml file for use by cloud-init")
	initImageCmd.Flags().String("ubuntu-img-url", "", "The direct URL to source the ubuntu cloud image")
}
