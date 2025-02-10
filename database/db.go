package database

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/Gerixmus/go-api/docs"

	"github.com/go-sql-driver/mysql"
)

func Connect() (*sql.DB, error) {
	mustGetenv := func(k string) string {
		v := os.Getenv(k)
		if v == "" {
			log.Fatalf("Fatal Error in connect_tcp.go: %s environment variable not set.", k)
		}
		return v
	}

	var (
		dbUser = mustGetenv("DB_USER")
		dbPwd  = mustGetenv("DB_PASSWORD")
		dbHost = mustGetenv("DB_HOST")
		dbPort = mustGetenv("DB_PORT")
		dbName = mustGetenv("DB_NAME")
	)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPwd, dbHost, dbPort, dbName)

	if dbCert, ok := os.LookupEnv("DB_CERT"); ok {
		pool := x509.NewCertPool()
		pem, err := os.ReadFile(dbCert)
		if err != nil {
			return nil, err
		}
		if ok := pool.AppendCertsFromPEM(pem); !ok {
			return nil, err
		}
		mysql.RegisterTLSConfig("config", &tls.Config{
			RootCAs:            pool,
			InsecureSkipVerify: true,
		})
		dsn += "?tls=config"
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	return db, nil
}
