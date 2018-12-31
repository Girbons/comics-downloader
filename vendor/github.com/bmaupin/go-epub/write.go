package epub

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// UnableToCreateEpubError is thrown by Write if it cannot create the destination EPUB file
type UnableToCreateEpubError struct {
	Path string // The path that was given to Write to create the EPUB
	Err  error  // The underlying error that was thrown
}

func (e *UnableToCreateEpubError) Error() string {
	return fmt.Sprintf("Error creating EPUB at %q: %+v", e.Path, e.Err)
}

var extensionMediaTypes = map[string]string{
	".css":  mediaTypeCSS,
	".gif":  "image/gif",
	".jpeg": mediaTypeJpeg,
	".jpg":  mediaTypeJpeg,
	".otf":  "application/x-font-otf",
	".png":  "image/png",
	".svg":  "image/svg+xml",
	".ttf":  "application/x-font-ttf",
}

const (
	containerFilename     = "container.xml"
	containerFileTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<container version="1.0" xmlns="urn:oasis:names:tc:opendocument:xmlns:container">
  <rootfiles>
    <rootfile full-path="%s/%s" media-type="application/oebps-package+xml" />
  </rootfiles>
</container>
`
	// This seems to be the standard based on the latest EPUB spec:
	// http://www.idpf.org/epub/31/spec/epub-ocf.html
	contentFolderName    = "EPUB"
	coverImageProperties = "cover-image"
	// Permissions for any new directories we create
	dirPermissions = 0755
	// Permissions for any new files we create
	filePermissions   = 0644
	mediaTypeCSS      = "text/css"
	mediaTypeEpub     = "application/epub+zip"
	mediaTypeJpeg     = "image/jpeg"
	mediaTypeNcx      = "application/x-dtbncx+xml"
	mediaTypeXhtml    = "application/xhtml+xml"
	metaInfFolderName = "META-INF"
	mimetypeFilename  = "mimetype"
	pkgFilename       = "package.opf"
	tempDirPrefix     = "go-epub"
	xhtmlFolderName   = "xhtml"
)

// Write writes the EPUB file. The destination path must be the full path to
// the resulting file, including filename and extension.
func (e *Epub) Write(destFilePath string) error {
	tempDir, err := ioutil.TempDir("", tempDirPrefix)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			panic(fmt.Sprintf("Error removing temp directory: %s", err))
		}
	}()
	if err != nil {
		panic(fmt.Sprintf("Error creating temp directory: %s", err))
	}

	writeMimetype(tempDir)
	createEpubFolders(tempDir)

	// Must be called after:
	// createEpubFolders()
	writeContainerFile(tempDir)

	// Must be called after:
	// createEpubFolders()
	err = e.writeCSSFiles(tempDir)
	if err != nil {
		return err
	}

	// Must be called after:
	// createEpubFolders()
	err = e.writeFonts(tempDir)
	if err != nil {
		return err
	}

	// Must be called after:
	// createEpubFolders()
	err = e.writeImages(tempDir)
	if err != nil {
		return err
	}

	// Must be called after:
	// createEpubFolders()
	e.writeSections(tempDir)

	// Must be called after:
	// createEpubFolders()
	// writeSections()
	e.writeToc(tempDir)

	// Must be called after:
	// createEpubFolders()
	// writeCSSFiles()
	// writeImages()
	// writeSections()
	// writeToc()
	e.writePackageFile(tempDir)

	// Must be called last
	err = e.writeEpub(tempDir, destFilePath)
	if err != nil {
		return err
	}

	return nil
}

// Create the EPUB folder structure in a temp directory
func createEpubFolders(tempDir string) {
	if err := os.Mkdir(
		filepath.Join(
			tempDir,
			contentFolderName,
		),
		dirPermissions); err != nil {
		// No reason this should happen if tempDir creation was successful
		panic(fmt.Sprintf("Error creating EPUB subdirectory: %s", err))
	}

	if err := os.Mkdir(
		filepath.Join(
			tempDir,
			contentFolderName,
			xhtmlFolderName,
		),
		dirPermissions); err != nil {
		panic(fmt.Sprintf("Error creating xhtml subdirectory: %s", err))
	}

	if err := os.Mkdir(
		filepath.Join(
			tempDir,
			metaInfFolderName,
		),
		dirPermissions); err != nil {
		panic(fmt.Sprintf("Error creating META-INF subdirectory: %s", err))
	}
}

// Write the contatiner file (container.xml), which mostly just points to the
// package file (package.opf)
//
// Sample: https://github.com/bmaupin/epub-samples/blob/master/minimal-v3plus2/META-INF/container.xml
// Spec: http://www.idpf.org/epub/301/spec/epub-ocf.html#sec-container-metainf-container.xml
func writeContainerFile(tempDir string) {
	containerFilePath := filepath.Join(tempDir, metaInfFolderName, containerFilename)
	if err := ioutil.WriteFile(
		containerFilePath,
		[]byte(
			fmt.Sprintf(
				containerFileTemplate,
				contentFolderName,
				pkgFilename,
			),
		),
		filePermissions,
	); err != nil {
		panic(fmt.Sprintf("Error writing container file: %s", err))
	}
}

// Write the CSS files to the temporary directory and add them to the package
// file
func (e *Epub) writeCSSFiles(tempDir string) error {
	err := e.writeMedia(tempDir, e.css, CSSFolderName)
	if err != nil {
		return err
	}

	// Clean up the cover temp file if one was created
	os.Remove(e.cover.cssTempFile)

	return nil
}

// Write the EPUB file itself by zipping up everything from a temp directory
func (e *Epub) writeEpub(tempDir string, destFilePath string) error {
	f, err := os.Create(destFilePath)
	if err != nil {
		return &UnableToCreateEpubError{
			Path: destFilePath,
			Err:  err,
		}
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	z := zip.NewWriter(f)
	defer func() {
		if err := z.Close(); err != nil {
			panic(err)
		}
	}()

	skipMimetypeFile := false

	var addFileToZip = func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get the path of the file relative to the folder we're zipping
		relativePath, err := filepath.Rel(tempDir, path)
		relativePath = filepath.ToSlash(relativePath)
		if err != nil {
			// tempDir and path are both internal, so we shouldn't get here
			panic(fmt.Sprintf("Error closing EPUB file: %s", err))
		}

		// Only include regular files, not directories
		if !info.Mode().IsRegular() {
			return nil
		}

		var w io.Writer
		if path == filepath.Join(tempDir, mimetypeFilename) {
			// Skip the mimetype file if it's already been written
			if skipMimetypeFile == true {
				return nil
			}
			// The mimetype file must be uncompressed according to the EPUB spec
			w, err = z.CreateHeader(&zip.FileHeader{
				Name:   relativePath,
				Method: zip.Store,
			})
		} else {
			w, err = z.Create(relativePath)
		}
		if err != nil {
			panic(fmt.Sprintf("Error creating zip writer: %s", err))
		}

		r, err := os.Open(path)
		if err != nil {
			panic(fmt.Sprintf("Error opening file being added to EPUB: %s", err))
		}
		defer func() {
			if err := r.Close(); err != nil {
				panic(err)
			}
		}()

		_, err = io.Copy(w, r)
		if err != nil {
			panic(fmt.Sprintf("Error copying contents of file being added EPUB: %s", err))
		}

		return nil
	}

	// Add the mimetype file first
	mimetypeFilePath := filepath.Join(tempDir, mimetypeFilename)
	mimetypeInfo, err := os.Lstat(mimetypeFilePath)
	if err != nil {
		panic(fmt.Sprintf("Unable to get FileInfo for mimetype file: %s", err))
	}
	err = addFileToZip(mimetypeFilePath, mimetypeInfo, nil)
	if err != nil {
		panic(fmt.Sprintf("Unable to add mimetype file to EPUB: %s", err))
	}

	skipMimetypeFile = true

	err = filepath.Walk(tempDir, addFileToZip)
	if err != nil {
		panic(fmt.Sprintf("Unable to add file to EPUB: %s", err))
	}

	return nil
}

// Get fonts from their source and save them in the temporary directory
func (e *Epub) writeFonts(tempDir string) error {
	return e.writeMedia(tempDir, e.fonts, FontFolderName)
}

// Get images from their source and save them in the temporary directory
func (e *Epub) writeImages(tempDir string) error {
	return e.writeMedia(tempDir, e.images, ImageFolderName)
}

// Get images from their source and save them in the temporary directory
func (e *Epub) writeMedia(tempDir string, mediaMap map[string]string, mediaFolderName string) error {
	if len(mediaMap) > 0 {
		mediaFolderPath := filepath.Join(tempDir, contentFolderName, mediaFolderName)
		if err := os.Mkdir(mediaFolderPath, dirPermissions); err != nil {
			panic(fmt.Sprintf("Unable to create directory: %s", err))
		}

		for mediaFilename, mediaSource := range mediaMap {
			// Get the media file from the source
			u, err := url.Parse(mediaSource)
			if err != nil {
				return &FileRetrievalError{Source: mediaSource, Err: err}
			}

			var r io.ReadCloser
			var resp *http.Response
			// If it's a URL
			if u.Scheme == "http" || u.Scheme == "https" {
				resp, err = http.Get(mediaSource)
				if err != nil {
					return &FileRetrievalError{Source: mediaSource, Err: err}
				}
				r = resp.Body

				// Otherwise, assume it's a local file
			} else {
				r, err = os.Open(mediaSource)
			}
			if err != nil {
				return &FileRetrievalError{Source: mediaSource, Err: err}
			}

			mediaFilePath := filepath.Join(
				mediaFolderPath,
				mediaFilename,
			)

			// Add the file to the EPUB temp directory
			w, err := os.Create(mediaFilePath)
			if err != nil {
				panic(fmt.Sprintf("Unable to create file: %s", err))
			}

			_, err = io.Copy(w, r)
			// Close the reader and writer manually. If we use a defer instead,
			// they won't close until the function exits.
			func() {
				if err := r.Close(); err != nil {
					panic(err)
				}
			}()
			func() {
				if err := w.Close(); err != nil {
					panic(err)
				}
			}()
			if err != nil {
				// There shouldn't be any problem with the writer, but the reader
				// might have an issue
				return &FileRetrievalError{Source: mediaSource, Err: err}
			}

			mediaType := extensionMediaTypes[strings.ToLower(filepath.Ext(mediaFilename))]
			if mediaType == "" {
				panic(fmt.Sprintf(
					"Unmatched file extension, media type not set for file: %s",
					mediaFilename))
			}

			// The cover image has a special value for the properties attribute
			mediaProperties := ""
			if mediaFilename == e.cover.imageFilename {
				mediaProperties = coverImageProperties
			}

			// Add the file to the OPF manifest
			e.pkg.addToManifest(mediaFilename, filepath.Join(mediaFolderName, mediaFilename), mediaType, mediaProperties)
		}
	}

	return nil
}

// Write the mimetype file
//
// Sample: https://github.com/bmaupin/epub-samples/blob/master/minimal-v3plus2/mimetype
// Spec: http://www.idpf.org/epub/301/spec/epub-ocf.html#sec-zip-container-mime
func writeMimetype(tempDir string) {
	mimetypeFilePath := filepath.Join(tempDir, mimetypeFilename)

	if err := ioutil.WriteFile(mimetypeFilePath, []byte(mediaTypeEpub), filePermissions); err != nil {
		panic(fmt.Sprintf("Error writing mimetype file: %s", err))
	}
}

func (e *Epub) writePackageFile(tempDir string) {
	e.pkg.write(tempDir)
}

// Write the section files to the temporary directory and add the sections to
// the TOC and package files
func (e *Epub) writeSections(tempDir string) {
	if len(e.sections) > 0 {
		// If a cover was set, add it to the package spine first so it shows up
		// first in the reading order
		if e.cover.xhtmlFilename != "" {
			e.pkg.addToSpine(e.cover.xhtmlFilename)
		}

		for i, section := range e.sections {
			// Set the title of the cover page XHTML to the title of the EPUB
			if section.filename == e.cover.xhtmlFilename {
				section.xhtml.setTitle(e.Title())
			}

			sectionFilePath := filepath.Join(tempDir, contentFolderName, xhtmlFolderName, section.filename)
			section.xhtml.write(sectionFilePath)

			relativePath := filepath.Join(xhtmlFolderName, section.filename)
			// Don't add pages without titles or the cover to the TOC
			if section.xhtml.Title() != "" && section.filename != e.cover.xhtmlFilename {
				e.toc.addSection(i, section.xhtml.Title(), relativePath)
			}
			// The cover page should have already been added to the spine first
			if section.filename != e.cover.xhtmlFilename {
				e.pkg.addToSpine(section.filename)
			}
			e.pkg.addToManifest(section.filename, relativePath, mediaTypeXhtml, "")
		}
	}
}

// Write the TOC file to the temporary directory and add the TOC entries to the
// package file
func (e *Epub) writeToc(tempDir string) {
	e.pkg.addToManifest(tocNavItemID, tocNavFilename, mediaTypeXhtml, tocNavItemProperties)
	e.pkg.addToManifest(tocNcxItemID, tocNcxFilename, mediaTypeNcx, "")

	e.toc.write(tempDir)
}
