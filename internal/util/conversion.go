package util

import (
	"os"
	"strconv"
	"strings"
)

func Read2int(file string) int {
	reads, _ := os.ReadFile(file)
	num, _ := strconv.Atoi(strings.ReplaceAll(string(reads), "\n", ""))
	return num
}
