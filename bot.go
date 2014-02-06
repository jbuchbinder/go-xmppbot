package xmppbot

import (
	"github.com/mattn/go-xmpp"
	"log"
	"os"
	"strings"
	"sync"
)

// XmppBot initialization parameters.
type XmppBotParams struct {
	Server   string
	Username string
	Password string
	UseTls   bool
	Debug    bool
}

// Container, instantiated to allow custom actions for the XMPP bot.
type XmppBotCommand struct {
	Name     string
	HelpText string
	Command  func(xmppbot *XmppBot, user string, args []string) error
}

// Main object instance for the XMPP chat bot.
type XmppBot struct {
	client   *xmpp.Client
	params   *XmppBotParams
	commands map[string]*XmppBotCommand
}

// CreateXmppBot() is the factory method to create a new XMPP chat bot and
// connect to the upstream XMPP server.
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

// AddCommand() adds additional capabilties to the bot. If this is never
// called, only "version" and "help" commands will be available to end
// users.
func (self *XmppBot) AddCommand(c *XmppBotCommand) {
	self.commands[c.Name] = c
}

// Run() launches the processing loop for the bot. If this is not called,
// the bot will connect to the XMPP server, but will not process any input.
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
					self.dealWithCmd(v.Remote, v.Text)
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

// Internal routine for processing commands from a user to the bot.
func (self *XmppBot) dealWithCmd(user string, raw string) {
	switch strings.TrimSpace(strings.ToLower(raw)) {
	case "help":
		// Iterate through commands, plus help and version
		helptxt := "\n"
		for k := range self.commands {
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
		for k := range self.commands {
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

// Convenience routine to send a response to a specified user.
func (self *XmppBot) SendClient(user string, msg string) {
	log.Print("[" + user + "] : " + msg)
	self.client.Send(xmpp.Chat{Remote: user, Type: "chat", Text: msg})
}
