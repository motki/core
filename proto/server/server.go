// Package server contains an implementation of the MOTKI GRPC server.
//
// Much of the Server interface is generated using the protocol buffer definitions in
// the proto package.
package server

import (
	"crypto/tls"
	"net"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/motki/motki/eveapi"
	"github.com/motki/motki/evedb"
	"github.com/motki/motki/log"
	"github.com/motki/motki/model"
	"github.com/motki/motki/proto"
)

var ErrBadCredentials = errors.New("username or password is incorrect")

// A Server represents the raw interface for a MOTKI protobuf server.
type Server interface {
	proto.AuthenticationServiceServer
	proto.ProductServiceServer
	proto.MarketPriceServiceServer
	proto.InfoServiceServer
	proto.EveDBServiceServer
	proto.CorporationServiceServer
	proto.InventoryServiceServer

	// Serve opens a listening socket for the GRPC server.
	Serve() error
	// Shutdown attempts to gracefully shutdown the GRPC server.
	Shutdown() error
}

// Ensure grpcServer implements the Server interface.
var _ Server = &grpcServer{}

type grpcServer struct {
	config proto.Config

	model  *model.Manager
	evedb  *evedb.EveDB
	eveapi *eveapi.EveAPI
	logger log.Logger

	grpc *grpc.Server

	server net.Listener
	local  net.Listener
}

// New creates a new Server using the given configuration and dependencies.
func New(conf proto.Config, m *model.Manager, edb *evedb.EveDB, api *eveapi.EveAPI, l log.Logger) (Server, error) {
	srv := &grpcServer{config: conf, model: m, evedb: edb, eveapi: api, logger: l, grpc: grpc.NewServer()}
	proto.RegisterAuthenticationServiceServer(srv.grpc, srv)
	proto.RegisterProductServiceServer(srv.grpc, srv)
	proto.RegisterMarketPriceServiceServer(srv.grpc, srv)
	proto.RegisterInfoServiceServer(srv.grpc, srv)
	proto.RegisterEveDBServiceServer(srv.grpc, srv)
	proto.RegisterCorporationServiceServer(srv.grpc, srv)
	proto.RegisterInventoryServiceServer(srv.grpc, srv)
	return srv, nil
}

func (srv *grpcServer) Shutdown() error {
	srv.grpc.GracefulStop()
	srv.server = nil
	srv.local = nil
	return nil
}

func (srv *grpcServer) Serve() error {
	if srv.config.ServerEnabled {
		srv.logger.Debugf("grpc server: listening on %s", srv.config.ServerGRPC.ServerAddr)
		if srv.config.ServerGRPC.InsecureSkipVerify {
			srv.logger.Warnf("insecure_skip_verify_ssl is enabled, grpc client will not verify server certificate")
		}
		tlsConf, err := srv.config.ServerGRPC.TLSConfig()
		if err != nil {
			return err
		}
		lis, err := tls.Listen("tcp", srv.config.ServerGRPC.ServerAddr, tlsConf)
		if err != nil {
			return err
		}
		srv.server = lis
		go func() {
			// TODO: Close returns an error
			defer lis.Close()
			err := srv.grpc.Serve(lis)
			if err != nil {
				srv.logger.Warnf("grpc server return with: %s", err.Error())
			}
		}()
	}
	if srv.config.Kind == proto.BackendLocalGRPC {
		srv.logger.Debugf("grpc server: starting local listener")
		lis := srv.config.LocalGRPC.Listener
		if lis == nil {
			return errors.New("expected listener, got nil in local grpc config")
		}
		srv.local = lis
		go func() {
			// TODO: Close returns an error
			defer lis.Close()
			err := srv.grpc.Serve(lis)
			if err != nil {
				srv.logger.Warnf("grpc server return with: %s", err.Error())
			}
		}()
	}
	return nil
}

var successResult = &proto.Result{Status: proto.Status_SUCCESS}

func errorResult(err error) *proto.Result {
	return &proto.Result{Status: proto.Status_FAILURE, Description: err.Error()}
}

func (srv *grpcServer) Authenticate(ctx context.Context, req *proto.AuthenticateRequest) (resp *proto.AuthenticateResponse, err error) {
	defer func() {
		if err != nil {
			resp = &proto.AuthenticateResponse{
				Result: errorResult(ErrBadCredentials),
			}
			err = nil
		}
	}()
	_, tok, err := srv.model.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	return &proto.AuthenticateResponse{Result: successResult, Token: &proto.Token{tok}}, nil
}
