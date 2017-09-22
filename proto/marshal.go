package proto

import (
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/motki/motki/evedb"
	"github.com/motki/motki/model"
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

func ProtoToIcon(p *Icon) evedb.Icon {
	return evedb.Icon{
		IconID:          int(p.IconId),
		IconFile:        p.ImageUrl,
		IconDescription: p.Description,
	}
}

func IconToProto(m evedb.Icon) *Icon {
	return &Icon{
		IconId:      int64(m.IconID),
		ImageUrl:    m.IconFile,
		Description: m.IconDescription,
	}
}

func ProtoToRace(p *Race) *evedb.Race {
	return &evedb.Race{
		ID:               int(p.RaceId),
		Name:             p.Name,
		Description:      p.Description,
		ShortDescription: p.ShortDesc,
		Icon:             ProtoToIcon(p.Icon),
	}
}

func RaceToProto(m *evedb.Race) *Race {
	return &Race{
		RaceId:      int64(m.ID),
		Name:        m.Name,
		Description: m.Description,
		ShortDesc:   m.ShortDescription,
		Icon:        IconToProto(m.Icon),
	}
}

func ProtoToAncestry(p *Ancestry) *evedb.Ancestry {
	return &evedb.Ancestry{
		ID:               int(p.AncestryId),
		Name:             p.Name,
		Description:      p.Description,
		ShortDescription: p.ShortDesc,
		BloodlineID:      int(p.BloodlineId),
		Charisma:         int(p.Charisma),
		Willpower:        int(p.Willpower),
		Perception:       int(p.Perception),
		Memory:           int(p.Memory),
		Intelligence:     int(p.Intelligence),
		Icon:             ProtoToIcon(p.Icon),
	}
}

func AncestryToProto(m *evedb.Ancestry) *Ancestry {
	return &Ancestry{
		AncestryId:   int64(m.ID),
		Name:         m.Name,
		Description:  m.Description,
		ShortDesc:    m.ShortDescription,
		BloodlineId:  int64(m.BloodlineID),
		Charisma:     int64(m.Charisma),
		Willpower:    int64(m.Willpower),
		Perception:   int64(m.Perception),
		Memory:       int64(m.Memory),
		Intelligence: int64(m.Intelligence),
		Icon:         IconToProto(m.Icon),
	}
}

func ProtoToBloodline(p *Bloodline) *evedb.Bloodline {
	return &evedb.Bloodline{
		ID:                     int(p.BloodlineId),
		Name:                   p.Name,
		Description:            p.Description,
		ShortDescription:       p.ShortDesc,
		RaceID:                 int(p.RaceId),
		CorporationID:          int(p.CorporationId),
		Charisma:               int(p.Charisma),
		Willpower:              int(p.Willpower),
		Perception:             int(p.Perception),
		Memory:                 int(p.Memory),
		Intelligence:           int(p.Intelligence),
		FemaleDescription:      p.FemaleDesc,
		ShortFemaleDescription: p.ShortFemaleDesc,
		MaleDescription:        p.MaleDesc,
		ShortMaleDescription:   p.ShortMaleDesc,
		Icon:                   ProtoToIcon(p.Icon),
	}
}

func BloodlineToProto(m *evedb.Bloodline) *Bloodline {
	return &Bloodline{
		BloodlineId:     int64(m.ID),
		Name:            m.Name,
		Description:     m.Description,
		ShortDesc:       m.ShortDescription,
		RaceId:          int64(m.RaceID),
		CorporationId:   int64(m.CorporationID),
		Charisma:        int64(m.Charisma),
		Willpower:       int64(m.Willpower),
		Perception:      int64(m.Perception),
		Memory:          int64(m.Memory),
		Intelligence:    int64(m.Intelligence),
		FemaleDesc:      m.FemaleDescription,
		ShortFemaleDesc: m.ShortFemaleDescription,
		MaleDesc:        m.MaleDescription,
		ShortMaleDesc:   m.ShortMaleDescription,
		Icon:            IconToProto(m.Icon),
	}
}

func ProtoToSystem(p *System) *evedb.System {
	return &evedb.System{
		SystemID:        int(p.SystemId),
		Name:            p.Name,
		ConstellationID: int(p.ConstellationId),
		RegionID:        int(p.RegionId),
		Security:        p.Security,
	}
}

func SystemToProto(m *evedb.System) *System {
	return &System{
		SystemId:        int64(m.SystemID),
		Name:            m.Name,
		ConstellationId: int64(m.ConstellationID),
		RegionId:        int64(m.RegionID),
		Security:        m.Security,
	}
}

func ProtoToConstellation(p *Constellation) *evedb.Constellation {
	return &evedb.Constellation{
		ConstellationID: int(p.ConstellationId),
		Name:            p.Name,
		RegionID:        int(p.RegionId),
	}
}

func ConstellationToProto(m *evedb.Constellation) *Constellation {
	return &Constellation{
		Name:            m.Name,
		ConstellationId: int64(m.ConstellationID),
		RegionId:        int64(m.RegionID),
	}
}

func ProtoToRegion(p *Region) *evedb.Region {
	return &evedb.Region{
		RegionID: int(p.RegionId),
		Name:     p.Name,
	}

}

func RegionToProto(m *evedb.Region) *Region {
	return &Region{
		RegionId: int64(m.RegionID),
		Name:     m.Name,
	}
}

func ProtoToItemType(p *ItemType) *evedb.ItemType {
	return &evedb.ItemType{
		ID:          int(p.TypeId),
		Name:        p.Name,
		Description: p.Description,
	}
}

func ItemTypeToProto(m *evedb.ItemType) *ItemType {
	return &ItemType{
		TypeId:      int64(m.ID),
		Name:        m.Name,
		Description: m.Description,
	}
}

func ProtoToItemTypeDetail(p *ItemTypeDetail) *evedb.ItemTypeDetail {
	var deriv []int
	for _, d := range p.DerivativeTypeId {
		deriv = append(deriv, int(d))
	}
	return &evedb.ItemTypeDetail{
		ItemType: &evedb.ItemType{
			ID:          int(p.TypeId),
			Name:        p.Name,
			Description: p.Description,
		},
		GroupID:           int(p.GroupId),
		GroupName:         p.GroupName,
		CategoryID:        int(p.CategoryId),
		CategoryName:      p.CategoryName,
		Mass:              decimal.NewFromFloat(p.Mass),
		Volume:            decimal.NewFromFloat(p.Volume),
		Capacity:          decimal.NewFromFloat(p.Capacity),
		PortionSize:       int(p.PortionSize),
		BasePrice:         decimal.NewFromFloat(p.BasePrice),
		ParentTypeID:      int(p.ParentTypeId),
		BlueprintID:       int(p.BlueprintId),
		DerivativeTypeIDs: deriv,
	}
}

func ItemTypeDetailToProto(m *evedb.ItemTypeDetail) *ItemTypeDetail {
	var deriv []int64
	for _, d := range m.DerivativeTypeIDs {
		deriv = append(deriv, int64(d))
	}
	mass, _ := m.Mass.Float64()
	volume, _ := m.Volume.Float64()
	capacity, _ := m.Capacity.Float64()
	basePrice, _ := m.BasePrice.Float64()
	return &ItemTypeDetail{
		TypeId:           int64(m.ID),
		Name:             m.Name,
		Description:      m.Description,
		GroupId:          int64(m.GroupID),
		GroupName:        m.GroupName,
		CategoryId:       int64(m.CategoryID),
		CategoryName:     m.CategoryName,
		Mass:             mass,
		Volume:           volume,
		Capacity:         capacity,
		PortionSize:      int64(m.PortionSize),
		BasePrice:        basePrice,
		ParentTypeId:     int64(m.ParentTypeID),
		BlueprintId:      int64(m.BlueprintID),
		DerivativeTypeId: deriv,
	}
}

func ProtoToMatSheet(p *MaterialSheet) *evedb.MaterialSheet {
	var mats []*evedb.Material
	for _, mat := range p.Materials {
		mats = append(mats, ProtoToMaterial(mat))
	}
	return &evedb.MaterialSheet{
		ItemType:    ProtoToItemType(p.Type),
		Materials:   mats,
		ProducesQty: int(p.ProducesQty),
	}
}

func MatSheetToProto(m *evedb.MaterialSheet) *MaterialSheet {
	var mats []*Material
	for _, mat := range m.Materials {
		mats = append(mats, MaterialToProto(mat))
	}
	return &MaterialSheet{
		Type:        ItemTypeToProto(m.ItemType),
		Materials:   mats,
		ProducesQty: int64(m.ProducesQty),
	}
}

func ProtoToMaterial(p *Material) *evedb.Material {
	return &evedb.Material{
		ItemType: ProtoToItemType(p.Type),
		Quantity: int(p.Quantity),
	}
}

func MaterialToProto(m *evedb.Material) *Material {
	return &Material{
		Type:     ItemTypeToProto(m.ItemType),
		Quantity: int64(m.Quantity),
	}
}

func ProtoToBlueprint(p *Blueprint) *model.Blueprint {
	kind := model.BlueprintOriginal
	if p.Kind == Blueprint_COPY {
		kind = model.BlueprintCopy
	}
	return &model.Blueprint{
		ItemID:             int(p.ItemId),
		LocationID:         int(p.LocationId),
		TypeID:             int(p.TypeId),
		TypeName:           p.TypeName,
		FlagID:             int(p.FlagId),
		TimeEfficiency:     int(p.TimeEff),
		MaterialEfficiency: int(p.MaterialEff),
		Kind:               kind,
		Quantity:           int(p.Quantity),
		Runs:               int(p.Runs),
	}
}

func BlueprintToProto(m *model.Blueprint) *Blueprint {
	kind := Blueprint_ORIGINAL
	if m.Kind == model.BlueprintCopy {
		kind = Blueprint_COPY
	}
	return &Blueprint{
		ItemId:      int64(m.ItemID),
		LocationId:  int64(m.LocationID),
		TypeId:      int64(m.TypeID),
		TypeName:    m.TypeName,
		FlagId:      int64(m.FlagID),
		TimeEff:     int64(m.TimeEfficiency),
		MaterialEff: int64(m.MaterialEfficiency),
		Kind:        kind,
		Quantity:    int64(m.Quantity),
		Runs:        int64(m.Runs),
	}
}
