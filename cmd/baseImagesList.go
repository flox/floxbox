/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// imagesListCmd represents the imagesList command
var baseImagesListCmd = &cobra.Command{
	Use:   "base-images-list",
	Short: "list base images based on disto type",
	Long:  `This command will list all the base  images for the distro you specify.`,
	Run: func(cmd *cobra.Command, args []string) {
		distro, _ := cmd.Flags().GetString("distro")
		if distro == "ubuntu-focal" {
			ubuntuFocalBaseImagesList()
		} else {
			fmt.Println(distro, " distro is not supported")
		}
	},
}

func ubuntuBaseImgFullPathStr() string {
	imgcfgdir := viper.Get("ubuntu-base-images-dir")
	imgdirstr := fmt.Sprintf("%v", imgcfgdir)
	return imgdirstr
}

func ubuntuFocalBaseImagesList() {
	path := ubuntuBaseImgFullPathStr()
	cmd := exec.Command("ls", "-lah", path)
	output, _ := cmd.CombinedOutput()
	fmt.Println(string(output))
}

func init() {
	rootCmd.AddCommand(baseImagesListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// imagesListCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// imagesListCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	baseImagesListCmd.Flags().String("distro", "ubuntu-focal", "Desired distro type to list images for. Currently supported distros include: ubuntu-focal")
}
