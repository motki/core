package model

import "github.com/motki/core/eveapi"

var (
	userScopes = []string{
		eveapi.ScopeESISkillsReadSkills,
		eveapi.ScopeESISkillsReadSkillqueue,
		eveapi.ScopeESIKillmailsReadKillmails,
		eveapi.ScopeESICharactersReadCorporationRoles,
	}
	logisticsScopes = []string{
		eveapi.ScopeESISkillsReadSkills,
		eveapi.ScopeESIUniverseReadStructures,
		eveapi.ScopeESIAssetsReadAssets,
		eveapi.ScopeESIWalletReadCharacterWallet,
		eveapi.ScopeESIMarketReadStructureMarkets,
		eveapi.ScopeESIMarketReadCharacterOrders,
		eveapi.ScopeESIMarketReadStructureMarkets,
		eveapi.ScopeESIIndustryReadCharacterJobs,
		eveapi.ScopeESICharactersReadBlueprints,
		eveapi.ScopeESICharactersReadCorporationRoles,
	}
	directorScopes = []string{
		eveapi.ScopeESIIndustryReadCorporationJobs,
		eveapi.ScopeESICorporationsReadStructures,
		eveapi.ScopeESIMarketReadCorporationOrders,
		eveapi.ScopeESICorporationsReadBlueprints,
		eveapi.ScopeESIAssetsReadCorporationAssets,
		eveapi.ScopeESICorporationsReadDivisions,
		eveapi.ScopeESIWalletReadCorporationWallet,
	}
)

func APIScopesForRole(r Role) []string {
	switch r {
	case RoleUser:
		return userScopes
	case RoleLogistics:
		return logisticsScopes
	case RoleDirector:
		s := make([]string, len(logisticsScopes))
		copy(s, logisticsScopes)
		return append(s, directorScopes...)
	default:
		return []string{}
	}
}
