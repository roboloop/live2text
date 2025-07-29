package text

import (
	"slices"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"
)

type lineBreak struct {
	segmentIndex int
	runeOffset   int
}

// Formatter is responsible for displaying real-time recognized text as
// subtitles, maintaining a fixed number of visible lines (default: 2).
// The format of incoming data: `text` + `isFinal` mark.
// If `isFinal` is true, the previous section is replaced with the current.
//
// Older content is discarded once the line limit is exceeded,
// ensuring only the most recent lines are shown.
type Formatter struct {
	lineCount     int
	lineRuneCount int

	prevIsFinal      bool
	lineSeparator    string
	segmentSeparator string
	segments         []string
	lineBreaks       []lineBreak
	mu               sync.RWMutex

	commitedLinesOffset          int
	lineBreaksOfPrevSegmentCount int
}

// NewSubtitleFormatter creates a new instance, where
// `displayLineLimit` emphasizes an exact number of output lines
// `lineRuneCount` describes line length in character count.
func NewSubtitleFormatter(lineCount int, lineCharCount int) *Formatter {
	if lineCount == 0 {
		lineCount = 2
	}
	if lineCharCount == 0 {
		lineCharCount = 80
	}
	return &Formatter{
		lineCount:     lineCount,
		lineRuneCount: lineCharCount,

		prevIsFinal:      false,
		lineSeparator:    "\n",
		segmentSeparator: " ",
		segments:         []string{},
		lineBreaks:       []lineBreak{},
	}
}

func (f *Formatter) Append(text string, isFinal bool) {
	f.mu.Lock()
	defer f.mu.Unlock()

	text = sanitize(text)
	if text == "" && f.prevIsFinal && isFinal {
		return
	}

	text = capitalize(text)
	if isFinal {
		text = finalize(text)
	}

	if len(f.segments) == 0 || f.prevIsFinal {
		f.segments = append(f.segments, text)
	} else {
		f.segments[len(f.segments)-1] = text
	}
	f.clearLineBreaksOfLastSegment()
	f.fillLineBreaksOfLastSegment()
	f.adjustSkippedSegments()
	f.prevIsFinal = isFinal

	f.dropOutdatedSegments()
}

func (f *Formatter) clearLineBreaksOfLastSegment() {
	activeSegmentIndex := len(f.segments) - 1
	index := slices.IndexFunc(f.lineBreaks, func(lb lineBreak) bool {
		return lb.segmentIndex == activeSegmentIndex
	})
	if index != -1 {
		f.lineBreaks = f.lineBreaks[:index]
	}
}

func (f *Formatter) fillLineBreaksOfLastSegment() {
	lineLength := f.commitedLineLength()
	activeSegmentIndex := len(f.segments) - 1
	segment := f.segments[activeSegmentIndex]
	runeOffset := 0
	for i, word := range strings.Fields(segment) {
		// Dirty hack for not implementing that case
		if utf8.RuneCountInString(word) > f.lineRuneCount {
			runes := []rune(word)
			word = string(runes[:f.lineRuneCount])
		}

		var wordLength int
		// Count the space before the word (not applied to the first word)
		if i > 0 {
			wordLength = utf8.RuneCountInString(f.segmentSeparator)
		}
		wordLength += utf8.RuneCountInString(word)
		lineLength += wordLength

		// If the line is longer than should be, then cut it
		if lineLength > f.lineRuneCount {
			f.lineBreaks = append(f.lineBreaks, lineBreak{segmentIndex: activeSegmentIndex, runeOffset: runeOffset})
			lineLength = utf8.RuneCountInString(word)
		}
		runeOffset += wordLength
	}
}

func (f *Formatter) adjustSkippedSegments() {
	if f.prevIsFinal {
		f.commitedLinesOffset = 0
		f.lineBreaksOfPrevSegmentCount = 0
		return
	}

	activeSegmentIndex := len(f.segments) - 1
	lineBreaksOfActiveSegmentCount := 0
	for _, lb := range f.lineBreaks {
		if lb.segmentIndex == activeSegmentIndex {
			lineBreaksOfActiveSegmentCount += 1
		}
	}

	f.commitedLinesOffset = max(0, f.lineBreaksOfPrevSegmentCount-lineBreaksOfActiveSegmentCount)
	f.lineBreaksOfPrevSegmentCount = lineBreaksOfActiveSegmentCount
}

