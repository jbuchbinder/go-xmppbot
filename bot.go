package xmppbot

import (
	"log"
	"os"
	"strings"
	"sync"

	"github.com/xmppo/go-xmpp"
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

// CreateXmppBot is the factory method to create a new XMPP chat bot and
// connect to the upstream XMPP server.
func CreateXMPPBot(parameters *XmppBotParams) (*XmppBot, error) {
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

// AddCommand adds additional capabilties to the bot. If this is never
// called, only "version" and "help" commands will be available to end
// users.
func (x *XmppBot) AddCommand(c *XmppBotCommand) {
	x.commands[c.Name] = c
}

// Run launches the processing loop for the bot. If this is not called,
// the bot will connect to the XMPP server, but will not process any input.
func (x *XmppBot) Run() {
	w := sync.WaitGroup{}
	w.Add(1)
	go func() {
		for {
			chat, err := x.client.Recv()
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			switch v := chat.(type) {
			case xmpp.Chat:
				if x.params.Debug {
					log.Print("RECV[" + v.Remote + "]: " + v.Text)
				}
				if strings.TrimSpace(v.Text) != "" {
					x.dealWithCmd(v.Remote, v.Text)
				}
			case xmpp.Presence:
				if x.params.Debug {
					log.Print("PRES[" + v.From + "]: " + v.Show)
				}
			}
		}
	}()

	// Easier than endless for... loop
	w.Wait()
}

// dealWithCmd is an internal routine for processing commands from a user to the bot.
func (x *XmppBot) dealWithCmd(user string, raw string) {
	switch strings.TrimSpace(strings.ToLower(raw)) {
	case "help":
		// Iterate through commands, plus help and version
		helptxt := "\n"
		for k := range x.commands {
			helptxt += x.commands[k].Name + " : " + x.commands[k].HelpText + "\n"
		}
		helptxt += "help : This help text\n" +
			"version : Display the ops bot version\n"
		x.SendClient(user, helptxt)
		break
	case "version":
		x.SendClient(user, "Version: "+Version)
		break
	default:
		for k := range x.commands {
			if strings.HasPrefix(strings.TrimSpace(strings.ToLower(raw)), strings.ToLower(x.commands[k].Name)) {
				x.commands[k].Command(
					x,
					user,
					strings.Split(
						strings.TrimPrefix(
							strings.TrimSpace(strings.ToLower(raw)),
							strings.ToLower(x.commands[k].Name)),
						" "),
				)
				break
			}
		}
		x.SendClient(user, "UNKNOWN COMMAND: '"+raw+"'\nTry 'help' for a list of valid commands.")
		break
	}
	return
}

// SendClient is a convenience routine to send a response to a specified user.
func (self *XmppBot) SendClient(user string, msg string) {
	log.Print("[" + user + "] : " + msg)
	self.client.Send(xmpp.Chat{Remote: user, Type: "chat", Text: msg})
}
