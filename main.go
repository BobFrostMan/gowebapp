package main

import (
	"log"
	"runtime"
	"gowebapp/pet/app/shared/config/jsonconfig"
	"encoding/json"
	"os"
	"gowebapp/pet/app/shared/server"
	"gowebapp/pet/app/shared/database"
)

// *****************************************************************************
// Application Logic
// *****************************************************************************

func init() {
	// Verbose logging with file name and line number
	log.SetFlags(log.Lshortfile)

	// Use all CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	// Parsing app configurations
	filepath := "config"+string(os.PathSeparator)+"config.json"
	jsonconfig.Load(filepath, config)

	// Connect to database
	database.Connect(config.Database)
}

// *****************************************************************************
//  Settings
// *****************************************************************************

// config the settings variable
var config = &configuration{}

// configuration contains the application settings
type configuration struct {
	Database  database.Info   `json:"Database"`
	Server    server.Server   `json:"Server"`
}

// ParseJSON unmarshals bytes to structs
func (c *configuration) ParseJSON(b []byte) error {
	return json.Unmarshal(b, &c)
}