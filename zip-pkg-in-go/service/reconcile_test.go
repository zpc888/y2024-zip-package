package service

import "testing"

func TestColIndexToString(t *testing.T) {
	tests := []struct {
		colIndex int
		want     string
	}{
		{0, "A"},
		{1, "B"},
		{25, "Z"},
		{26, "AA"},
		{27, "AB"},
		{51, "AZ"},
		{52, "BA"},
		{53, "BB"},
		{77, "BZ"},
		{78, "CA"},
		{79, "CB"},
		{701, "ZZ"},
		{702, "AAA"},
		{703, "AAB"},
		{727, "AAZ"},
		{728, "ABA"},
		{729, "ABB"},
		{1352, "AZA"},
		{1353, "AZB"},
		{1354, "AZC"},
		{1378, "BAA"},
		{1379, "BAB"},
		{1380, "BAC"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			actualMap := colIndexToString(tt.colIndex)
			if got := actualMap[tt.colIndex]; got != tt.want {
				t.Errorf("ColIndexToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
