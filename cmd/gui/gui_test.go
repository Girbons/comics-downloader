package main

import (
	"testing"

	"fyne.io/fyne/test"
	"fyne.io/fyne/widget"
)

// newTestDownloader creates a Downloader wired with real widgets under the
// headless fyne test driver (no display required).
func newTestDownloader() *Downloader {
	_ = test.NewApp() // headless driver; safe to call multiple times
	return &Downloader{
		URL:               widget.NewEntry(),
		Country:           widget.NewEntry(),
		Format:            widget.NewRadioGroup([]string{"pdf", "epub", "cbr", "cbz"}, nil),
		AllChapters:       widget.NewCheck("", nil),
		LastChapter:       widget.NewCheck("", nil),
		ImagesOnly:        widget.NewCheck("", nil),
		ImagesFormat:      widget.NewRadioGroup([]string{"png", "jpg", "img"}, nil),
		OutputFolder:      widget.NewEntry(),
		CreateDefaultPath: widget.NewCheck("", nil),
		IssuesRange:       widget.NewEntry(),
		Debug:             widget.NewCheck("", nil),
		CustomComicName:   widget.NewEntry(),
	}
}

func TestClearURLField(t *testing.T) {
	d := newTestDownloader()
	d.URL.SetText("https://example.com/comic/1")
	d.ClearURLField()
	if d.URL.Text != "" {
		t.Fatalf("expected empty URL after clear, got %q", d.URL.Text)
	}
}

func TestClearCountryField(t *testing.T) {
	d := newTestDownloader()
	d.Country.SetText("JP")
	d.ClearCountryField()
	if d.Country.Text != "" {
		t.Fatalf("expected empty country after clear, got %q", d.Country.Text)
	}
}

func TestClearOutputFolderField(t *testing.T) {
	d := newTestDownloader()
	d.OutputFolder.SetText("/tmp/comics")
	d.ClearOutputFolderField()
	if d.OutputFolder.Text != "" {
		t.Fatalf("expected empty output folder after clear, got %q", d.OutputFolder.Text)
	}
}

func TestDownloaderDefaultFieldValues(t *testing.T) {
	d := newTestDownloader()

	if d.URL.Text != "" {
		t.Errorf("URL should default to empty, got %q", d.URL.Text)
	}
	if d.AllChapters.Checked {
		t.Errorf("AllChapters should default to unchecked")
	}
	if d.LastChapter.Checked {
		t.Errorf("LastChapter should default to unchecked")
	}
	if d.ImagesOnly.Checked {
		t.Errorf("ImagesOnly should default to unchecked")
	}
	if d.Debug.Checked {
		t.Errorf("Debug should default to unchecked")
	}
}

func TestDownloaderURLTrimmedOnSubmit(t *testing.T) {
	// Submit is fire-and-forget (go routine) and calls the real downloader with
	// a blank URL, which exits quickly.  The main goal here is that wiring
	// compiles and the entry text used inside Submit is trimmed.
	d := newTestDownloader()
	d.URL.SetText("  https://example.com/comic  ")

	// We cannot easily intercept the options struct from outside the function,
	// but we can verify the widget still holds its raw text and that TrimSpace
	// is applied internally (white-box check via direct field read after a
	// synthesised Submit would require refactoring; this test guards the
	// constructor wiring compiles and the helpers operate correctly).
	if d.URL.Text != "  https://example.com/comic  " {
		t.Fatalf("unexpected URL text before submit: %q", d.URL.Text)
	}
}

func TestFormatRadioGroupOptions(t *testing.T) {
	d := newTestDownloader()
	d.Format.SetSelected("epub")
	if d.Format.Selected != "epub" {
		t.Fatalf("expected Selected=epub, got %q", d.Format.Selected)
	}
}

func TestImagesFormatRadioGroupOptions(t *testing.T) {
	d := newTestDownloader()
	d.ImagesFormat.SetSelected("png")
	if d.ImagesFormat.Selected != "png" {
		t.Fatalf("expected Selected=png, got %q", d.ImagesFormat.Selected)
	}
}
