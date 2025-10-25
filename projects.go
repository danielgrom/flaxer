package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/karrick/godirwalk"
)

type ProjectFile struct {
	Name  string
	Image string
	Path  string
}

func (app *Config) CreateProjectsTab() *fyne.Container {
	title := widget.NewRichText()
	title.ParseMarkdown("# Projects")
	separator := widget.NewSeparator()

	projectList, err := GetProjectList(app.Settings.ProjectsDirectory)

	if err != nil {
		dialog.NewError(err, app.MainWindow)
	}

	grid := widget.NewGridWrapWithData(*projectList,
		func() fyne.CanvasObject {
			text := widget.NewLabel("Title")
			text.Alignment = fyne.TextAlignCenter
			defaultIcon := canvas.NewImageFromResource(resourceIconPng)
			defaultIcon.FillMode = canvas.ImageFillOriginal
			content := container.NewBorder(nil, text, nil, nil, defaultIcon)
			return container.NewGridWrap(fyne.NewSquareSize(250), content)
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			value, _ := i.(binding.Untyped).Get()
			item := value.(ProjectFile)

			image := o.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*canvas.Image)
			image.File = item.Image
			image.Resource = nil

			label := o.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label)
			label.SetText(item.Name)
		})
	var selectedItem *ProjectFile
	newBtn := widget.NewButton("Create New Project", func() {
		//launch the flax engine
		pname := widget.NewEntry()
		pname.SetText("newproject")
		dlg := dialog.NewCustomConfirm(
			"Choose a new project name",
			"Confirm",
			"Cancel",
			pname,
			func(b bool) {
				if b {
					projectFile := app.Settings.ProjectsDirectory + "/" + pname.Text + "/" + pname.Text + ".flaxproj"
					println(projectFile)
					if !doesFileExist(projectFile) {
						// Check if FlaxLocation is set
						if app.Settings.FlaxLocation == "" {
							dialog.ShowInformation("Error", "Flax executable location not set. Please configure it in Settings.", app.MainWindow)
							return
						}
						// Properly split arguments for exec.Command
						cmd := exec.Command(app.Settings.FlaxLocation, "-new", "-project", app.Settings.ProjectsDirectory+"/"+pname.Text)
						if err := cmd.Start(); err != nil {
							log.Printf("Error starting Flax: %v", err)
							dialog.ShowError(err, app.MainWindow)
							return
						}
						app.MainWindow.Close()
					} else {
						dialog.ShowInformation("Error", "Project already exists", app.MainWindow)
					}
				}
			},
			app.MainWindow)
		dlg.Show()
	})

	openBtn := widget.NewButton("Open Selected Project", func() {
		//launch the flax engine
		if selectedItem != nil {
			// Check if FlaxLocation is set
			if app.Settings.FlaxLocation == "" {
				dialog.ShowInformation("Error", "Flax executable location not set. Please configure it in Settings.", app.MainWindow)
				return
			}
			// Properly split arguments for exec.Command
			cmd := exec.Command(app.Settings.FlaxLocation, "-std", "-project", selectedItem.Path)
			if err := cmd.Start(); err != nil {
				log.Printf("Error starting Flax: %v", err)
				dialog.ShowError(err, app.MainWindow)
				return
			}
			app.MainWindow.Close()
		}
	})
	openBtn.Disable()
	grid.OnSelected = func(id widget.ListItemID) {
		if value, err := (*projectList).GetValue(id); err == nil {
			item := value.(ProjectFile)
			selectedItem = &item
			openBtn.Enable()
		} else {
			fmt.Println("Error getting project:", err)
		}
	}

	grid.OnUnselected = func(id widget.ListItemID) {
		selectedItem = nil
		openBtn.Disable()
	}

	titl := container.NewVBox(title, newBtn, openBtn, separator)
	grd := container.NewStack(grid)
	ret := container.NewBorder(titl, nil, nil, nil, grd)

	return ret
}

func GetProjectList(directory string) (*binding.UntypedList, error) {
	ret := binding.NewUntypedList()
	err := godirwalk.Walk(directory, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if de.IsRegular() {
				if strings.Contains(osPathname, ".flaxproj") {
					path := filepath.Dir(osPathname)
					image := path + "/Preview.png"
					if !doesFileExist(image) {
						image = "icon.png"
					}
					file := ProjectFile{
						Path:  osPathname,
						Image: image,
						Name:  strings.TrimSuffix(de.Name(), ".flaxproj")}
					ret.Append(file)
				}
			}
			return nil
		},
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			return godirwalk.SkipNode
		},
	})
	return &ret, err
}

func doesFileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	// check if error is "file not exists"
	return !os.IsNotExist(err)
}
