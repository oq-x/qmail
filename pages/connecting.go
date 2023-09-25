package pages

import (
	"fmt"
	"qmail/client"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func ConnectingPage(window fyne.Window) {
	title := widget.NewRichTextFromMarkdown("# Connecting...")

	window.SetContent(container.NewCenter(title))

	err := client.ConnectIMAP(imapconfig.address, imapconfig.username, imapconfig.password)
	if err != nil {
		ConnectionFailure(window, fmt.Sprintf("IMAP: %s", err))
		return
	}
	title.ParseMarkdown("# Syncing...")
	Save()
	client.GetMailboxes()
	Main(window, 1)
}

func ConnectionFailure(window fyne.Window, err string) {
	title := container.NewCenter(widget.NewRichTextFromMarkdown("# Error: " + err))
	retry := widget.NewButton("Retry", func() {
		LoginPage(window)
	})
	window.SetContent(container.NewVBox(layout.NewSpacer(), title, container.NewGridWithColumns(3, layout.NewSpacer(), retry), layout.NewSpacer()))
}
