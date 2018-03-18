package client_test

import (
	"fmt"

	"github.com/motki/core/log"
	"github.com/motki/core/proto"
	"github.com/motki/core/proto/client"
)

// ExampleNew shows how one might use the client package directly.
func ExampleNew() {
	// This configuration connects to the public application server at motki.org.
	c := proto.Config{
		Kind:       proto.BackendRemoteGRPC,
		RemoteGRPC: proto.RemoteConfig{ServerAddr: "motki.org:18443"},
	}
	// Gotta have a logger.
	l := log.New(log.Config{Level: "warn"})

	// Create the client or panic.
	cl, err := client.New(c, l)
	if err != nil {
		panic(err)
	}

	// This method call will always fail when ran under the godoc sandbox without network access.
	it, err := cl.GetItemType(2281)
	if err != nil {
		panic(err)
	}

	fmt.Println(it.Name)

	// Output:
	// Adaptive Invulnerability Field II
}
