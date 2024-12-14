package db

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/quinn-tao/hmis/v1/config"
	"github.com/quinn-tao/hmis/v1/internal/chrono"
	"github.com/quinn-tao/hmis/v1/internal/coins"
	"github.com/quinn-tao/hmis/v1/internal/debug"
	"github.com/quinn-tao/hmis/v1/internal/profile"
	"github.com/quinn-tao/hmis/v1/internal/record"
	"github.com/quinn-tao/hmis/v1/internal/util"
)

var Persistor struct {
	Path string
	db   *sql.DB
}

func mkstmt(stmts ...string) string {
	return strings.Join(stmts, " ")
}

func PersistorInit(path string) {
	// Create db if not exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		debug.Tracef("db at %v does not exist. Creating a new db", path)
		_, err := os.Create(path)
		util.CheckErrorf(err, "Cannot create db at %v", path)
	}
	debug.Tracef("Located db file at %v", path)

	// Open db connection
	db, err := sql.Open("sqlite3", config.StorageLocation())
	util.CheckErrorf(err, "Cannot open db at provided location: %v", config.StorageLocation())
	Persistor.db = db
	debug.Trace("DB connection established")

	// Create necessary tables if not exists
	stmt := `create table if not exists rec (
        id integer not null primary key, 
        cents bigint not null, 
        name text not null,
        category text not null,
        recordDate text not null
    )`
	_, err = db.Exec(stmt)
	util.CheckErrorf(err, "Cannot create/locate rec table in db at %v: %v", path, err)
	debug.Trace("Necessary tables prepared")
}

func PersistorClose() {
	Persistor.db.Close()
	defer debug.Tracef("Connection to %v closed", Persistor.Path)
}

func InsertRec(rec record.Record) error {
	stmt, err := Persistor.db.Prepare("insert into rec(cents, name, category, recordDate) " +
		"values (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(rec.Amount, rec.Name, rec.Category, fmt.Sprintf("%v", rec.Date))
	return err
}

type recordSearchOptions struct {
	fromDate  string
	toDate    string
	categroy  string
	nameMatch string
}

type RecordSearchOption interface {
	apply(*recordSearchOptions)
}

type funcRecordSearchOption struct {
	op func(*recordSearchOptions)
}

func (fo *funcRecordSearchOption) apply(option *recordSearchOptions) {
	fo.op(option)
}

func RecordSearchWithFromDate(date chrono.Date) RecordSearchOption {
	return &funcRecordSearchOption{
		op: func(rso *recordSearchOptions) {
			rso.fromDate = date.String()
		},
	}
}

func RecordSearchWithToDate(date chrono.Date) RecordSearchOption {
	return &funcRecordSearchOption{
		op: func(rso *recordSearchOptions) {
			rso.toDate = date.String()
		},
	}
}

func RecordSearchWithCategory(c profile.Category) RecordSearchOption {
	return &funcRecordSearchOption{
		op: func(rso *recordSearchOptions) {
			rso.categroy = c.Name
		},
	}
}

func RecordSearchWithNameRg(rg string) RecordSearchOption {
	return &funcRecordSearchOption{
		op: func(rso *recordSearchOptions) {
			rso.nameMatch = ConvertRegexToLikePattern(rg)
		},
	}
}

func GetRecords(options ...RecordSearchOption) ([]record.Record, error) {
	// Interpret options
	var opt recordSearchOptions
	for _, option := range options {
		option.apply(&opt)
	}

	specified := func(optString string) bool {
		return optString != ""
	}

	var selectors []string
	if specified(opt.fromDate) {
		selectors = append(selectors, fmt.Sprintf("recordDate >= %v", opt.fromDate))
	}

	if specified(opt.toDate) {
		selectors = append(selectors, fmt.Sprintf("recordDate <= %v", opt.toDate))
	}

	if specified(opt.categroy) {
		selectors = append(selectors, fmt.Sprintf("category = %v", opt.categroy))
	}

	if specified(opt.categroy) {
		selectors = append(selectors, fmt.Sprintf("name like %v", opt.nameMatch))
	}

	stmt := SearchStmt{
		From:  "rec",
		Where: selectors,
	}

	debug.Tracef("Formatting record search string is %v", stmt)

	rows, err := Persistor.db.Query(stmt.String())
	if err != nil {
		return nil, err
	}
	retv := make([]record.Record, 0)
	for rows.Next() {
		var rec record.Record
		var dateStr string
		err = rows.Scan(&rec.Id, &rec.Amount, &rec.Name, &rec.Category, &dateStr)
		if err != nil {
			return nil, err
		}
		debug.Tracef("Date:%v", dateStr)
		date, err := chrono.NewDate(dateStr)
		if err != nil {
			return nil, err
		}

		rec.Date = date
		retv = append(retv, rec)
	}
	return retv, nil
}

// Get sum of all records as a record
func GetSumRecord() (record.Record, error) {
	stmt := "select sum(cents), count(*) from rec"
	var sum int64
	var cnt int
	err := Persistor.db.QueryRow(stmt).Scan(&sum, &cnt)
	if err != nil {
		return record.Record{}, err
	}

	return record.Record{
		Id:       cnt,
		Name:     "Sum",
		Amount:   coins.RawAmountVal(sum),
		Category: "/",
		Date:     chrono.Today(),
	}, nil
}

func RemoveRec(id int) error {
	stmt := fmt.Sprintf("delete from rec where id = %v", id)
	_, err := Persistor.db.Exec(stmt)
	return err
}
