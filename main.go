package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "net/url"
    "time"

    "h12.io/socks"
    "sync"
)

var threads = 0
var wg sync.WaitGroup

func main() {
    fmt.Print("\033[H\033[2J")
    go func() {
        for {
            fmt.Println(threads)
            time.Sleep(time.Second)
        }
    }()
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        wg.Add(1)
        threads++
        proxy := scanner.Text()
        go func() {
            defer func() {
                threads--
            }()
            checkProxy(proxy, os.Args[1])
        }()
    }
    wg.Wait()
}


func checkProxy(proxy string, port string) {
    defer wg.Done()

    httpClient := &http.Client{Timeout: time.Second * 8}
    httpClient.Transport = &http.Transport{}
    proxyandport := fmt.Sprintf("%s:%s", proxy, port)
    defer httpClient.Transport.(*http.Transport).CloseIdleConnections()
    var protocols = []string{"http", "socks4", "socks5"}
    for _, protocol := range protocols {
        url := "http://ip-api.com/json/"
        proxyUrl := fmt.Sprintf("%s://%s:%s", protocol, proxy, port)
        switch protocol {
        case "http":
            httpClient.Transport = &http.Transport{Proxy: http.ProxyURL(MustParseURL(proxyUrl))}
        case "socks4":
            httpClient.Transport = &http.Transport{Dial: socks.Dial(proxyUrl)}
        case "socks5":
            httpClient.Transport = &http.Transport{Dial: socks.Dial(proxyUrl)}
        }
        resp, err := httpClient.Get(url)
        if err != nil {
            //fmt.Println(err)
            continue
        }
        defer resp.Body.Close()

        if resp.StatusCode == 200 {
            var data map[string]interface{}
            bodyBytes, err := ioutil.ReadAll(resp.Body)
            if err != nil {
                //fmt.Println(err)
                continue
            }
            if err := json.Unmarshal(bodyBytes, &data); err != nil {
                //fmt.Println(err)
                continue
            }
            if data["status"] == "success" {
                switch protocol {
                case "http":
                    writeToFile("http.txt", proxyandport)
                    fmt.Printf("%s://%s:%s - %s\n" , protocol, proxy, port, data["as"])
                case "socks4":
                    writeToFile("socks4.txt", proxyandport)
                    fmt.Printf("%s://%s:%s - %s\n" , protocol, proxy, port, data["as"])
                case "socks5":
                    writeToFile("socks5.txt", proxyandport)
                    fmt.Printf("%s://%s:%s - %s\n" , protocol, proxy, port, data["as"])
                }
            }
        }
    }
}


func writeToFile(filename string, content string) {
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer file.Close()

    if _, err := file.WriteString(content + "\n"); err != nil {
        fmt.Println(err)
    }
}

func MustParseURL(rawurl string) *url.URL {
    u, err := url.Parse(rawurl)
    if err != nil {
        panic(err)
    }
    return u
}
