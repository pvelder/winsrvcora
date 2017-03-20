package oracle

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/fsnotify/fsnotify"
	_ "github.com/mattn/go-oci8"
)

type tomlConfig struct {
	DB database   `toml:"database"`
	FS filesystem `toml:"filesystem"`
}

type database struct {
	Type           string
	Connectstring  string
	ApplyStatement string
}

type filesystem struct {
	WatchDirectory string
	LogFile        string
}

var config tomlConfig

func Apply(filename string) {
	// should we use a dedicated connection to the database?
	// or only open the connection when archive logs are detected?
	db, err := sql.Open(config.DB.Type, config.DB.Connectstring)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		fmt.Printf("Error connecting to the database: %s\n", err)
		return
	}

	applyTemplate := config.DB.ApplyStatement
	cmd := fmt.Sprintf(applyTemplate, filename)
	log.Println("Start Applying: " + cmd)
	result, error := db.Exec(cmd)
	if error != nil {
		log.Println("Error:")
		log.Println(error)
		log.Println("Result:")
		log.Println(result)
		return
	}
	log.Println("Applied ...")

}

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

	// log.Println("Three")

	go func() {
		for {
			select {
			case event := <-watcher.Events:

				log.Println("event ouside :", event)

				if event.Op&fsnotify.Create == fsnotify.Create {

					filename := event.Name

					// wait until file is fully written
					// http://stackoverflow.com/questions/13434555/hotfolder-in-go-wait-for-file-to-be-written
					for {
						timer := time.NewTimer(1 * time.Second)
						select {
						case ev := <-watcher.Events:
							log.Println("event inside", ev)
						case err := <-watcher.Errors:
							log.Println("error inside:", err)
						case <-timer.C:
							Apply(filename)
						}
						timer.Stop()
					}
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
