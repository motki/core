package eveapi

import (
	"strconv"

	"github.com/antihax/goesi/esi"
	"github.com/antihax/goesi/optional"
	"golang.org/x/net/context"
)

//[ AssetSafety, AutoFit, Bonus, Booster, BoosterBay, Capsule, Cargo, CorpDeliveries, CorpSAG1, CorpSAG2, CorpSAG3, CorpSAG4, CorpSAG5, CorpSAG6, CorpSAG7, CrateLoot, Deliveries, DroneBay, DustBattle, DustDatabank, FighterBay, FighterTube0, FighterTube1, FighterTube2, FighterTube3, FighterTube4, FleetHangar, Hangar, HangarAll, HiSlot0, HiSlot1, HiSlot2, HiSlot3, HiSlot4, HiSlot5, HiSlot6, HiSlot7, HiddenModifers, Implant, Impounded, JunkyardReprocessed, JunkyardTrashed, LoSlot0, LoSlot1, LoSlot2, LoSlot3, LoSlot4, LoSlot5, LoSlot6, LoSlot7, Locked, MedSlot0, MedSlot1, MedSlot2, MedSlot3, MedSlot4, MedSlot5, MedSlot6, MedSlot7, OfficeFolder, Pilot, PlanetSurface, QuafeBay, Reward, RigSlot0, RigSlot1, RigSlot2, RigSlot3, RigSlot4, RigSlot5, RigSlot6, RigSlot7, SecondaryStorage, ServiceSlot0, ServiceSlot1, ServiceSlot2, ServiceSlot3, ServiceSlot4, ServiceSlot5, ServiceSlot6, ServiceSlot7, ShipHangar, ShipOffline, Skill, SkillInTraining, SpecializedAmmoHold, SpecializedCommandCenterHold, SpecializedFuelBay, SpecializedGasHold, SpecializedIndustrialShipHold, SpecializedLargeShipHold, SpecializedMaterialBay, SpecializedMediumShipHold, SpecializedMineralHold, SpecializedOreHold, SpecializedPlanetaryCommoditiesHold, SpecializedSalvageHold, SpecializedShipHold, SpecializedSmallShipHold, StructureActive, StructureFuel, StructureInactive, StructureOffline, SubSystemSlot0, SubSystemSlot1, SubSystemSlot2, SubSystemSlot3, SubSystemSlot4, SubSystemSlot5, SubSystemSlot6, SubSystemSlot7, SubsystemBay, Unlocked, Wallet, Wardrobe ]

type Blueprint struct {
	ItemID             int64
	LocationID         int64
	TypeID             int64
	LocationFlag       string
	TimeEfficiency     int64
	MaterialEfficiency int64

	// -2 = BPC (and always qty 1), else BPO
	Quantity int64

	// -1 = infinite runs (a BPO)
	Runs int64
}

func (api *EveAPI) GetCorporationBlueprints(ctx context.Context, corpID int) ([]*Blueprint, error) {
	_, err := TokenFromContext(ctx)
	if err != nil {
		return nil, err
	}
	var bps []*Blueprint
	var max int
	for p := 0; p <= max; p++ {
		res, resp, err := api.client.ESI.CorporationApi.GetCorporationsCorporationIdBlueprints(
			ctx,
			int32(corpID),
			&esi.GetCorporationsCorporationIdBlueprintsOpts{Page: optional.NewInt32(int32(p))})
		if err != nil {
			return nil, err
		}
		max, err = strconv.Atoi(resp.Header.Get("X-Pages"))
		if err != nil {
			api.logger.Debugf("error reading X-Pages header: ", err.Error())
		}
		for _, bp := range res {
			bps = append(bps, &Blueprint{
				ItemID:             bp.ItemId,
				LocationID:         bp.LocationId,
				LocationFlag:       bp.LocationFlag,
				TypeID:             int64(bp.TypeId),
				Quantity:           int64(bp.Quantity),
				TimeEfficiency:     int64(bp.TimeEfficiency),
				MaterialEfficiency: int64(bp.MaterialEfficiency),
				Runs:               int64(bp.Runs),
			})
		}
	}
	return bps, nil
}
