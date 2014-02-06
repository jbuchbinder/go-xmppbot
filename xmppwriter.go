package xmppbot

import (
	"github.com/mattn/go-xmpp"
)

// Writer type which wraps the ability to push traffic to an xmpp.Client
// instance.
type XmppWriter struct {
	Client *xmpp.Client
	User   string
}

// Instantiate Writer from XmppBot object for a specified XMPP user.
func (self *XmppBot) GetWriter(user string) *XmppWriter {
	return &XmppWriter{Client: self.client, User: user}
}

func (self *XmppWriter) Write(p []byte) (n int, err error) {
	self.Client.Send(xmpp.Chat{Remote: self.User, Type: "chat", Text: string(p)})
	return len(p), nil
}
