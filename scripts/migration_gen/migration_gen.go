package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/pkg/errors"
)

const (
	dbPkgPath          = "./internal/database"
	migrationsListPath = dbPkgPath + "/migrations.go"
)

var (
	fileName = flag.String("name", "", "Set migration file name")
	tmplText = `package database

import (
	"database/sql"

	"github.com/lopezator/migrator"
	"github.com/pkg/errors"
)

//nolint // to bypass gosec sql concat warning
func {{.MigrationVarName}}(schema string) *migrator.Migration {
	return &migrator.Migration{
		Name: "{{.MigrationFileName}}",
		Func: func(tx *sql.Tx) error {
			qs := []string{}

			for k, query := range qs {
				if _, err := tx.Exec(query); err != nil {
					return errors.Wrapf(err, "applying {{.MigrationFileName}} migration #%d", k)
				}
			}
			return nil
		},

	}
}

/* ROLLBACK SQL

*/
`
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

type migration struct {
	MigrationVarName  string
	MigrationFileName string
}

func main() {
	flag.Parse()
	if *fileName == "" {
		log.Fatal("ERROR: not set filename.")
	}

	migrationName := fmt.Sprintf("%v_%s", time.Now().Format("200601021504"), toSnakeCase(*fileName))
	newMigrationFile, err := os.Create(filepath.Join(dbPkgPath, filepath.Base(migrationName+".go")))
	if err != nil {
		log.Fatal(errors.Wrap(err, "create new migration error"))
	}
	defer newMigrationFile.Close()

	newMgr := migration{
		MigrationVarName:  getVarName(migrationName),
		MigrationFileName: migrationName,
	}

	tpl, err := template.New("migration template").Parse(tmplText)
	if err != nil {
		log.Fatal(errors.Wrap(err, "parse template error"))
	}

	err = tpl.ExecuteTemplate(newMigrationFile, "migration template", newMgr)
	if err != nil {
		log.Fatal(errors.Wrap(err, "exec template error"))
	}

	updateMigrationsList(newMgr)
}

func updateMigrationsList(newMgr migration) {
	b, err := ioutil.ReadFile(migrationsListPath)
	if err != nil {
		log.Fatal(errors.Wrap(err, "migrations list file read error"))
	}

	insertIndex := strings.LastIndex(string(b), "),")
	addString := "\t\t" + newMgr.MigrationVarName + "(schema)," + "\n\t\t"
	migrationsList, err := os.OpenFile(migrationsListPath, os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(errors.Wrap(err, "migrations list open file error"))
	}
	defer migrationsList.Close()

	if _, err = migrationsList.WriteAt([]byte(addString), int64(insertIndex-1)); err != nil {
		log.Fatal(errors.Wrap(err, "write file error"))
	}
	if _, err = migrationsList.WriteAt(b[insertIndex:], int64(insertIndex-1+len(addString))); err != nil {
		log.Fatal(errors.Wrap(err, "write file error"))
	}
}

func getVarName(migrationName string) string {
	parts := strings.Split(migrationName, "_")

	res := make([]string, 0, len(parts))
	for i := 1; i < len(parts); i++ {
		res = append(res, strings.Title(strings.ToLower(parts[i])))
	}

	return "migration" + strings.Join(res, "")
}

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
