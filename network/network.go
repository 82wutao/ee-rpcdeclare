package network

import "fmt"

type HostPort struct {
	Host  string
	Port  int16
	Proto string // tcp/udp/http
}

func (hp HostPort) ProtoHostPort() string {
	return fmt.Sprintf("%s@%s:%d",
		hp.Proto, hp.Host, hp.Port)
}
func (hp HostPort) HostPort() string {
	return fmt.Sprintf("%s:%d",
		hp.Host, hp.Port)
}
func (hp HostPort) SchemaHostPort() string {
	return fmt.Sprintf("%s://%s:%d",
		hp.Proto, hp.Host, hp.Port)
}
