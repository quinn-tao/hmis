package db_test

import (
	"testing"

	"github.com/quinn-tao/hmis/v1/internal/db"
)

func TestSqlSearch(t *testing.T) {
	tcs := []struct {
		Name                 string
		ExpectedSearchString string
		SearchParams         db.SearchStmt
	}{
		{
			Name: "select all statement",
			SearchParams: db.SearchStmt{
				From: "test",
			},
			ExpectedSearchString: "select * from test",
		},
		{
			Name: "select column statement",
			SearchParams: db.SearchStmt{
				Select: []string{"testCol1", "testCol2"},
				From:   "test",
			},
			ExpectedSearchString: "select testCol1, testCol2 from test",
		},
		{
			Name: "select all statement with single where clause",
			SearchParams: db.SearchStmt{
				From:  "test",
				Where: []string{"testCol1=1"},
			},
			ExpectedSearchString: "select * from test\n" +
				"where\n" + "testCol1=1",
		}, {
			Name: "select all statement with more where clause",
			SearchParams: db.SearchStmt{
				From:  "test",
				Where: []string{"testCol1=1", "testCol2=2"},
			},
			ExpectedSearchString: "select * from test\n" +
				"where\n" + "testCol1=1\n" + "and testCol2=2",
		},
	}

	for _, tc := range tcs {
		t.Logf("[TestSqlSearch] running %v", tc.Name)
		actualSearchString := tc.SearchParams.String()
		if actualSearchString != tc.ExpectedSearchString {
			t.Fatalf("Expected %v, Got %v", tc.ExpectedSearchString, actualSearchString)
		}
	}
}
