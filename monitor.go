package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
)

var db, _ = sql.Open("sqlite3", "./foo.db")

func periodicDump() {
	for _ = range time.Tick(30 * time.Second) {
		dump()
	}
}

func parseDate(str string) time.Time {
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		fmt.Println(err)
		return time.Now()
	}
	return t
}

func lp(l Log) {
	//fmt.Printf("%s [<---] %s %-15s %s\n", l.DtDisc.Format(time.RFC822), l.Wlan, getHostName(l.Mac), l.Duration)
	addLog(l)
}

func connected(mac string, wlan string, dt time.Time) {
	disconnected(mac, wlan, dt)
	client_map[mac+wlan] = Client{dt, wlan, mac}
	//	fmt.Printf("%s [--->] %s %-15s\n", dt.Format(time.RFC822), wlanToHuman(wlan), getHostName(mac))
}

func disconnected(mac string, wlan string, dt time.Time) {
	if _, ok := client_map[mac+wlan]; ok {
		client := client_map[mac+wlan]
		l := Log{client.DtConn, dt, wlanToHuman(client.Wlan), client.Mac, calcDate(client.DtConn, dt)}
		lp(l)
		delete(client_map, mac+wlan)
	}
}

func dump() {
	for _, client := range client_map {
		fmt.Printf("%s [DUMP] %s %-15s %s\n", time.Now().Format(time.RFC822), wlanToHuman(client.Wlan), getHostName(client.Mac), calcDate(client.DtConn, time.Now()))
	}
}

func main() {
	readHosts()
	go periodicDump()
	go readDhcp()
	go readHostAp()

	http.HandleFunc("/", hello)
	http.ListenAndServe(":3000", nil)
	defer db.Close()
}
