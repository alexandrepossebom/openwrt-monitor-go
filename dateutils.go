package main

import (
	"fmt"
	"time"
)

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
