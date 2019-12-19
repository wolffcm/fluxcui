package view

import (
	"fmt"
	
	"github.com/jroimartin/gocui"
)

func (c *cui) setEditorKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding(editorView, gocui.MouseLeft, gocui.ModNone, setcurrentView); err != nil {
		return err
	}
	return nil
}

func (c *cui) doEditorView(g *gocui.Gui, x0, y0, x1, y1 int) error {
	editor, err := g.SetView(editorView, x0, y0, x1, y1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		editor.Title = "editor"
		editor.Editable = true
		editor.Wrap = true

		if _, err = fmt.Fprintf(editor, c.cfg.EditorText); err != nil {
			return err
		}

		if _, err := setCurrentViewOnTop(g, editorView); err != nil {
			return err
		}
	}

	return nil
}
