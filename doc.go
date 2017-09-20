// Package motkid is an EVE Online corporation management tool.
package motkid

//go:generate go-bindata -prefix "./public" -pkg template -tags "release" -ignore .DS_Store -o "./http/template/bindata_release.go" ./views/...
//go:generate go-bindata -prefix "./views" -pkg assets -tags "release" -ignore .DS_Store -o "./http/module/assets/bindata_release.go" ./public/fonts/... ./public/images/ ./public/scripts/... ./public/styles/... ./public/
//go:generate protoc -I proto/ -I vendor/github.com/golang/protobuf/ptypes/timestamp/ proto/motki.proto proto/model.proto proto/evedb.proto --go_out=plugins=grpc:proto
