package main

import (
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	downloader "github.com/Girbons/comics-downloader/cmd/app"
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

var release string // sha1 revision used to build the program

func main() {
	options := []string{"pdf", "epub", "cbr", "cbz"}
	imagesFormat := []string{"png", "jpg", "img"}

	app := app.New()
	w := app.NewWindow(fmt.Sprintf("Comics Downloader %s", release))

	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("Comic URL or URLs separated by a comma")

	countryEntry := widget.NewEntry()
	countryEntry.SetPlaceHolder("Country param used by mangarock")
	countryEntry.Hide()

	choices := widget.NewRadio(options, nil)
	choices.SetSelected("pdf")

	imagesFormatChoices := widget.NewRadio(imagesFormat, nil)
	imagesFormatChoices.SetSelected("jpg")

	allChaptersCheck := widget.NewCheck("", nil)
	lastChapterCheck := widget.NewCheck("", nil)
	imagesOnlyCheck := widget.NewCheck("", nil)

	outputFolderEntry := widget.NewEntry()
	outputFolderEntry.SetPlaceHolder("Folder where the comics will be saved")

	d := &Downloader{
		URL:          urlEntry,
		Country:      countryEntry,
		Format:       choices,
		AllChapters:  allChaptersCheck,
		LastChapter:  lastChapterCheck,
		ImagesOnly:   imagesOnlyCheck,
		ImagesFormat: imagesFormatChoices,
		OutputFolder: outputFolderEntry,
	}

	clearCountryFieldButton := widget.NewButton("Clear Country", func() {
		d.ClearCountryField()
	})
	clearCountryFieldButton.Hide()

	showCountryOption := widget.NewCheck("Show Country Options", func(on bool) {
		if on {
			countryEntry.Show()
			clearCountryFieldButton.Show()
		} else {
			countryEntry.Hide()
			clearCountryFieldButton.Hide()
		}
	})

	form := widget.NewForm()
	form.Append("URL", d.URL)
	form.Append("Country", d.Country)
	form.Append("", showCountryOption)
	form.Append("Output", d.Format)
	form.Append("All chapters", d.AllChapters)
	form.Append("Last chapter", d.LastChapter)
	form.Append("Images Only", d.ImagesOnly)
	form.Append("Images Format", d.ImagesFormat)
	form.Append("Output Folder", d.OutputFolder)

	box := widget.NewVBox()

	clearURLFieldButton := widget.NewButton("Clear URL", func() {
		d.ClearURLField()
	})

	clearLogsButton := widget.NewButton("Clear Logs", func() {
		box.Children = make([]fyne.CanvasObject, 0)
		widget.Refresh(box)
	})

	clearOutputFolderButton := widget.NewButton("Clear Output Folder", func() {
		d.ClearOutputFolderField()
	})

	submitButton := widget.NewButton("Download", func() {
		d.Submit()
	})
	submitButton.Style = widget.PrimaryButton

	buttons := widget.NewHBox(
		clearURLFieldButton,
		clearLogsButton,
		clearCountryFieldButton,
		clearOutputFolderButton,
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
