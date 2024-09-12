package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/quinn-tao/hmis/v1/config"
	"github.com/quinn-tao/hmis/v1/internal/debug"
	"github.com/quinn-tao/hmis/v1/internal/record"
	"github.com/quinn-tao/hmis/v1/internal/util"
)

var Persistor struct {
    Path string 
    db *sql.DB
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
        cents int not null, 
        name text not null,
        category text not null
    )` 
    _, err = db.Exec(stmt)
    util.CheckErrorf(err, "Cannot create/locate rec table in db at %v: %v", path, err)
    debug.Trace("Necessary tables prepared")
}

func PersistorClose() {
    Persistor.db.Close()
    defer debug.Tracef("Connection to %v closed", Persistor.Path)
}

func InsertRec(cents int, name string, category string) error {
    stmt, err := Persistor.db.Prepare("insert into rec(cents, name, category) values (?, ?, ?)")
    if err != nil {
        return err
    }
    defer stmt.Close()
    _, err = stmt.Exec(cents, name, category)
    return err
}

func GetAllRec() ([]record.Record, error) {
    stmt := "select * from rec" 
    rows, err := Persistor.db.Query(stmt)
    if err != nil {
        return nil, err
    }
    retv := make([]record.Record, 0)
    for rows.Next() {
        var rec record.Record
        err = rows.Scan(&rec.Id, &rec.Cents, &rec.Name, &rec.Category)
        if err != nil {
            return nil, err
        }
        retv = append(retv, rec)
    }
    return retv, nil
}

func RemoveRec(id int) error {
    stmt := fmt.Sprintf("delete from rec where id = %v", id)
    _, err := Persistor.db.Exec(stmt)
    return err 
}
