// Package core contains libraries that provide various functionality, shared by all MOTKI applications.
//
// This package is made up of subpackages, and currently serves informational purposes only.
package core // import "github.com/motki/core"

//go:generate protoc -I proto/ proto/motki.proto proto/model.proto proto/evedb.proto --go_out=plugins=grpc:proto
