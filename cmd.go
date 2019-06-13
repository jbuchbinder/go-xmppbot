package xmppbot

import (
	"io"
	"log"
	"os/exec"
)

// System "run command and report" functions

// RunCmd is a simplistic method of running a local command with arguments.
// The cmd argument must be a fully qualified path to the destination
// executable file, otherwise it will not function.
func (x *XmppBot) RunCmd(user string, cmd string, argv []string) {
	log.Print("For user '" + user + "' executing command : " + cmd)

	proc := exec.Command(cmd, argv...)
	go func() {
		stdout, err := proc.StdoutPipe()
		if err != nil {
			x.SendClient(user, cmd+": Failed to create stdout pipe")
			return
		}
		stderr, err := proc.StderrPipe()
		if err != nil {
			x.SendClient(user, cmd+": Failed to create stderr pipe")
			return
		}

		err = proc.Start()
		if err != nil {
			x.SendClient(user, cmd+": Failed to start process : "+err.Error())
			return
		}
		defer proc.Wait()

		// Hack to write everything back to the user
		w := &XmppWriter{Client: x.client, User: user}
		go io.Copy(w, stdout)
		go io.Copy(w, stderr)
	}()
}
