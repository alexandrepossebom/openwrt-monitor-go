package main

import (
	"fmt"
	"github.com/ActiveState/tail"
	"regexp"
	"strings"
	"time"
)

var re = regexp.MustCompile("(.*) (.*) hostapd: ([a-z]+[0-9]+[-]*[0-9]*): STA (?i)(([0-9A-F]{2}[:-]){5}([0-9A-F]{2}))")

func readHostAp() {
	t, err := tail.TailFile("/tmp/hostapd.log", tail.Config{Follow: true, ReOpen: true})
	for line := range t.Lines {
		if strings.Contains(line.Text, "deauthenticated") {
			disconnected(parse(line.Text))
		} else if strings.Contains(line.Text, "authenticated") {
			connected(parse(line.Text))
		}
	}
	fmt.Println(err)
}

func parse(text string) (string, string, time.Time) {
	wlan, mac, str := "", "", ""
	match := re.FindStringSubmatch(text)

	if len(match) == 7 {
		mac, wlan, str = match[4], match[3], match[1]
	}

	return mac, wlan, parseDate(str)
}
