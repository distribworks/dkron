package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/distribworks/dkron/v3/dkron"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var config = dkron.DefaultConfig()

var rpcAddr string
var ip string

// dkronCmd represents the dkron command
var dkronCmd = &cobra.Command{
	Use:   "dkron",
	Short: "Open source distributed job scheduling system",
	Long: `Dkron is a system service that runs scheduled jobs at given intervals or times,
just like the cron unix service but distributed in several machines in a cluster.
If a machine fails (the leader), a follower will take over and keep running the scheduled jobs without human intervention.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := dkronCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	dkronCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("dkron")        // name of config file (without extension)
		viper.AddConfigPath("/etc/dkron")   // call multiple times to add many search paths
		viper.AddConfigPath("$HOME/.dkron") // call multiple times to add many search paths
		viper.AddConfigPath("./config")     // call multiple times to add many search paths
	}

	viper.SetEnvPrefix("dkron")
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv() // read in environment variables that match

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		logrus.WithError(err).Info("No valid config found: Applying default values.")
	}

	if err := viper.Unmarshal(config); err != nil {
		logrus.WithError(err).Fatal("config: Error unmarshalling config")
	}

	cliTags := viper.GetStringSlice("tag")
	var tags map[string]string

	if len(cliTags) > 0 {
		tags, err = UnmarshalTags(cliTags)
		if err != nil {
			logrus.WithError(err).Fatal("config: Error unmarshalling cli tags")
		}
	} else {
		tags = viper.GetStringMapString("tags")
	}

	config.Tags = tags

	dkron.InitLogger(viper.GetString("log-level"), config.NodeName)
}
