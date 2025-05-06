package photoorganizer

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type XmpMeta struct {
	XMLName xml.Name `xml:"x:xmpmeta"`
	XmlnsX  string   `xml:"xmlns:x,attr"`
	RDF     RDF      `xml:"rdf:RDF"`
}

type RDF struct {
	XMLName     xml.Name    `xml:"rdf:RDF"`
	XmlnsRdf    string      `xml:"xmlns:rdf,attr"`
	XmlnsDc     string      `xml:"xmlns:dc,attr"`
	Description Description `xml:"rdf:Description"`
}

type Description struct {
	Description DescriptionText `xml:"dc:description"`
	Subject     Subject         `xml:"dc:subject"`
}

type DescriptionText struct {
	Alt Alt `xml:"rdf:Alt"`
}

type Alt struct {
	LangString LangString `xml:"rdf:li"`
}

type LangString struct {
	Lang  string `xml:"xml:lang,attr"`
	Value string `xml:",chardata"`
}

type Subject struct {
	Bag Bag `xml:"rdf:Bag"`
}

type Bag struct {
	Items []string `xml:"rdf:li"`
}

func generateXMPFile(path, description string, tags []string) error {
	base := strings.TrimSuffix(path, filepath.Ext(path))
	output := base + ".xmp"

	xmp := XmpMeta{
		XmlnsX: "adobe:ns:meta/",
		RDF: RDF{
			XmlnsRdf: "http://www.w3.org/1999/02/22-rdf-syntax-ns#",
			XmlnsDc:  "http://purl.org/dc/elements/1.1/",
			Description: Description{
				Description: DescriptionText{
					Alt: Alt{
						LangString: LangString{
							Lang:  "x-default",
							Value: description,
						},
					},
				},
				Subject: Subject{
					Bag: Bag{Items: tags},
				},
			},
		},
	}

	file, err := os.Create(output)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")
	if err := encoder.Encode(xmp); err != nil {
		return err
	}
	fmt.Println("XMP written to:", output)
	return nil
}
