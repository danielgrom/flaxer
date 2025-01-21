package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (app *Config) CreateLearnTab() *fyne.Container {
	title := widget.NewRichText()
	title.ParseMarkdown("# Learn")
	separator := widget.NewSeparator()
	ret := container.NewVBox(title, separator)
	return ret
}
