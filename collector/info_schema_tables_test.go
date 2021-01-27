// Scrape `information_schema.tables`.

package collector

import (
	"context"
	"database/sql"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartystreets/goconvey/convey"
)

func TestScrapeTableSchema(t *testing.T) { //nolint:unused
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	dbName := "test_cache_db"
	tableName := "test_cache_table"
	*tableSchemaDatabases = dbName

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { //nolint:wsl
		_, err = db.Exec("DROP DATABASE " + dbName)
		if err != nil {
			t.Fatal(err)
		}
	}()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS " + dbName + "." + tableName + " (id int(64))")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { //nolint:wsl
		_, err = db.Exec("DROP TABLE " + dbName + "." + tableName)
		if err != nil {
			t.Fatal(err)
		}
	}()
	_, err = db.Exec("TRUNCATE " + dbName + "." + tableName)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	addRowAndCheckRowsCount(t, ctx, db, dbName, tableName, 1)
	addRowAndCheckRowsCount(t, ctx, db, dbName, tableName, 2)
}

func addRowAndCheckRowsCount(t *testing.T, ctx context.Context, db *sql.DB, dbName, tableName string, expectedRowsCount float64) { //nolint:go-lint
	_, err := db.Exec("INSERT INTO " + dbName + "." + tableName + " VALUES(50)")
	if err != nil {
		t.Fatal(err)
	}
	ch := make(chan prometheus.Metric)
	go func() { //nolint:wsl
		if err = (ScrapeTableSchema{}).Scrape(ctx, db, ch); err != nil {
			t.Errorf("error calling function on test: %s", err)
		}
		close(ch)
	}()

	// For test is important only second receive to channel.
	// Others can be ignored.
	<-ch
	got := readMetric(<-ch)
	<-ch
	<-ch
	<-ch

	expect := MetricResult{
		labels: labelMap{
			"schema": dbName,
			"table":  tableName,
		},
		value:      expectedRowsCount,
		metricType: 1,
	}
	// Variable got.value contains actual rows count in table.
	// Should be equal to count of calling this method.
	convey.ShouldEqual(got, expect)
}
