package manage

import (
	"os"
	"strconv"
	"strings"
)

var (
	tempFile string
)

func init() {
	baseDir := "/sys/class/thermal"
	thermals, _ := os.ReadDir(baseDir)
	for _, zone := range thermals {
		name := zone.Name()
		if strings.Contains(name, "thermal_zone") {
			zoneType, _ := os.ReadFile(baseDir + "/" + name + "/type")
			if strings.Contains(string(zoneType), "x86_pkg_temp") {
				tempFile = baseDir + "/" + name + "/temp"
				return
			}
		}
	}
	panic("x86_pkg_temp not found!")
}

func CPUTemp() int {
	temp, _ := os.ReadFile(tempFile)
	temp_int, _ := strconv.Atoi(string(temp)[:len(temp)-1])
	return temp_int / 1000
}
