/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// runImageCmd represents the runImage command
var runImageCmd = &cobra.Command{
	Use:   "run-image",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		image, _ := cmd.Flags().GetString("image-name")
		if image != "" {
			runImage(image)
		} else {
			fmt.Println("you must supply an --image-name value")
		}
	},
}

func ubuntuRunImgDirStr() string {
	imgcfgdir := viper.Get("ubuntu-images-dir")
	imgdirstr := fmt.Sprintf("%v", imgcfgdir)
	return imgdirstr
}

func runImage(image string) {
	home, _ := os.UserHomeDir()
	imgdirstr := ubuntuRunImgDirStr()
	fullpath := home + "/" + imgdirstr + "/" + image
	if _, err := os.Stat(fullpath); os.IsNotExist(err) {
		if err != nil && !os.IsExist(err) {
			log.Print(err)
			fmt.Println(fullpath, " image not found")
		}
	} else {
		runcmd := exec.Command("qemu-system-x86_64", "-drive", "file="+fullpath+",format=qcow2", "-cpu", "host", "-enable-kvm", "-m", "10G", "-smp", "2", "-net", "user,hostfwd=tcp::10022-:22", "-net", "nic", "-display", "none")
		fmt.Println("*** Running " + image + " ***")
		runerr := runcmd.Run()
		if runerr != nil {
			log.Fatal(runerr)
		}
	}
}

func init() {
	rootCmd.AddCommand(runImageCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runImageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runImageCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	runImageCmd.Flags().String("image-name", "", "Desired image to run. you must include the full name of the image as retrived from itest images-list")
}
