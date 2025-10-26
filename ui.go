package main

import (
	"flaxer/repository"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (app *Config) makeUI() {
	img := canvas.NewImageFromResource(resourceIconPng)
	img.SetMinSize(fyne.Size{Width: 150, Height: 150})

	listEntries := [4]string{"News", "Learn", "Projects", "Settings"}
	list := widget.NewList(
		func() int { return len(listEntries) },
		func() fyne.CanvasObject { return widget.NewLabel("Template Item") },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			entry := listEntries[i]
			o.(*widget.Label).SetText(entry)
		},
	)
	s := repository.ProjectSettings{ProjectsDirectory: "", FlaxLocation: ""}
	var cont *fyne.Container = container.NewStack(nil)

	// Initialize the RefreshProjects function
	app.RefreshProjects = func() {
		// Only refresh if projects tab is currently displayed
		if app.ProjectsTab != nil && len(cont.Objects) > 0 {
			// Check if the current tab is the projects tab
			currentTab := cont.Objects[0]
			if currentTab == app.ProjectsTab {
				// Recreate the projects tab
				app.ProjectsTab.RemoveAll()
				app.ProjectsTab = nil
				app.ProjectsTab = app.CreateProjectsTab()
				cont.Objects = []fyne.CanvasObject{app.ProjectsTab}
				cont.Refresh()
			}
		}
	}
	list.Select(0)
	app.NewsTab = app.CreateNewsTab()
	cont.Objects = []fyne.CanvasObject{app.NewsTab}
	cont.Refresh()
	list.OnSelected = func(id widget.ListItemID) {
		switch id {
		case 0:
			//News
			if app.NewsTab != nil {
				app.NewsTab.RemoveAll()
				app.NewsTab = nil
			}
			app.NewsTab = app.CreateNewsTab()
			cont.Objects = nil
			cont.Objects = []fyne.CanvasObject{app.NewsTab}
			cont.Refresh()
		case 1:
			//Learn
			if app.LearnTab != nil {
				app.LearnTab.RemoveAll()
				app.LearnTab = nil
			}
			app.LearnTab = app.CreateLearnTab()
			cont.Objects = nil
			cont.Objects = []fyne.CanvasObject{app.LearnTab}
			cont.Refresh()
		case 2:
			//Projects
			if app.ProjectsTab != nil {
				app.ProjectsTab.RemoveAll()
				app.ProjectsTab = nil
			}
			app.ProjectsTab = app.CreateProjectsTab()
			cont.Objects = nil
			cont.Objects = []fyne.CanvasObject{app.ProjectsTab}
			cont.Refresh()
		case 3:
			// Settings
			if app.SettingsTab != nil {
				app.SettingsTab.RemoveAll()
				app.SettingsTab = nil
			}
			app.SettingsTab = app.CreateSettingsTab(s)
			cont.Objects = nil
			cont.Objects = []fyne.CanvasObject{app.SettingsTab}
			cont.Refresh()
		default:
		}
	}

	listContainer := container.NewVScroll(list)
	listContainer.SetMinSize(fyne.NewSize(200, 600)) // Set a minimum height to ensure it fills the space

	leftGrid := container.NewVBox(
		container.NewCenter(img),
		listContainer,
	)

	rightGrid := container.NewVScroll(
		cont,
	)

	ncolor := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	registerBackground := canvas.NewRectangle(color.RGBA{R: 0, G: 126, B: 200, A: 255})
	registerBackground.SetMinSize(fyne.Size{Width: 1, Height: 5})
	statusBar := canvas.NewText("Ok", ncolor)

	content := container.NewBorder(nil, container.NewStack(registerBackground, statusBar), leftGrid, nil, rightGrid)

	app.MainWindow.SetContent(content)
}
