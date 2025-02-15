package main

import (
	"database/sql"
	"testing"
	"time"

	"github.com/jsandberg07/clitest/internal/database"
)

func TestCareDays(t *testing.T) {
	max := 20
	start := normalizeDate(time.Now())
	end := normalizeDate(time.Now().AddDate(0, 0, max))
	ccMax := make([]database.GetCardsDateRangeRow, max)
	date := normalizeDate(time.Now())
	for i := 0; i < max; i++ {
		cc := database.GetCardsDateRangeRow{
			CcID:          int32(i),
			ActivatedOn:   sql.NullTime{Valid: true, Time: date},
			DeactivatedOn: sql.NullTime{Valid: true, Time: date.AddDate(0, 0, i)},
		}
		ccMax[i] = cc
	}

	num := careDaysQuery(start, end, ccMax)

	expected := 0
	for i := 0; i < max; i++ {
		expected += i + 1
	}
	if num != expected {
		t.Fatalf("Max Test: Care days did not match E: %v -- N: %v\n", expected, num)
	}

	// cage cards are still the same, but end date matches start date ie 1 each
	start = normalizeDate(time.Now())
	end = normalizeDate(time.Now())
	num = careDaysQuery(start, end, ccMax)
	expected = 0
	for i := 0; i < max; i++ {
		expected += 1
	}
	if num != expected {
		t.Fatalf("One Test: Care days did not match E: %v -- N: %v\n", expected, num)
	}
}
