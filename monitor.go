package main

import (
	"fmt"
	"github.com/ActiveState/tail"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type Log struct {
	DtConn   time.Time
	DtDisc   time.Time
	Wlan     string
	Mac      string
	Duration string
}

type Client struct {
	DtConn time.Time
	Wlan   string
	Mac    string
}

func periodicDump() {
	for _ = range time.Tick(30 * time.Second) {
		dump()
	}
}

var client_map = map[string]Client{}
var host_map = map[string]string{}
var logs []Log

var re = regexp.MustCompile("(.*) (.*) hostapd: ([a-z]+[0-9]+[-]*[0-9]*): STA (?i)(([0-9A-F]{2}[:-]){5}([0-9A-F]{2}))")
var reHosts = regexp.MustCompile("((?i)([0-9A-F]{2}[:-]){5}(?i)([0-9A-F]{2})) (.*)$")

func parseDate(str string) time.Time {
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		fmt.Println(err)
		return time.Now()
	}
	return t
}

func lp(l Log) {
	fmt.Printf("%s [<---] %s %-15s %s\n", l.DtDisc.Format(time.RFC822), l.Wlan, getHostName(l.Mac), l.Duration)
	if strings.Contains(l.Duration, "m") {
		logs = append(logs, l)
	}
}

func connected(mac string, wlan string, dt time.Time) {
	disconnected(mac, wlan, dt)
	client_map[mac+wlan] = Client{dt, wlan, mac}
	fmt.Printf("%s [--->] %s %-15s\n", dt.Format(time.RFC822), wlanToHuman(wlan), getHostName(mac))
}

func disconnected(mac string, wlan string, dt time.Time) {
	if _, ok := client_map[mac+wlan]; ok {
		client := client_map[mac+wlan]
		l := Log{client.DtConn, dt, wlanToHuman(client.Wlan), client.Mac, calcDate(client.DtConn, dt)}
		lp(l)
		delete(client_map, mac+wlan)
	}
}

func parse(text string) (string, string, time.Time) {
	wlan, mac, str := "", "", ""
	match := re.FindStringSubmatch(text)

	if len(match) == 7 {
		mac, wlan, str = match[4], match[3], match[1]
	}

	return mac, wlan, parseDate(str)
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
				// fmt.Printf("host updated!\n")
			}
		} else {
			host_map[mac] = host
			//  fmt.Printf("new host : %s = %s !\n", mac, host)
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

func calcDate(connect, disconnect time.Time) string {
	duration := disconnect.Sub(connect)
	elapsed := int(duration.Seconds())

	// fmt.Printf("duration %s %s %d !\n" , connect.Format(time.RFC822),disconnect.Format(time.RFC822), duration)

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
	if val, ok := host_map[mac]; ok {
		return val
	}
	return mac
}

func dump() {
	for _, client := range client_map {
		fmt.Printf("%s [DUMP] %s %-15s %s\n", time.Now().Format(time.RFC822), wlanToHuman(client.Wlan), getHostName(client.Mac), calcDate(client.DtConn, time.Now()))
	}
}

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

func readDhcp() {
	t, err := tail.TailFile("/tmp/dnsmasq-dhcp.log", tail.Config{Follow: true, ReOpen: true})
	for line := range t.Lines {
		if strings.Contains(line.Text, "DHCPACK") {
			parseHost(line.Text)
		}
	}
	fmt.Println(err)
}

func hello(w http.ResponseWriter, r *http.Request) {

	io.WriteString(w, "<!DOCTYPE html>")
	io.WriteString(w, "<html lang=\"en\">")
	io.WriteString(w, "<head>")
	io.WriteString(w, "<title>Wifi Stats</title>")
	io.WriteString(w, "<meta charset=\"utf-8\">")
	io.WriteString(w, "<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">")
	io.WriteString(w, "<link rel=\"stylesheet\" href=\"http://maxcdn.bootstrapcdn.com/bootstrap/3.3.5/css/bootstrap.min.css\">")
	io.WriteString(w, "<script src=\"https://ajax.googleapis.com/ajax/libs/jquery/1.11.3/jquery.min.js\"></script>")
	io.WriteString(w, "<script src=\"http://maxcdn.bootstrapcdn.com/bootstrap/3.3.5/js/bootstrap.min.js\"></script>")
	io.WriteString(w, "</head>")
	io.WriteString(w, "<body>")

	io.WriteString(w, "<div class=\"container\">")

	io.WriteString(w, "<h2>Hosts connected</h2>")
	const initTablec = `
	<table class="table table-striped">
<thead>
<tr>
<th>Connected</th>
<th>Interface</th>
<th>Hostname</th>
<th>Duration</th>
</tr>
</thead>
<tbody>
`
	io.WriteString(w, initTablec)

	for _, client := range client_map {

		io.WriteString(w, "<tr>")
		io.WriteString(w, "<td>"+client.DtConn.Format("02/01 15:04")+"</td>")
		io.WriteString(w, "<td>"+wlanToHuman(client.Wlan)+"</td>")
		io.WriteString(w, "<td>"+getHostName(client.Mac)+"</td>")
		io.WriteString(w, "<td>"+calcDate(client.DtConn, time.Now())+"</td>")
		io.WriteString(w, "</tr>")
	}

	io.WriteString(w, "</tbody>")
	io.WriteString(w, "</table>")

	io.WriteString(w, "<h2>Hosts history</h2>")

	const initTableh = `
<table class="table table-striped">
<thead>
<tr>
<th>Connected</th>
<th>Disconnected</th>
<th>Interface</th>
<th>Hostname</th>
<th>Duration</th>
</tr>
</thead>
<tbody>
`

	io.WriteString(w, initTableh)

	for _, l := range logs {
		io.WriteString(w, "<tr>")

		// currentTime.Format("01-02 15:00")
		//
		// strTimeConn := fmt.Sprintf("%d/%d %d:%d", l.DtConn..Format("01-02 15:00")Day(), l.DtConn.Month(), l.DtConn.Hour(), l.DtConn.Minute())
		// strTimeDisc := fmt.Sprintf("%d/%d %d:%d", l.DtDisc.Day(), l.DtDisc.Month(), l.DtDisc.Hour(), l.DtDisc.Minute())

		io.WriteString(w, "<td>"+l.DtConn.Format("02/01 15:04")+"</td>")
		io.WriteString(w, "<td>"+l.DtDisc.Format("02/01 15:04")+"</td>")
		io.WriteString(w, "<td>"+l.Wlan+"</td>")
		io.WriteString(w, "<td>"+getHostName(l.Mac)+"</td>")
		io.WriteString(w, "<td>"+l.Duration+"</td>")
		io.WriteString(w, "</tr>")
	}

	io.WriteString(w, "</tbody>")
	io.WriteString(w, "</table>")
	io.WriteString(w, "</div>")

	io.WriteString(w, "</body>")
	io.WriteString(w, "</html>")

}

func main() {
	go periodicDump()
	go readDhcp()
	go readHostAp()

	http.HandleFunc("/", hello)
	http.ListenAndServe(":3000", nil)
}
