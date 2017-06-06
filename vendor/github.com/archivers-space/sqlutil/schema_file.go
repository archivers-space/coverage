package sqlutil

import (
	"github.com/gchaincl/dotsql"
)

// LoadSchemaFile takes a filepath to a sql file with create & drop table commands
// and returns a SchemaFile
func LoadSchemaFile(sqlFilePath string) (*SchemaFile, error) {
	f, err := dotsql.LoadFromFile(sqlFilePath)
	if err != nil {
		return nil, err
	}

	return &SchemaFile{
		file: f,
	}, nil
}

// SchemaFile is an sql file that defines a database schema
type SchemaFile struct {
	file *dotsql.DotSql
}

// InitializeDatabase drops everything and calls create on all tables
// WARNING - THIS ZAPS WHATEVER DB IT'S GIVEN. DO NOT CALL THIS SHIT.
// used for testing only, returns a teardown func
func (s *SchemaFile) InitializeDatabase(db Execable) error {
	// TODO - infer table names from de-prefixed create commands,
	// use this to check for data existence
	// // test query to check for database schema existence
	// var exists bool
	// if err = db.QueryRow("select exists(select * from primers limit 1)").Scan(&exists); err == nil {
	//   return nil
	// }

	if err := s.DropAll(db); err != nil {
		return err
	}

	if err := s.CreateAll(db); err != nil {
		return err
	}

	return nil
}

// DropAll executes the command named "drop-all" from the sql file
// this should be a command in the form:
// DROP TABLE IF EXISTS foo, bar, baz ...
func (s *SchemaFile) DropAll(db Execable) error {
	_, err := s.file.Exec(db, "drop-all")
	return err
}

// CreateAll executes all commands that have the prefix "create"
func (s *SchemaFile) CreateAll(db Execable) error {
	for _, cmd := range commandsWithPrefix(s.file, "create") {
		if _, err := s.file.Exec(db, cmd); err != nil {
			return err
		}
	}
	return nil
}
