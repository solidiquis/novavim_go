package main

import (
	"fmt"
	ansi "github.com/solidiquis/ansigo"
	"log"
	"math"
	"os"
	"os/signal"
	"syscall"
)

type session struct {
	// Window dimensions
	Width  int
	Height int

	// Current mode
	Mode string

	// Text and associated line number
	Lines map[int]string

	// Last line number
	LastLine int

	// Cursor Position
	CursorRow int
	CursorCol int

	// Offsets for cursor boundaries
	ColOffset int
	RowOffset int

	// Whitespace from left of window to line num
	OffsetChars string
}

func InitSession() *session {
	cols, rows, err := ansi.TerminalDimensions()
	if err != nil {
		log.Fatalln(err)
	}

	return &session{
		Width:     cols,
		Height:    rows,
		Mode:      MD_NORMAL,
		Lines:     make(map[int]string),
		CursorRow: CURSOR_COL_START,
		CursorCol: CURSOR_ROW_START,
		ColOffset: CURSOR_COL_START,
		RowOffset: rows - 1,

		// TODO: Should be determined by highest
		// line number in a given file.
		OffsetChars: "  ",
	}
}

func (s *session) InitWindow() {
	ansi.EraseScreen()
	ansi.CursorSetPos(s.Height, 0)
	fmt.Print(ansi.Bright(MD_NORMAL))
	ansi.CursorSetPos(0, 0)
	fmt.Print(ansi.FgYellow("  1 "))
}

func (s *session) WinResizeListener() {
	for {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGWINCH)
		<-sig

		c, r, err := ansi.TerminalDimensions()
		if err != nil {
			log.Fatalln(err)
		}

		s.Height = c
		s.Width = r
	}
}

func (s *session) AddLine() {
	s.LastLine++

	// How much initial whitespace before line number
	offset := int(math.Floor(math.Log10(float64(s.LastLine))+1)) - 1

	ws := s.OffsetChars
	if offsetCol > 0 {
		trim := len(ws) - offset
		ws = s.OffsetChars[:trim]
	}

	// print the new line number
	ansi.CursorSetPos(s.LastLine, 0)
	ln := ansi.FgYellow(fmt.Sprintf("%s%d ", ws, s.LastLine))
	fmt.Print(ln)

	// place cursor in correct position
	switch key {
	case VI_ENTER, VI_o:
		(*cursorPos)["row"]++
	}

	ansi.CursorSetPos((*cursorPos)["row"], CURSOR_COL_START)
}
