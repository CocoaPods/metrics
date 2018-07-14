package aggregator

import (
	"database/sql"
	"fmt"

	"github.com/cocoapods/metrics/internal/config"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	// Importing postgres sql driver
	_ "github.com/lib/pq"
)

type Aggregator struct {
	warehouseConn *sql.DB
	metricsConn   *sql.DB
}

func openDBConn(dbs *config.DBConfig, sslMode string) (*sql.DB, error) {
	conStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%d sslmode=%s", dbs.Username, dbs.DatabaseName, dbs.Password, dbs.Host, dbs.Port, sslMode)
	db, err := sql.Open("postgres", conStr)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to connect to database")
	}

	return db, nil
}

const (
	extractDownloadStatsQuery = `SELECT dependency_name, dependency_version, to_char(sent_at, 'YYYY-MM') AS rollup_date,
	count(CASE WHEN pod_try THEN 1 ELSE null END) AS pod_tries,
	count(CASE WHEN pod_try THEN null ELSE 1 END) AS downloads
	FROM cocoapods_stats_production.install
	GROUP BY rollup_date, dependency_name, dependency_version;`
)

func NewAggregator(c *config.Config) (*Aggregator, error) {
	warehouseConn, err := openDBConn(c.WarehouseDB, "verify-full")
	if err != nil {
		return nil, errors.Wrap(err, "Failed to connect to warehouse")
	}
	metricsConn, err := openDBConn(c.MetricsDB, "disable")
	if err != nil {
		return nil, errors.Wrap(err, "Failed to connect to app db")
	}
	return &Aggregator{warehouseConn: warehouseConn, metricsConn: metricsConn}, nil
}

type DownloadRow struct {
	Name      string
	Version   string
	Date      string
	Tries     int
	Downloads int
}

func flushRows(rows []DownloadRow, tx *sql.Tx) error {
	stmt, err := tx.Prepare("INSERT INTO cocoapods_stats_rolluped.pod_monthly_metrics(pod_name, pod_version, rollup_datestamp, tries, downloads) VALUES($1, $2, $3, $4, $5);")
	if err != nil {
		return errors.Wrap(err, "Failed to prepare statement")
	}
	defer stmt.Close()
	for _, r := range rows {
		_, err := stmt.Exec(r.Name, r.Version, r.Date, r.Tries, r.Downloads)
		if err != nil {
			return errors.Wrap(err, "Failed to execute insert")
		}
	}

	return nil
}

func (a *Aggregator) aggregateDownloadStats() error {
	logrus.Info("Fetching download stats")

	rows, err := a.warehouseConn.Query(extractDownloadStatsQuery)
	if err != nil {
		return errors.Wrap(err, "Failed to extract download stats")
	}
	defer rows.Close()

	rowBuf := []DownloadRow{}

	logrus.Info("Preparing download stats")

	for rows.Next() {
		_, _ = rows.Columns()

		// Extract
		r := &DownloadRow{}
		err := rows.Scan(&r.Name, &r.Version, &r.Date, &r.Tries, &r.Downloads)
		if err != nil {
			logrus.WithError(err).Info("Failed to scan row")
			continue
		}

		// Load

		rowBuf = append(rowBuf, *r)

	}

	logrus.WithField("row_count", len(rowBuf)).Info("Inserting download stats")

	tx, err := a.metricsConn.Begin()
	if err != nil {
		return errors.Wrap(err, "Failed to start transaction")
	}

	if err := flushRows(rowBuf, tx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logrus.WithError(rollbackErr).Error("Failed to rollback transaction")
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "Failed to commit transaction")
	}

	logrus.Info("Inserted download stats")

	return nil
}

func (a *Aggregator) Aggregate() error {
	if err := a.aggregateDownloadStats(); err != nil {
		logrus.WithError(err).Warn("Aggregation failed")
	}
	return nil
}
