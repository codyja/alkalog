package main

import (
	"flag"
	"github.com/codyja/alkatronic/api"
	"log"
	"os"
	"sync"
)

type DbClient struct {
	Host     string
	Port     string
	Name     string
	Username string
	Password string
}

func NewDbClient(host, port, name, username, password string) *DbClient {
	return &DbClient{
		Host: host,
		Port: port,
		Name: name,
		Username: username,
		Password: password,
	}
}

func main() {

	username, ok := os.LookupEnv("ALKATRONIC_USERNAME")
	if !ok {
		log.Fatalf("ALKATRONIC_USERNAME not set")
	}
	password, ok := os.LookupEnv("ALKATRONIC_PASSWORD")
	if !ok {
		log.Fatalf("ALKATRONIC_PASSWORD not set")
	}
	dbConn, ok := os.LookupEnv("DB_CONNECTION_STRING")
	if !ok {
		log.Fatalf("DB_CONNECTION_STRING not set")
	}

	// read flags
	flagDaemon := flag.Bool("d", false, "Run in Daemon mode to keep polling for new Alkatronic data")
	flagDays := flag.Int("days", 7, "Number of days worth of records to retrieve. Valid days: 7,30, or 90")
	flag.Parse()


	// Initialize new Alkatronic client
	c, err := api.NewAlkatronicClient()
	if err != nil {
		log.Fatalf("error initializing new Alkatronic Client: %s", err)
	}

	// Initialize new Postgresql client
	pg, err := NewPostgresAlkatronic(dbConn)
	if err != nil {
		log.Fatalf("error initializing new Postgresql Client: %s", err )
	}

	// Create wait group so that alkatronic authentication completes before requesting data
	var wg sync.WaitGroup
	wg.Add(1)
	go alkAuth(c, username, password, &wg)
	wg.Wait()


	if *flagDaemon {
		alkLoop(c, pg)
	} else {
		GetAllAlkData(c, pg, *flagDays)
	}

}
