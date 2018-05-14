package eveapi

const (
	ScopeESIAlliancesReadContacts             = "esi-alliances.read_contacts.v1"                  // EVE SSO scope esi-alliances.read_contacts.v1
	ScopeESIAssetsReadAssets                  = "esi-assets.read_assets.v1"                       // EVE SSO scope esi-assets.read_assets.v1
	ScopeESIAssetsReadCorporationAssets       = "esi-assets.read_corporation_assets.v1"           // EVE SSO scope esi-assets.read_corporation_assets.v1
	ScopeESIBookmarksReadCharacterBookmarks   = "esi-bookmarks.read_character_bookmarks.v1"       // EVE SSO scope esi-bookmarks.read_character_bookmarks.v1
	ScopeESIBookmarksReadCorporationBookmarks = "esi-bookmarks.read_corporation_bookmarks.v1"     // EVE SSO scope esi-bookmarks.read_corporation_bookmarks.v1
	ScopeESICalendarReadCalendarEvents        = "esi-calendar.read_calendar_events.v1"            // EVE SSO scope esi-calendar.read_calendar_events.v1
	ScopeESICalendarRespondCalendarEvents     = "esi-calendar.respond_calendar_events.v1"         // EVE SSO scope esi-calendar.respond_calendar_events.v1
	ScopeESICharactersReadAgentsResearch      = "esi-characters.read_agents_research.v1"          // EVE SSO scope esi-characters.read_agents_research.v1
	ScopeESICharactersReadBlueprints          = "esi-characters.read_blueprints.v1"               // EVE SSO scope esi-characters.read_blueprints.v1
	ScopeESICharactersReadChatChannels        = "esi-characters.read_chat_channels.v1"            // EVE SSO scope esi-characters.read_chat_channels.v1
	ScopeESICharactersReadContacts            = "esi-characters.read_contacts.v1"                 // EVE SSO scope esi-characters.read_contacts.v1
	ScopeESICharactersReadCorporationRoles    = "esi-characters.read_corporation_roles.v1"        // EVE SSO scope esi-characters.read_corporation_roles.v1
	ScopeESICharactersReadFatique             = "esi-characters.read_fatigue.v1"                  // EVE SSO scope esi-characters.read_fatigue.v1
	ScopeESICharactersReadFWStats             = "esi-characters.read_fw_stats.v1"                 // EVE SSO scope esi-characters.read_fw_stats.v1
	ScopeESICharactersReadLoyalty             = "esi-characters.read_loyalty.v1"                  // EVE SSO scope esi-characters.read_loyalty.v1
	ScopeESICharactersReadMedals              = "esi-characters.read_medals.v1"                   // EVE SSO scope esi-characters.read_medals.v1
	ScopeESICharactersReadNotifications       = "esi-characters.read_notifications.v1"            // EVE SSO scope esi-characters.read_notifications.v1
	ScopeESICharactersReadOpportunities       = "esi-characters.read_opportunities.v1"            // EVE SSO scope esi-characters.read_opportunities.v1
	ScopeESICharactersReadStandings           = "esi-characters.read_standings.v1"                // EVE SSO scope esi-characters.read_standings.v1
	ScopeESICharactersReadTitles              = "esi-characters.read_titles.v1"                   // EVE SSO scope esi-characters.read_titles.v1
	ScopeESICharactersWriteContacts           = "esi-characters.write_contacts.v1"                // EVE SSO scope esi-characters.write_contacts.v1
	ScopeESICharactersStats                   = "esi-characterstats.read.v1"                      // EVE SSO scope esi-characterstats.read.v1
	ScopeESIClonesReadClones                  = "esi-clones.read_clones.v1"                       // EVE SSO scope esi-clones.read_clones.v1
	ScopeESIClonesReadImplants                = "esi-clones.read_implants.v1"                     // EVE SSO scope esi-clones.read_implants.v1
	ScopeESIContractsCharacterContracts       = "esi-contracts.read_character_contracts.v1"       // EVE SSO scope esi-contracts.read_character_contracts.v1
	ScopeESIContractsCorporationContracts     = "esi-contracts.read_corporation_contracts.v1"     // EVE SSO scope esi-contracts.read_corporation_contracts.v1
	ScopeESICorporationsReadBlueprints        = "esi-corporations.read_blueprints.v1"             // EVE SSO scope esi-corporations.read_blueprints.v1
	ScopeESICorporationsReadContacts          = "esi-corporations.read_contacts.v1"               // EVE SSO scope esi-corporations.read_contacts.v1
	ScopeESICorporationsReadContainerLogs     = "esi-corporations.read_container_logs.v1"         // EVE SSO scope esi-corporations.read_container_logs.v1
	ScopeESICorporationsReadMembership        = "esi-corporations.read_corporation_membership.v1" // EVE SSO scope esi-corporations.read_corporation_membership.v1
	ScopeESICorporationsReadDivisions         = "esi-corporations.read_divisions.v1"              // EVE SSO scope esi-corporations.read_divisions.v1
	ScopeESICorporationsReadFacilities        = "esi-corporations.read_facilities.v1"             // EVE SSO scope esi-corporations.read_facilities.v1
	ScopeESICorporationsReadFWStats           = "esi-corporations.read_fw_stats.v1"               // EVE SSO scope esi-corporations.read_fw_stats.v1
	ScopeESICorporationsReadMedals            = "esi-corporations.read_medals.v1"                 // EVE SSO scope esi-corporations.read_medals.v1
	ScopeESICorporationsReadOutposts          = "esi-corporations.read_outposts.v1"               // EVE SSO scope esi-corporations.read_outposts.v1
	ScopeESICorporationsReadStandings         = "esi-corporations.read_standings.v1"              // EVE SSO scope esi-corporations.read_standings.v1
	ScopeESICorporationsReadStarbases         = "esi-corporations.read_starbases.v1"              // EVE SSO scope esi-corporations.read_starbases.v1
	ScopeESICorporationsReadStructures        = "esi-corporations.read_structures.v1"             // EVE SSO scope esi-corporations.read_structures.v1
	ScopeESICorporationsReadTitles            = "esi-corporations.read_titles.v1"                 // EVE SSO scope esi-corporations.read_titles.v1
	ScopeESICorporationsTrackMembers          = "esi-corporations.track_members.v1"               // EVE SSO scope esi-corporations.track_members.v1
	ScopeESIFittingsReadFittings              = "esi-fittings.read_fittings.v1"                   // EVE SSO scope esi-fittings.read_fittings.v1
	ScopeESIFittingsWriteFittings             = "esi-fittings.write_fittings.v1"                  // EVE SSO scope esi-fittings.write_fittings.v1
	ScopeESIFleetsReadFleets                  = "esi-fleets.read_fleet.v1"                        // EVE SSO scope esi-fleets.read_fleet.v1
	ScopeESIFleetsWriteFleets                 = "esi-fleets.write_fleet.v1"                       // EVE SSO scope esi-fleets.write_fleet.v1
	ScopeESIIndustryReadCharacterJobs         = "esi-industry.read_character_jobs.v1"             // EVE SSO scope esi-industry.read_character_jobs.v1
	ScopeESIIndustryReadCharacterMining       = "esi-industry.read_character_mining.v1"           // EVE SSO scope esi-industry.read_character_mining.v1
	ScopeESIIndustryReadCorporationJobs       = "esi-industry.read_corporation_jobs.v1"           // EVE SSO scope esi-industry.read_corporation_jobs.v1
	ScopeESIIndustryReadCorporationMining     = "esi-industry.read_corporation_mining.v1"         // EVE SSO scope esi-industry.read_corporation_mining.v1
	ScopeESIKillmailsReadCorporationKillmails = "esi-killmails.read_corporation_killmails.v1"     // EVE SSO scope esi-killmails.read_corporation_killmails.v1
	ScopeESIKillmailsReadKillmails            = "esi-killmails.read_killmails.v1"                 // EVE SSO scope esi-killmails.read_killmails.v1
	ScopeESILocationsReadLocation             = "esi-location.read_location.v1"                   // EVE SSO scope esi-location.read_location.v1
	ScopeESILocationsReadOnline               = "esi-location.read_online.v1"                     // EVE SSO scope esi-location.read_online.v1
	ScopeESILocationsReadShipType             = "esi-location.read_ship_type.v1"                  // EVE SSO scope esi-location.read_ship_type.v1
	ScopeESIMailOrganizeMail                  = "esi-mail.organize_mail.v1"                       // EVE SSO scope esi-mail.organize_mail.v1
	ScopeESIMailReadMail                      = "esi-mail.read_mail.v1"                           // EVE SSO scope esi-mail.read_mail.v1
	ScopeESIMailSendMail                      = "esi-mail.send_mail.v1"                           // EVE SSO scope esi-mail.send_mail.v1
	ScopeESIMarketReadCharacterOrders         = "esi-markets.read_character_orders.v1"            // EVE SSO scope esi-markets.read_character_orders.v1
	ScopeESIMarketReadCorporationOrders       = "esi-markets.read_corporation_orders.v1"          // EVE SSO scope esi-markets.read_corporation_orders.v1
	ScopeESIMarketReadStructureMarkets        = "esi-markets.structure_markets.v1"                // EVE SSO scope esi-markets.structure_markets.v1
	ScopeESIPlanetsManagePlanets              = "esi-planets.manage_planets.v1"                   // EVE SSO scope esi-planets.manage_planets.v1
	ScopeESIPlanetsReadCustomsOffices         = "esi-planets.read_customs_offices.v1"             // EVE SSO scope esi-planets.read_customs_offices.v1
	ScopeESISearchSearchStructures            = "esi-search.search_structures.v1"                 // EVE SSO scope esi-search.search_structures.v1
	ScopeESISkillsReadSkillqueue              = "esi-skills.read_skillqueue.v1"                   // EVE SSO scope esi-skills.read_skillqueue.v1
	ScopeESISkillsReadSkills                  = "esi-skills.read_skills.v1"                       // EVE SSO scope esi-skills.read_skills.v1
	ScopeESIUIOpenWindow                      = "esi-ui.open_window.v1"                           // EVE SSO scope esi-ui.open_window.v1
	ScopeESIUIWriteWaypoint                   = "esi-ui.write_waypoint.v1"                        // EVE SSO scope esi-ui.write_waypoint.v1
	ScopeESIUniverseReadStructures            = "esi-universe.read_structures.v1"                 // EVE SSO scope esi-universe.read_structures.v1
	ScopeESIWalletReadCharacterWallet         = "esi-wallet.read_character_wallet.v1"             // EVE SSO scope esi-wallet.read_character_wallet.v1
	ScopeESIWalletReadCorporationWallet       = "esi-wallet.read_corporation_wallets.v1"          // EVE SSO scope esi-wallet.read_corporation_wallets.v1
)

