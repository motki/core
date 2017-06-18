package eveapi

const (
	ScopeCorporationContactsRead                  = "corporationContactsRead"
	ScopePublicData                               = "publicData"
	ScopeCharacterStatsRead                       = "characterStatsRead"
	ScopeCharacterFittingsRead                    = "characterFittingsRead"
	ScopeCharacterFittingsWrite                   = "characterFittingsWrite"
	ScopeCharacterContactsRead                    = "characterContactsRead"
	ScopeCharacterContactsWrite                   = "characterContactsWrite"
	ScopeCharacterLocationRead                    = "characterLocationRead"
	ScopeCharacterNavigationWrite                 = "characterNavigationWrite"
	ScopeCharacterWalletRead                      = "characterWalletRead"
	ScopeCharacterAssetsRead                      = "characterAssetsRead"
	ScopeCharacterCalendarRead                    = "characterCalendarRead"
	ScopeCharacterFactionalWarfareRead            = "characterFactionalWarfareRead"
	ScopeCharacterIndustryJobsRead                = "characterIndustryJobsRead"
	ScopeCharacterKillsRead                       = "characterKillsRead"
	ScopeCharacterMailRead                        = "characterMailRead"
	ScopeCharacterMarketOrdersRead                = "characterMarketOrdersRead"
	ScopeCharacterMedalsRead                      = "characterMedalsRead"
	ScopeCharacterNotificationsRead               = "characterNotificationsRead"
	ScopeCharacterResearchRead                    = "characterResearchRead"
	ScopeCharacterSkillsRead                      = "characterSkillsRead"
	ScopeCharacterAccountRead                     = "characterAccountRead"
	ScopeCharacterContractsRead                   = "characterContractsRead"
	ScopeCharacterBookmarksRead                   = "characterBookmarksRead"
	ScopeCharacterChatChannelsRead                = "characterChatChannelsRead"
	ScopeCharacterClonesRead                      = "characterClonesRead"
	ScopeCharacterOpportunitiesRead               = "characterOpportunitiesRead"
	ScopeCharacterLoyaltyPointsRead               = "characterLoyaltyPointsRead"
	ScopeCorporationWalletRead                    = "corporationWalletRead"
	ScopeCorporationAssetsRead                    = "corporationAssetsRead"
	ScopeCorporationMedalsRead                    = "corporationMedalsRead"
	ScopeCorporationFactionalWarfareRead          = "corporationFactionalWarfareRead"
	ScopeCorporationIndustryJobsRead              = "corporationIndustryJobsRead"
	ScopeCorporationKillsRead                     = "corporationKillsRead"
	ScopeCorporationMembersRead                   = "corporationMembersRead"
	ScopeCorporationMarketOrdersRead              = "corporationMarketOrdersRead"
	ScopeCorporationStructuresRead                = "corporationStructuresRead"
	ScopeCorporationShareholdersRead              = "corporationShareholdersRead"
	ScopeCorporationContractsRead                 = "corporationContractsRead"
	ScopeCorporationBookmarksRead                 = "corporationBookmarksRead"
	ScopeFleetRead                                = "fleetRead"
	ScopeFleetWrite                               = "fleetWrite"
	ScopeStructureVulnUpdate                      = "structureVulnUpdate"
	ScopeRemoteClientUI                           = "remoteClientUI"
	ScopeESICalendarRespondCalendarEvents         = "esi-calendar.respond_calendar_events.v1"
	ScopeESICalendarReadCalendarEvents            = "esi-calendar.read_calendar_events.v1"
	ScopeESILocationReadLocation                  = "esi-location.read_location.v1"
	ScopeESILocationReadShipType                  = "esi-location.read_ship_type.v1"
	ScopeESIMailOrganizeEmail                     = "esi-mail.organize_mail.v1"
	ScopeESIMailReadMail                          = "esi-mail.read_mail.v1"
	ScopeESIMailSendMail                          = "esi-mail.send_mail.v1"
	ScopeESISkillsReadSkills                      = "esi-skills.read_skills.v1"
	ScopeESISkillsReadSkillQueue                  = "esi-skills.read_skillqueue.v1"
	ScopeESIWalletReadCharacterWallet             = "esi-wallet.read_character_wallet.v1"
	ScopeESISearchSearchStructures                = "esi-search.search_structures.v1"
	ScopeESIClonesReadClones                      = "esi-clones.read_clones.v1"
	ScopeESICharactersReadContacts                = "esi-characters.read_contacts.v1"
	ScopeESIUniverseReadStructures                = "esi-universe.read_structures.v1"
	ScopeESIBookmarksReadCharacterBookmarks       = "esi-bookmarks.read_character_bookmarks.v1"
	ScopeESIKillmailsReadKillmails                = "esi-killmails.read_killmails.v1"
	ScopeESICorporationsReadCorporationMembership = "esi-corporations.read_corporation_membership.v1"
	ScopeESIAssetsReadAssets                      = "esi-assets.read_assets.v1"
	ScopeESIPlanetsManagePlanets                  = "esi-planets.manage_planets.v1"
	ScopeESIFleetsReadFleet                       = "esi-fleets.read_fleet.v1"
	ScopeESIFleetsWriteFleet                      = "esi-fleets.write_fleet.v1"
	ScopeESIUIOpenWindow                          = "esi-ui.open_window.v1"
	ScopeESIUIWriteWaypoint                       = "esi-ui.write_waypoint.v1"
	ScopeESICharactersWriteContacts               = "esi-characters.write_contacts.v1"
	ScopeESIFittingsReadFittings                  = "esi-fittings.read_fittings.v1"
	ScopeESIFittingsWriteFittings                 = "esi-fittings.write_fittings.v1"
	ScopeESIMarketsStructureMarkets               = "esi-markets.structure_markets.v1"
	ScopeESICorporationsReadStructures            = "esi-corporations.read_structures.v1"
	ScopeESICorporationsWriteStructures           = "esi-corporations.write_structures.v1"
	ScopeESICharactersReadLoyalty                 = "esi-characters.read_loyalty.v1"
	ScopeESICharactersReadOpportunities           = "esi-characters.read_opportunities.v1"
	ScopeESICharactersReadChatChannels            = "esi-characters.read_chat_channels.v1"
	ScopeESICharactersReadMedals                  = "esi-characters.read_medals.v1"
	ScopeESICharactersReadStandings               = "esi-characters.read_standings.v1"
	ScopeESICharactersReadAgentsResearch          = "esi-characters.read_agents_research.v1"
	ScopeESIIndustryReadCharacterJobs             = "esi-industry.read_character_jobs.v1"
	ScopeESIMarketsReadCharacterOrders            = "esi-markets.read_character_orders.v1"
	ScopeESICharactersReadBlueprints              = "esi-characters.read_blueprints.v1"
)

