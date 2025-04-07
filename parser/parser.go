package parser

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Shape struct {
	ID   string
	Text string
	PinX float64
	PinY float64
}

type Connect struct {
	FromID string
	ToID   string
}

type Page struct {
	XMLName  xml.Name     `xml:"PageContents"`
	Shapes   []RawShape   `xml:"Shapes>Shape"`
	Connects []RawConnect `xml:"Connects>Connect"`
}

type RawShape struct {
	ID   string `xml:"ID,attr"`
	Text string `xml:"Text"`
	PinX string `xml:"XForm>PinX"`
	PinY string `xml:"XForm>PinY"`
}

type RawConnect struct {
	FromSheet string `xml:"FromSheet,attr"`
	ToSheet   string `xml:"ToSheet,attr"`
}

func ExtractVSDX(filename string) ([]Shape, []Connect, error) {
	r, err := zip.OpenReader(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("opening VSDX: %w", err)
	}
	defer r.Close()

	var pageXML io.ReadCloser
	for _, f := range r.File {
		if strings.HasSuffix(f.Name, "pages/page1.xml") {
			pageXML, err = f.Open()
			if err != nil {
				return nil, nil, fmt.Errorf("opening page1.xml: %w", err)
			}
			defer pageXML.Close()
			break
		}
	}

	if pageXML == nil {
		return nil, nil, fmt.Errorf("page1.xml not found in VSDX")
	}

	var page Page
	decoder := xml.NewDecoder(pageXML)
	if err := decoder.Decode(&page); err != nil {
		return nil, nil, fmt.Errorf("decoding page1.xml: %w", err)
	}

	var shapes []Shape
	for _, rs := range page.Shapes {
		pinX, _ := strconv.ParseFloat(strings.TrimSpace(rs.PinX), 64)
		pinY, _ := strconv.ParseFloat(strings.TrimSpace(rs.PinY), 64)
		shape := Shape{
			ID:   rs.ID,
			Text: strings.TrimSpace(rs.Text),
			PinX: pinX,
			PinY: pinY,
		}
		shapes = append(shapes, shape)
	}

	var connectors []Connect
	for _, rc := range page.Connects {
		connectors = append(connectors, Connect{
			FromID: rc.FromSheet,
			ToID:   rc.ToSheet,
		})
	}

	return shapes, connectors, nil
}
