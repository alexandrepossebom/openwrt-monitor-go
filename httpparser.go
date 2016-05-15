package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var host_map = map[string]string{}

func parseHostname(str string){
		dhcp :=  strings.Split(str, " ")

		mac := dhcp[1]
		hostname := dhcp[3]

		if len(hostname) <= 1{
			hostname = mac
		}
		host_map[mac] = hostname
		// fmt.Printf("mac: %s, hostname: %s\n", mac, hostname)
}

func parseClient(str string){
		client :=  strings.Split(strings.Split(str, "\"")[1]," ")[0]
		fmt.Printf("client: %s %s\n", client, host_map[client])
}

func main() {
	fmt.Println("init...")

	form := url.Values{}
	form.Add("password", "la01ks92")

	req, err := http.NewRequest("POST", "http://coxande.no-ip.info/utility/get_password_cookie.sh", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	cookies := strings.Replace(string(body), "Set-Cookie:", "", -1)
	cookies = strings.Replace(cookies, "Path=/;", "", -1)
	cookies = strings.Replace(cookies, "\n", "", -1)

	resp.Body.Close()

	if strings.Contains(cookies, "hash") {

		req, err = http.NewRequest("POST", "http://coxande.no-ip.info/hosts.sh", strings.NewReader(form.Encode()))
		req.Header.Set("Cookie", cookies)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		client = &http.Client{}
		resp, err = client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		body, _ = ioutil.ReadAll(resp.Body)

		for _, line := range strings.Split(string(body), "\n") {
			if strings.Contains(line, "dhcpLeaseLines.push") {
				parseHostname(line)
			}
			if strings.Contains(line, "wifiLines.push") {
				parseClient(line)
			}
		}

	}

}
