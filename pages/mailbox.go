package pages

import (
	"encoding/base64"
	"fmt"
	"math"
	"os"
	"qmail/client"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/sqweek/dialog"

	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var page uint32 = 1
var selected = 0
var selectedname string

type attachment struct {
	name string
	data []byte
}

func Main(window fyne.Window, p uint32) {
	page = p
	cont := container.NewHSplit(Mailboxes(window), Mails(p, window))
	cont.SetOffset(0)

	window.SetContent(cont)
}

func resync(window fyne.Window) {
	cont := container.NewHSplit(Mailboxes(window), container.NewCenter(widget.NewRichTextFromMarkdown("# Syncing...")))
	cont.SetOffset(0)
	window.SetContent(cont)

	Main(window, page)
}

func CalculateMaxPage(totalEntries uint32) uint32 {
	return uint32(math.Ceil(float64(totalEntries) / 25))
}

func OpenMail(mail *imapclient.FetchMessageBuffer, window fyne.Window) {
	var attachments []*attachment
	back := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		Main(window, page)
	})

	var content string
	var body string
	for _, b := range mail.BodySection {
		body = string(b)
	}

	for _, section := range strings.Split(body, "--") {
		if section == "" {
			continue
		}
		sp := strings.Split(section, "\n")
		section = strings.TrimSpace(strings.TrimPrefix(section, sp[0]))
		sp = strings.Split(section, "Content-Type")
		for i, s := range sp {
			if s == "" {
				continue
			}
			spp := strings.Split(s, "\n")
			c := spp[0]
			c = strings.TrimPrefix(c, ": ")
			if strings.HasPrefix(c, "text/plain") {
				content += strings.TrimPrefix(strings.Join(spp[1:], "\n"), "Content-Transfer-Encoding: quoted-printable") + "\n"
			} else if strings.HasPrefix(c, "image/") {
				var data string
				var attr = make(map[string]string)
				for i, l := range spp {
					if i == 0 {
						continue
					}
					sp := strings.Split(l, ":")
					if len(sp) != 2 {
						data += l
					} else {
						//fmt.Println(sp[0])
						attr[sp[0]] = strings.TrimSpace(sp[1])
					}
				}
				var filename string
				encoding, ok := attr["Content-Transfer-Encoding"]
				if !ok {
					continue
				}
				if cd, ok := attr["Content-Disposition"]; ok {
					sp := strings.Split(cd, " ")
					for _, s := range sp {
						if strings.HasPrefix(s, "filename=") {
							filename = strings.TrimSuffix(strings.TrimPrefix(s, `filename="`), `"`)
						}
					}
				} else {
					continue
				}
				if filename == "" {
					filename = fmt.Sprintf("File %d", i)
				}
				switch encoding {
				case "base64":
					{
						d, err := base64.StdEncoding.DecodeString(data)
						if err != nil {
							continue
						}
						attachments = append(attachments, &attachment{name: filename, data: d})
					}
				}
			}
		}
	}

	text := widget.NewLabel(content)

	from := mail.Envelope.Sender[0].Addr()
	if mail.Envelope.Sender[0].Name != "" && mail.Envelope.Sender[0].Name != mail.Envelope.Sender[0].Addr() {
		from += fmt.Sprintf(" (%s)", mail.Envelope.Sender[0].Name)
	}
	if mail.Envelope.Sender[0].Addr() == imapconfig.username {
		from += " (*You*)"
	}

	to := mail.Envelope.To[0].Addr()
	if mail.Envelope.To[0].Name != "" && mail.Envelope.To[0].Name != mail.Envelope.To[0].Addr() {
		to += fmt.Sprintf(" (%s)", mail.Envelope.To[0].Name)
	}
	if mail.Envelope.To[0].Addr() == imapconfig.username {
		to += " (*You*)"
	}

	s := mail.Envelope.Subject
	if s == "" {
		s = "No Subject"
	}

	sender := widget.NewRichTextFromMarkdown(fmt.Sprintf("From: **%s**", from))
	reciever := widget.NewRichTextFromMarkdown(fmt.Sprintf("To: **%s**", to))
	subject := widget.NewRichTextFromMarkdown(fmt.Sprintf("## %s", s))

	envelope := container.NewVBox(sender, reciever)

	attachs := container.NewHBox()

	for _, attachment := range attachments {
		a := widget.NewButtonWithIcon(attachment.name, theme.DownloadIcon(), func() {
			d := dialog.File()
			d.StartFile = attachment.name
			p, err := d.Save()
			if err == nil {
				os.WriteFile(p, attachment.data, 0755)
			}
		})
		a.Importance = widget.LowImportance
		attachs.Add(a)
	}

	co := container.NewBorder(envelope, nil, nil, nil, container.NewBorder(nil, attachs, nil, nil, text))

	cont := container.NewBorder(
		container.NewHBox(
			subject,
			layout.NewSpacer(),
			back,
		),
		nil,
		nil,
		nil,
		container.NewVScroll(container.NewHScroll(co)),
	)

	conte := container.NewHSplit(Mailboxes(window), cont)
	conte.SetOffset(0)

	window.SetContent(conte)
}

