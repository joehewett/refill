package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/unidoc/unipdf/v3/extractor"
	pdfModel "github.com/unidoc/unipdf/v3/model"
)

type FileType int

const (
	txt FileType = iota
	pdf
)

type File interface {
	Name() string
	Path() string
	Type() FileType
	Load() (string, error)
}

func NewFile(path string, extension string) (File, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("failed to find file %s: %s", path, err)
	}

	if extension == ".pdf" {
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

func (p *PDF) Name() string {
	return filepath.Base(p.Path())
}

func (p *PDF) Type() FileType {
	return pdf
}

func (p *PDF) Load() (string, error) {
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
	}

	// fmt.Printf("%s\n----------------------------------------\n", text)

	return text, nil
}

type TXT struct {
	path string
}

func (t *TXT) Path() string {
	return t.path
}

func (t *TXT) Name() string {
	return filepath.Base(t.Path())
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
