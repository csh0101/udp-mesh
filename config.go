package main

type Config struct {
	LocalAddr string
	LocalPort int
}

func NewConfig(laddr, raddr string, lport, rport int) Config {
	return Config{
		LocalAddr: laddr,
		LocalPort: lport,
	}
}
