package main

import (
	"fyne.io/fyne/widget"

	downloader "github.com/Girbons/comics-downloader/cmd/app"
	"github.com/Girbons/comics-downloader/pkg/config"
)

type Downloader struct {
	URL          *widget.Entry
	Country      *widget.Entry
	Format       *widget.Radio
	AllChapters  *widget.Check
	LastChapter  *widget.Check
	ImagesOnly   *widget.Check
	ImagesFormat *widget.Radio
	OutputFolder *widget.Entry
	IssuesRange  *widget.Entry
	Debug        *widget.Check
}

func (d *Downloader) ClearURLField() {
	d.URL.SetText("")
}

func (d *Downloader) ClearCountryField() {
	d.Country.SetText("")
}

func (d *Downloader) ClearOutputFolderField() {
	d.OutputFolder.SetText("")
}

func (d *Downloader) Submit() {
	opts := &config.Options{
		Debug:        d.Debug.Checked,
		All:          d.AllChapters.Checked,
		Last:         d.LastChapter.Checked,
		Url:          d.URL.Text,
		Format:       d.Format.Selected,
		Country:      d.Country.Text,
		ImagesFormat: d.ImagesFormat.Selected,
		ImagesOnly:   d.ImagesOnly.Checked,
		OutputFolder: d.OutputFolder.Text,
		IssuesRange:  d.IssuesRange.Text,
	}

	downloader.GuiRun(opts)
}
