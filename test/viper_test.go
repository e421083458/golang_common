package test

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"testing"
)

func Test_ViperConf(t *testing.T) {
	viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")
	// any approach to require this configuration into your program.
	var yamlExample = []byte(`
Hacker: true
name: steve
user_name: jobs
hobbies:
- skateboarding
- snowboarding
- go
clothing:
  jacket: leather
  trousers: denim
age: 35
eyes : brown
beard: true
`)
	viper.ReadConfig(bytes.NewBuffer(yamlExample))
	type YamlTConf struct{
		Name string `mapstructure:"user_name"`
	}
	yt:=YamlTConf{}
	if err:=viper.Unmarshal(&yt);err!=nil{
		t.Fatal(err)
	}
	viper.WriteConfigAs("viper.toml")
	fmt.Println()
	fmt.Println(yt)
	fmt.Println(viper.Get("name")) // this would be "steve"
}