package main

import (
	"flag"
	"fmt"
	bot "github.com/jbuchbinder/go-xmppbot"
)

var (
	username = flag.String("user", "", "GTalk username")
	password = flag.String("pass", "", "GTalk password")
)

func main() {
	b, err := bot.CreateXmppBot(&bot.XmppBotParams{
		Server:   "talk.google.com:443",
		Username: *username,
		Password: *password,
		UseTls:   true,
	})
	if err != nil {
		panic(err)
	}
	b.AddCommand(&bot.XmppBotCommand{
		Name:     "test",
		HelpText: "Bot test routine",
		Command: func(xmppbot *bot.XmppBot, user string, args []string) error {
			xmppbot.SendClient(user, "This is the result of a test\n")
			return nil
		},
	})
	go b.Run()
	fmt.Println("*** Hit ENTER to terminate ***")
	var x string
	_, _ = fmt.Scanf(x)
}
