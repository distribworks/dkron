package dcron

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	RPCAddr  string `mapstructure:"rpc_addr"`
	NodeName string `mapstructure:"node_name"`
	BindAddr string `mapstructure:"bind_addr"`
	HTTPAddr string `mapstructure:"http_addr"`
	Discover string `mapstructure:"discover"`
}

func init() {
	viper.SetConfigName("dcron")    // name of config file (without extension)
	viper.AddConfigPath("./config") // call multiple times to add many search paths
	err := viper.ReadInConfig()     // Find and read the config file
	if err != nil {                 // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	viper.SetEnvPrefix("dcr") // will be uppercased automatically
	viper.AutomaticEnv()
}
