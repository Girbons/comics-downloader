package main

import (
	"fyne.io/fyne/widget"

	downloader "github.com/Girbons/comics-downloader/cmd/app"
	"github.com/Girbons/comics-downloader/pkg/config"
)

// Downloader represents a downloader instance.
type Downloader struct {
	URL               *widget.Entry
	Country           *widget.Entry
	Format            *widget.Radio
	AllChapters       *widget.Check
	LastChapter       *widget.Check
	ImagesOnly        *widget.Check
	ImagesFormat      *widget.Radio
	OutputFolder      *widget.Entry
	CreateDefaultPath *widget.Check
	IssuesRange       *widget.Entry
	Debug             *widget.Check
	CustomComicName   *widget.Entry
}

// ClearURLField resets the url text field.
func (d *Downloader) ClearURLField() {
	d.URL.SetText("")
}

// ClearCountryField resets the country text field.
func (d *Downloader) ClearCountryField() {
	d.Country.SetText("")
}

// ClearOutputFolderField resets the output text field.
func (d *Downloader) ClearOutputFolderField() {
	d.OutputFolder.SetText("")
}

// Submit calls the downloader api with the given options.
func (d *Downloader) Submit() {
	opts := &config.Options{
		Debug:             d.Debug.Checked,
		All:               d.AllChapters.Checked,
		Last:              d.LastChapter.Checked,
		URL:               d.URL.Text,
		Format:            d.Format.Selected,
		Country:           d.Country.Text,
		ImagesFormat:      d.ImagesFormat.Selected,
		ImagesOnly:        d.ImagesOnly.Checked,
		OutputFolder:      d.OutputFolder.Text,
		CreateDefaultPath: d.CreateDefaultPath.Checked,
		IssuesRange:       d.IssuesRange.Text,
		CustomComicName:   d.CustomComicName.Text,
	}

	downloader.GuiRun(opts)
}
