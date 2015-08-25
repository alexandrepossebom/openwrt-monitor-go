package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

type Client struct {
	Connect time.Time
	Wlan    string
	Mac     string
}

func periodicFree() {
	for _ = range time.Tick(5 * time.Second) {
		dump()
	}
}

var client_map = map[string]Client{}
var host_map = map[string]string{}

var re = regexp.MustCompile("hostapd: ([a-z]+[0-9]+[-]*[0-9]*): STA (?i)(([0-9A-F]{2}[:-]){5}([0-9A-F]{2}))")
var reHosts = regexp.MustCompile("((?i)([0-9A-F]{2}[:-]){5}(?i)([0-9A-F]{2})) (.*)$")

func connected(mac, wlan string) {
	disconnected(mac, wlan)
	client_map[mac] = Client{time.Now(), wlan, mac}
	fmt.Printf("Connected: %s %s\n", wlan, mac)
}

func disconnected(mac, wlan string) {
	if _, ok := client_map[mac]; ok {
		delete(client_map, mac)
		fmt.Printf("Disconnected: %s %s\n", wlan, mac)
	}
}

func parse(text string) (string, string) {
	wlan, mac := "", ""
	match := re.FindStringSubmatch(text)

	if len(match) == 5 {
		mac, wlan = match[2], match[1]
	}

	return mac, wlan
}

func parseHost(text string) {
	match := reHosts.FindStringSubmatch(text)
	if len(match) == 5 {
		host := match[4]
		mac := match[1]
		fmt.Printf("match : %s -> %s\n", host, mac)

		if val, ok := host_map[mac]; ok {
			if val != host {
				host_map[mac] = host
				fmt.Printf("host updated!\n")
			}
		} else {
			host_map[mac] = host
			fmt.Printf("new host!\n")
		}
	}
}

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

func calcDate(connect time.Time) string {
	duration := time.Since(connect)
	elapsed := int(duration.Seconds())

	days := elapsed / 86400
	elapsed = elapsed % 86400
	hours := elapsed / 3600
	elapsed = elapsed % 3600
	minutes := elapsed / 60
	seconds := elapsed % 60

	strTime := ""

	if days > 0 {
		strTime = fmt.Sprintf("%02dd%02dh%02dm%02ds", days, hours, minutes, seconds)
	} else if hours > 0 {
		strTime = fmt.Sprintf("   %02dh%02dm%02ds", hours, minutes, seconds)
	} else if minutes > 0 {
		strTime = fmt.Sprintf("      %02dm%02ds", minutes, seconds)
	} else {
		strTime = fmt.Sprintf("         %02ds", seconds)
	}
	return strTime
}

func getHostName(mac string) string {
	str := mac
	if val, ok := host_map[mac]; ok {
		str = val
	}
	return str
}

func dump() {
	fmt.Printf("DUMP\n")
	for _, client := range client_map {
		fmt.Printf("%s [-->] %s %-25s\n", calcDate(client.Connect), wlanToHuman(client.Wlan), getHostName(client.Mac))
	}
}

func main() {
	go periodicFree()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "deauthenticated") {
			disconnected(parse(line))
		} else if strings.Contains(line, "authenticated") {
			connected(parse(line))
		} else if strings.Contains(line, "DHCPACK") {
			parseHost(line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
