package pdftotext

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

type Getter struct {
}

func NewGetter() *Getter {
	return &Getter{}
}

func (g *Getter) GetTextFromPDF(fileName string, customOffset int) (string, error) {
	args := []string{
		"-layout",  // Maintain (as best as possible) the original physical layout of the text.
		"-nopgbrk", // Don't insert page breaks (form feed characters) between pages.
		fileName,   // The input file.
		"-",        // Send the output to stdout.
	}

	cmd := exec.Command("pdftotext", args...)

	var buf bytes.Buffer
	cmd.Stdout = &buf

	if err := cmd.Run(); err != nil {
		return "", err
	}

	text := buf.String()

	noNewLines := strings.Split(text, "\n")
	if len(noNewLines) <= customOffset {
		return "", errors.New("no new lines")
	}

	return strings.Join(noNewLines[customOffset:], "\n"), nil
}
