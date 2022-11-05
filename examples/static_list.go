package main

import (
	"github.com/Minagoroshi/SyzProxy"
	"log"
	"net/http"
)

func main() {
	PManager := &SyzProxy.ProxyManager{}
	num, err := PManager.LoadFromFile("proxies.txt", "http")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(PManager.ProxyType)
	log.Println("Loaded", num, "proxies")

	transport, err := PManager.GetRandomTransport()
	if err != nil {
		log.Fatal(err)
	}
	proxyClient := http.Client{Transport: transport}

	resp, err := proxyClient.Get("https://httpbin.org/ip")
	if err != nil {
		log.Fatalln("Error getting response: ", err)
	}

	defer resp.Body.Close()

	log.Println(resp)
}
