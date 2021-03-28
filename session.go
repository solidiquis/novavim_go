package main

import (
	"fmt"
	ansi "github.com/solidiquis/ansigo"
	"log"

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

	// Cursor Position
	CursorRow int
	CursorCol int

	// Offsets for cursor boundaries
	ColOffset int
	RowOffset int
}

func initSession() *session {
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
	}
}

func (s *session) initWindow() {
	ansi.EraseScreen()
	ansi.UnbufferStdin()
	ansi.UnechoStdin()
	ansi.CursorSetPos(s.Height, 0)
	fmt.Print(ansi.Bright(MD_NORMAL))
	ansi.CursorSetPos(0, 0)
	fmt.Print(ansi.FgYellow("  1 "))
}

func (s *session) winResizeListener() {
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
