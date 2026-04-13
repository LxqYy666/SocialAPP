package utils

import (
	"fmt"

	"github.com/spf13/viper"
)

func Init() {
	// Implement configuration reading logic here, e.g., using Viper
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".") // look for config in the current directory
	viper.SetDefault("JWT_SECRET", "mysecret")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
	}
}
