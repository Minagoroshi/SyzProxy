package main

import (
	"github.com/Minagoroshi/SyzProxy"
	"io/ioutil"
	"log"
	"time"
)

func main() {
	proxyClient, err := SyzProxy.ClientFromProxy(SyzProxy.ReturnProxy("68.1.210.163", 4145, "", ""), "socks5")
	if err != nil {
		log.Fatal(err)
	}
	proxyClient.Timeout = 10 * time.Second

	resp, err := proxyClient.Get("https://httpbin.org/ip")
	if err != nil {
		log.Fatalln("Error getting response: ", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Error reading body: ", err)
	}

	log.Println(string(body))
}
