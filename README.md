# SyzProxy
## A simple and modern proxy management library for Go

## Table of Contents
- [Installation](#installation)
- [Examples](#examples)
- [Documentation](#documentation)

## Installation
```bash
go get github.com/minagoroshi/syzproxy@v0.1.1
```

## Examples
### Single/Rotating Proxy
```go
package main

import (
	"github.com/Minagoroshi/SyzProxy"
	"io/ioutil"
	"log"
	"time"
)

func main() {

	// If your proxy does not require authentication leave the username and password empty ("")
	proxyClient, err := SyzProxy.ClientFromProxy(SyzProxy.ReturnProxy("proxy.example.host", 12345, "username", "password"))
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

```

### Proxy List
```go
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

```

## Documentation
https://pkg.go.dev/github.com/Minagoroshi/SyzProxy#section-readme

