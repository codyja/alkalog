package main

import (
	"fmt"
	"github.com/codyja/alkatronic/api"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

// alkLoop starts up a ticker that calls GetAlkData every X minutes
func alkLoop(c *api.AlkatronicClient, pg *PostgresAlkatronic) {
	ticker := time.NewTicker(30 * time.Minute)

	for {
		GetAlkData(c, pg)
		<-ticker.C
	}
}

// alkAuth performs initial authentication or reauthentication as needed
func alkAuth(c *api.AlkatronicClient, username string, password string, wg *sync.WaitGroup) {
	// Read in home directory to read and write token file to
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error getting home directory: %s", err)
	}

	// Location to store token inside file
	tokenFileLocation := fmt.Sprintf("%s/.alkatronic", home)

	// Read existing token from token storage file
	existingToken, err := ioutil.ReadFile(tokenFileLocation)
	if err != nil {
		log.Printf("Checking for existing token file: %s", err)
	}

	// Authenticate if there's no existing token file, otherwise call SetAccessToken to reuse existing token on subsequent calls
	if string(existingToken) == "" {
		c.Authenticate(username, password)

		// write token locally
		tokenBytes := []byte(fmt.Sprintf("%s",c.AccessToken()))
		err = ioutil.WriteFile(tokenFileLocation, tokenBytes, 0600)
		if err != nil {
			log.Fatalf("Error writting token file: %s", err)
		}
	} else {
		log.Printf("Using existing access token read from ~/.alkatronic")
		c.SetAccessToken(string(existingToken))
	}

	wg.Done()
}

// GetAlkData calls GetLatestResult to get the latest record then calls InsertRecord to log to the DB
func GetAlkData(c *api.AlkatronicClient, pg *PostgresAlkatronic) {
	devices, err := c.GetDevices()
	if err != nil {
		log.Fatalf("Error getting devices: %s", err)
	}

	for _, d := range devices.Data {
		r, err := c.GetLatestResult(d.DeviceID)
		if err != nil {
			log.Fatalf("Error getting last test result: %s", err)
		}

		log.Printf(
			"Latest test result: Device Name: %s, Record ID: %d, KH_Value: %.2f, Create Time: %s\n",
			d.FriendlyName,
			r.RecordID,
			api.ConvertKh(r.KhValue),
			time.Unix(r.CreateTime, 0).Format(time.RFC822Z))

		r.KhValue   = api.ConvertKh(r.KhValue)

		device := api.Device {
			UpperKh: api.ConvertKh(d.UpperKh),
			LowerKh: api.ConvertKh(d.LowerKh),
		}

		err = pg.InsertRecord(r, device)
		if err != nil {
			log.Fatalf("error calling InsertRecord(): %s", err)
		}

	}
}

func GetAllAlkData(c *api.AlkatronicClient, pg *PostgresAlkatronic, days int) {
	devices, err := c.GetDevices()
	if err != nil {
		log.Fatalf("Error getting devices: %s", err)
	}

	for _, d := range devices.Data {

		records, err := c.GetRecords(d.DeviceID, days)
		if err != nil {
			log.Fatalf("Error getting device records: %s", err)
		}

		for _, r := range records.Data {
			r.KhValue   = api.ConvertKh(r.KhValue)

			device := api.Device {
				UpperKh: api.ConvertKh(d.UpperKh),
				LowerKh: api.ConvertKh(d.LowerKh),
			}

			log.Printf(
				"Latest test result: Device Name: %s, Record ID: %d, KH_Value: %.2f, Create Time: %s\n",
				d.FriendlyName,
				r.RecordID,
				api.ConvertKh(r.KhValue),
				time.Unix(r.CreateTime, 0).Format(time.RFC822Z))

			err = pg.InsertRecord(r, device)
			if err != nil {
				log.Fatalf("error calling InsertRecord(): %s", err)
			}
		}

	}

}

func boolConverter(i int) bool {
	if i == 0 {
		return false
	} else {
		return true
	}
}
