package epub

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"
)

const (
	pkgAuthorID       = "role"
	pkgAuthorData     = "aut"
	pkgAuthorProperty = "role"
	pkgAuthorRefines  = "#creator"
	pkgAuthorScheme   = "marc:relators"
	pkgCreatorID      = "creator"
	pkgFileTemplate   = `<?xml version="1.0" encoding="UTF-8"?>
<package version="3.0" unique-identifier="pub-id" xmlns="http://www.idpf.org/2007/opf">
  <metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:identifier id="pub-id"></dc:identifier>
    <dc:title></dc:title>
    <dc:language></dc:language>
  </metadata>
  <manifest>
  </manifest>
  <spine toc="ncx">
  </spine>
</package>
`
	pkgModifiedProperty = "dcterms:modified"
	pkgUniqueIdentifier = "pub-id"

	xmlnsDc = "http://purl.org/dc/elements/1.1/"
)

// pkg implements the package document file (package.opf), which contains
// metadata about the EPUB (title, author, etc) as well as a list of files the
// EPUB contains.
//
// Sample: https://github.com/bmaupin/epub-samples/blob/master/minimal-v3plus2/EPUB/package.opf
// Spec: http://www.idpf.org/epub/301/spec/epub-publications.html
type pkg struct {
	xml          *pkgRoot
	authorMeta   *pkgMeta
	modifiedMeta *pkgMeta
}

// This holds the actual XML for the package file
type pkgRoot struct {
	XMLName          xml.Name    `xml:"http://www.idpf.org/2007/opf package"`
	UniqueIdentifier string      `xml:"unique-identifier,attr"`
	Version          string      `xml:"version,attr"`
	Metadata         pkgMetadata `xml:"metadata"`
	ManifestItems    []pkgItem   `xml:"manifest>item"`
	Spine            pkgSpine    `xml:"spine"`
}

// <dc:creator>, e.g. the author
type pkgCreator struct {
	XMLName xml.Name `xml:"dc:creator"`
	ID      string   `xml:"id,attr"`
	Data    string   `xml:",chardata"`
}

// <dc:identifier>, where the unique identifier is stored
// Ex: <dc:identifier id="pub-id">urn:uuid:fe93046f-af57-475a-a0cb-a0d4bc99ba6d</dc:identifier>
type pkgIdentifier struct {
	ID   string `xml:"id,attr"`
	Data string `xml:",chardata"`
}

// <item> elements, one per each file stored in the EPUB
// Ex: <item id="nav" href="nav.xhtml" media-type="application/xhtml+xml" properties="nav" />
//     <item id="ncx" href="toc.ncx" media-type="application/x-dtbncx+xml" />
//     <item id="section0001.xhtml" href="xhtml/section0001.xhtml" media-type="application/xhtml+xml" />
type pkgItem struct {
	ID         string `xml:"id,attr"`
	Href       string `xml:"href,attr"`
	MediaType  string `xml:"media-type,attr"`
	Properties string `xml:"properties,attr,omitempty"`
}

// <itemref> elements, which define the reading order
// Ex: <itemref idref="section0001.xhtml" />
type pkgItemref struct {
	Idref string `xml:"idref,attr"`
}

// The <meta> element, which contains modified date, role of the creator (e.g.
// author), etc
// Ex: <meta refines="#creator" property="role" scheme="marc:relators" id="role">aut</meta>
//     <meta property="dcterms:modified">2011-01-01T12:00:00Z</meta>
type pkgMeta struct {
	Refines  string `xml:"refines,attr,omitempty"`
	Property string `xml:"property,attr"`
	Scheme   string `xml:"scheme,attr,omitempty"`
	ID       string `xml:"id,attr,omitempty"`
	Data     string `xml:",chardata"`
}

// The <metadata> element
type pkgMetadata struct {
	XmlnsDc    string        `xml:"xmlns:dc,attr"`
	Identifier pkgIdentifier `xml:"dc:identifier"`
	// Ex: <dc:title>Your title here</dc:title>
	Title string `xml:"dc:title"`
	// Ex: <dc:language>en</dc:language>
	Language string `xml:"dc:language"`
	Creator  *pkgCreator
	Meta     []pkgMeta `xml:"meta"`
}

