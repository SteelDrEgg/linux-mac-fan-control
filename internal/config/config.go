package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"strconv"
)

type modes struct {
	FixedRPM     int
	FixedPercent int
	TempPercent  int
	TempRPM      int
}

var Modes = modes{
	FixedRPM:     0,
	FixedPercent: 1,
	TempPercent:  2,
	TempRPM:      3,
}

var (
	Interval     int
	Mode         int
	FixedPercent int
	FixedRPM     int
	TempPercent  map[string]map[int]int
	TempRPM      map[string]map[int]int
)

func init() {
	viper.SetConfigType("yaml")
	viper.SetConfigFile("./config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	loadConfig()
}

func loadConfig() {
	if viper.InConfig("Interval") {
		Interval = viper.Get("Interval").(int)
	} else {
		Interval = 1000
	}

	if viper.InConfig("Mode") {
		mode := viper.Get("Mode").(string)
		if mode == "FixedRPM" {
			Mode = Modes.FixedRPM
		} else if mode == "FixedPercent" {
			Mode = Modes.FixedPercent
		} else if mode == "TempPercent" {
			Mode = Modes.TempPercent
		} else if mode == "TempRPM" {
			Mode = Modes.TempRPM
		} else {
			panic(errors.New("Incorrect config at 'Mode'"))
		}
	} else {
		panic("Incorrect config, missing 'Mode'")
	}

	if viper.InConfig("FixedPercent") {
		FixedPercent = viper.Get("FixedPercent").(int)
	} else {
		FixedPercent = 100
	}
	if viper.InConfig("FixedRPM") {
		FixedRPM = viper.Get("FixedRPM").(int)
	} else {
		FixedRPM = 2500
	}

	if viper.InConfig("TempPercent") {
		tp := viper.Get("TempPercent").([]interface{})
		TempPercent = make(map[string]map[int]int)
		parseTempConfig(TempPercent, tp)
	}

	if viper.InConfig("TempRPM") {
		tr := viper.Get("TempRPM").([]interface{})
		TempRPM = make(map[string]map[int]int)
		parseTempConfig(TempRPM, tr)
	}

}

func parseTempConfig(target map[string]map[int]int, set []interface{}) {
	for _, m := range set {
		temp := make(map[int]int)
		fanName := ""
		for k, v := range m.(map[string]interface{}) {
			temperature, err := strconv.Atoi(k)
			if err != nil {
				fanName = fmt.Sprintf("%s%v", k, v)
			} else {
				temp[temperature] = v.(int)
			}
		}
		target[fanName] = temp
	}
}
