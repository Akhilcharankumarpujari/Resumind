package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	ledpdf "github.com/ledongthuc/pdf"
)


func ExtractTextFromPDF(r io.ReaderAt, size int64) (string, error) {
	reader, err := ledpdf.NewReader(r, size)
	if err != nil {
		return "", fmt.Errorf("failed to open PDF reader: %w", err)
	}

	var fullText bytes.Buffer
	numPages := reader.NumPage()

	if numPages == 0 {
		return "", fmt.Errorf("PDF has no pages")
	}

	for pageNum := 1; pageNum <= numPages; pageNum++ {
		page := reader.Page(pageNum)
		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)
		if err != nil {

			continue
		}

		fullText.WriteString(text)
		fullText.WriteString("\n\n--- Page Break ---\n\n")
	}

	result := strings.TrimSpace(fullText.String())
	if result == "" {
		return "", fmt.Errorf("no readable text found in PDF (it may be a scanned image-only PDF)")
	}

	return result, nil
}
