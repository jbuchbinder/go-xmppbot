package xmppbot

import (
	"github.com/mattn/go-xmpp"
)

type XmppWriter struct {
	Client *xmpp.Client
	User   string
}

func (self *XmppWriter) Write(p []byte) (n int, err error) {
	self.Client.Send(xmpp.Chat{Remote: self.User, Type: "chat", Text: string(p)})
	return len(p), nil
}
