package main

import (
	"qmail/pages"

	"fyne.io/fyne/v2/app"
)

func main() {
	app := app.NewWithID("com.oq.qmail")
	window := app.NewWindow("QMail")
	pages.StartPage(window)
	window.ShowAndRun()
}
