package client

import (
	"github.com/pkg/errors"

	"github.com/motki/motki/evedb"
	"github.com/motki/motki/log"
	"github.com/motki/motki/model"
)

var ErrNotAuthenticated = errors.New("not authenticated")
var ErrBadCredentials = errors.New("username or password is incorrect")

// A Client provides a remote interface to the motki model package.
//
// A Client connects to a process-local or remote GRPC server to facilitate
// remote procedure calls. When used with a remote GRPC server, the Client allows
// client applications to consume MOTKI and EVESDE data without storing anything
// on the local machine.
type Client interface {
	// Authenticate attempts to authenticate with the GRPC server, returning a session token if successful.
	//
	// If the authentication is successful, a new session token is returned. This token is used internally
	// by the client to authorize subsequent GRPC calls.
	Authenticate(username, password string) (string, error)

	// CharacterForRole returns the current session's associated character for the given role.
	CharacterForRole(model.Role) (*model.Character, error)
	// GetCharacter returns a populated Character for the given character ID.
	GetCharacter(charID int) (*model.Character, error)
	// GetCorporation returns a populated Corporation for the given corporation ID.
	GetCorporation(corpID int) (*model.Corporation, error)
	// GetAlliance returns a populated Alliance for the given alliance ID.
	GetAlliance(allianceID int) (*model.Alliance, error)

	// NewProduct creates a new Production Chain for the given type ID.
	//
	// If a production chain already exists for the given type ID, it will be returned.
	NewProduct(typeID int) (*model.Product, error)
	// GetProduct attempts to load an existing production chain using its unique product ID.
	GetProduct(productID int) (*model.Product, error)
	// SaveProduct attempts to save the given production chain to the backend database.
	SaveProduct(product *model.Product) error
	// GetProducts gets all production chains for the current session's corporation.
	GetProducts() ([]*model.Product, error)
	// UpdateProductPrices updates all items in a production chain with the latest market sell price.
	UpdateProductPrices(*model.Product) (*model.Product, error)

	// GetInventory returns all inventory items for the current session's corporation.
	GetInventory() ([]*model.InventoryItem, error)
	// NewInventoryItem creates a new inventory item for the given type ID and location ID.
	//
	// If an inventory item already exists for the given type and location ID, it will be returned.
	NewInventoryItem(typeID, locationID int) (*model.InventoryItem, error)
	// SaveInventoryItem attempts to save the given inventory item to the backend database.
	SaveInventoryItem(*model.InventoryItem) error

	// GetCorpBlueprints returns the current session's corporation's blueprints.
	GetCorpBlueprints() ([]*model.Blueprint, error)
	// GetMarketPrice returns the current market price for the given type ID.
	GetMarketPrice(typeID int) (*model.MarketPrice, error)
	// GetMarketPrices returns a slice of market prices for each of the given type IDs.
	GetMarketPrices(typeID int, typeIDs ...int) ([]*model.MarketPrice, error)

	// GetItemType returns information about the given type ID.
	GetItemType(typeID int) (*evedb.ItemType, error)
	// GetItemTypeDetail returns detailed information about the given type ID.
	GetItemTypeDetail(typeID int) (*evedb.ItemTypeDetail, error)

	// QueryItemTypes searches for matching item types given a search query and optional category IDs.
	QueryItemTypes(query string, catIDs ...int) ([]*evedb.ItemType, error)
	// QueryItemTypeDetails searches for matching item types, returning detail type information for each match.
	QueryItemTypeDetails(query string, catIDs ...int) ([]*evedb.ItemTypeDetail, error)
	// GetMaterialSheet returns manufacturing information about the given type ID.
	GetMaterialSheet(typeID int) (*evedb.MaterialSheet, error)

	// GetRace returns information about the given race ID.
	GetRace(raceID int) (*evedb.Race, error)
	// GetRaces returns information about all races in the EVE universe.
	GetRaces() ([]*evedb.Race, error)
	// GetBloodline returns information about the given bloodline ID.
	GetBloodline(bloodlineID int) (*evedb.Bloodline, error)
	// GetAncestry returns information about the given ancestry ID.
	GetAncestry(ancestryID int) (*evedb.Ancestry, error)

	// GetRegion returns information about the given region ID.
	GetRegion(regionID int) (*evedb.Region, error)
	// GetRegions returns a slice containing information about all regions in the EVE universe.
	GetRegions() ([]*evedb.Region, error)
	// GetConstellation returns information about the given constellation ID.
	GetConstellation(constellationID int) (*evedb.Constellation, error)
	// GetSystem returns information about the given system ID.
	GetSystem(systemID int) (*evedb.System, error)
}

// New creates a new Client using the given model configuration.
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
