// pkg/utils/xml_converter.go
package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/clbanning/mxj/v2"
)

type XMLConverter struct {
	RemoveNamespaces bool
	RemoveLineBreaks bool
	CleanOutput      bool
}

func NewXMLConverter(removeNamespaces, cleanOutput bool) *XMLConverter {
	return &XMLConverter{
		RemoveNamespaces: removeNamespaces,
		RemoveLineBreaks: true,
		CleanOutput:      cleanOutput,
	}
}

// From *bytes.Buffer
func (c *XMLConverter) ToJSONFromBuffer(buffer *bytes.Buffer) ([]byte, error) {
	return c.processXML(buffer.Bytes())
}

// From bytes.Buffer (value)
func (c *XMLConverter) ToJSONFromBufferValue(buffer bytes.Buffer) ([]byte, error) {
	return c.processXML(buffer.Bytes())
}

// From io.Reader
func (c *XMLConverter) ToJSONFromReader(reader io.Reader) ([]byte, error) {
	xmlData, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read XML data: %w", err)
	}
	return c.processXML(xmlData)
}

// From string
func (c *XMLConverter) ToJSONFromString(xmlString string) ([]byte, error) {
	return c.processXML([]byte(xmlString))
}

// Core processing logic
func (c *XMLConverter) processXML(xmlData []byte) ([]byte, error) {
	if len(xmlData) == 0 {
		return nil, fmt.Errorf("empty XML data")
	}

	if c.CleanOutput {
		xmlData = c.cleanXML(xmlData)
	}

	mv, err := mxj.NewMapXml(xmlData)
	if err != nil {
		return nil, fmt.Errorf("XML parsing failed: %w", err)
	}

	return json.MarshalIndent(mv, "", "  ")
}

func (c *XMLConverter) ToMap(xmlData []byte) (map[string]any, error) {
	if len(xmlData) == 0 {
		return nil, fmt.Errorf("empty XML data")
	}

	if c.CleanOutput {
		xmlData = c.cleanXML(xmlData)
	}

	mv, err := mxj.NewMapXml(xmlData)
	if err != nil {
		return nil, fmt.Errorf("XML parsing failed: %w", err)
	}

	return map[string]any(mv), nil
}

func (c *XMLConverter) cleanXML(xmlData []byte) []byte {
	cleaned := string(xmlData)
	if c.RemoveNamespaces {
		// Simple namespace removal
		cleaned = strings.ReplaceAll(cleaned, ` xmlns="[^"]*"`, "")
		cleaned = strings.ReplaceAll(cleaned, ` xmlns:[^=]*="[^"]*"`, "")
	}

	if c.RemoveLineBreaks {
		cleaned = c.removeLineBreaks(cleaned)
	}

	return []byte(strings.TrimSpace(cleaned))
}

// Zeilenumbrüche und überflüssige Whitespaces entfernen
func (c *XMLConverter) removeLineBreaks(xml string) string {
	// Verschiedene Zeilenumbruch-Typen
	xml = strings.ReplaceAll(xml, "\r\n", " ") // Windows
	xml = strings.ReplaceAll(xml, "\r", " ")   // Mac (alt)
	xml = strings.ReplaceAll(xml, "\n", " ")   // Unix/Linux

	// Tabs entfernen
	xml = strings.ReplaceAll(xml, "\t", " ")

	// Mehrfache Leerzeichen zu einem zusammenfassen
	re := regexp.MustCompile(`\s+`)
	xml = re.ReplaceAllString(xml, " ")

	// Leerzeichen zwischen Tags bereinigen
	re2 := regexp.MustCompile(`>\s+<`)
	xml = re2.ReplaceAllString(xml, "><")

	return xml
}