var AllScopes = []string{
	//ScopeCorporationContactsRead,
	//ScopePublicData,
	//ScopeCharacterStatsRead,
	//ScopeCharacterFittingsRead,
	//ScopeCharacterFittingsWrite,
	//ScopeCharacterContactsRead,
	//ScopeCharacterContactsWrite,
	//ScopeCharacterLocationRead,
	//ScopeCharacterNavigationWrite,
	//ScopeCharacterWalletRead,
	ScopeCharacterAssetsRead,
	//ScopeCharacterCalendarRead,
	//ScopeCharacterFactionalWarfareRead,
	ScopeCharacterIndustryJobsRead,
	//ScopeCharacterKillsRead,
	//ScopeCharacterMailRead,
	//ScopeCharacterMarketOrdersRead,
	//ScopeCharacterMedalsRead,
	//ScopeCharacterNotificationsRead,
	//ScopeCharacterResearchRead,
	//ScopeCharacterSkillsRead,
	//ScopeCharacterAccountRead,
	//ScopeCharacterContractsRead,
	//ScopeCharacterBookmarksRead,
	//ScopeCharacterChatChannelsRead,
	//ScopeCharacterClonesRead,
	//ScopeCharacterOpportunitiesRead,
	//ScopeCharacterLoyaltyPointsRead,
	//ScopeCorporationWalletRead,
	ScopeCorporationAssetsRead,
	//ScopeCorporationMedalsRead,
	//ScopeCorporationFactionalWarfareRead,
	ScopeCorporationIndustryJobsRead,
	//ScopeCorporationKillsRead,
	//ScopeCorporationMembersRead,
	//ScopeCorporationMarketOrdersRead,
	//ScopeCorporationStructuresRead,
	//ScopeCorporationShareholdersRead,
	//ScopeCorporationContractsRead,
	//ScopeCorporationBookmarksRead,
	//ScopeFleetRead,
	//ScopeFleetWrite,
	//ScopeStructureVulnUpdate,
	//ScopeRemoteClientUI,
	ScopeESICalendarRespondCalendarEvents,
	ScopeESICalendarReadCalendarEvents,
	ScopeESILocationReadLocation,
	ScopeESILocationReadShipType,
	//ScopeESIMailOrganizeEmail,
	//ScopeESIMailReadMail,
	//ScopeESIMailSendMail,
	ScopeESISkillsReadSkills,
	ScopeESISkillsReadSkillQueue,
	ScopeESIWalletReadCharacterWallet,
	ScopeESISearchSearchStructures,
	//ScopeESIClonesReadClones,
	//ScopeESICharactersReadContacts,
	ScopeESIUniverseReadStructures,
	//ScopeESIBookmarksReadCharacterBookmarks,
	ScopeESIKillmailsReadKillmails,
	ScopeESICorporationsReadCorporationMembership,
	ScopeESIAssetsReadAssets,
	//ScopeESIPlanetsManagePlanets,
	//ScopeESIFleetsReadFleet,
	//ScopeESIFleetsWriteFleet,
	//ScopeESIUIOpenWindow,
	//ScopeESIUIWriteWaypoint,
	//ScopeESICharactersWriteContacts,
	ScopeESIFittingsReadFittings,
	ScopeESIFittingsWriteFittings,
	ScopeESIMarketsStructureMarkets,
	ScopeESICorporationsReadStructures,
	ScopeESICorporationsWriteStructures,
	//ScopeESICharactersReadLoyalty,
	//ScopeESICharactersReadOpportunities,
	//ScopeESICharactersReadChatChannels,
	//ScopeESICharactersReadMedals,
	//ScopeESICharactersReadStandings,
	//ScopeESICharactersReadAgentsResearch,
	ScopeESIIndustryReadCharacterJobs,
	ScopeESIMarketsReadCharacterOrders,
	ScopeESICharactersReadBlueprints,
}
