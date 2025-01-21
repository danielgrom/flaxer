package main

import (
	"database/sql"
	"flaxer/repository"
	"log"
	"net/http"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	_ "github.com/glebarez/go-sqlite"
)

type Config struct {
	App         fyne.App
	InfoLog     *log.Logger
	ErrorLog    *log.Logger
	MainWindow  fyne.Window
	SettingsTab *fyne.Container
	NewsTab     *fyne.Container
	LearnTab    *fyne.Container
	ProjectsTab *fyne.Container
	Settings    *repository.ProjectSettings
	DB          repository.Repository
	HTTPClient  *http.Client
}

func main() {
	// ... rest of the program ...
	var myApp Config
	fyneApp := app.NewWithID("project.go.flaxer")
	myApp.App = fyneApp
	myApp.HTTPClient = &http.Client{}
	myApp.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	myApp.ErrorLog = log.New(os.Stdout, "Error\t", log.Ldate|log.Lshortfile)

	//open a connection to a database
	sqlDB, err := myApp.connectSQL()
	if err != nil {
		log.Panic(err)
	}
	//create a database repository
	myApp.setupDB(sqlDB)

	//load settings
	myApp.Settings = &repository.ProjectSettings{ProjectsDirectory: "", FlaxLocation: ""}
	savedSettings := GetSavedFlaxerSettings(&myApp)
	myApp.Settings.FlaxLocation = savedSettings.FlaxLocation
	myApp.Settings.ProjectsDirectory = savedSettings.ProjectsDirectory
	myApp.MainWindow = fyneApp.NewWindow("Flaxer")
	myApp.MainWindow.Resize(fyne.NewSize(1000, 500))
	myApp.MainWindow.SetFixedSize(false)
	myApp.MainWindow.SetMaster()

	myApp.makeUI()
	myApp.MainWindow.ShowAndRun()
}

func (app *Config) connectSQL() (*sql.DB, error) {
	path := ""

	if os.Getenv("FLAXERDB_PATH") != "" {
		path = os.Getenv("FLAXERDB_PATH")
	} else {
		path = app.App.Storage().RootURI().Path() + "/sql.db"
		app.InfoLog.Println("db in:", path)
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (app *Config) setupDB(sqlDB *sql.DB) {
	app.DB = repository.NewSQLiteRepository(sqlDB)

	err := app.DB.Migrate()
	if err != nil {
		app.ErrorLog.Println(err)
		log.Panic()
	}
}
