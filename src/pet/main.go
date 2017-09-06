package main

import (
	"log"
	"runtime"
	"pet/app/shared/config/jsonconfig"
	"encoding/json"
	"os"
	"pet/app/shared/server"
	"pet/app/shared/database"
	"pet/app/model"
	"pet/app/route"
	"pet/app/executor"
	"pet/app/shared/context"
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
	filepath := "config" + string(os.PathSeparator) + "config.json"
	jsonconfig.Load(filepath, config)

	// Connect to database
	database.Connect(config.Database)

	// Create initial DB entities
	createDefaultDBEntities()

	//TODO: on application init:
	//TODO: fetch all methods from db through dao
	//TODO: all subsequent db operations to be done via methods (whats that?)
	//Loading API methods to App context
	loadApiMethodsToCtx()

	//Loading api methods to FSM
	loadApiMethodsToFsm()

	// Configure API endpoint, register handlers
	route.ConfigRoutes()

	// Starting server using configuration from config
	route.StartServer(&config.Server)
}


// *****************************************************************************
//  Settings
// *****************************************************************************

// config the settings variable
var config = &configuration{}

// configuration contains the application settings
type configuration struct {
	Database database.Info   `json:"Database"`
	Server   server.Server   `json:"Server"`
}

// ParseJSON unmarshals bytes to structs
func (c *configuration) ParseJSON(b []byte) error {
	return json.Unmarshal(b, &c)
}


// *****************************************************************************
//  DB preparations
// *****************************************************************************

// createDefaultDBEntities
// Creates first permissions, group and user and saves them to DB (if not created yet)
func createDefaultDBEntities() {
	if _, err := model.MethodByName("auth"); err != nil {
		err = model.CreateMethod("auth", []model.Parameter{
			model.Parameter{Name: "login", Required: true, Type: model.ParamTypeString},
			model.Parameter{Name: "pass", Required: true, Type: model.ParamTypeString },
		}, []byte{});
	}

	initialPerms := getInitPermissions()
	for _, p := range initialPerms {
		if _, err := model.PermissionByName(p.Name); err != nil {
			//if permission not found create them
			model.PermissionCreate(p.Name, p.Type, p.Value, p.Read, p.Update, p.Execute)
		}
	}
	var permissions []model.Permission
	for _, p := range initialPerms {
		if permission, err := model.PermissionByName(p.Name); err == nil {
			permissions = append(permissions, *permission)
		} else {
			log.Printf("Error occured during retieving permission group '%s'", permission)
		}
	}

	// Create initial permission group
	groupName := "initial"
	group, err := model.GroupByName(groupName);
	if err != nil {
		model.GroupCreate(groupName, permissions)
		//filling group again
		group, _ = model.GroupByName(groupName);
	}
	//appending group to groups array
	groups := []model.Group{}
	groups = append(groups, *group)

	// Creating initial user with initial permission groups
	login := "Fluggegecheimen"
	password := login
	userName := "The Bandit"
	if _, err := model.UserByLogin(login); err != nil {
		model.UserCreate(login, userName, password, groups)
	}
}
// *****************************************************************************
//  Initial context configuration
// *****************************************************************************

func loadApiMethodsToFsm() {
	executor.LoadFSM(context.AppContext.GetAllMethods())
}

func loadApiMethodsToCtx() {
	context.AppContext.InitContext()
	methods := *model.GetAllMethods()
	for _, method := range methods{
		context.AppContext.Put(method.Name, method)
	}
}

// *****************************************************************************
//  Initial database values
// *****************************************************************************
func getInitPermissions() []model.Permission {
	return []model.Permission{
		{Name: "readMethod", Type: model.TypeMethod,
			Value: "readSomeMethod",
			Read:true,
			Update:false,
			Execute:false,
		},
		{Name: "executeMethod",
			Type: model.TypeMethod,
			Value: "executeSomeMethod",
			Read:true,
			Update:false,
			Execute:true,
		}, {
			Name: "updateField",
			Type: model.TypeField,
			Value: "updateSomeField",
			Read:true,
			Update:true,
			Execute:false,
		},
	}
}