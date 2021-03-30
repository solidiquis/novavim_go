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

func switchMode(mode *string, newMode string) {
	_, row, err := ansi.TerminalDimensions()
	if err != nil {
		panic(err)
	}
	ansi.CursorSavePos()

	*mode = newMode

	ansi.CursorSetPos(row, 0)
	ansi.EraseLine()
	fmt.Print(ansi.Bright(newMode))
	ansi.CursorRestorePos()
}

func newLine(cursorPos *map[string]int, lastLine int, key byte) {
	// offset lastLine numbers from column 1
	offsetWs := "  "
	offsetCol := int(math.Floor(math.Log10(float64(lastLine))+1)) - 1
	if offsetCol > 0 {
		offsetWs = offsetWs[:len(offsetWs)-offsetCol]
	}

	// print the new line number
	ansi.CursorSetPos(lastLine, 0)
	fmt.Print(ansi.FgYellow(fmt.Sprintf("%s%d ", offsetWs, lastLine)))

	// place cursor in correct position
	switch key {
	case VI_ENTER, VI_o:
		(*cursorPos)["row"]++
	}

	ansi.CursorSetPos((*cursorPos)["row"], CURSOR_COL_START)
}

func main() {
	sn := InitSession()
	stdin := make(chan string, 1)

	go sn.WinResizeListener()
	go ansi.GetChar(stdin)

	for {
		select {
		case ch := <-stdin:
			if ch[0] == VI_ESC {
				switchMode(&mode, MD_NORMAL)
				continue
			}

			// NORMAL MODE
			if mode == MD_NORMAL {
				switch ch[0] {
				// movement
				case VI_h:
					if sn.CursorCol > colOffset {
						ansi.CursorBackward(1)
						sn.CursorCol["col"]--
					}
				case VI_j:
					currentLine++
					sn.CursorRow++
					ansi.CursorDown(1)
				case VI_k:
					currentLine--
					sn.CursorRow--
					ansi.CursorUp(1)
				case VI_l:
					if sn.CursorCol < len(lines[currentLine])+CURSOR_COL_START { // should be using offset, not CURSOR_COL_START
						ansi.CursorForward(1)
					}

				// delete
				case VI_d:
					subCh := <-stdin
					switch subCh[0] {
					case VI_d:
						ansi.EraseLine()
						ansi.CursorBackward(col)
						ansi.CursorUp(1)
					default:
						continue
					}

				// insert
				case VI_O, VI_o:
					currentLine++
					newLine(&cursorPos, currentLine, ch[0])
					switchMode(&mode, MD_INSERT)
				case VI_i:
					switchMode(&mode, MD_INSERT)
				}
				continue
			}

			// INSERT MODE
			if mode == MD_INSERT {
				switch ch[0] {
				case VI_BACKSPACE:
					lines[currentLine] = lines[currentLine][:len(lines[currentLine])-1]
					ansi.Backspace()
				case VI_ENTER:
					currentLine++
					newLine(&cursorPos, currentLine, ch[0])
				default:
					fmt.Print(string(ch))
					lines[currentLine] += string(ch)
					sn.CursorCol++
				}
			}
		}
	} // for
}
