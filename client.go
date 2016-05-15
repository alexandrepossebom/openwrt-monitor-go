package main

import (
	"time"
)

type Client struct {
	DtConn time.Time
	Wlan   string
	Mac    string
}

var client_map = map[string]Client{}
