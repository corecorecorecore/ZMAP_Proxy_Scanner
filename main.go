package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "net/url"
    "sync"
    "time"
    "strconv"

    "h12.io/socks"
)

func main() {
    fmt.Print("\033[H\033[2J")

    maxThreads, _ := strconv.Atoi(os.Args[1])
    sem := make(chan struct{}, maxThreads)

    go printThreads(sem)

    var wg sync.WaitGroup
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        proxy := scanner.Text()

        sem <- struct{}{}

        wg.Add(1)
        go func() {
            defer func() {
                <-sem
            }()
            checkProxy(proxy, &wg, os.Args[2])
        }()
    }
    wg.Wait()
}


func checkProxy(proxy string, wg *sync.WaitGroup, port string) {
    defer wg.Done()

    httpClient := &http.Client{Timeout: time.Second * 8}
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
                    writeToFile("http.txt", proxy)
                    fmt.Printf("%s://%s:%s - %s\n" , protocol, proxy, port, data["as"])
                case "socks4":
                    writeToFile("socks4.txt", proxy)
                    fmt.Printf("%s://%s:%s - %s\n" , protocol, proxy, port, data["as"])
                case "socks5":
                    writeToFile("socks5.txt", proxy)
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

func contains(arr []string, elem string) bool {
    for _, a := range arr {
        if a == elem {
            return true
        }
    }
    return false
}


func MustParseURL(rawurl string) *url.URL {
    u, err := url.Parse(rawurl)
    if err != nil {
        panic(err)
    }
    return u
}

func printThreads(sem chan struct{}) {
    for {
        fmt.Print("\033[H\033[2J") // clear the shell
        fmt.Printf("Running threads: %d\n", len(sem))
        time.Sleep(time.Second)
    }
}

