package main

import (
	"time"
)

type Log struct {
	DtConn   time.Time
	DtDisc   time.Time
	Wlan     string
	Mac      string
	Duration string
}
