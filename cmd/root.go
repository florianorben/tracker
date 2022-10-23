package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/user"
	"path"
	"github.com/florianorben/tracker/helpers"
	"github.com/florianorben/tracker/tracker"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const AppVersion = "1.0.0"

var (
	configFile = "config"
	appDir     string
)

var RootCmd = &cobra.Command{
	Use:   "tracker",
	Short: "Tracker is a tool aimed at helping you monitoring your time.",
	Long: `Tracker is a tool aimed at helping you monitoring your time.

  You just have to tell Watson when you start working on your project with
  the 'start' command, and you can stop the timer when you're done with the
  'stop' command.

  Heavily inspired by Watson from TailorDev (https://github.com/TailorDev/Watson/)`,
	Run: root,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.Flags().BoolP("version", "v", false, "Show the version and exit.")

	u, err := user.Current()
	if err != nil {
		fmt.Printf("Error: %s\n", helpers.PrintRed(err.Error()))
		os.Exit(1)
	}
	appDir = path.Join(u.HomeDir, ".tracker")

	cobra.OnInitialize(createAppFolder, createConfigFile, initConfig)
}

func main() {
}

func root(cmd *cobra.Command, args []string) {
	version, err := cmd.Flags().GetBool("version")

	if err == nil && version != false {
		fmt.Println("tracker " + AppVersion)
		os.Exit(0)
	}

	cmd.Help()
}

func createAppFolder() {
	err := os.Mkdir(appDir, 0744)
	if err != nil && !os.IsExist(err) {
		fmt.Printf("Error: %s\n", helpers.PrintRed(err.Error()))
		os.Exit(1)
	}
}

func createConfigFile() {
	filename := path.Join(appDir, configFile+".toml")

	if file, err := os.Open(filename); err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(filename)
			if err != nil {
				fmt.Printf("Error: %s\n", helpers.PrintRed(err.Error()))
				os.Exit(1)
			}
			defer file.Close()

			_, err := file.Write(bytes.TrimSpace(tracker.DefaultConfig))
			if err != nil {
				fmt.Printf("Error: %s\n", helpers.PrintRed(err.Error()))
				os.Exit(1)
			}

			fmt.Printf("Created config file at %s\n", filename)
		} else {
			fmt.Printf("Error: %s\n", helpers.PrintRed(err.Error()))
			os.Exit(1)
		}
	} else {
		file.Close()
	}
}

func initConfig() {
	viper.SetConfigName(configFile) // name of config file (without extension)
	viper.AddConfigPath(appDir)     // adding home directory as first search path
	viper.AutomaticEnv()            // read in environment variables that match

	viper.SetDefault("appDir", appDir)
	viper.SetDefault("framesFile", path.Join(appDir, "frames"))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error: %s\n", helpers.PrintRed("configuration file could not be read. Using default configuration."))
	}
}
