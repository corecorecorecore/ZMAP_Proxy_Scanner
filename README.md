# ZMAP_Proxy_Scanner
Supports HTTP/SOCKS4/SOCKS5

```go mod init scanner
go build
zmap -p port | ./scanner port threads
zmap -p 3128 | ./scanner 3128 1000
```
