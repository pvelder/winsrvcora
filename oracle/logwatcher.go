package oracle

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/fsnotify/fsnotify"
	_ "github.com/mattn/go-oci8"
)

type tomlConfig struct {
	DB database   `toml:"database"`
	FS filesystem `toml:"filesystem"`
}

type database struct {
	Type          string
	Connectstring string
}

type filesystem struct {
	WatchDirectory string
	LogFile        string
}

var config tomlConfig

func DoWatchLogs(directoryToWatch string) {

	// logging to file
	// TODO configure
	f, err := os.OpenFile(config.FS.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		// TODO : log to service default output?
		// t.Fatalf("error opening file: %v", err)
		fmt.Printf("error opening file: %v", err)

	}
	defer f.Close()

	log.SetOutput(f)

	// log.Println("One")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// log.Println("Two")

	done := make(chan bool)

	// should we use a dedicated connection to the database?
	// or only open the connection when archive logs are detected?
	db, err := sql.Open(config.DB.Type, config.DB.Connectstring)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	// log.Println("Three")

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				// log.Println("event:", event)
				if event.Op&fsnotify.Create == fsnotify.Create {
					cmd := fmt.Sprintf("alter database register logfile '%s'", event.Name)
					log.Println(cmd)
					// db.Exec(cmd)
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	// log.Println("Four")

	err = watcher.Add(directoryToWatch)
	log.Println("Watching directory: " + directoryToWatch)
	if err != nil {
		log.Println("Error watching directory: " + directoryToWatch)

		log.Fatal(err)
	}
	<-done
}

func WatchOracleArchiveLogs() {

	// TODO refer to toml.config using command line argument,
	// falling back to default "current" directory, or default e.g ${home}/.../config.toml ?
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		fmt.Println(err)
		return
	}

	DoWatchLogs(config.FS.WatchDirectory)

}
