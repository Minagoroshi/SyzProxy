package SyzProxy

import (
	"bufio"
	"errors"
	"fmt"
	"h12.io/socks"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// Todo: add IsProxyAlive value to Proxy and its functionality

var (
	ProxyTypes = []string{"http", "socks5", "socks4", "socks4a"}
)

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

// NewProxyManager returns a new ProxyManager
func NewProxyManager() ProxyManager {
	return ProxyManager{}
}

// GetRandomTransport returns a transport using the GetRandomProxy function
func (pm *ProxyManager) GetRandomTransport() (*http.Transport, error) {
	proxyType := strings.ToLower(pm.ProxyType)
	switch proxyType {
	case "http":
		return TransportFromProxy(pm.GetRandomProxy(), proxyType)
	case "socks5":
		return TransportFromProxy(pm.GetRandomProxy(), proxyType)
	case "socks4":
		return TransportFromProxy(pm.GetRandomProxy(), proxyType)
	case "socks4a":
		return TransportFromProxy(pm.GetRandomProxy(), proxyType)
	default:
		return nil, errors.New("Invalid proxy type")
	}
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
// Proxy types are http, socks4, socks4a, socks5
func (pm *ProxyManager) LoadFromFile(filename string, proxyType string) (int, error) {
	proxyType = strings.ToLower(proxyType)
	if !arrayContains(ProxyTypes, proxyType) {
		return 0, errors.New("Invalid proxy type")
	} else if proxyType == "https" {
		proxyType = "http"
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
// Proxy types are http, socks4, socks4a, socks5
func TransportFromProxy(proxy Proxy, proxyType string) (*http.Transport, error) {

	// Validate the host url
	_, err := url.Parse("http://" + proxy.host + ":" + strconv.Itoa(proxy.port))
	if err != nil {
		return nil, err
	}

	proxyType = strings.ToLower(proxyType)
	switch proxyType {
	case "http":
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
	case "socks4":
		userpass := proxy.username + ":" + proxy.password + "@"
		if proxy.username == "" && proxy.password == "" {
			userpass = ""
		}
		dialSocksProxy := socks.Dial(fmt.Sprintf("socks4://%s%s:%s?timeout=10s", userpass, proxy.host, strconv.Itoa(proxy.port)))
		return &http.Transport{Dial: dialSocksProxy}, nil
	case "socks4a":
		userpass := proxy.username + ":" + proxy.password + "@"
		if proxy.username == "" && proxy.password == "" {
			userpass = ""
		}
		dialSocksProxy := socks.Dial(fmt.Sprintf("socks4a://%s%s:%s?timeout=10s", userpass, proxy.host, strconv.Itoa(proxy.port)))
		return &http.Transport{Dial: dialSocksProxy}, nil
	case "socks5":
		userpass := proxy.username + ":" + proxy.password + "@"
		if proxy.username == "" && proxy.password == "" {
			userpass = ""
		}
		dialSocksProxy := socks.Dial(fmt.Sprintf("socks5://%s%s:%s?timeout=10s", userpass, proxy.host, strconv.Itoa(proxy.port)))
		return &http.Transport{Dial: dialSocksProxy}, nil
	default:
		return nil, errors.New("Invalid proxy type")
	}
}

// ClientFromProxy returns a http.Client with the transport set
// Proxy types are http, socks4, socks4a, socks5
func ClientFromProxy(proxy Proxy, proxyType string) (*http.Client, error) {
	if !arrayContains(ProxyTypes, proxyType) {
		return nil, errors.New("Invalid proxy type")
	}
	proxyType = strings.ToLower(proxyType)

	transport, err := TransportFromProxy(proxy, proxyType)
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

func arrayContains(array []string, value string) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}