func (f *Formatter) dropOutdatedSegments() {
	lineBreaksCount := 0
	if !f.prevIsFinal {
		for _, line := range f.lineBreaks {
			if line.segmentIndex == len(f.segments)-1 {
				lineBreaksCount++
			}
		}
	}

	if len(f.lineBreaks)-lineBreaksCount <= f.lineCount {
		return
	}

	linesToDrop := len(f.lineBreaks) - f.lineCount - lineBreaksCount
	f.lineBreaks = f.lineBreaks[linesToDrop:]

	segmentIndexOffset := f.lineBreaks[0].segmentIndex
	f.segments = f.segments[segmentIndexOffset:]
	for i := range f.lineBreaks {
		f.lineBreaks[i].segmentIndex -= segmentIndexOffset
	}
}

func (f *Formatter) commitedLineLength() int {
	commitedLineLength := 0
	startSegmentIndex := 0
	if len(f.lineBreaks) > 0 {
		lastLineBreak := f.lineBreaks[len(f.lineBreaks)-1]
		segmentIndex := lastLineBreak.segmentIndex
		runeOffset := lastLineBreak.runeOffset
		commitedLineLength = utf8.RuneCountInString(f.segments[segmentIndex]) - runeOffset

		startSegmentIndex = segmentIndex + 1
	}
	for i := startSegmentIndex; i < len(f.segments)-1; i++ {
		commitedLineLength += utf8.RuneCountInString(f.segmentSeparator) + utf8.RuneCountInString(f.segments[i])
	}

	return commitedLineLength
}

func (f *Formatter) Format() string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if len(f.segments) == 0 {
		return ""
	}

	var lines [][]string
	// At least one line output is expected
	lines = append(lines, []string{})

	for i, s := range f.segments {
		var runeOffsets []int

		// Find all line breaks that apply to this section
		for _, lb := range f.lineBreaks {
			if lb.segmentIndex == i {
				runeOffsets = append(runeOffsets, lb.runeOffset)
			}
		}

		// If no line breaks in this section, add entire text to current line
		if len(runeOffsets) == 0 {
			lines[len(lines)-1] = append(lines[len(lines)-1], s)
			continue
		}

		// combine all segments according to their offsets
		runes := []rune(s)
		prevRuneOffset := 0
		for _, runeOffset := range runeOffsets {
			if runeOffset != 0 {
				fragment := strings.TrimLeft(string(runes[prevRuneOffset:runeOffset]), f.segmentSeparator)
				lines[len(lines)-1] = append(lines[len(lines)-1], fragment)
			}

			lines = append(lines, []string{})
			prevRuneOffset = runeOffset
		}

		// Add the final segment (from last offset to end of the segment)
		lines[len(lines)-1] = append(
			lines[len(lines)-1],
			strings.TrimLeft(string(runes[prevRuneOffset:]), f.segmentSeparator),
		)
	}

	// Ensure we only show the most recent lines if exceeding our display limit
	start := 0
	if len(lines) > f.lineCount {
		start = len(lines) - f.lineCount
	}
	start += f.commitedLinesOffset

	var committedLines []string
	for _, line := range lines[start:] {
		committedLines = append(committedLines, strings.Join(line, f.segmentSeparator))
	}

	// Padding the list of committed lines
	for i := len(committedLines); i < f.lineCount; i++ {
		committedLines = append(committedLines, "")
	}

	return strings.Join(committedLines, f.lineSeparator)
}

func sanitize(text string) string {
	return strings.TrimSpace(text)
}

func capitalize(text string) string {
	if len(text) == 0 {
		return text
	}
	runes := []rune(text)
	if r := runes[0]; unicode.IsLetter(r) && unicode.IsLower(r) {
		runes[0] = unicode.ToUpper(r)
		return string(runes)
	}
	return text
}

func finalize(text string) string {
	return text + "."
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
