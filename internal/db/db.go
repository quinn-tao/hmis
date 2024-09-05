package db

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/quinn-tao/hmis/v1/config"
	"github.com/quinn-tao/hmis/v1/internal/debug"
	"github.com/quinn-tao/hmis/v1/internal/display"
)

var Persistor struct {
    Path string 
    db *sql.DB
} 

func PersistorInit(path string) {
    // Create db if not exist
    if _, err := os.Stat(path); os.IsNotExist(err) {
        debug.TraceF("db at %v does not exist. Creating a new db", path)
        _, err := os.Create(path)
        if err != nil {
            display.Errorf("Cannot create db at %v", path)
        }
    }
    debug.TraceF("Located db file at %v", path)
    
    // Open db connection
    db, err := sql.Open("sqlite3", config.StorageLocation()) 
    if err != nil {
        display.Errorf("Cannot open db at provided location: %v", config.StorageLocation())
    }
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
    if err != nil {
        display.Errorf("Cannot create/locate rec table in db at %v: %v", path, err)
    }
    debug.Trace("Necessary tables prepared")
}

func PersistorClose() {
    Persistor.db.Close()
    defer debug.TraceF("Connection to %v closed", Persistor.Path)
}


