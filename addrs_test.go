package addrs_test

import (
	"flag"
	"fmt"
	"net"
	"sync"
	"testing"

	"github.com/facebookgo/flag.addrs"
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

func TestFlagOneInvalidNetwork(t *testing.T) {
	t.Parallel()
	name := genName()
	var a1 net.Addr
	addrs.FlagOneVar(&a1, name, "", "")

	const network = "foo"
	const addr = "127.0.0.1:1234"

	err := flag.Set(name, network+":"+addr)
	if err == nil {
		t.Fatal("was expecting an error")
	}
	if err.Error() != "unknown network foo" {
		t.Fatal("did not get expected error, got", err)
	}
}

func TestFlagOneDefaultValue(t *testing.T) {
	t.Parallel()
	name := genName()
	var a1 net.Addr
	const network = "udp"
	const addr = "127.0.0.1:1234"
	addrs.FlagOneVar(&a1, name, network+":"+addr, "")
	if a1.Network() != network {
		t.Fatal("did not find expected network")
	}
	if a1.String() != addr {
		t.Fatal("did not find expected addr")
	}
}

func TestFlagOneInvalidDefaultValue(t *testing.T) {
	t.Parallel()
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("was expecting panic")
		}
	}()
	name := genName()
	var a1 net.Addr
	const network = "foo"
	const addr = "127.0.0.1:1234"
	addrs.FlagOneVar(&a1, name, network+":"+addr, "")
}

func TestFlagOneInvalidFormat(t *testing.T) {
	t.Parallel()
	name := genName()
	var a1 net.Addr
	addrs.FlagOneVar(&a1, name, "", "")

	err := flag.Set(name, "foo")
	if err == nil {
		t.Fatal("was expecting an error")
	}
	if err.Error() != `invalid address format, must be "net:host:port": foo` {
		t.Fatal("did not get expected error, got", err)
	}
}

func TestFlagMany(t *testing.T) {
	t.Parallel()
	name := genName()
	var a1 []net.Addr
	addrs.FlagManyVar(&a1, name, "", "")

	const network0 = "tcp"
	const addr0 = "127.0.0.1:1234"
	const network1 = "tcp"
	const addr1 = "127.0.0.1:5678"

	if err := flag.Set(name, network0+":"+addr0+","+network1+":"+addr1); err != nil {
		t.Fatal(err)
	}
	if a1[0].Network() != network0 {
		t.Fatal("did not find expected network")
	}
	if a1[0].String() != addr0 {
		t.Fatal("did not find expected addr")
	}
	if a1[1].Network() != network1 {
		t.Fatal("did not find expected network")
	}
	if a1[1].String() != addr1 {
		t.Fatal("did not find expected addr")
	}
}

func TestFlagManyDefaultValue(t *testing.T) {
	t.Parallel()
	name := genName()
	var a1 []net.Addr
	const network0 = "tcp"
	const addr0 = "127.0.0.1:1234"
	const network1 = "tcp"
	const addr1 = "127.0.0.1:5678"
	addrs.FlagManyVar(&a1, name, network0+":"+addr0+","+network1+":"+addr1, "")

	if a1[0].Network() != network0 {
		t.Fatal("did not find expected network")
	}
	if a1[0].String() != addr0 {
		t.Fatal("did not find expected addr")
	}
	if a1[1].Network() != network1 {
		t.Fatal("did not find expected network")
	}
	if a1[1].String() != addr1 {
		t.Fatal("did not find expected addr")
	}
}

func TestFlagManyInvalidDefaultValue(t *testing.T) {
	t.Parallel()
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("was expecting panic")
		}
	}()
	name := genName()
	var a1 []net.Addr
	const network = "foo"
	const addr = "127.0.0.1:1234"
	addrs.FlagManyVar(&a1, name, network+":"+addr, "")
}
