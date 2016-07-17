package main

import (
	"os"
	"fmt"
	"time"
	"strconv"
	"encoding/json"

	"github.com/jhoonb/archivex"
	"github.com/jasonlvhit/gocron"
)

type Configuration struct {
	Backup   string
	Interval int
}

func main() {

	file, _ := os.Open("conf.json")
	decoder := json.NewDecoder(file)
	config := Configuration{}
	err := decoder.Decode(&config)

	if err != nil {
		fmt.Println("[Backup] Unable to load conf.json, are you sure it's present? Error:", err.Error())
		os.Exit(1)
	}

	fmt.Println("[Backup] Started up... Preparing to archive:", config.Backup, "Every:", config.Interval, "seconds.")

	backup := func() {
		zip := new(archivex.ZipFile)
		zipName := "world@" + strconv.Itoa(int(time.Now().Unix()))
		zip.Create(zipName)
		zip.AddAll(config.Backup, true)
		zip.Close()
		fmt.Println("[Backup] Successfully archived:", config.Backup, "Id:", zipName)
	}

	s := gocron.NewScheduler()
	s.Every(uint64(config.Interval)).Seconds().Do(backup)
	<-s.Start()

}