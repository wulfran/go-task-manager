package db

import (
	"database/sql"
	"fmt"
	"sync"
	"task-manager/internal/config"
	"task-manager/internal/helpers"
	"time"

	_ "github.com/lib/pq"
)

var (
	once      sync.Once
	dbConnect *DB
)

const (
	dbDir = "./internal/db/"
	mDir  = dbDir + "migrations/"
)

type DB struct {
	*sql.DB
}

type dsnConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

func createDsn(c config.DBConfig) string {
	cfg := dsnConfig{
		Username: c.User,
		Password: c.Password,
		Host:     c.Host,
		Port:     c.Port,
		Database: c.Name,
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
}

func InitDb(c config.DBConfig) (*DB, error) {
	fmt.Println("initializing db")
	var err error
	dsn := createDsn(c)
	once.Do(func() {
		var sqlDB *sql.DB
		sqlDB, err = sql.Open("postgres", dsn)
		if err != nil {
			return
		}
		dbConnect = &DB{sqlDB}
	})

	if err != nil {
		return nil, err
	}

	return dbConnect, nil
}

func (d DB) RunMigrations() error {
	exists, err := d.tableExists("migrations")
	if err != nil {
		return fmt.Errorf("db: RunMigrations: error while checking for migrations table: %v", err)
	}
	if !exists {
		if err = d.createTable("0001_create_migrations_table.sql"); err != nil {
			return fmt.Errorf("db: RunMigrations: error while initializing migrations table: %v", err)
		}
	}

	mDone, err := d.getMigrated()
	if err != nil {
		return fmt.Errorf("db: RunMigrations: %v", err)
	}

	m, err := SQLFiles.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("db: RunMigrations: error while reading migrations dir: %s", err)
	}

	for _, item := range m {
		if !item.IsDir() {
			if !helpers.SliceContains(mDone, item.Name()) {
				if err := d.createTable(item.Name()); err != nil {
					return fmt.Errorf("db: RunMigrations: error while running %s, %v", item.Name(), err)
				}
			}
		}
	}

	return nil
}

func GetQuery(path string) (string, error) {
	data, err := SQLFiles.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("GetQuery: unable to read embedded query from path %s: %v", path, err)
	}

	return string(data), nil
}

func (d DB) tableExists(n string) (bool, error) {
	q, err := GetQuery("queries/utils/tableExists.sql")
	if err != nil {
		fmt.Println(q)
		return false, fmt.Errorf("tableExists: error while reading query: %v", err)
	}
	q = fmt.Sprintf(q, n)
	var exists bool
	err = d.QueryRow(q).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("tableExists: error while executing query: %v", err)
	}
	return exists, nil
}

func (d DB) createTable(n string) error {
	q, err := GetQuery("migrations/" + n)
	if err != nil {
		return fmt.Errorf("createTable: error reading query from %s, %v", "migrations/"+n, err)
	}

	_, err = d.Exec(q)
	if err != nil {
		return fmt.Errorf("createTable: error while executing create table query: %v", err)
	}

	if err = d.registerMigration(n); err != nil {
		return fmt.Errorf("createTable: %v", err)
	}

	return nil
}

func (d DB) registerMigration(n string) error {
	q, err := GetQuery("queries/utils/insertMigration.sql")
	if err != nil {
		return fmt.Errorf("registerMigration: error while reading query: %v", err)
	}
	q = fmt.Sprintf(q, n, time.Now().Format(time.RFC3339))
	_, err = d.Exec(q)
	if err != nil {
		return fmt.Errorf("registerMigration: error while executing query: %v", err)
	}
	return nil
}

func (d DB) getMigrated() ([]string, error) {
	q, err := GetQuery("queries/utils/readMigrationsTable.sql")
	if err != nil {
		return nil, fmt.Errorf("db: getMigrated: error while reading migrations table: %v", err)
	}

	var n string
	rows, err := d.Query(q)
	if err != nil {
		return nil, fmt.Errorf("getMigrated: error while executing query: %v", err)
	}
	defer rows.Close()

	var list []string
	for rows.Next() {
		if err = rows.Scan(&n); err != nil {
			return nil, fmt.Errorf("getMigrated: error while reading rows: %v", err)
		}
		list = append(list, n)
	}

	return list, nil
}
