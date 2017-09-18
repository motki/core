package client

import (
	"github.com/pkg/errors"

	"fmt"

	"github.com/motki/motkid/log"
	"github.com/motki/motkid/model"
)

var ErrNotAuthenticated = errors.New("not authenticated")

type Client interface {
	Authenticate(username, password string) (string, error)

	CharacterForRole(model.Role) (*model.Character, error)
	GetCharacter(charID int) (*model.Character, error)
	GetCorporation(corpID int) (*model.Corporation, error)
	GetAlliance(allianceID int) (*model.Alliance, error)

	NewProduct(typeID int) (*model.Product, error)
	GetProduct(productID int) (*model.Product, error)
	SaveProduct(product *model.Product) error
	GetProducts() ([]*model.Product, error)
	UpdateProductPrices(*model.Product) (*model.Product, error)
}

func New(conf model.Config, logger log.Logger) (Client, error) {
	var cl Client
	var err error
	switch conf.Kind {
	case model.BackendLocalGRPC:
		logger.Debugf("grpc client: initializing local client.")
		cl, err = NewLocalGRPC(conf.LocalGRPC.Listener, logger)
		if err != nil {
			return nil, errors.Wrap(err, "app: unable to initialize backend")
		}
		fmt.Println("doneeee")

	case model.BackendRemoteGRPC:
		logger.Debugf("grpc client: initializing remote client, server address: %s", conf.RemoteGRPC.ServerAddr)
		conf := conf.RemoteGRPC
		if conf.InsecureSkipVerify {
			logger.Warnf("insecure_skip_verify_ssl is enabled, grpc client will not verify server certificate")
		}
		tc, err := conf.TLSConfig()
		if err != nil {
			return nil, errors.Wrap(err, "app: unable to initialize backend")
		}
		cl, err = NewRemoteGRPC(conf.ServerAddr, logger, tc)
		if err != nil {
			return nil, errors.Wrap(err, "app: unable to initialize backend")
		}

	default:
		return nil, errors.Errorf("unsupported backend kind %s", conf.Kind)
	}
	return cl, nil
}
