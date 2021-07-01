package database

import (
	"database/sql"

	"github.com/lopezator/migrator"
	"github.com/pkg/errors"
)

//nolint // to bypass gosec sql concat warning
func migrationUser(schema string) *migrator.Migration {
	return &migrator.Migration{
		Name: "202106291342_user",
		Func: func(tx *sql.Tx) error {
			qs := []string{
				`CREATE TABLE IF NOT EXISTS user (
					id BIGINT AUTO_INCREMENT PRIMARY KEY,
					email VARCHAR(255) NOT NULL UNIQUE,
					hashed_password VARBINARY(255) NOT NULL,
					type VARCHAR(20) NOT NULL DEFAULT "guest",
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
				)`,
			}

			for k, query := range qs {
				if _, err := tx.Exec(query); err != nil {
					return errors.Wrapf(err, "applying 202106291342_user migration #%d", k)
				}
			}
			return nil
		},
	}
}

/* ROLLBACK SQL
DROP TABLE IF EXISTS user;
*/
