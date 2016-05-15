package main

import (
	_ "github.com/mattn/go-sqlite3"
	"time"
)

func updateHost(host, mac string) {
	stmt, _ := db.Prepare("update hosts set hostname=? where mac=?")
	stmt.Exec(host, mac)
	stmt.Close()
}

func addHost(host, mac string) {
	stmt, _ := db.Prepare("INSERT INTO hosts(mac, hostname) values(?,?)")
	stmt.Exec(mac, host)
	stmt.Close()
}

func addLog(log Log) {
	stmt, _ := db.Prepare("INSERT INTO history VALUES(?,?,?,?,?)")
	stmt.Exec(log.DtConn, log.DtDisc, log.Wlan, log.Mac, log.Duration)
	stmt.Close()
}

func readLogs() []Log {
	var logs []Log
	rows, err := db.Query("SELECT * FROM history where duration like '%m%' order by disconnected limit 10")

	for rows.Next() {
		var connected time.Time
		var disconnected time.Time
		var wlan string
		var mac string
		var duration string
		err = rows.Scan(&connected, &disconnected,&wlan,&mac,&duration)
		if err == nil {
			l := Log{connected, disconnected, wlan, mac, duration}
			logs = append(logs, l)
		}
	}
	return logs
}

func readHosts() {
	rows, err := db.Query("SELECT mac,hostname FROM hosts")

	for rows.Next() {
		var mac string
		var hostname string
		err = rows.Scan(&mac, &hostname)
		if err != nil {
			host_map[mac] = hostname
		}
	}
}
