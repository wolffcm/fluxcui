package main

import (
	"bytes"
	"errors"
	"github.com/chzyer/readline"
	"github.com/jroimartin/gocui"
	"io"
	"log"
	"os"
	"sync"
	"unicode"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	//g.Mouse = true
	g.SelFgColor = gocui.ColorGreen
	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) (err error) {
	maxX, maxY := g.Size()
	if v, err := g.SetView("hello", 0, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = true
		v.Autoscroll = true
		v.Wrap = true

		editor, err := newReadlineEditor(g, v)
		v.Editor = editor
		if err != nil {
			return err
		}

		if _, err := g.SetCurrentView("hello"); err != nil {
			return err
		}
		if _, err := g.SetViewOnTop("hello"); err != nil {
			return err
		}
		go readLines(g, editor)

		return nil

	}
	return nil
}

type readlineEditor struct {
	keyboard *os.File

	rl *readline.Instance
}

func newReadlineEditor(g *gocui.Gui, v *gocui.View) (*readlineEditor, error) {
	stdin, in, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	stdout := &guiWriter{g: g, w: v}
	getWidth := func() int {
		x, _ := v.Size()
		return x
	}
	cfg := &readline.Config{
		Prompt:                 "> ",
		HistoryFile:            "",
		HistoryLimit:           0,
		DisableAutoSaveHistory: false,
		HistorySearchFold:      false,
		AutoComplete:           nil,
		Listener:               nil,
		Painter:                nil,
		VimMode:                false,
		InterruptPrompt:        "",
		EOFPrompt:              "",
		FuncGetWidth:           getWidth,
		Stdin:                  stdin,
		StdinWriter:            nil,
		Stdout:                 stdout,
		Stderr:                 stdout,
		EnableMask:             false,
		MaskRune:               0,
		UniqueEditLine:         false,
		FuncFilterInputRune:    nil,
		FuncIsTerminal:         func() bool { return false },
		FuncMakeRaw:            func() error { return nil },
		FuncExitRaw:            func() error { return nil },
		FuncOnWidthChanged:     nil,
		ForceUseInteractive:    true,
	}
	rl, err := readline.NewEx(cfg)
	if err != nil {
		return nil, err
	}
	return &readlineEditor{
		keyboard: in,
		rl:       rl,
	}, nil
}

func readLines(g *gocui.Gui, e *readlineEditor) {
	for {
		_, err := e.rl.Readline()
		if err != nil {
			return
		}
	}
}

