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
		LastLine:  1,
		CursorRow: CURSOR_ROW_START,
		CursorCol: CURSOR_COL_START,
		ColOffset: CURSOR_COL_START,
		RowOffset: rows - 1,

		// TODO: Should be determined by highest
		// line number in a given file.
		OffsetChars: "  ",
	}
}

func (sn *session) InitWindow() {
	ansi.EraseScreen()
	ansi.CursorSetPos(sn.Height, 0)
	fmt.Print(ansi.Bright(MD_NORMAL))
	ansi.CursorSetPos(0, 0)
	fmt.Print(ansi.FgYellow("  1 "))
}

func (sn *session) WinResizeListener() {
	for {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGWINCH)
		<-sig

		c, r, err := ansi.TerminalDimensions()
		if err != nil {
			log.Fatalln(err)
		}

		sn.Height = c
		sn.Width = r
	}
}

func (sn *session) AddLine(key byte) {
	sn.LastLine++

	// How much initial whitespace before line number
	offset := int(math.Floor(math.Log10(float64(sn.LastLine))+1)) - 1

	ws := sn.OffsetChars
	if sn.ColOffset > 0 {
		trim := len(ws) - offset
		ws = sn.OffsetChars[:trim]
	}

	// print the new line number
	ansi.CursorSetPos(sn.LastLine, 0)
	ln := ansi.FgYellow(fmt.Sprintf("%s%d ", ws, sn.LastLine))
	fmt.Print(ln)

	// place cursor in correct position
	switch key {
	case VI_ENTER, VI_o:
		sn.CursorRow++
	}

	ansi.CursorSetPos(sn.CursorRow, CURSOR_COL_START)
}

func (sn *session) SetMode(mode string) {
	ansi.CursorSavePos()
	sn.Mode = mode
	ansi.CursorSetPos(sn.Height, 0)
	ansi.EraseLine()
	fmt.Print(ansi.Bright(mode))
	ansi.CursorRestorePos()
}