func Mails(page uint32, window fyne.Window) fyne.CanvasObject {
	newMail := widget.NewButtonWithIcon("New Mail", theme.ContentAddIcon(), nil)
	newMail.Importance = widget.HighImportance

	sync := widget.NewButtonWithIcon("", theme.MediaReplayIcon(), func() {
		resync(window)
	})
	prev := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		Main(window, page-1)
	})
	if page-1 == 0 {
		prev.Disable()
	}
	next := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
		Main(window, page+1)
	})
	if page == CalculateMaxPage(client.Selected.NumMessages) {
		next.Disable()
	}
	mails, _ := client.GetMails(page)

	list := widget.NewList(func() int {
		return len(mails)
	}, func() fyne.CanvasObject {
		return widget.NewRichTextWithText("")
	}, func(lii widget.ListItemID, co fyne.CanvasObject) {
		text := mails[lii].Envelope.Subject
		if text == "" {
			text = "No Subject"
		}
		co.(*widget.RichText).ParseMarkdown(fmt.Sprintf("### %s", text))
	})

	list.OnSelected = func(id widget.ListItemID) {
		OpenMail(mails[id], window)
	}

	if len(mails) == 0 {
		return container.NewCenter(widget.NewRichTextFromMarkdown("# No Mails"))
	}
	return container.NewBorder(
		container.NewHBox(
			widget.NewRichTextFromMarkdown("## "+selectedname),
			layout.NewSpacer(),
			newMail,
			sync,
			prev,
			next,
		),
		nil,
		nil,
		nil,
		list,
	)
}

func Mailboxes(window fyne.Window) fyne.CanvasObject {
	cont := container.NewVBox()
	for i, mb := range client.Mailboxes {
		d := i
		mailbox := mb
		name := strings.TrimPrefix(mailbox.Mailbox, "[Gmail]")
		if name == "" {
			continue
		}
		name = strings.TrimPrefix(name, "/")
		button := widget.NewButton(name, nil)
		button.Importance = widget.LowImportance
		button.OnTapped = func() {
			cont.Objects[selected].(*widget.Button).Importance = widget.LowImportance
			cont.Objects[selected].(*widget.Button).Refresh()
			selected = d
			selectedname = name
			button.Importance = widget.HighImportance
			button.Refresh()
			client.Select(mailbox.Mailbox)
			page = 1
			Main(window, 1)
		}
		if d == selected {
			button.Importance = widget.HighImportance
			client.Select(mailbox.Mailbox)
		}
		cont.Add(button)
	}
	settings := widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
		Settings(window)
	})
	settings.Importance = widget.LowImportance
	cont.Add(layout.NewSpacer())
	cont.Add(settings)
	return cont
}
