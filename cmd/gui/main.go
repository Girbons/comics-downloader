package main

import (
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	downloader "github.com/Girbons/comics-downloader/cmd/app"
	"github.com/Girbons/comics-downloader/internal/version"
)

func watchLogs(logSection *widget.ScrollContainer, box *widget.Box) {
	for {
		box.Append(widget.NewLabel(<-downloader.Messages))
		logSection.Resize(logSection.Size())
	}
}

func appStatus(downloadButton *widget.Button) {
	for {
		if <-downloader.AppStatus {
			downloadButton.Disable()
		} else {
			downloadButton.Enable()
		}
	}
}

func main() {
	options := []string{"pdf", "epub", "cbr", "cbz"}
	imagesFormat := []string{"png", "jpg", "img"}

	app := app.New()
	w := app.NewWindow(fmt.Sprintf("Comics Downloader %s", version.Tag))

	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("Comic URL or URLs separated by a comma")

	countryEntry := widget.NewEntry()
	countryEntry.SetPlaceHolder("Country param used by mangadex which uses ISO 3166-1 codes")

	choices := widget.NewRadioGroup(options, nil)
	choices.SetSelected("pdf")

	imagesFormatChoices := widget.NewRadioGroup(imagesFormat, nil)
	imagesFormatChoices.SetSelected("jpg")

	allChaptersCheck := widget.NewCheck("", nil)
	lastChapterCheck := widget.NewCheck("", nil)
	imagesOnlyCheck := widget.NewCheck("", nil)

	createDefaultPath := widget.NewCheck("", nil)
	createDefaultPath.Checked = true

	debugCheck := widget.NewCheck("", nil)

	outputFolderEntry := widget.NewEntry()
	outputFolderEntry.SetPlaceHolder("Folder where the comics will be saved")

	customComicName := widget.NewEntry()
	customComicName.SetPlaceHolder("Custom comic name")

	issuesRange := widget.NewEntry()
	issuesRange.SetPlaceHolder("1-10")

	d := &Downloader{
		URL:               urlEntry,
		Country:           countryEntry,
		Format:            choices,
		AllChapters:       allChaptersCheck,
		LastChapter:       lastChapterCheck,
		ImagesOnly:        imagesOnlyCheck,
		ImagesFormat:      imagesFormatChoices,
		CreateDefaultPath: createDefaultPath,
		OutputFolder:      outputFolderEntry,
		IssuesRange:       issuesRange,
		Debug:             debugCheck,
		CustomComicName:   customComicName,
	}

	form := widget.NewForm()
	form.Append("URL", d.URL)
	form.Append("Country", d.Country)
	form.Append("Custom comic name", d.CustomComicName)
	form.Append("Output", d.Format)
	form.Append("All chapters", d.AllChapters)
	form.Append("Last chapter", d.LastChapter)
	form.Append("Debug Mode", d.Debug)
	form.Append("Issues Range", d.IssuesRange)
	form.Append("Images Only", d.ImagesOnly)
	form.Append("Images Format", d.ImagesFormat)
	form.Append("Output Folder", d.OutputFolder)
	form.Append("Create Default Download Path", d.CreateDefaultPath)

	box := widget.NewVBox()

	clearLogsButton := widget.NewButton("Clear Logs", func() {
		box.Children = make([]fyne.CanvasObject, 0)
		widget.Refresh(box)
	})

	submitButton := widget.NewButton("Download", func() {
		d.Submit()
	})
	submitButton.Style = widget.PrimaryButton

	buttons := widget.NewHBox(
		clearLogsButton,
		layout.NewSpacer(),
		submitButton,
	)

	logSection := widget.NewScrollContainer(box)

	go watchLogs(logSection, box)
	go appStatus(submitButton)

	w.SetContent(fyne.NewContainerWithLayout(layout.NewBorderLayout(form, buttons, nil, nil), form, buttons, logSection))
	w.Resize(fyne.NewSize(800, 400))
	w.ShowAndRun()
}
