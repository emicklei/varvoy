package api

import (
	"os"
	"strconv"
	"strings"
)

func IsTraceEnabled() bool {
	return getLogFlag()
}

func flagValueString(f string) string {
	if !strings.Contains(f, "=") {
		return ""
	}
	return strings.Split(f, "=")[1]
}

func flagValueInt(f string) int {
	if !strings.Contains(f, "=") {
		return 0
	}
	is := strings.Split(f, "=")[1]
	i, err := strconv.Atoi(is)
	if err != nil {
		return 0
	}
	return i
}

func getListenFlag() string {
	for _, each := range os.Args {
		if strings.HasPrefix(each, "--listen") {
			return each
		}
	}
	return ""
}
func getLogDestFlag() string {
	for _, each := range os.Args {
		if strings.HasPrefix(each, "--log-dest") {
			return each
		}
	}
	return ""
}
func getLogFlag() bool {
	for _, each := range os.Args {
		if each == "--log" {
			return true
		}
	}
	return false
}
