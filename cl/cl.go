package cl

import (
	"net"
	"net/http"
	"net/http/cookiejar"
	"time"
	. "tools"
)

func Client(t time.Duration) http.Client {
	jar, err := cookiejar.New(nil)
	ErrHandle(err, `x`, `obtain_cookiejar`)
	return http.Client{&http.Transport{Dial: func(network string, address string) (net.Conn, error) {
		return net.DialTimeout(network, address, t*time.Millisecond)
	}}, nil, jar}
}
