package main

import (
	"fmt"
	"github.com/ActiveState/tail"
	"regexp"
	"strings"
)


var host_map = map[string]string{}
var reHosts = regexp.MustCompile("((?i)([0-9A-F]{2}[:-]){5}(?i)([0-9A-F]{2})) (.*)$")

func readDhcp() {
	t, err := tail.TailFile("/tmp/dnsmasq-dhcp.log", tail.Config{Follow: true, ReOpen: true})
	for line := range t.Lines {
		if strings.Contains(line.Text, "DHCPACK") {
			parseHost(line.Text)
		}
	}
	fmt.Println(err)
}

func parseHost(text string) {
	match := reHosts.FindStringSubmatch(text)
	if len(match) == 5 {
		host := match[4]
		mac := match[1]

		if len(host) == 0 {
			return
		}

		if val, ok := host_map[mac]; ok {
			if val != host {
				host_map[mac] = host
				updateHost(host, mac)
				// fmt.Printf("host updated!\n")
			}
		} else {
			host_map[mac] = host
			addHost(host, mac)
			//  fmt.Printf("new host : %s = %s !\n", mac, host)
		}
	}
}

func getHostName(mac string) string {
	if val, ok := host_map[mac]; ok {
		return val
	}
	return mac
}
