package model

import "github.com/motki/core/eveapi"

var (
	userScopes = []string{
		eveapi.ScopePublicData,
		eveapi.ScopeESISkillsReadSkills,
		eveapi.ScopeESISkillsReadSkillQueue,
		eveapi.ScopeESIKillmailsReadKillmails,
	}
	logisticsScopes = []string{
		eveapi.ScopeCharacterAssetsRead,
		eveapi.ScopeCharacterIndustryJobsRead,
		eveapi.ScopeCharacterMarketOrdersRead,
		eveapi.ScopeCharacterWalletRead,
		eveapi.ScopeCorporationMarketOrdersRead,
		eveapi.ScopeCorporationIndustryJobsRead,
		eveapi.ScopeCorporationWalletRead,
		eveapi.ScopeESISkillsReadSkills,
		eveapi.ScopeESIUniverseReadStructures,
		eveapi.ScopeESIAssetsReadAssets,
		eveapi.ScopeESIWalletReadCharacterWallet,
		eveapi.ScopeESIMarketsStructureMarkets,
		eveapi.ScopeESIIndustryReadCharacterJobs,
		eveapi.ScopeESIMarketsReadCharacterOrders,
		eveapi.ScopeESICharactersReadBlueprints,
	}
	directorScopes = []string{
		eveapi.ScopeESICorporationsReadStructures,
		eveapi.ScopeESICorporationsWriteStructures,
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
