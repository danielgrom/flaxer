package main

import (
	"io"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (app *Config) CreateNewsTab() *fyne.Container {
	title := widget.NewRichText()
	title.ParseMarkdown("# News")
	separator := widget.NewSeparator()
	news := widget.NewRichTextFromMarkdown("")
	siteNews, err := getNews(app)
	if err != nil {
		dialog.NewError(err, app.MainWindow)
	}
	news.AppendMarkdown(siteNews)
	ret := container.NewVBox(title, separator, news)
	return ret
}

func getNews(app *Config) (string, error) {

	if app.HTTPClient == nil {
		app.HTTPClient = &http.Client{}
	}
	ret := ""
	url := "https://danielgrom.github.io/Flaxer/news.md"
	req, _ := http.NewRequest("GET", url, nil)

	resp, err := app.HTTPClient.Do(req)
	if err != nil {
		app.ErrorLog.Println("error contacting github: ", err)
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		app.ErrorLog.Println("error reading file", err)
		return "", err
	}

	ret = string(body[:])
	return ret, err
}
