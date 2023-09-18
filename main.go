package main

import (
	"fmt"
	"linux-mac-fan-control/internal/config"
	"linux-mac-fan-control/internal/manage"
	"os/user"
	"time"
)

func main() {
	curUser, _ := user.Current()
	if curUser.Uid != "0" {
		panic("running in root privilege")
	}

	fans := manage.Fans
	// Enable manual control
	for _, fan := range fans {
		if !fan.ControlEnabled() {
			fan.ToggleControl()
		}
	}
	var destRPM int
	for true {
		fmt.Println("Temp:", manage.CPUTemp())
		for _, fan := range fans {
			switch config.Mode {

			case config.Modes.FixedPercent:
				destRPM = fan.MaxRPM() * config.FixedPercent / 100
				setIfNotSet(fan, destRPM)

			case config.Modes.FixedRPM:
				if fan.MaxRPM() <= config.FixedRPM {
					destRPM = fan.MaxRPM()
				} else {
					destRPM = config.FixedRPM
				}
				setIfNotSet(fan, destRPM)

			case config.Modes.TempRPM:
				thisConfig, err := config.TempPercent[fan.Name()]
				if !err {
					panic("No config for" + fan.Name())
				}
				destRPM = best_level_according_to_temp(thisConfig)
				setIfNotSet(fan, destRPM)

			case config.Modes.TempPercent:
				destPercent := best_level_according_to_temp(config.TempPercent[fan.Name()])
				destRPM = fan.MaxRPM() * destPercent / 100
				if fan.MaxRPM() <= destRPM {
					destRPM = fan.MaxRPM()
				}
				setIfNotSet(fan, destRPM)
			}
			fmt.Println("Fan:", fan.Name(), "DestRPM:", destRPM, "CurRPM", fan.CurrentRPM())
		}
		time.Sleep(time.Duration(config.Interval) * time.Millisecond)
	}
}

func setIfNotSet(fan manage.TheFan, destRPM int) {
	if fan.DestRPM() != destRPM {
		fan.SetRPM(destRPM)
	}
}

func best_level_according_to_temp(thisConfig map[int]int) int {
	curTemp := manage.CPUTemp()
	lowestAcceptable := 65535 // Just a big number
	destLevel := 65535 // Just a big number
	for k, v := range thisConfig {
		if curTemp <= k && k <= lowestAcceptable {
			destLevel = v
			lowestAcceptable = k
		}
	}
	fmt.Println("BestLevel:", destLevel)
	return destLevel
}
