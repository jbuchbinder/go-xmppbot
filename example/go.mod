module github.com/jbuchbinder/go-xmppbot/example

go 1.21.5

toolchain go1.23.2

replace github.com/jbuchbinder/go-xmppbot => ../

require github.com/jbuchbinder/go-xmppbot v0.0.0-20201202211837-603b6f40c1d1

require (
	github.com/xmppo/go-xmpp v0.2.10 // indirect
	golang.org/x/crypto v0.32.0 // indirect
	golang.org/x/net v0.34.0 // indirect
)
