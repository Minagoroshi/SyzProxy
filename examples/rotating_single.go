package main

import (
	"github.com/Minagoroshi/SyzProxy"
	"io/ioutil"
	"log"
	"time"
)

func main() {

	// If your proxy does not require authentication leave the username and password empty ("")

	proxyClient, err := SyzProxy.ClientFromProxy(SyzProxy.ReturnProxy("proxy.example.host", 12345, "username", "password"), "http")
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
