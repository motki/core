package proto_test

import (
	"testing"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/motki/core/proto"
)

func TestMarshalCharacter(t *testing.T) {
	char := proto.ProtoToCharacter(&proto.Character{
		Id:            1,
		CorporationId: 2,
		AllianceId:    3,
		Name:          "Test",
		BloodlineId:   4,
		RaceId:        5,
		AncestryId:    6,
		BirthDate:     &timestamp.Timestamp{Seconds: 15000000},
		Description:   "My bio",
	})

	if char.CharacterID != 1 {
		t.Errorf("expected model character ID to be 1, got %d", char.CharacterID)
	}
	if char.CorporationID != 2 {
		t.Errorf("expected model corporation ID to be 2, got %d", char.CorporationID)
	}
	if char.AllianceID != 3 {
		t.Errorf("expected model alliance ID to be 3, got %d", char.AllianceID)
	}
	if char.Name != "Test" {
		t.Errorf("expected model character name to be 'Test', got %s", char.Name)
	}
	if char.BloodlineID != 4 {
		t.Errorf("expected model bloodline ID to be 4, got %d", char.BloodlineID)
	}
	if char.RaceID != 5 {
		t.Errorf("expected model race ID to be 5, got %d", char.RaceID)
	}
	if char.AncestryID != 6 {
		t.Errorf("expected model ancestry ID to be 6, got %d", char.AncestryID)
	}
	if char.BirthDate.Unix() != 15000000 {
		t.Errorf("expected model birth date to be 15000000, got %d", char.BirthDate.Unix())
	}
	if char.Description != "My bio" {
		t.Errorf("expected model description to be 'My bio', got %s", char.Description)
	}

	pchar := proto.CharacterToProto(char)
	if pchar.Id != 1 {
		t.Errorf("expected proto character ID to be 1, got %d", pchar.Id)
	}
	if pchar.CorporationId != 2 {
		t.Errorf("expected proto corporation ID to be 2, got %d", pchar.CorporationId)
	}
	if pchar.AllianceId != 3 {
		t.Errorf("expected proto alliance ID to be 3, got %d", pchar.AllianceId)
	}
	if pchar.Name != "Test" {
		t.Errorf("expected proto character name to be 'Test', got %s", pchar.Name)
	}
	if pchar.BloodlineId != 4 {
		t.Errorf("expected proto bloodline ID to be 4, got %d", pchar.BloodlineId)
	}
	if pchar.RaceId != 5 {
		t.Errorf("expected proto race ID to be 5, got %d", pchar.RaceId)
	}
	if pchar.AncestryId != 6 {
		t.Errorf("expected proto ancestry ID to be 6, got %d", pchar.AncestryId)
	}
	if pchar.BirthDate.Seconds != 15000000 {
		t.Errorf("expected proto birth date to be 15000000, got %d", pchar.BirthDate.Seconds)
	}
	if pchar.Description != "My bio" {
		t.Errorf("expected proto description to be 'My bio', got %s", pchar.Description)
	}
}

func TestMarshalCorporation(t *testing.T) {
	corp := proto.ProtoToCorporation(&proto.Corporation{
		Id:           4,
		AllianceId:   3,
		Name:         "Taste",
		Ticker:       "WOTKI",
		CreationDate: &timestamp.Timestamp{Seconds: 15000000},
		Description:  "A bio",
	})

	if corp.CorporationID != 4 {
		t.Errorf("expected model corporation ID to be 4, got %d", corp.CorporationID)
	}
	if corp.AllianceID != 3 {
		t.Errorf("expected model alliance ID to be 3, got %d", corp.AllianceID)
	}
	if corp.Name != "Taste" {
		t.Errorf("expected model corporation name to be 'Taste', got %s", corp.Name)
	}
	if corp.Ticker != "WOTKI" {
		t.Errorf("expected model ticker to be 'WOTKI', got %s", corp.Ticker)
	}
	if corp.CreationDate.Unix() != 15000000 {
		t.Errorf("expected model creation date to be 15000000, got %d", corp.CreationDate.Unix())
	}
	if corp.Description != "A bio" {
		t.Errorf("expected model description to be 'A bio', got %s", corp.Description)
	}

	pcorp := proto.CorporationToProto(corp)
	if pcorp.Id != 4 {
		t.Errorf("expected proto corporation ID to be 4, got %d", pcorp.Id)
	}
	if pcorp.AllianceId != 3 {
		t.Errorf("expected proto alliance ID to be 3, got %d", pcorp.AllianceId)
	}
	if pcorp.Name != "Taste" {
		t.Errorf("expected proto corporation name to be 'Taste', got %s", pcorp.Name)
	}
	if pcorp.Ticker != "WOTKI" {
		t.Errorf("expected proto ticker to be 'WOTKI', got %s", pcorp.Ticker)
	}
	if pcorp.CreationDate.Seconds != 15000000 {
		t.Errorf("expected proto creation date to be 15000000, got %d", pcorp.CreationDate.Seconds)
	}
	if pcorp.Description != "A bio" {
		t.Errorf("expected proto description to be 'A bio', got %s", pcorp.Description)
	}
}

func TestMarshalAlliance(t *testing.T) {
	alliance := proto.ProtoToAlliance(&proto.Alliance{
		Id:          10,
		Name:        "Toast",
		Ticker:      "TRST",
		DateFounded: &timestamp.Timestamp{Seconds: 15000000},
	})

	if alliance.AllianceID != 10 {
		t.Errorf("expected model alliance ID to be 10, got %d", alliance.AllianceID)
	}
	if alliance.Name != "Toast" {
		t.Errorf("expected model character name to be 'Toast', got %s", alliance.Name)
	}
	if alliance.DateFounded.Unix() != 15000000 {
		t.Errorf("expected model date founded to be 15000000, got %d", alliance.DateFounded.Unix())
	}
	if alliance.Ticker != "TRST" {
		t.Errorf("expected model ticker to be 'TRST', got %s", alliance.Ticker)
	}

	palliance := proto.AllianceToProto(alliance)
	if palliance.Id != 10 {
		t.Errorf("expected proto alliance ID to be 10, got %d", palliance.Id)
	}
	if palliance.Name != "Toast" {
		t.Errorf("expected proto character name to be 'Toast', got %s", palliance.Name)
	}
	if palliance.DateFounded.Seconds != 15000000 {
		t.Errorf("expected proto date founded to be 15000000, got %d", palliance.DateFounded.Seconds)
	}
	if palliance.Ticker != "TRST" {
		t.Errorf("expected proto ticker to be 'TRST', got %s", palliance.Ticker)
	}
}