func (r *readlineEditor) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	var bb bytes.Buffer
	if ch != 0 {
		if mod == gocui.ModAlt {
			bb.WriteRune('\033')
		}
		bb.WriteRune(ch)
	} else {
		if mod == gocui.ModAlt {
			bb.WriteRune('\033')
		}
		switch key {
		case gocui.KeyF1:
		case gocui.KeyF2:
		case gocui.KeyF3:
		case gocui.KeyF4:
		case gocui.KeyF5:
		case gocui.KeyF6:
		case gocui.KeyF7:
		case gocui.KeyF8:
		case gocui.KeyF9:
		case gocui.KeyF10:
		case gocui.KeyF11:
		case gocui.KeyF12:
		case gocui.KeyInsert:
		case gocui.KeyDelete:
			bb.WriteRune(readline.CharBackspace)
		case gocui.KeyHome:
		case gocui.KeyEnd:
		case gocui.KeyPgup:
		case gocui.KeyPgdn:
		case gocui.KeyArrowUp:
			bb.WriteRune(readline.CharPrev)
		case gocui.KeyArrowDown:
			bb.WriteRune(readline.CharNext)
		case gocui.KeyArrowLeft:
			bb.WriteRune(readline.CharBackward)
		case gocui.KeyArrowRight:
			bb.WriteRune(readline.CharForward)
		case gocui.MouseLeft, gocui.MouseMiddle, gocui.MouseRight,
			gocui.MouseRelease, gocui.MouseWheelUp, gocui.MouseWheelDown:
			// Do nothing

		//case gocui.KeyCtrlTilde:
		//case gocui.KeyCtrl2:
		//case gocui.KeyCtrlSpace:
		case gocui.KeyCtrlA:
			bb.WriteRune(readline.CharLineStart)
		//case gocui.KeyCtrlB:
		//case gocui.KeyCtrlC:
		//case gocui.KeyCtrlD:
		case gocui.KeyCtrlE:
			bb.WriteRune(readline.CharLineEnd)
		//case gocui.KeyCtrlF:
		//case gocui.KeyCtrlG:
		//case gocui.KeyBackspace:
		//case gocui.KeyCtrlH:
		//case gocui.KeyTab:
		//case gocui.KeyCtrlI:
		//case gocui.KeyCtrlJ:
		case gocui.KeyCtrlK:
			bb.WriteRune(readline.CharKill)
		//case gocui.KeyCtrlL:
		case gocui.KeyEnter:
			bb.WriteRune(readline.CharEnter)
		//case gocui.KeyCtrlM:
		//case gocui.KeyCtrlN:
		//case gocui.KeyCtrlO:
		//case gocui.KeyCtrlP:
		//case gocui.KeyCtrlQ:
		//case gocui.KeyCtrlR:
		//case gocui.KeyCtrlS:
		//case gocui.KeyCtrlT:
		//case gocui.KeyCtrlU:
		//case gocui.KeyCtrlV:
		//case gocui.KeyCtrlW:
		//case gocui.KeyCtrlX:
		case gocui.KeyCtrlY:
			bb.WriteRune(readline.CharCtrlY)
		//case gocui.KeyCtrlZ:
		//case gocui.KeyEsc:
		//case gocui.KeyCtrlLsqBracket:
		//case gocui.KeyCtrl3:
		//case gocui.KeyCtrl4:
		//case gocui.KeyCtrlBackslash:
		//case gocui.KeyCtrl5:
		//case gocui.KeyCtrlRsqBracket:
		//case gocui.KeyCtrl6:
		//case gocui.KeyCtrl7:
		//case gocui.KeyCtrlSlash:
		//case gocui.KeyCtrlUnderscore:
		case gocui.KeySpace:
			bb.WriteRune(' ')
		case gocui.KeyBackspace2:
			bb.WriteRune(readline.CharBackspace)
			//case gocui.KeyCtrl8:
		}
	}
	bb.WriteTo(r.keyboard)
}

type guiWriter struct {
	m sync.Mutex

	g *gocui.Gui
	w *gocui.View
}

func handleEscapeSequence(v *gocui.View, bb *bytes.Buffer) error {
	var seq string
	for {
		rn, _, err := bb.ReadRune()
		if err != nil {
			return err
		}

		seq += string(rn)
		if unicode.IsLetter(rn) {
			break
		}
	}
	switch seq {
	case "[2K":
		v.EditClearLine()
	case "[J":
		v.EditTruncateBuffer()
	default:
		// ignore
	}
	return nil
}

func (gw *guiWriter) Write(p []byte) (n int, err error) {
	bb := bytes.NewBuffer(p)
	gw.m.Lock()
	gw.g.Update(func(*gocui.Gui) error {
		defer gw.m.Unlock()
		for {
			var rn rune
			rn, _, err = bb.ReadRune()
			if err != nil && err != io.EOF {
				return err
			} else if err == io.EOF {
				break
			}
			switch rn {
			case '\a':
				// bell; do nothing
			case '\033':
				if err := handleEscapeSequence(gw.w, bb); err != nil {
					return err
				}
			case '\n':
				gw.w.EditNewLine()
			case '\r':
				gw.w.EditCarriageReturn()
			case '\b':
				gw.w.MoveCursor(-1, 0, true)
			default:
				if rn < ' ' {
					// Some control character we do not yet handle
					return errors.New("unrecognized control char")
				}
				gw.w.EditWrite(rn)
			}
		}
		return nil
	})

	return n, err
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
