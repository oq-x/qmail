package pages

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var d *dialog.CustomDialog

type config struct {
	address  string
	username string
	password string
}

var smtpconfig config
var imapconfig config

func Save() {
	app := fyne.CurrentApp()
	app.Preferences().SetString("IMAP_ADDRESS", imapconfig.address)
	app.Preferences().SetString("IMAP_USERNAME", imapconfig.username)
	app.Preferences().SetString("IMAP_PASSWORD", imapconfig.password)

	app.Preferences().SetString("SMTP_ADDRESS", smtpconfig.address)
	app.Preferences().SetString("SMTP_USERNAME", smtpconfig.username)
	app.Preferences().SetString("SMTP_PASSWORD", smtpconfig.password)
}

func Reset() {
	app := fyne.CurrentApp()
	app.Preferences().RemoveValue("IMAP_ADDRESS")
	app.Preferences().RemoveValue("IMAP_USERNAME")
	app.Preferences().RemoveValue("IMAP_PASSWORD")

	app.Preferences().RemoveValue("SMTP_ADDRESS")
	app.Preferences().RemoveValue("SMTP_USERNAME")
	app.Preferences().RemoveValue("SMTP_PASSWORD")
}

func Load() bool {
	app := fyne.CurrentApp()
	ia := app.Preferences().String("IMAP_ADDRESS")
	iu := app.Preferences().String("IMAP_USERNAME")
	ip := app.Preferences().String("IMAP_PASSWORD")

	sa := app.Preferences().String("SMTP_ADDRESS")
	su := app.Preferences().String("SMTP_USERNAME")
	sp := app.Preferences().String("SMTP_PASSWORD")

	if ia == "" || iu == "" || ip == "" || sa == "" || su == "" || sp == "" {
		return false
	}
	imapconfig.address = ia
	imapconfig.username = iu
	imapconfig.password = ip

	smtpconfig.address = sa
	smtpconfig.username = su
	smtpconfig.password = sp

	return true
}

func sendingForm(window fyne.Window) *fyne.Container {
	form := widget.NewForm()
	address := widget.NewFormItem("Address", &widget.Entry{
		Wrapping: fyne.TextTruncate,
		Text:     smtpconfig.address,
	})
	username := widget.NewFormItem("Username", &widget.Entry{
		Wrapping: fyne.TextTruncate,
		Text:     smtpconfig.username,
	})
	p := widget.NewPasswordEntry()
	p.SetText(smtpconfig.password)
	password := widget.NewFormItem("Password", p)
	form.AppendItem(address)
	form.AppendItem(username)
	form.AppendItem(password)

	title := widget.NewRichTextFromMarkdown("# Sending (SMTP)")
	subtitle := widget.NewLabel("one or more inputs are empty")
	subtitle.Hide()

	next := widget.NewButton("Next", func() {
		address := address.Widget.(*widget.Entry).Text
		username := username.Widget.(*widget.Entry).Text
		password := password.Widget.(*widget.Entry).Text
		if address == "" || username == "" || password == "" {
			subtitle.Show()
			return
		}
		smtpconfig.address = address
		smtpconfig.username = username
		smtpconfig.password = password
		d.Hide()
		loginPage(window, recievingForm(window))
	})
	next.Importance = widget.HighImportance

	return container.NewVBox(
		title,
		subtitle,
		form,
		container.NewGridWithColumns(3, layout.NewSpacer(), next),
	)
}

func recievingForm(window fyne.Window) *fyne.Container {
	form := widget.NewForm()
	address := widget.NewFormItem("Address", &widget.Entry{
		Wrapping: fyne.TextTruncate,
		Text:     imapconfig.address,
	})
	username := widget.NewFormItem("Username", &widget.Entry{
		Wrapping: fyne.TextTruncate,
		Text:     imapconfig.username,
	})
	p := widget.NewPasswordEntry()
	p.SetText(imapconfig.password)
	password := widget.NewFormItem("Password", p)

	form.AppendItem(address)
	form.AppendItem(username)
	form.AppendItem(password)

	title := widget.NewRichTextFromMarkdown("# Recieving (IMAP)")
	subtitle := widget.NewLabel("one or more inputs are empty")
	subtitle.Hide()

	next := widget.NewButton("Login", func() {
		address := address.Widget.(*widget.Entry).Text
		username := username.Widget.(*widget.Entry).Text
		password := password.Widget.(*widget.Entry).Text
		if address == "" || username == "" || password == "" {
			subtitle.Show()
			return
		}
		imapconfig.address = address
		imapconfig.username = username
		imapconfig.password = password
		d.Hide()
		ConnectingPage(window)
	})
	next.Importance = widget.HighImportance

	return container.NewVBox(
		title,
		subtitle,
		form,
		container.NewGridWithColumns(3, layout.NewSpacer(), next),
	)
}
func loginPage(window fyne.Window, cont *fyne.Container) {
	d = dialog.NewCustomWithoutButtons("Login", cont, window)
	d.Resize(fyne.NewSize(400, 200))
	d.Show()
}
func LoginPage(window fyne.Window) {
	loginPage(window, sendingForm(window))
}

func StartPage(window fyne.Window) {
	window.Resize(fyne.NewSize(600, 400))
	title := widget.NewRichTextFromMarkdown("# Welcome to QMail")
	start := widget.NewButton("Get Started", func() {
		LoginPage(window)
	})
	start.Importance = widget.HighImportance

	if Load() {
		ConnectingPage(window)
	} else {
		window.SetContent(container.NewVBox(container.NewHBox(layout.NewSpacer(), title, layout.NewSpacer()), layout.NewSpacer(), container.NewGridWithColumns(3, layout.NewSpacer(), start), layout.NewSpacer()))
	}
}
