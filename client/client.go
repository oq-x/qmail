package client

import (
	"mime"

	"slices"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/emersion/go-message"

	_ "github.com/emersion/go-message/charset"
)

var imapc *imapclient.Client
var Mailboxes []*imap.ListData
var Selected *imapclient.SelectedMailbox

func ConnectIMAP(address, username, password string) error {
	options := &imapclient.Options{
		WordDecoder: &mime.WordDecoder{CharsetReader: message.CharsetReader},
	}
	client, err := imapclient.DialTLS(address, options)
	if err != nil {
		return err
	}
	err = client.Login(username, password).Wait()
	imapc = client
	return err
}

func GetMailboxes() {
	Mailboxes, _ = imapc.List("", "*", nil).Collect()
}

func Select(mailbox string) {
	if Selected != nil && Selected.Name == mailbox {
		return
	}
	imapc.Select(mailbox, nil).Wait()
	Selected = imapc.Mailbox()
}

func GetMails(page uint32) ([]*imapclient.FetchMessageBuffer, error) {
	pageSize := uint32(25)

	endMessage := Selected.NumMessages - (page-1)*pageSize
	startMessage := endMessage - pageSize + 1

	if startMessage < 1 {
		startMessage = 1
	}
	if endMessage > Selected.NumMessages {
		endMessage = Selected.NumMessages
	}

	seqset := imap.SeqSet{}
	seqset.AddRange(startMessage, endMessage)
	mails, err := imapc.Fetch(seqset, &imap.FetchOptions{
		Envelope: true,
		Flags:    true,
		BodySection: []*imap.FetchItemBodySection{
			{Specifier: imap.PartSpecifierText},
		},
	}).Collect()
	slices.Reverse(mails)
	return mails, err
}
