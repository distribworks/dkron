package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/victorcoder/dkron/dkron"
)

var cfgFile string
var config = &dkron.Config{}

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

	dkronCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is /etc/dkron/dkron.yml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigName("dkron")        // name of config file (without extension)
	viper.AddConfigPath("/etc/dkron")   // call multiple times to add many search paths
	viper.AddConfigPath("$HOME/.dkron") // call multiple times to add many search paths
	viper.AddConfigPath("./config")     // call multiple times to add many search paths
	viper.SetEnvPrefix("dkron")         // will be uppercased automatically
	viper.AutomaticEnv()

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		logrus.WithError(err).Info("No valid config found: Applying default values.")
	}

	cliTags := viper.GetStringSlice("tag")
	var tags map[string]string

	if len(cliTags) > 0 {
		tags, err = unmarshalTags(cliTags)
		if err != nil {
			logrus.Fatal("config: Error unmarshaling cli tags")
		}
	} else {
		tags = viper.GetStringMapString("tags")
	}

	server := viper.GetBool("server")
	nodeName := viper.GetString("node_name")

	if server {
		tags["dkron_server"] = "true"
	} else {
		tags["dkron_server"] = "false"
	}
	tags["dkron_version"] = dkron.Version

	dkron.InitLogger(viper.GetString("log_level"), nodeName)

	dkronCmd.Flags().VisitAll(func(f *pflag.Flag) {
		fmt.Println(f.Value.String())
		v := strings.Replace(f.Name, "-", "_", -1)
		if f.Value.String() != f.DefValue {
			viper.Set(v, f.Value.String())
		} else {
			viper.SetDefault(v, f.Value.String())
		}
	})

	viper.Unmarshal(config)
	spew.Dump(config)
}

// unmarshalTags is a utility function which takes a slice of strings in
// key=value format and returns them as a tag mapping.
func unmarshalTags(tags []string) (map[string]string, error) {
	result := make(map[string]string)
	for _, tag := range tags {
		parts := strings.SplitN(tag, "=", 2)
		if len(parts) != 2 || len(parts[0]) == 0 {
			return nil, fmt.Errorf("Invalid tag: '%s'", tag)
		}
		result[parts[0]] = parts[1]
	}
	return result, nil
}
