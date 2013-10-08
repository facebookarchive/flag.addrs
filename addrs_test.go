package addrs_test

import (
	"flag"
	"fmt"
	"net"
	"sync"
	"testing"

	"github.com/daaku/go.flag.addrs"
)

var (
	genNameCount int
	genNameMutex sync.Mutex
)

func genName() string {
	genNameMutex.Lock()
	defer genNameMutex.Unlock()
	defer func() { genNameCount++ }()
	return fmt.Sprintf("flag-%d", genNameCount)
}

func TestFlagOne(t *testing.T) {
	t.Parallel()
	name := genName()
	var a1 net.Addr
	addrs.FlagOneVar(&a1, name, "", "")

	const network = "udp"
	const addr = "127.0.0.1:1234"

	if err := flag.Set(name, network+":"+addr); err != nil {
		t.Fatal(err)
	}
	if a1.Network() != network {
		t.Fatal("did not find expected network")
	}
	if a1.String() != addr {
		t.Fatal("did not find expected addr")
	}
}
