package main

import (
	"fmt"
	ansi "github.com/solidiquis/ansigo"
)

func main() {
	sn := InitSession()
	stdin := make(chan string, 1)

	sn.InitWindow()

	go sn.WinResizeListener()
	go ansi.GetChar(stdin)

	wr := Debug("ttys014")

	for {
		wr([]byte(fmt.Sprint(sn.CursorCol)))
		wr([]byte(fmt.Sprint(sn.CursorRow)))
		select {
		case ch := <-stdin:
			if ch[0] == VI_ESC {
				sn.SetMode(MD_NORMAL)
				continue
			}

			// NORMAL MODE
			if sn.Mode == MD_NORMAL {
				switch ch[0] {
				// movement
				case VI_h:
					sn.CursorLeft(1)
				case VI_j:
					sn.CursorDown(1)
				case VI_k:
					sn.CursorRow--
					ansi.CursorUp(1)
				case VI_l:
					sn.CursorRight(1)

				// delete
				case VI_d:
					subCh := <-stdin
					switch subCh[0] {
					case VI_d:
						ansi.EraseLine()
						ansi.CursorBackward(1) // Delete to beginning of column offset
						ansi.CursorUp(1)
					default:
						continue
					}

				// insert
				case VI_O, VI_o:
					sn.AddLine(ch[0])
					sn.SetMode(MD_INSERT)
				case VI_i:
					sn.SetMode(MD_INSERT)
				}
				continue
			}

			// INSERT MODE
			if sn.Mode == MD_INSERT {
				switch ch[0] {
				case VI_BACKSPACE:
					sn.Backspace()
				case VI_ENTER:
					sn.AddLine(ch[0])
				default:
					fmt.Print(string(ch))
					sn.Lines[sn.CursorRow] += string(ch)
					sn.CursorCol++
				}
			}
		}
	} // for
}
