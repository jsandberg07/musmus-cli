package main

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/jsandberg07/clitest/internal/database"
)

// learning a lot about generics thats for sure
func TestNormalizeCCExport(t *testing.T) {
	date := time.Now()

	expected := []CageCardExport{
		{
			CcID:          69,
			IName:         "Beans Johnson",
			PNumber:       "12-24-32",
			SName:         sql.NullString{Valid: true, String: "000664"},
			ActivatedOn:   sql.NullTime{Valid: true, Time: date},
			DeactivatedOn: sql.NullTime{Valid: false},
		},
		{
			CcID:          420,
			IName:         "Beans Johnson",
			PNumber:       "12-24-32",
			SName:         sql.NullString{Valid: true, String: "022"},
			ActivatedOn:   sql.NullTime{Valid: true, Time: date},
			DeactivatedOn: sql.NullTime{Valid: true, Time: date},
		},
	}

	invTest := []database.GetCageCardsInvestigatorRow{
		{
			CcID:          69,
			IName:         "Beans Johnson",
			PNumber:       "12-24-32",
			SName:         sql.NullString{Valid: true, String: "000664"},
			ActivatedOn:   sql.NullTime{Valid: true, Time: date},
			DeactivatedOn: sql.NullTime{Valid: false},
		},
		{
			CcID:          420,
			IName:         "Beans Johnson",
			PNumber:       "12-24-32",
			SName:         sql.NullString{Valid: true, String: "022"},
			ActivatedOn:   sql.NullTime{Valid: true, Time: date},
			DeactivatedOn: sql.NullTime{Valid: true, Time: date},
		},
	}

	dateTest := []database.GetCardsDateRangeRow{
		{
			CcID:          69,
			IName:         "Beans Johnson",
			PNumber:       "12-24-32",
			SName:         sql.NullString{Valid: true, String: "000664"},
			ActivatedOn:   sql.NullTime{Valid: true, Time: date},
			DeactivatedOn: sql.NullTime{Valid: false},
		},
		{
			CcID:          420,
			IName:         "Beans Johnson",
			PNumber:       "12-24-32",
			SName:         sql.NullString{Valid: true, String: "022"},
			ActivatedOn:   sql.NullTime{Valid: true, Time: date},
			DeactivatedOn: sql.NullTime{Valid: true, Time: date},
		},
	}

	allTest := []database.GetCageCardsAllRow{
		{
			CcID:          69,
			IName:         "Beans Johnson",
			PNumber:       "12-24-32",
			SName:         sql.NullString{Valid: true, String: "000664"},
			ActivatedOn:   sql.NullTime{Valid: true, Time: date},
			DeactivatedOn: sql.NullTime{Valid: false},
		},
		{
			CcID:          420,
			IName:         "Beans Johnson",
			PNumber:       "12-24-32",
			SName:         sql.NullString{Valid: true, String: "022"},
			ActivatedOn:   sql.NullTime{Valid: true, Time: date},
			DeactivatedOn: sql.NullTime{Valid: true, Time: date},
		},
	}

	activeTest := []database.GetCageCardsActiveRow{
		{
			CcID:          69,
			IName:         "Beans Johnson",
			PNumber:       "12-24-32",
			SName:         sql.NullString{Valid: true, String: "000664"},
			ActivatedOn:   sql.NullTime{Valid: true, Time: date},
			DeactivatedOn: sql.NullTime{Valid: false},
		},
		{
			CcID:          420,
			IName:         "Beans Johnson",
			PNumber:       "12-24-32",
			SName:         sql.NullString{Valid: true, String: "022"},
			ActivatedOn:   sql.NullTime{Valid: true, Time: date},
			DeactivatedOn: sql.NullTime{Valid: true, Time: date},
		},
	}

	protocolTest := []database.GetCageCardsProtocolRow{
		{
			CcID:          69,
			IName:         "Beans Johnson",
			PNumber:       "12-24-32",
			SName:         sql.NullString{Valid: true, String: "000664"},
			ActivatedOn:   sql.NullTime{Valid: true, Time: date},
			DeactivatedOn: sql.NullTime{Valid: false},
		},
		{
			CcID:          420,
			IName:         "Beans Johnson",
			PNumber:       "12-24-32",
			SName:         sql.NullString{Valid: true, String: "022"},
			ActivatedOn:   sql.NullTime{Valid: true, Time: date},
			DeactivatedOn: sql.NullTime{Valid: true, Time: date},
		},
	}

	// var tests []interface{}
	// tests = append(tests, invTest, dateTest, allTest, activeTest, protocolTest)
	// i really wanted a cool array BUT []interface{} messes with generic type deduction

	// should pass
	output := NormalizeCCExport(invTest)
	if !reflect.DeepEqual(expected, output) {
		t.Fatal("Inv test failed")
	}

	output = NormalizeCCExport(dateTest)
	if !reflect.DeepEqual(expected, output) {
		t.Fatal("Date test failed")
	}

	output = NormalizeCCExport(allTest)
	if !reflect.DeepEqual(expected, output) {
		t.Fatal("All test failed")
	}

	output = NormalizeCCExport(activeTest)
	if !reflect.DeepEqual(expected, output) {
		t.Fatal("Active test failed")
	}

	output = NormalizeCCExport(protocolTest)
	if !reflect.DeepEqual(expected, output) {
		t.Fatal("Protocol test failed")
	}

	// empty
	empty := []database.GetCageCardsActiveRow{}
	expected = []CageCardExport{}
	output = NormalizeCCExport(empty)
	if !reflect.DeepEqual(expected, output) {
		t.Fatal("empty test failed")
	}

}
