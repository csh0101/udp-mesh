package main

type Config struct {
	RemoteAddr string
	RemotePort int
	LocalAddr  string
	LocalPort  int
}

var (
	DefaultBroadCast = Config{
		RemoteAddr: "255.255.255.255",
		RemotePort: 18888,
		LocalAddr:  "",
		LocalPort:  -1,
	}
)
