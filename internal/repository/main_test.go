package repository

import (
	"log"
	"os"
	"task-manager/internal/config"
	"task-manager/internal/db"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var (
	dsn = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5"
	cfg = config.DBConfig{
		Name:     "test_db",
		User:     "postgres",
		Password: "postgres",
		Host:     "localhost",
		Port:     "5435",
	}
)

var resource *dockertest.Resource
var pool *dockertest.Pool
var testDB *db.DB

func TestMain(m *testing.M) {
	p, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}
	pool = p

	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14.5",
		Env: []string{
			"POSTGRES_USER=" + cfg.User,
			"POSTGRES_PASSWORD=" + cfg.Password,
			"POSTGRES_DB=" + cfg.Name,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: cfg.Port},
			},
		},
	}

	resource, err = pool.RunWithOptions(&opts)
	if err != nil {
		if resource != nil {
			_ = pool.Purge(resource)
		}
		log.Fatalf("could not start the resource: %s", err)
	}

	if err := pool.Retry(func() error {
		var err error
		testDB, err = db.InitDb(cfg)
		if err != nil {
			log.Println(err)
			return err
		}
		return testDB.Ping()
	}); err != nil {
		if resource != nil {
			_ = pool.Purge(resource)
		}

		log.Fatalf("could not connect to database: %s", err)
	}

	err = testDB.RunMigrations()
	if err != nil {
		if resource != nil {
			_ = pool.Purge(resource)
		}
		log.Fatalf("could not run migrations: %s", err)
	}

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("could not purge the resource: %s", err)
	}

	os.Exit(code)
}

func TestPingDB(t *testing.T) {
	err := testDB.Ping()
	if err != nil {
		t.Errorf("could not ping the database: %s", err)
	}
}
