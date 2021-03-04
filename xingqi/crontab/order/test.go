package main

import (
	"fmt"
	"github.com/spf13/viper"
)

func main() {

	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")    // optionally look for config in the working directory
	viper.ReadInConfig() // Find and read the config file
	//if err != nil { // Handle errors reading the config file
	//	panic(fmt.Errorf("Fatal error config file: %s \n", err))
	//}

	fmt.Println(viper.Get("last_chufangid"))
	fmt.Println(viper.Get("token"))
	viper.Set("last_chufangid", 1111)
	fmt.Println(viper.Get("last_chufangid"))
	fmt.Println(viper.Get("token"))
	viper.WriteConfig()

	//viper.WriteConfig() // writes current config to predefined path set by 'viper.AddConfigPath()' and 'viper.SetConfigName'
	//viper.SafeWriteConfig()
	//viper.WriteConfigAs("/path/to/my/.config")
	//viper.SafeWriteConfigAs("/path/to/my/.config") // will error since it has already been written
	//viper.SafeWriteConfigAs("/path/to/my/.other_config")

	//func initTOML() {
	//	//Reset()
	//	SetConfigType("toml")
	//	r := bytes.NewReader(tomlExample)
	//
	//	unmarshalReader(r, v.config)
	//}

}
