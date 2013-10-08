// Package addrs provides a flags to define one or an array of net.Addr.
package addrs

import (
	"flag"
	"fmt"
	"net"
	"strings"
)

type flagOne struct {
	addr *net.Addr
}

func (f *flagOne) Set(addr string) error {
	a, err := resolveAddr(addr)
	if err != nil {
		return err
	}
	*f.addr = a
	return nil
}

func (f *flagOne) String() string {
	if f.addr == nil || (*f.addr) == nil {
		return ""
	}
	return (*f.addr).Network() + ":" + (*f.addr).String()
}

// Set a single net.Addr by a flag.
func FlagOneVar(dest *net.Addr, name string, addr string, usage string) {
	if addr != "" {
		ra, err := resolveAddr(addr)
		if err != nil {
			panic(err)
		}
		*dest = ra
	}
	flag.Var(&flagOne{addr: dest}, name, usage)
}

func resolveAddr(addr string) (net.Addr, error) {
	parts := strings.Split(addr, ":")
	if len(parts) != 3 {
		return nil, fmt.Errorf(
			`invalid address format, must be "net:host:port": %s`,
			addr,
		)
	}

	hp := parts[1] + ":" + parts[2]
	switch parts[0] {
	default:
		return nil, net.UnknownNetworkError(parts[0])
	case "ip", "ip4", "ip6":
		return net.ResolveIPAddr(parts[0], hp)
	case "tcp", "tcp4", "tcp6":
		return net.ResolveTCPAddr(parts[0], hp)
	case "udp", "udp4", "udp6":
		return net.ResolveUDPAddr(parts[0], hp)
	case "unix", "unixgram", "unixpacket":
		return net.ResolveUnixAddr(parts[0], hp)
	}
}
