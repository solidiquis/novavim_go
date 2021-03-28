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

func terminalResizeEvent(col, row *int) {
	for {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGWINCH)
		<-sig

		c, r, err := ansi.TerminalDimensions()
		if err != nil {
			log.Fatalln(err)
		}

		*col = c
		*row = r
	}
}

// TODO: Make work for 4+ digit line nums
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
	// for logging
	f, _ := os.OpenFile("/dev/ttys015", os.O_WRONLY, 0755)

	// Initialize window size
	col, row, err := ansi.TerminalDimensions()
	if err != nil {
		log.Fatalln(err)
	}

	// Initialize initial offsets
	colOffset := CURSOR_COL_START

	// Init mode
	mode := MD_NORMAL

	// Text and associated line number
	lines := make(map[int]string)

	// What line is the cursor on
	currentLine := 1

	// Listen for window resize, update col and row vals
	go terminalResizeEvent(&col, &row)

	// Set cursor start pos and monitor
	cursorPos := map[string]int{
		"col": CURSOR_COL_START,
		"row": CURSOR_ROW_START,
	}

	// Initialize screen
	ansi.EraseScreen()
	ansi.UnbufferStdin()
	ansi.UnechoStdin()
	ansi.CursorSetPos(row, 0)
	fmt.Print(ansi.Bright(MD_NORMAL))
	ansi.CursorSetPos(0, 0)
	fmt.Print(ansi.FgYellow("  1 "))

	stdin := make(chan string, 1)
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
					f.Write([]byte(fmt.Sprint(cursorPos["col"], "\n")))
					if cursorPos["col"] > colOffset {
						ansi.CursorBackward(1)
						cursorPos["col"]--
					}
				case VI_j:
					currentLine++
					cursorPos["row"]++
					ansi.CursorDown(1)
				case VI_k:
					currentLine--
					cursorPos["row"]--
					ansi.CursorUp(1)
				case VI_l:
					if cursorPos["col"] < len(lines[currentLine])+CURSOR_COL_START { // should be using offset, not CURSOR_COL_START
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
					cursorPos["col"]++
				}
			}
		}
	} // for
}
