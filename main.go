package main

import (
	"encoding/xml"
	"fmt"
	"os"
)

func main() {
	data := map[string]interface{}{
		"id":   1,
		"name": "john",
		"sampleList": []interface{}{
			"text", "2",
		},
		"sampleMap": map[string]interface{}{
			"key": "value",
		},
	}

	// create new writer
	dst, _ := os.Create("/tmp/dat2")

	// write default XML header
	dst.WriteString(xml.Header)

	// create new encoder with writer
	xmlEncoder := xml.NewEncoder(dst)
	xmlEncoder.Indent("", "\t")

	// create root element
	rootElement := xml.StartElement{
		Name: xml.Name{
			Local: "profiles",
		},
		Attr: []xml.Attr{
			xml.Attr{
				Name:  xml.Name{Local: "language"},
				Value: "En-US",
			},
		},
	}
	if err := xmlEncoder.EncodeToken(rootElement); err != nil {
		fmt.Println(err)
		return
	}

	// write buffer to file
	xmlEncoder.Flush()

	for k, value := range data {
		switch v := value.(type) {
		case string, int:
			// add attributes to file
			if err := xmlEncoder.EncodeElement(v, xml.StartElement{Name: xml.Name{Local: k}}); err != nil {
				fmt.Println(err)
				return
			}

		case map[string]interface{}:
			// create root element
			nodeElement := xml.StartElement{
				Name: xml.Name{
					Local: "node",
				},
			}

			if err := xmlEncoder.EncodeToken(nodeElement); err != nil {
				fmt.Println(err)
				return
			}

			for z, xalue := range v {
				if err := xmlEncoder.EncodeElement(xalue, xml.StartElement{Name: xml.Name{Local: z}}); err != nil {
					fmt.Println(err)
					return
				}
			}

			nodeCloser := xml.EndElement{Name: nodeElement.Name}
			if err := xmlEncoder.EncodeToken(nodeCloser); err != nil {
				fmt.Println(err)
				return
			}
		}
	}

	// close root element
	rootCloser := xml.EndElement{Name: rootElement.Name}
	if err := xmlEncoder.EncodeToken(rootCloser); err != nil {
		fmt.Println(err)
		return
	}

	xmlEncoder.Flush()
}
