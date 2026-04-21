package main

import (
	"fmt"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	downloader "github.com/Girbons/comics-downloader/cmd/app"
	"github.com/Girbons/comics-downloader/internal/version"
)

func watchLogs(logSection *container.Scroll, box *widget.Box, statusLabel *widget.Label) {
	for message := range downloader.Messages {
		trimmed := strings.TrimSpace(message)
		if trimmed == "" {
			continue
		}

		level := ""
		text := trimmed
		if parts := strings.SplitN(trimmed, ":", 2); len(parts) == 2 {
			level = strings.ToUpper(strings.TrimSpace(parts[0]))
			text = strings.TrimSpace(parts[1])
		}

		switch level {
		case "INFO", "DEBUG":
			if text == "" {
				statusLabel.SetText(level)
			} else {
				statusLabel.SetText(text)
			}
		case "WARNING", "ERROR":
			statusLabel.SetText(trimmed)
			box.Append(widget.NewLabel(trimmed))
			logSection.ScrollToBottom()
		default:
			statusLabel.SetText(trimmed)
		}
	}
}

func appStatus(downloadButton *widget.Button, progress *widget.ProgressBarInfinite, statusLabel *widget.Label) {
	for running := range downloader.AppStatus {
		if running {
			downloadButton.Disable()
			if !progress.Visible() {
				progress.Show()
			}
			if !progress.Running() {
				progress.Start()
			}
			statusLabel.SetText("Downloading...")
		} else {
			progress.Hide()
			downloadButton.Enable()
			if statusLabel.Text == "Downloading..." {
				statusLabel.SetText("Ready")
			}
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

	statusLabel := widget.NewLabel("Ready")
	statusLabel.Wrapping = fyne.TextWrapWord

	progress := widget.NewProgressBarInfinite()
	progress.Hide()

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
		box.Refresh()
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

	// logSection := widget.NewScrollContainer(box)
	logSection := container.NewScroll(box)

	footer := container.NewVBox(statusLabel, progress, buttons)

	go watchLogs(logSection, box, statusLabel)
	go appStatus(submitButton, progress, statusLabel)

	w.SetContent(fyne.NewContainerWithLayout(layout.NewBorderLayout(form, footer, nil, nil), form, footer, logSection))
	w.Resize(fyne.NewSize(800, 400))
	w.ShowAndRun()
}
