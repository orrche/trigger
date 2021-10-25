package triggerkasa

import (
	"io/ioutil"
	"log"

	"github.com/ivanbeldad/kasa-go"
	yaml "gopkg.in/yaml.v2"
)

type kasaConf struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

var c kasaConf
var api kasa.API

func turnOff() {
	hs100, err := api.GetHS100("AC")
	if err != nil {
		log.Fatal(err)
	}

	err = hs100.TurnOff()
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	yamlFile, err := ioutil.ReadFile("configs/kasa.yml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("Unmarhsal: %v", err)
	}
	api, err = kasa.Connect(c.User, c.Password)
	if err != nil {
		log.Fatalf("Uhm: %v", err)
	}
}
