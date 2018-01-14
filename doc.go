// Package motki contains libraries that integrate various functionality, shared by all MOTKI applications.
package motki

//go:generate protoc -I proto/ proto/motki.proto proto/model.proto proto/evedb.proto --go_out=plugins=grpc:proto
