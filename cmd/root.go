/*

 */
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "flox-qemu",
	Short: "A linux qemu/kvm management tool for flox",
	Long:  `This tool can configure and manage qemu/kvm images for testing or general use with flox`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func intializeDirectory() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	dotdir := home + "/" + ".flox-qemu"

	if _, err := os.Stat(dotdir); os.IsNotExist(err) {
		err := os.MkdirAll(dotdir, 0750)
		if err != nil && !os.IsExist(err) {
			log.Print(err)
		} else {
			fmt.Println("Created ", dotdir)
		}
	} else {
		fmt.Println(dotdir, " already exists, moving on")
	}

}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(intializeDirectory)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.flox-qemu.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".flox-qemu" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".flox-qemu")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else {

		home, err := os.UserHomeDir()
		viper.Set("ubuntu-base-images-dir", home+"/.flox-qemu/ubuntu-base-images")
		viper.Set("ubuntu-snapshot-images-dir", home+"/.flox-qemu/ubuntu-snapshot-images")
		cobra.CheckErr(err)
		cnfpath := home + "/" + ".flox-qemu.yaml"
		viper.SafeWriteConfigAs(cnfpath)
	}
}
