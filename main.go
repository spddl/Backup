package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jasonlvhit/gocron"
	"github.com/jhoonb/archivex"
)

type configuration struct {
	Path     string `json:"path"`
	Interval int    `json:"interval"`
	Files    []struct {
		Name string `json:"name"`
		Path string `json:"path"`
	} `json:"files"`
}

func main() {

	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	config := configuration{}
	err := decoder.Decode(&config)

	if err != nil {
		fmt.Println("[Backup] Unable to load config.json, are you sure it's present? Error:", err.Error())
		os.Exit(1)
	}

	names := make([]string, 0)
	for _, json := range config.Files {
		names = append(names, json.Name)
	}

	fmt.Println("[Backup] Started up... Every:", config.Interval, "seconds. Preparing to archive:", strings.Join(names[:], ", "))

	backup := func() {
		for _, json := range config.Files {
			zip := new(archivex.ZipFile)
			zipName := config.Path + json.Name + "@" + strconv.Itoa(int(time.Now().Unix()))
			zip.Create(zipName)
			zip.AddAll(json.Path, true)
			zip.Close()
			fmt.Println("[Backup] Successfully archived:", zipName+".zip")
		}
	}

	backup()
	s := gocron.NewScheduler()
	s.Every(uint64(config.Interval)).Seconds().Do(backup)
	<-s.Start()
}
