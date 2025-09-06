package cmd

import (
	s_tracker "Apartment-Tracker/s-tracker"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	client     s_tracker.Client
	configPath = "/.config/s-tracker"
	rootCmd    = &cobra.Command{
		Use:   "s-tracker",
		Short: "It finds student apartments from s.dk and tracks them and shit idk",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			client, err := s_tracker.NewClient()
			if err != nil {
				return err
			}
			userInfo, err := client.GetUserInfo()
			if err != nil {
				return err
			}

			fmt.Printf("\nLogged in as: %s, (ApplicantId: %s)\n", userInfo.Username, client.GetApplicantId())
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("we do be testing")
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// check if config directory exists, and makes one if it doesn't
	home, err := os.UserHomeDir()
	if _, err := os.Stat(home + configPath); err != nil {
		err = os.Mkdir(home+configPath, 0755)
	}
	cobra.CheckErr(err)

	// config setup stuff
	viper.AddConfigPath(home + configPath)
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	// check if we can read and write
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalln("Error reading config: ", err) // cannot read config
		}

		if err := viper.SafeWriteConfig(); err != nil {
			log.Fatalln("Error writing to config: ", err) // can read config, but cannot write
		}
	}
}
