package main

import (
	"io"
	"net/http"
	"time"
)

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

	for _, l := range readLogs() {
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
