package dcron

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var config *viper.Viper

func init() {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	config = viper.New()

	config.SetDefault("rpc_addr", "127.0.0.1:7373")
	config.SetDefault("discover", "dcron")
	config.SetDefault("node_name", hostname)
	config.SetDefault("node", config.GetString("node_name"))
	config.SetDefault("bind", "0.0.0.0:7946")
	config.SetDefault("http_addr", ":8080")

	config.SetConfigName("dcron")    // name of config file (without extension)
	config.AddConfigPath("./config") // call multiple times to add many search paths
	err = config.ReadInConfig()      // Find and read the config file
	if err != nil {                  // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	config.SetEnvPrefix("dcr") // will be uppercased automatically
	config.AutomaticEnv()
}
