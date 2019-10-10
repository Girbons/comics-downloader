package main

import (
	"fyne.io/fyne/widget"
	downloader "github.com/Girbons/comics-downloader/cmd/app"
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
	downloader.GuiRun(
		d.URL.Text,
		d.Format.Selected,
		d.Country.Text,
		d.ImagesFormat.Selected,
		d.AllChapters.Checked,
		d.LastChapter.Checked,
		d.ImagesOnly.Checked,
		d.OutputFolder.Text,
	)
}
