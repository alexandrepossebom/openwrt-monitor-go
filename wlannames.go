package main

import (
	"strings"
)

func wlanToHuman(wlan string) string {
	str := ""
	if strings.Contains(wlan, "wlan0-1") {
		str = "[FREE]"
	} else if strings.Contains(wlan, "wlan0") {
		str = "[ 2G ]"
	} else if strings.Contains(wlan, "wlan1") {
		str = "[ 5G ]"
	} else {
		str = wlan
	}
	return str
}
