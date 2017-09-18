package server

import (
	"context"
	"crypto/tls"
	"net"

	"github.com/motki/motkid/log"
	"github.com/motki/motkid/model"
	"github.com/motki/motkid/model/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type Server interface {
	proto.AuthenticationServiceServer
	proto.ProductServiceServer
	proto.MarketPriceServiceServer
	proto.InfoServiceServer

	Serve() error
	Shutdown() error
}

var (
	_ Server = &GRPCServer{}
)

type GRPCServer struct {
	config model.Config

	model  *model.Manager
	logger log.Logger
	grpc   *grpc.Server

	server net.Listener
	local  net.Listener
}

func New(conf model.Config, m *model.Manager, l log.Logger) (Server, error) {
	srv := &GRPCServer{config: conf, model: m, logger: l, grpc: grpc.NewServer()}
	proto.RegisterAuthenticationServiceServer(srv.grpc, srv)
	proto.RegisterProductServiceServer(srv.grpc, srv)
	proto.RegisterMarketPriceServiceServer(srv.grpc, srv)
	proto.RegisterInfoServiceServer(srv.grpc, srv)
	return srv, nil
}

func (srv *GRPCServer) Shutdown() error {
	srv.grpc.GracefulStop()
	srv.server = nil
	srv.local = nil
	return nil
}

func (srv *GRPCServer) Serve() error {
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
	if srv.config.Kind == model.BackendLocalGRPC {
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

func (srv *GRPCServer) Authenticate(ctx context.Context, req *proto.AuthenticateRequest) (resp *proto.AuthenticateResponse, err error) {
	defer func() {
		if err != nil {
			resp = &proto.AuthenticateResponse{
				Result: errorResult(err),
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
