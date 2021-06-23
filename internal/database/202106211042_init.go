package database

import (
	"database/sql"

	"github.com/lopezator/migrator"
	"github.com/pkg/errors"
)

//nolint // to bypass gosec sql concat warning
func migrationInit(schema string) *migrator.Migration {
	return &migrator.Migration{
		Name: "202106211042_init",
		Func: func(tx *sql.Tx) error {
			qs := []string{
				`CREATE TABLE IF NOT EXISTS city (
					id INT AUTO_INCREMENT PRIMARY KEY,
					name VARCHAR(255) NOT NULL UNIQUE
				)`,
				`CREATE TABLE IF NOT EXISTS bus (
					id BIGINT AUTO_INCREMENT PRIMARY KEY,
					city_id INT NOT NULL,
					-- num can be like '35A'.
					num VARCHAR(255) NOT NULL,
					UNIQUE(city_id, num),
					FOREIGN KEY(city_id) REFERENCES city(id) ON UPDATE CASCADE ON DELETE CASCADE
				)`,
				`CREATE TABLE IF NOT EXISTS stop (
					id BIGINT AUTO_INCREMENT PRIMARY KEY,
					city_id INT NOT NULL,
					address VARCHAR(255) NOT NULL,
					UNIQUE(city_id, address),
					FOREIGN KEY(city_id) REFERENCES city(id) ON UPDATE CASCADE ON DELETE CASCADE
				)`,
				`CREATE TABLE IF NOT EXISTS route (
					bus_id BIGINT NOT NULL,
					stop_id BIGINT NOT NULL,
					-- step is the number of a bus stop in its route.
					step TINYINT NOT NULL,
					PRIMARY KEY(bus_id, step),
					FOREIGN KEY(bus_id) REFERENCES bus(id) ON UPDATE CASCADE ON DELETE CASCADE,
					FOREIGN KEY(stop_id) REFERENCES stop(id) ON UPDATE CASCADE
				)`,
			}

			for k, query := range qs {
				if _, err := tx.Exec(query); err != nil {
					return errors.Wrapf(err, "applying 202106211042_init migration #%d", k)
				}
			}
			return nil
		},
	}
}

/* ROLLBACK SQL
DROP TABLE IF EXISTS city;
DROP TABLE IF EXISTS bus;
DROP TABLE IF EXISTS stop;
DROP TABLE IF EXISTS route;
*/
