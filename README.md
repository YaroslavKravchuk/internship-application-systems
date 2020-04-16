# Cloudflare Internship Application: Systems
4/15/20
This is an CLI application I have built which takes a web address or an Ipv4 IP address and pings messages to the target, keeping track of package loss and RTT time.
This project was written in Golang.
Currently supports Ipv4 only

## How to run

```
go build Yping.go
```
#### Input an address:
```
sudo ./Yping www.google.com
```

#### or input an IP address:
```
sudo ./Yping 8.8.8.8
```