var AllScopes = []string{
	ScopeESICalendarRespondCalendarEvents,
	ScopeESICalendarReadCalendarEvents,
	ScopeESILocationsReadLocation,
	ScopeESILocationsReadShipType,
	//ScopeESIMailOrganizeEmail,
	//ScopeESIMailReadMail,
	//ScopeESIMailSendMail,
	ScopeESISkillsReadSkillqueue,
	ScopeESIWalletReadCharacterWallet,
	ScopeESISearchSearchStructures,
	//ScopeESIClonesReadClones,
	//ScopeESICharactersReadContacts,
	ScopeESIUniverseReadStructures,
	//ScopeESIBookmarksReadCharacterBookmarks,
	ScopeESIKillmailsReadKillmails,
	ScopeESICorporationsReadMembership,
	ScopeESIAssetsReadAssets,
	//ScopeESIPlanetsManagePlanets,
	//ScopeESIFleetsReadFleet,
	//ScopeESIFleetsWriteFleet,
	//ScopeESIUIOpenWindow,
	//ScopeESIUIWriteWaypoint,
	//ScopeESICharactersWriteContacts,
	ScopeESIFittingsReadFittings,
	ScopeESIFittingsWriteFittings,
	ScopeESIMarketReadStructureMarkets,
	ScopeESICorporationsReadStructures,
	//ScopeESICharactersReadLoyalty,
	//ScopeESICharactersReadOpportunities,
	//ScopeESICharactersReadChatChannels,
	//ScopeESICharactersReadMedals,
	//ScopeESICharactersReadStandings,
	//ScopeESICharactersReadAgentsResearch,
	ScopeESIIndustryReadCharacterJobs,
	ScopeESIMarketReadCharacterOrders,
	ScopeESICharactersReadBlueprints,
}
