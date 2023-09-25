package pages

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func account(window fyne.Window) fyne.CanvasObject {
	title := container.NewHBox(layout.NewSpacer(), widget.NewRichTextFromMarkdown("# Account"), layout.NewSpacer())
	imap := widget.NewForm(
		widget.NewFormItem("Address", widget.NewLabel(imapconfig.address)),
		widget.NewFormItem("Username", widget.NewLabel(imapconfig.username)),
		widget.NewFormItem("Password", widget.NewLabel(imapconfig.password)),
	)
	imap.Hide()
	smtp := widget.NewForm(
		widget.NewFormItem("Address", widget.NewLabel(smtpconfig.address)),
		widget.NewFormItem("Username", widget.NewLabel(smtpconfig.username)),
		widget.NewFormItem("Password", widget.NewLabel(smtpconfig.password)),
	)
	smtp.Hide()

	imapbutton := &widget.Button{
		Importance: widget.LowImportance,
		Icon:       theme.MenuDropDownIcon(),
	}
	imapbutton.OnTapped = func() {
		if imap.Hidden {
			imap.Show()
			imapbutton.SetIcon(theme.MenuDropUpIcon())
		} else {
			imap.Hide()
			imapbutton.SetIcon(theme.MenuDropDownIcon())
		}
	}

	smtpbutton := &widget.Button{
		Importance: widget.LowImportance,
		Icon:       theme.MenuDropDownIcon(),
	}
	smtpbutton.OnTapped = func() {
		if smtp.Hidden {
			smtp.Show()
			smtpbutton.SetIcon(theme.MenuDropUpIcon())
		} else {
			smtp.Hide()
			smtpbutton.SetIcon(theme.MenuDropDownIcon())
		}
	}

	return container.NewVScroll(container.NewVBox(
		title,
		container.NewHBox(widget.NewRichTextFromMarkdown("## SMTP"), smtpbutton),
		smtp,
		container.NewHBox(widget.NewRichTextFromMarkdown("## IMAP"), imapbutton),
		imap,
		container.NewHBox(&widget.Button{
			Text:       "Log Out",
			Importance: widget.DangerImportance,
			OnTapped: func() {
				Reset()
				settingsdialog.Hide()
				StartPage(window)
			},
		}, &widget.Button{
			Text:       "Edit",
			Importance: widget.WarningImportance,
			OnTapped: func() {
				settingsdialog.Hide()
				LoginPage(window)
			},
		}),
	))
}

var settingsdialog *dialog.CustomDialog
var buttons = []fyne.CanvasObject{
	&widget.Button{
		Text:       "Account",
		Importance: widget.HighImportance,
		OnTapped:   nil,
	},
	layout.NewSpacer(),
	&widget.Button{
		Text: "Close",
		OnTapped: func() {
			settingsdialog.Hide()
		},
		Importance: widget.DangerImportance,
	},
}

func Settings(window fyne.Window) {
	left := container.NewVBox(buttons...)

	cont := container.NewBorder(nil, nil, left, nil, account(window))
	settingsdialog = dialog.NewCustomWithoutButtons("Settings", cont, window)
	settingsdialog.Resize(fyne.NewSize(500, 400))
	settingsdialog.Show()
}
