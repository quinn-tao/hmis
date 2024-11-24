package cli

import (
	"errors"
	"os"
	"reflect"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

var (
    RowInsertionError = errors.New("Cannot insert row to table")
)

type Table struct {
    title string
    columns []string
    mustSetColunms int64
    writer table.Writer
}

type Column struct {
    Name string 
    Required bool
}

// Create a new table renderer 
func NewTable(title string, columns ...Column) *Table {
    tbl := Table{}
    tbl.title = title
    tbl.mustSetColunms = 0
    tbl.columns = make([]string, len(columns), len(columns))
    for i, column := range columns {
        tbl.columns[i] = column.Name
        tbl.mustSetColunms &= 1 << i
    }
    tbl.writer = table.NewWriter()
    tbl.writer.SetOutputMirror(os.Stdout)
    tbl.writer.Style().Options.SeparateRows = false
    tbl.writer.Style().Options.SeparateColumns = true
    tbl.writer.SetTitle(tbl.title)
    header := make([]interface{}, len(columns))
    for i, column := range columns {
        header[i] = column.Name
    }
    tbl.writer.AppendHeader(header)
    tbl.writer.Style().Color.Header = append(tbl.writer.Style().Color.Header, 
        text.FgGreen)
    tbl.writer.Style().Title.Colors = append(tbl.writer.Style().Title.Colors,
        text.FgGreen)
    return &tbl
}


// Append an row mapped from an arbitrary object of which 
// the fields is a superset of the column names
// Each field value will then be mapped to the corresponding column
func (tbl *Table) AppendRow(obj interface{}) error {
	val := reflect.ValueOf(obj)
	if val.Kind() != reflect.Struct {
		return RowInsertionError
	}

    var setColumns int64
    numColumns := len(tbl.columns)
    row := make([]interface{}, numColumns, numColumns)
	for field := 0; field < val.NumField(); field++ {
		fieldName := val.Type().Field(field).Name
        fieldName = strings.ToLower(fieldName)
        for i, column := range tbl.columns {
            column = strings.ToLower(column)
            if column == fieldName {
		        fieldValue := val.Field(field).Interface()
                row[i] = fieldValue
                setColumns &= 1 << i
                break
            }
        }
	}

    if setColumns & tbl.mustSetColunms != tbl.mustSetColunms {
        return RowInsertionError
    }

    tbl.writer.AppendRow(row)
    return nil
}

// Renders the table to Stdout 
func (tbl *Table) Render() {
    tbl.writer.Render()
}

