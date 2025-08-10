//go:build ignore
// +build ignore

package main

import (
	"client"
	"fmt"
	"log"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func AddMessageToList(l *widgets.List, ws_client *client.Client) {
	for {
		message := ws_client.Read()
		l.Rows = append(l.Rows, message)
		ui.Render(l)
	}
}

func main() {

	var name string
	fmt.Println("Enter your name")
	fmt.Scan(&name)

	ws_client := client.NewClient("localhost:8008", name)

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	l := widgets.NewList()
	l.Title = "Chat"
	l.Rows = []string{}
	l.TextStyle = ui.NewStyle(ui.ColorBlue)
	l.WrapText = false
	l.SetRect(0, 0, 60, 16)
	p := widgets.NewParagraph()
	par_text := "Text here: "
	p.Text = par_text
	p.SetRect(0, 17, 60, 20)

	ui.Render(l, p)

	go AddMessageToList(l, ws_client)
	uiEvents := ui.PollEvents()
	text_data := ""
	for {
		e := <-uiEvents
		switch e.ID {
		case "<C-c>":
			return
		case "<Down>":
			l.ScrollDown()
		case "<Up>":
			l.ScrollUp()
		case "<C-d>":
			l.ScrollHalfPageDown()
		case "<C-u>":
			l.ScrollHalfPageUp()
		case "<C-f>":
			l.ScrollPageDown()
		case "<C-b>":
			l.ScrollPageUp()
		case "<Enter>":
			{
				ws_client.Write(text_data)
				p.Text = par_text
				ui.Render(p)
				text_data = ""
			}
		case "<Space>":
			text_data += " "
		case "<C-<Backspace>>", "<Resize>":
			// nothing for now
		default:
			text_data += e.ID
		}
		p.Text = par_text + text_data
		ui.Render(p)
		ui.Render(l)
	}
}