// The <spine> element
type pkgSpine struct {
	Items []pkgItemref `xml:"itemref"`
	Toc   string       `xml:"toc,attr"`
	Ppd   string       `xml:"page-progression-direction,attr,omitempty"`
}

// Constructor for pkg
func newPackage() *pkg {
	p := &pkg{
		xml: &pkgRoot{
			Metadata: pkgMetadata{
				XmlnsDc: xmlnsDc,
				Identifier: pkgIdentifier{
					ID: pkgUniqueIdentifier,
				},
			},
		},
	}

	err := xml.Unmarshal([]byte(pkgFileTemplate), &p.xml)
	if err != nil {
		panic(fmt.Sprintf(
			"Error unmarshalling package file XML: %s\n"+
				"\tp.xml=%#v\n"+
				"\tpkgFileTemplate=%s",
			err,
			*p.xml,
			pkgFileTemplate))
	}

	return p
}

func (p *pkg) addToManifest(id string, href string, mediaType string, properties string) {
	href = filepath.ToSlash(href)
	i := &pkgItem{
		ID:         id,
		Href:       href,
		MediaType:  mediaType,
		Properties: properties,
	}
	p.xml.ManifestItems = append(p.xml.ManifestItems, *i)
}

func (p *pkg) addToSpine(id string) {
	i := &pkgItemref{
		Idref: id,
	}

	p.xml.Spine.Items = append(p.xml.Spine.Items, *i)
}

func (p *pkg) setAuthor(author string) {
	p.xml.Metadata.Creator = &pkgCreator{
		Data: author,
		ID:   pkgCreatorID,
	}
	p.authorMeta = &pkgMeta{
		Data:     pkgAuthorData,
		ID:       pkgAuthorID,
		Property: pkgAuthorProperty,
		Refines:  pkgAuthorRefines,
		Scheme:   pkgAuthorScheme,
	}

	p.xml.Metadata.Meta = updateMeta(p.xml.Metadata.Meta, p.authorMeta)
}

func (p *pkg) setIdentifier(identifier string) {
	p.xml.Metadata.Identifier.Data = identifier
}

func (p *pkg) setLang(lang string) {
	p.xml.Metadata.Language = lang
}

func (p *pkg) setPpd(direction string) {
	p.xml.Spine.Ppd = direction
}

func (p *pkg) setModified(timestamp string) {
	p.modifiedMeta = &pkgMeta{
		Data:     timestamp,
		Property: pkgModifiedProperty,
	}

	p.xml.Metadata.Meta = updateMeta(p.xml.Metadata.Meta, p.modifiedMeta)
}

func (p *pkg) setTitle(title string) {
	p.xml.Metadata.Title = title
}

// Update the <meta> element
func updateMeta(a []pkgMeta, m *pkgMeta) []pkgMeta {
	indexToReplace := -1

	if len(a) > 0 {
		// If we've already added the modified meta element to the meta array
		for i, meta := range a {
			if meta == *m {
				indexToReplace = i
				break
			}
		}
	}

	// If the array is empty or the meta element isn't in it
	if indexToReplace == -1 {
		// Add the meta element to the array of meta elements
		a = append(a, *m)

		// If the meta element is found
	} else {
		// Replace it
		a[indexToReplace] = *m
	}

	return a
}

// Write the package file to the temporary directory
func (p *pkg) write(tempDir string) {
	now := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	p.setModified(now)

	pkgFilePath := filepath.Join(tempDir, contentFolderName, pkgFilename)

	output, err := xml.MarshalIndent(p.xml, "", "  ")
	if err != nil {
		panic(fmt.Sprintf(
			"Error marshalling XML for package file: %s\n"+
				"\tXML=%#v",
			err,
			p.xml))
	}
	// Add the xml header to the output
	pkgFileContent := append([]byte(xml.Header), output...)
	// It's generally nice to have files end with a newline
	pkgFileContent = append(pkgFileContent, "\n"...)

	if err := ioutil.WriteFile(pkgFilePath, []byte(pkgFileContent), filePermissions); err != nil {
		panic(fmt.Sprintf("Error writing package file: %s", err))
	}
}
