package xmppbot

import (
	"github.com/mattn/go-xmpp"
	"log"
	"os"
	"strings"
	"sync"
)

type XmppBotParams struct {
	Server   string
	Username string
	Password string
	UseTls   bool
	Debug    bool
}

type XmppBotCommand struct {
	Name     string
	HelpText string
	Command  func(xmppbot *XmppBot, user string, args []string) error
}

type XmppBot struct {
	client   *xmpp.Client
	params   *XmppBotParams
	commands map[string]*XmppBotCommand
}

func CreateXmppBot(parameters *XmppBotParams) (*XmppBot, error) {
	obj := &XmppBot{params: parameters, commands: make(map[string]*XmppBotCommand)}
	var err error
	log.Print("Initializing XMPP connection")
	if parameters.UseTls {
		obj.client, err = xmpp.NewClient(parameters.Server, parameters.Username, parameters.Password, parameters.Debug)
	} else {
		obj.client, err = xmpp.NewClientNoTLS(parameters.Server, parameters.Username, parameters.Password, parameters.Debug)
	}
	if err != nil {
		return nil, err
	}
	log.Print("Successfully connected")
	return obj, nil
}

func (self *XmppBot) AddCommand(c *XmppBotCommand) {
	self.commands[c.Name] = c
}

func (self *XmppBot) Run() {
	w := sync.WaitGroup{}
	w.Add(1)
	go func() {
		for {
			chat, err := self.client.Recv()
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			switch v := chat.(type) {
			case xmpp.Chat:
				if self.params.Debug {
					log.Print("RECV[" + v.Remote + "]: " + v.Text)
				}
				if strings.TrimSpace(v.Text) != "" {
					self.DealWithCmd(v.Remote, v.Text)
				}
			case xmpp.Presence:
				if self.params.Debug {
					log.Print("PRES[" + v.From + "]: " + v.Show)
				}
			}
		}
	}()

	// Easier than endless for... loop
	w.Wait()
}

func (self *XmppBot) DealWithCmd(user string, raw string) {
	switch strings.TrimSpace(strings.ToLower(raw)) {
	case "help":
		// Iterate through commands, plus help and version
		helptxt := "\n"
		for k, _ := range self.commands {
			helptxt += self.commands[k].Name + " : " + self.commands[k].HelpText + "\n"
		}
		helptxt += "help : This help text\n" +
			"version : Display the ops bot version\n"
		self.SendClient(user, helptxt)
		break
	case "version":
		self.SendClient(user, "Version: "+VERSION)
		break
	default:
		for k, _ := range self.commands {
			if strings.HasPrefix(strings.TrimSpace(strings.ToLower(raw)), strings.ToLower(self.commands[k].Name)) {
				self.commands[k].Command(self, user, strings.Split(strings.TrimPrefix(strings.TrimSpace(strings.ToLower(raw)), strings.ToLower(self.commands[k].Name)), " "))
				break
			}
		}
		self.SendClient(user, "UNKNOWN COMMAND: '"+raw+"'\nTry 'help' for a list of valid commands.")
		break
	}
	return
}

func (self *XmppBot) SendClient(user string, msg string) {
	log.Print("[" + user + "] : " + msg)
	self.client.Send(xmpp.Chat{Remote: user, Type: "chat", Text: msg})
}
