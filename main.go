package main

import (
	"bufio"
	"fmt"
	"linux-mac-fan-control/internal/config"
	"linux-mac-fan-control/internal/manage"
	"os"
	"os/user"
	"time"
)

func main() {
	curUser, _ := user.Current()
	if curUser.Uid != "0" {
		panic("please running in root privilege")
	}
	fans := manage.Fans
	reader := bufio.NewReader(os.Stdin)
	quit := make(chan bool)
	fmt.Println("Linux Mac FanControl")
	helpMsg := "1. Status \t 2. Reload \t 3. Exit"
	fmt.Println(helpMsg)
	go run(quit, &fans)
	var ipt string
	mainLoop:
	for true {
		fmt.Print("> ")
		ipt, _ = reader.ReadString('\n')
		switch ipt[:len(ipt)-1] {
		case "1":
			fmt.Println("CPU Package: ", manage.CPUTemp(), "°C")
			for _, fan := range fans {
				fmt.Print(fan.Name(), " RPM: ", fan.CurrentRPM(), "\t")
			}
			fmt.Println()
		case "2":
			quit <- true
			for _, _ = range [3]int{} {
				fmt.Print(".")
				time.Sleep(500 * time.Millisecond)
			}
			go run(quit, &fans)
			fmt.Println("\nSuccess!")
		case "3":
			quit <- true
			for _, fan := range fans {
				if fan.ControlEnabled() {
					fan.ToggleControl()
				}
			}
			fmt.Println("Goodbye")
			break mainLoop
		case "h":
			fmt.Println(helpMsg)
		default:
			fmt.Println("invalid input")
		}
	}
}

func run(quit chan bool, fans *[]manage.TheFan) {
	// Enable manual control
	/*for _, fan := range *fans {
		if !fan.ControlEnabled() {
			fan.ToggleControl()
		}
	}*/
	var destRPM int
	for true {
		select {
		case <-quit:
			break
		default:
			for _, fan := range *fans {
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
					thisConfig, err := config.TempRPM[fan.Name()]
					if !err {
						panic("No config for" + fan.Name())
					}
					destRPM = best_level_according_to_temp(thisConfig, config.TempLevelRPM[fan.Name()])
					setIfNotSet(fan, destRPM)

				case config.Modes.TempPercent:
					destPercent := best_level_according_to_temp(config.TempPercent[fan.Name()], config.TempLevelPercent[fan.Name()])
					destRPM = fan.MaxRPM() * destPercent / 100
					if fan.MaxRPM() <= destRPM {
						destRPM = fan.MaxRPM()
					}
					setIfNotSet(fan, destRPM)
				}
			}
			time.Sleep(time.Duration(config.Interval) * time.Millisecond)
		}
	}
}

func setIfNotSet(fan manage.TheFan, destRPM int) {
	if fan.DestRPM() != destRPM {
		fan.SetRPM(destRPM)
	}
}

func best_level_according_to_temp(thisConfig map[int]int, tempLevels []int) int {
	curTemp := manage.CPUTemp()
	destTempLevelIndex := -1
	for i:=0;i<len(tempLevels);i++{
		if curTemp <= tempLevels[i] {
			destTempLevelIndex=i
			break
		}
	}
	
	if destTempLevelIndex == -1 {
		
	}
	
	output:=thisConfig[tempLevels[destTempLevelIndex]]
	
	if tempLevels[destTempLevelIndex] != curTemp && destTempLevelIndex != 0 {
		//output = thisConfig[tempLevels[destTempLevelIndex]]+((thisConfig[tempLevels[destTempLevelIndex+1]] - thisConfig[tempLevels[destTempLevelIndex]]) / (tempLevels[destTempLevelIndex+1] - tempLevels[destTempLevelIndex]) * (tempLevels[destTempLevelIndex+1] - tempLevels[destTempLevelIndex]))
		hiTemp:=tempLevels[destTempLevelIndex]
		loTemp:=tempLevels[destTempLevelIndex-1]
		hiLevel:=thisConfig[hiTemp]
		loLevel:=thisConfig[loTemp]
		
		output = loLevel + ((curTemp - loTemp) * ((hiLevel-loLevel) / (hiTemp - loTemp)))
	}
	
	// fmt.Println( destLevel+((secondLowLevel - destLevel) / (secondLowTempLevel - lowestAcceptable) * (curTemp - lowestAcceptable)))
	return output
}
