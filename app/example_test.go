package app_test

import (
	"fmt"

	"github.com/motki/core/app"
	"github.com/motki/core/log"
	"github.com/motki/core/proto"
	"github.com/motki/core/proto/client"
)

// ExampleNewClientEnv shows the bare-minimum to connect to the public MOTKI
// application and start an interactive CLI session.
//
// Under go test and godoc sandboxes this example will always fail to create
// the ClientEnv because there is no network available.
func ExampleNewClientEnv() {
	conf := &app.Config{
		Backend: proto.Config{
			Kind:       proto.BackendRemoteGRPC,
			RemoteGRPC: proto.RemoteConfig{ServerAddr: "motki.org"},
		},
		Logging: log.Config{
			Level: "debug",
		},
	}

	// Create the application environment.
	env, err := app.NewClientEnv(conf)
	if err != nil {
		if err == client.ErrBadCredentials {
			fmt.Println("Invalid username or password.")
		}
		panic("motki: error initializing application environment: " + err.Error())
	}

	// This method call will always fail when ran under the godoc sandbox without network access.
	_, err = env.Client.GetCorporation(98513229)
	if err != nil {
		fmt.Println("motki: error getting corporation: " + err.Error())
	}

	// Output:
	// motki: error getting corporation: rpc error: code = Unavailable desc = all SubConns are in TransientFailure
}
