package SyzProxy

import (
	"bufio"
	"errors"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// Todo: add support for socks
// Todo: add IsProxyAlive value to Proxy and its functionality

type ProxyManager struct {
	ProxyList []Proxy
	ProxyType string
}

type Proxy struct {
	host     string
	port     int
	username string
	password string
}

// GetRandomTransport returns a transport using the GetRandomProxy function
func (pm *ProxyManager) GetRandomTransport() (*http.Transport, error) {
	return TransportFromProxy(pm.GetRandomProxy())
}

// GetRandomProxy returns a random proxy from the list
func (pm *ProxyManager) GetRandomProxy() Proxy {
	if len(pm.ProxyList) == 0 {
		return Proxy{}
	} else if len(pm.ProxyList) == 1 {
		return pm.ProxyList[0]
	}
	var proxy Proxy
	rand.Seed(time.Now().UnixNano())
	proxy = pm.ProxyList[rand.Intn(len(pm.ProxyList)-1)]
	return proxy
}

// LoadFromFile loads a list of proxies from a file
func (pm *ProxyManager) LoadFromFile(filename string, proxyType string) (int, error) {
	proxyType = strings.ToLower(proxyType)
	if proxyType != "http" && proxyType != "https" && proxyType != "socks5" && proxyType != "socks4" && proxyType != "socks4a" {
		if strings.Contains(proxyType, "socks") {
			return 0, errors.New("Unsupported proxy type: " + proxyType)
		}
		return 0, errors.New("Invalid proxy type")
	}

	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer func(file *os.File) error {
		err := file.Close()
		if err != nil {
			return err
		}
		return nil
	}(file)

	pm.ProxyList = []Proxy{}
	pm.ProxyType = proxyType
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Count(line, ":") == 3 {
			proxySplit := strings.Split(line, ":")
			port, err := strconv.Atoi(proxySplit[1])
			if err != nil {
				return 0, err
			}
			pm.ProxyList = append(pm.ProxyList, ReturnProxy(proxySplit[0], port, proxySplit[2], proxySplit[3]))

		} else {
			proxySplit := strings.Split(line, ":")
			port, err := strconv.Atoi(proxySplit[1])
			if err != nil {
				return 0, err
			}
			pm.ProxyList = append(pm.ProxyList, ReturnProxy(proxySplit[0], port, "", ""))
		}
	}
	return len(pm.ProxyList), nil
}

// ReturnProxy returns a filled Proxy Struct
// If the username and password are empty, it will return a proxy without authentication
func ReturnProxy(host string, port int, username string, password string) Proxy {
	return Proxy{
		host:     host,
		port:     port,
		username: username,
		password: password,
	}
}

// TranportFromProxy returns a http.Transport with the proxy set
func TransportFromProxy(proxy Proxy) (*http.Transport, error) {

	// Validate the host url
	_, err := url.Parse("http://" + proxy.host + ":" + strconv.Itoa(proxy.port))
	if err != nil {
		return nil, err
	}
	if proxy.username == "" && proxy.password == "" {
		return &http.Transport{
			Proxy: http.ProxyURL(&url.URL{
				Scheme: "http",
				Host:   proxy.host + ":" + strconv.Itoa(proxy.port),
			}),
		}, nil
	} else {
		return &http.Transport{
			Proxy: http.ProxyURL(&url.URL{
				Scheme: "http",
				Host:   proxy.host + ":" + strconv.Itoa(proxy.port),
				User:   url.UserPassword(proxy.username, proxy.password),
			}),
		}, nil
	}
}

// ClientFromProxy returns a http.Client with the transport set
func ClientFromProxy(proxy Proxy) (*http.Client, error) {
	transport, err := TransportFromProxy(proxy)
	if err != nil {
		return nil, err
	}

	return &http.Client{
		Transport: transport,
	}, nil
}

// ClientFromTransport returns a http.Client with the transport set
func ClientFromTransport(transport *http.Transport) *http.Client {
	return &http.Client{
		Transport: transport,
	}
}
