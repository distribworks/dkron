package dcron

import (
	"fmt"
	"github.com/spf13/viper"
)

var config *viper.Viper

func init() {
	config = viper.New()
	config.SetConfigName("dcron")    // name of config file (without extension)
	config.AddConfigPath("./config") // call multiple times to add many search paths
	err := config.ReadInConfig()     // Find and read the config file
	if err != nil {                  // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	config.SetEnvPrefix("dcr") // will be uppercased automatically
	config.AutomaticEnv()
}
