package main

import (
	"fmt"
	"io"
	"os"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/extractor"
	pdfModel "github.com/unidoc/unipdf/v3/model"
)

type FileType int

const (
	txt FileType = iota
	pdf
)

type File interface {
	Path() string
	Type() FileType
	Load() (string, error)
}

func NewFile(path string, extension string) (File, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("failed to find file %s: %s", path, err)
	}

	if extension == "pdf" {
		return &PDF{path: path}, nil
	} else {
		return &TXT{path: path}, nil
	}
}

type PDF struct {
	path string
}

func (p *PDF) Path() string {
	return p.path
}

func (p *PDF) Type() FileType {
	return pdf
}

func (p *PDF) Load() (string, error) {
	initPDFParser()

	f, err := os.Open(p.Path())
	if err != nil {
		return "", err
	}

	defer f.Close()

	pdfReader, err := pdfModel.NewPdfReader(f)
	if err != nil {
		return "", err
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return "", err
	}

	fmt.Printf("--------------------\n")
	fmt.Printf("PDF to text extraction:\n")
	fmt.Printf("--------------------\n")
	var text string
	for i := 0; i < numPages; i++ {
		pageNum := i + 1

		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return "", err
		}

		ex, err := extractor.New(page)
		if err != nil {
			return "", err
		}

		pageText, err := ex.ExtractText()
		if err != nil {
			return "", err
		}

		text += pageText

		fmt.Println("------------------------------")
		fmt.Printf("Page %d:\n", pageNum)
		fmt.Printf("\"%s\"\n", text)
		fmt.Println("------------------------------")
	}

	return text, nil
}

type TXT struct {
	path string
}

func (t *TXT) Path() string {
	return t.path
}

func (t *TXT) Type() FileType {
	return txt
}

func (t *TXT) Load() (string, error) {
	dataFile, err := os.Open(t.Path())
	if err != nil {
		return "", fmt.Errorf("failed to open file %s: %s", t.Path(), err)
	}

	defer dataFile.Close()

	bytes, err := io.ReadAll(dataFile)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %s", t.Path(), err)
	}

	data := string(bytes)
	if *verbose {
		fmt.Println("Reading JSON from file")
	}

	return data, nil
}

func initPDFParser() {
	// Make sure to load your metered License API key prior to using the library.
	// If you need a key, you can sign up and create a free one at https://cloud.unidoc.io
	err := license.SetMeteredKey(os.Getenv(`UNIDOC_LICENSE_API_KEY`))
	if err != nil {
		panic(err)
	}
}
