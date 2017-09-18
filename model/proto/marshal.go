package proto

import (
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/motki/motkid/model"
	"github.com/shopspring/decimal"
)

func ProtoToCharacter(char *Character) *model.Character {
	return &model.Character{
		CharacterID:   int(char.Id),
		Name:          char.Name,
		CorporationID: int(char.CorporationId),
		AllianceID:    int(char.AllianceId),
		RaceID:        int(char.RaceId),
		AncestryID:    int(char.AncestryId),
		BloodlineID:   int(char.BloodlineId),
		Description:   char.Description,
		BirthDate:     time.Unix(char.BirthDate.Seconds, int64(char.BirthDate.Nanos)),
	}
}

func CharacterToProto(char *model.Character) *Character {
	return &Character{
		Id:            int64(char.CharacterID),
		Name:          char.Name,
		CorporationId: int64(char.CorporationID),
		AllianceId:    int64(char.AllianceID),
		AncestryId:    int32(char.AncestryID),
		RaceId:        int32(char.RaceID),
		BloodlineId:   int32(char.BloodlineID),
		BirthDate: &timestamp.Timestamp{
			Seconds: char.BirthDate.Unix(),
			Nanos:   int32(char.BirthDate.Nanosecond()),
		},
		Description: char.Description,
	}
}

func ProtoToCorporation(corp *Corporation) *model.Corporation {
	return &model.Corporation{
		Name:          corp.Name,
		CorporationID: int(corp.Id),
		AllianceID:    int(corp.AllianceId),
		CreationDate:  time.Unix(corp.CreationDate.Seconds, int64(corp.CreationDate.Nanos)),
		Description:   corp.Description,
		Ticker:        corp.Ticker,
	}
}

func CorporationToProto(corp *model.Corporation) *Corporation {
	return &Corporation{
		Id:         int64(corp.CorporationID),
		Name:       corp.Name,
		AllianceId: int64(corp.AllianceID),
		Ticker:     corp.Ticker,
		CreationDate: &timestamp.Timestamp{
			Seconds: corp.CreationDate.Unix(),
			Nanos:   int32(corp.CreationDate.Nanosecond()),
		},
		Description: corp.Description,
	}
}

func ProtoToAlliance(alliance *Alliance) *model.Alliance {
	return &model.Alliance{
		AllianceID:  int(alliance.Id),
		Name:        alliance.Name,
		Ticker:      alliance.Ticker,
		DateFounded: time.Unix(alliance.DateFounded.Seconds, int64(alliance.DateFounded.Nanos)),
	}
}

func AllianceToProto(alliance *model.Alliance) *Alliance {
	return &Alliance{
		Id:     int64(alliance.AllianceID),
		Name:   alliance.Name,
		Ticker: alliance.Ticker,
		DateFounded: &timestamp.Timestamp{
			Seconds: alliance.DateFounded.Unix(),
			Nanos:   int32(alliance.DateFounded.Nanosecond()),
		},
	}
}

func ProtoToProduct(m *Product) *model.Product {
	kind := model.ProductBuild
	if m.Kind == Product_BUY {
		kind = model.ProductBuy
	}
	prod := &model.Product{
		ProductID:          int(m.Id),
		TypeID:             int(m.TypeId),
		Quantity:           int(m.Quantity),
		MarketPrice:        decimal.NewFromFloat(m.MarketPrice),
		MarketRegionID:     int(m.MarketRegionId),
		MaterialEfficiency: decimal.NewFromFloat(m.MaterialEfficiency),
		BatchSize:          int(m.BatchSize),
		Kind:               kind,
		ParentID:           int(m.ParentId),
		Materials:          []*model.Product{},
	}
	for _, p := range m.Material {
		prod.Materials = append(prod.Materials, ProtoToProduct(p))
	}
	return prod
}

func ProductToProto(p *model.Product) *Product {
	marketPrice, _ := p.MarketPrice.Float64()
	materialEfficiency, _ := p.MaterialEfficiency.Float64()
	kind := Product_BUILD
	if p.Kind == model.ProductBuy {
		kind = Product_BUY
	}
	prod := &Product{
		Id:                 int32(p.ProductID),
		TypeId:             int64(p.TypeID),
		Quantity:           int32(p.Quantity),
		MarketPrice:        marketPrice,
		MarketRegionId:     int32(p.MarketRegionID),
		MaterialEfficiency: materialEfficiency,
		BatchSize:          int32(p.BatchSize),
		Kind:               kind,
		ParentId:           int32(p.ParentID),
		Material:           []*Product{},
	}
	for _, mat := range p.Materials {
		prod.Material = append(prod.Material, ProductToProto(mat))
	}
	return prod
}
