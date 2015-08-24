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

var client_list = make([]Client, 0)

func connected(mac, wlan string) {
	disconnected(mac, wlan)
	client_list = append(client_list, Client{time.Now(), wlan, mac})
	fmt.Printf("Connected: %s %s\n", wlan, mac)
}

func disconnected(mac, wlan string) {
	for i, client := range client_list {
		if client.Mac == mac && client.Wlan == wlan {
			client_list = append(client_list[:i], client_list[i+1:]...)
			fmt.Printf("Disconnect: %s %s\n", wlan, mac)
			break
		}
	}
}

func parse(text string) (string, string) {
	wlan, mac := "", ""
	re := regexp.MustCompile("hostapd: ([a-z]+[0-9]+[-]*[0-9]*): STA (?i)(([0-9A-F]{2}[:-]){5}([0-9A-F]{2}))")
	match := re.FindStringSubmatch(text)

	if len(match) == 5 {
		mac, wlan = match[2], match[1]
	}
	return mac, wlan
}

func dump() {
	for i, client := range client_list {
		fmt.Printf("%d - %s\n", i, client.Mac)
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "deauthenticated") {
			disconnected(parse(line))
		} else if strings.Contains(line, "authenticated") {
			connected(parse(line))
		} else if strings.Contains(line, "DUMP") {
			dump()
		}

	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
