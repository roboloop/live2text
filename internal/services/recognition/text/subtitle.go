package text

import (
	"strings"
	"time"
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

type Formatter struct {
	totalLines int
	perLine    int
	sections   []section
	newLines   []lineBreak
}

func NewSubtitleFormatter(totalLines, perLine int) *Formatter {
	if totalLines == 0 {
		totalLines = 2
	}
	if perLine == 0 {
		perLine = 80
	}
	return &Formatter{
		totalLines: totalLines,
		perLine:    perLine,
	}
}

//nolint:gocognit // TODO: simplify
func (f *Formatter) Append(text string, isFinal bool) {
	if isFinal {
		text += "."
	}

	separatorLength := 1

	if len(f.sections) == 0 || f.sections[len(f.sections)-1].isFinal {
		f.sections = append(f.sections, section{text: text, isFinal: isFinal})
	} else {
		f.sections[len(f.sections)-1] = section{text: text, isFinal: isFinal}
		// Remove line breaks related to last section
		var newLines []lineBreak
		for _, l := range f.newLines {
			if l.sectionID != len(f.sections)-1 {
				newLines = append(newLines, l)
			}
		}
		f.newLines = newLines
	}

	lineLength := 0
	startSectionID := 0
	if len(f.newLines) > 0 {
		sectionID := f.newLines[len(f.newLines)-1].sectionID
		offset := f.newLines[len(f.newLines)-1].offset
		lineLength = utf8.RuneCountInString(f.sections[sectionID].text) - offset
		startSectionID = sectionID + 1
	}
	for i := startSectionID; i < len(f.sections)-1; i++ {
		lineLength += separatorLength + utf8.RuneCountInString(f.sections[i].text)
	}

	words := strings.Fields(text)
	newSectionID := len(f.sections) - 1
	offset := 0
	for i, word := range words {
		if utf8.RuneCountInString(word) > f.perLine {
			wordRunes := []rune(word)
			word = string(wordRunes[:f.perLine])
		}
		partLength := utf8.RuneCountInString(word)
		if i > 0 {
			partLength += separatorLength
		}
		lineLength += partLength
		if lineLength > f.perLine {
			f.newLines = append(f.newLines, lineBreak{sectionID: newSectionID, offset: offset})
			lineLength = utf8.RuneCountInString(word)
		}
		offset += partLength
	}

	// Normalize if lines exceed limit
	lastSectionID := len(f.sections) - 1
	newLinesBeforeNonFinal := 0
	if !f.sections[lastSectionID].isFinal {
		for _, line := range f.newLines {
			if line.sectionID >= lastSectionID {
				break
			}
			newLinesBeforeNonFinal++
		}
	}

	if len(f.newLines)-newLinesBeforeNonFinal > f.totalLines {
		extracts := len(f.newLines) - f.totalLines
		f.newLines = f.newLines[extracts:]

		firstSectionID := f.newLines[0].sectionID
		f.sections = f.sections[firstSectionID:]
		for i := range f.newLines {
			f.newLines[i].sectionID -= firstSectionID
		}
	}
}

func (f *Formatter) Format() string {
	lines := [][]string{{}}
	lineIndex := 0

	for i, section := range f.sections {
		text := section.text
		var offsets []int

		for _, l := range f.newLines {
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

		for j := range len(offsets) - 1 {
			segment := strings.TrimLeft(string(runes[offsets[j]:offsets[j+1]]), " ")
			lines[lineIndex] = append(lines[lineIndex], segment)
			lineIndex++
		}
		lines[lineIndex] = append(lines[lineIndex], strings.TrimLeft(string(runes[offsets[len(offsets)-1]:]), " "))
	}

	start := 0
	if len(lines) > f.totalLines {
		start = len(lines) - f.totalLines
	}
	var formatted []string
	for _, line := range lines[start:] {
		formatted = append(formatted, strings.Join(line, " "))
	}
	return strings.Join(formatted, "\n")
}

type subtitleWriter struct {
	formatter *Formatter
}

func NewSubtitleWriter(formatter *Formatter) Writer {
	return &subtitleWriter{formatter}
}

func (sw *subtitleWriter) PrintFinal(_ time.Duration, text string) error {
	sw.formatter.Append(text, true)

	return nil
}

func (sw *subtitleWriter) PrintCandidate(_ time.Duration, text string) error {
	sw.formatter.Append(text, false)

	return nil
}

func (sw *subtitleWriter) Finalize() error {
	return nil
}
