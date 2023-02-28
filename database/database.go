package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

func CreateIfNotExistsDatabase(name string) (*Database, error) {
	database, err := sql.Open("sqlite3", name)
	if err != nil {
		return nil, err
	}
	defer database.Close()
	return &Database{db: database}, nil
}

type SQLStatement = string

func (self *Database) Execute(statement SQLStatement) error {
	_, err := self.db.Exec(statement)
	if err != nil {
		return err
	}
	return nil
}

func (self *Database) Fetch(statement SQLStatement) ([][]string, error) {
	rows, err := self.db.Query(statement)
	if err != nil {
		return [][]string{}, err
	}
	defer rows.Close()
	result := [][]string{}
	for rows.Next() {
		var IP string
		var Type_banned string
		err = rows.Scan(&IP, &Type_banned)
		if err != nil {
			panic(err)
		}
		result = append(result, []string{IP, Type_banned})
	}
	return result, nil
}
