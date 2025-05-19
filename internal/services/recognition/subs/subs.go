package subs

import (
	"strings"
	"unicode/utf8"
)

type section struct {
	text    string
	isFinal bool
}

type lineBreak struct {
	sectionID int
	offset    int
}

type Writer struct {
	totalLines int
	perLine    int
	sections   []section
	newLines   []lineBreak
}

func NewWriter(totalLines, perLine int) *Writer {
	if totalLines == 0 {
		totalLines = 2
	}
	if perLine == 0 {
		perLine = 80
	}
	return &Writer{
		totalLines: totalLines,
		perLine:    perLine,
	}
}

func (w *Writer) AddSection(text string, isFinal bool) {
	if isFinal {
		text += "."
	}

	separatorLength := 1

	if len(w.sections) == 0 || w.sections[len(w.sections)-1].isFinal {
		w.sections = append(w.sections, section{text: text, isFinal: isFinal})
	} else {
		w.sections[len(w.sections)-1] = section{text: text, isFinal: isFinal}
		// Remove line breaks related to last section
		var newLines []lineBreak
		for _, l := range w.newLines {
			if l.sectionID != len(w.sections)-1 {
				newLines = append(newLines, l)
			}
		}
		w.newLines = newLines
	}

	lineLength := 0
	startSectionID := 0
	if len(w.newLines) > 0 {
		sectionID := w.newLines[len(w.newLines)-1].sectionID
		offset := w.newLines[len(w.newLines)-1].offset
		lineLength = utf8.RuneCountInString(w.sections[sectionID].text) - offset
		startSectionID = sectionID + 1
	}
	for i := startSectionID; i < len(w.sections)-1; i++ {
		lineLength += separatorLength + utf8.RuneCountInString(w.sections[i].text)
	}

	words := strings.Fields(text)
	newSectionID := len(w.sections) - 1
	offset := 0
	for i, word := range words {
		if utf8.RuneCountInString(word) > w.perLine {
			wordRunes := []rune(word)
			word = string(wordRunes[:w.perLine])
		}
		partLength := utf8.RuneCountInString(word)
		if i > 0 {
			partLength += separatorLength
		}
		lineLength += partLength
		if lineLength > w.perLine {
			w.newLines = append(w.newLines, lineBreak{sectionID: newSectionID, offset: offset})
			lineLength = utf8.RuneCountInString(word)
		}
		offset += partLength
	}

	// Normalize if lines exceed limit
	lastSectionID := len(w.sections) - 1
	newLinesBeforeNonFinal := 0
	if !w.sections[lastSectionID].isFinal {
		for _, line := range w.newLines {
			if line.sectionID >= lastSectionID {
				break
			}
			newLinesBeforeNonFinal++
		}
	}

	if len(w.newLines)-newLinesBeforeNonFinal > w.totalLines {
		extracts := len(w.newLines) - w.totalLines
		w.newLines = w.newLines[extracts:]

		firstSectionID := w.newLines[0].sectionID
		w.sections = w.sections[firstSectionID:]
		for i := range w.newLines {
			w.newLines[i].sectionID -= firstSectionID
		}
	}
}

func (w *Writer) Format() string {
	lines := [][]string{{}}
	lineIndex := 0

	for i, section := range w.sections {
		text := section.text
		var offsets []int

		for _, l := range w.newLines {
			if l.sectionID == i {
				offsets = append(offsets, l.offset)
				lines = append(lines, []string{})
			}
		}

		runes := []rune(text)

		if len(offsets) == 0 {
			lines[lineIndex] = append(lines[lineIndex], text)
			continue
		}

		if offsets[0] == 0 {
			lineIndex++
		} else {
			offsets = append([]int{0}, offsets...)
		}

		for j := 0; j < len(offsets)-1; j++ {
			segment := strings.TrimLeft(string(runes[offsets[j]:offsets[j+1]]), " ")
			lines[lineIndex] = append(lines[lineIndex], segment)
			lineIndex++
		}
		lines[lineIndex] = append(lines[lineIndex], strings.TrimLeft(string(runes[offsets[len(offsets)-1]:]), " "))
	}

	start := 0
	if len(lines) > w.totalLines {
		start = len(lines) - w.totalLines
	}
	var formatted []string
	for _, line := range lines[start:] {
		formatted = append(formatted, strings.Join(line, " "))
	}
	return strings.Join(formatted, "\n")
}
