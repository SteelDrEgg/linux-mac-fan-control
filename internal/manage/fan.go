package manage

import (
	"linux-mac-fan-control/internal/util"
	"os"
	"strconv"
	"strings"
)

var (
	fans    []string
	Fans    []TheFan
	baseDir = "/sys/devices/platform/applesmc.768"
)

func init() {
	items, _ := os.ReadDir(baseDir)
	for _, item := range items {
		name := item.Name()
		if strings.Contains(name, "fan") && strings.Contains(name, "manual") {
			fan := strings.Split(name, "_")
			fanExist := false
			for _, v := range fans {
				if fan[0] == v {
					fanExist = true
					break
				}
			}
			if !fanExist {
				fans = append(fans, fan[0])
			}
		}
	}
	loadAllFans()
}

func loadAllFans() {
	for _, fanName := range fans {
		util.Read2int(baseDir + "/" + fanName + "_max")
		var newFan TheFan = &fan{
			name: fanName,
		}
		newFan.UpdateStatus()
		Fans = append(Fans, newFan)
	}
}

type TheFan interface {
	ToggleControl()
	UpdateStatus()
	SetRPM(int)
	CurrentRPM() int
	MaxRPM() int
	MinRPM() int
	ControlEnabled() bool
	DestRPM() int
	Name() string
}

type fan struct {
	manual  bool
	name    string
	maxRPM  int
	minRPM  int
	destRPM int
}

func (self *fan) ToggleControl() {
	fileName := baseDir + "/" + self.name + "_manual"
	var toWrite string
	status, _ := os.ReadFile(fileName)
	if string(status)[0] == '0' {
		toWrite = "1"
		self.manual = true
	} else {
		toWrite = "0"
		self.manual = false
	}
	_ = os.WriteFile(fileName, []byte(toWrite), 0644)
}

func (self *fan) UpdateStatus() {
	self.maxRPM = util.Read2int(baseDir + "/" + self.name + "_max")
	self.minRPM = util.Read2int(baseDir + "/" + self.name + "_min")

	status, _ := os.ReadFile(baseDir + "/" + self.name + "_manual")
	if string(status)[0] == '0' {
		self.manual = false
	} else {
		self.manual = true
	}
}

func (self *fan) SetRPM(rpm int) {
	fileName := baseDir + "/" + self.name + "_output"
	self.destRPM = rpm
	_ = os.WriteFile(fileName, []byte(strconv.Itoa(rpm)), 0644)
}

func (self fan) CurrentRPM() int {
	fileName := baseDir + "/" + self.name + "_input"
	return util.Read2int(fileName)
}

func (self fan) MaxRPM() int {
	return self.maxRPM
}
func (self fan) MinRPM() int {
	return self.minRPM
}
func (self fan) ControlEnabled() bool {
	return self.manual
}
func (self fan) DestRPM() int {
	return self.destRPM
}
func (self fan) Name() string {
	return self.name
}
