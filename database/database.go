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
	return &Database{db: database}, nil
}

type SQLStatement = string

func (self *Database) Execute(statement SQLStatement) error {
	ctx, err := self.db.Begin()
	if err != nil {
		panic(err)
	}
	s, err := ctx.Prepare(statement)
	if err != nil {
		panic(err)
	}
	defer s.Close()
	_, err = s.Exec()
	if err != nil {
		panic(err)
	}
	err = ctx.Commit()
	if err != nil {
		panic(err)
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

func (self *Database) Close() {
	defer self.db.Close()
}
