package main

import (
	"flaxer/repository"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func (app *Config) CreateSettingsTab(settings repository.ProjectSettings) *fyne.Container {
	title := widget.NewRichText()
	title.ParseMarkdown("# Settings")
	separator := widget.NewSeparator()
	projLocationEntry := widget.NewEntry()
	savedSettings := GetSavedFlaxerSettings(app)
	projLocationEntry.SetText(savedSettings.ProjectsDirectory)
	projLocationEntry.PlaceHolder = "Location of Flax Projects"
	projLocation := widget.NewFormItem("Flax Projects Directory...", projLocationEntry)
	projdialog := dialog.NewFolderOpen(func(file fyne.ListableURI, err error) {
		if file != nil {
			savedSettings = GetSavedFlaxerSettings(app)
			projLocationEntry.SetText(file.Path())
			savedSettings.ProjectsDirectory = file.Path()
			SaveFlaxerSettings(app, savedSettings)
		}
	}, app.MainWindow)
	projButton := widget.NewButton("Browse Projects Directory...", func() { projdialog.Show() })
	projFormButton := widget.NewFormItem("", projButton)

	flaxLocationEntry := widget.NewEntry()
	flaxLocationEntry.Text = savedSettings.FlaxLocation
	flaxLocation := widget.NewFormItem("Flax Executable Location", flaxLocationEntry)
	flaxdialog := dialog.NewFileOpen(func(file fyne.URIReadCloser, err error) {
		if file != nil {
			savedSettings = GetSavedFlaxerSettings(app)
			flaxLocationEntry.SetText(file.URI().Path())
			savedSettings.FlaxLocation = file.URI().Path()
			SaveFlaxerSettings(app, savedSettings)
		}
	}, app.MainWindow)
	extensions := []string{".exe", ".bin", ""}
	flaxdialog.SetFilter(storage.NewExtensionFileFilter(extensions))
	flaxButton := widget.NewButton("Browse Flax Executable", func() { flaxdialog.Show() })
	flaxForButton := widget.NewFormItem("", flaxButton)

	form := widget.NewForm(projLocation, projFormButton, flaxLocation, flaxForButton)

	ret := container.NewVBox(title, separator, form)
	return ret
}

func GetSavedFlaxerSettings(app *Config) *repository.FlaxerSettings {
	settings, err := app.DB.GetFlaxerSettings()
	if err != nil {
		// Only show dialog if MainWindow is initialized
		if app.MainWindow != nil {
			dialog.ShowError(err, app.MainWindow)
		}
		// Return default settings on error
		defaultSettings := &repository.FlaxerSettings{ID: 0, ProjectsDirectory: "", FlaxLocation: ""}
		return defaultSettings
	}

	if settings == nil {
		// No settings found in database - create default settings
		defaultSettings := &repository.FlaxerSettings{ID: 0, ProjectsDirectory: "", FlaxLocation: ""}
		savedSettings, err := app.DB.InsertFlaxerSettings(*defaultSettings)
		if err != nil {
			// If insert fails, just return the default settings
			if app.MainWindow != nil {
				dialog.ShowError(err, app.MainWindow)
			}
			return defaultSettings
		}
		return savedSettings
	}
	return settings
}

func SaveFlaxerSettings(app *Config, fs *repository.FlaxerSettings) {
	app.DB.UpdateFlaxerSettings(fs.ID, *fs)

	// Update the app settings
	oldProjectsDir := app.Settings.ProjectsDirectory
	app.Settings.ProjectsDirectory = fs.ProjectsDirectory
	app.Settings.FlaxLocation = fs.FlaxLocation

	// If projects directory changed and projects tab exists, refresh it
	if oldProjectsDir != fs.ProjectsDirectory && app.RefreshProjects != nil {
		app.RefreshProjects()
	}
}
