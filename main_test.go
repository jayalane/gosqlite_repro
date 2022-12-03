// -*- tab-width: 2 -*-

package main

import (
	"database/sql"
	"fmt"
	_ "github.com/glebarez/go-sqlite"
	"os"
	"testing"
)

const perfTune = `
pragma journal_mode = WAL;
pragma synchronous = normal;
pragma temp_store = memory;
pragma mmap_size = 30000000000;`

var dbNameTemplate = "persist_%d.sqlite"
var dbPath = "./dbs/"

// var sqliteMutex = sync.Mutex{}

// createDB just creates the DB and puts the table in and runs the tuning SQL.
func createDB(index int) error {

	//	sqliteMutex.Lock()
	//	defer sqliteMutex.Unlock()

	dbName := fmt.Sprintf(dbNameTemplate, index)

	db, err := sql.Open("sqlite", dbPath+dbName)
	if err != nil {
		fmt.Println("Error opening", dbPath+dbName, err)
		return err
	}
	_, err = db.Exec(perfTune)
	if err != nil {
		fmt.Println("Error tuning", err)
		return err
	}
	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS ips (
  key TEXT,
  value INT,
  ttl INT)`)
	if err != nil {
		fmt.Println("Error creating if not", err)
		return err
	}
	return nil
}

func TestCreateDBs(t *testing.T) {

	err := os.Mkdir(dbPath, 0770)

	if err != nil && !os.IsExist(err) {
		fmt.Println("Can't make state DB ./dbs/ for SQLITE", err)
		t.Fail()
	}

	for N := 0; N < 255; N++ {
		go func(n int) {
			err := createDB(n)
			if err != nil {
				fmt.Println("DB open got error", err)
				t.Fail()
			}
		}(N)
	}
}
