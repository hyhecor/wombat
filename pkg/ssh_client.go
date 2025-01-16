package pkg

import (
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

var terminalModes = ssh.TerminalModes{
	ssh.ECHO:          1, // enable echoing (different from the example in docs)
	ssh.ECHOCTL:       1,
	ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
	ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
}

func NewShell() {

	config := ssh.ClientConfig{
		User: "h",
		Auth: []ssh.AuthMethod{
			ssh.PasswordCallback(func() (secret string, err error) {
				print("password: ")

				fileDescriptor := int(os.Stdin.Fd())
				b, err := term.ReadPassword(fileDescriptor)

				println("")
				return string(b), err
			}),
			ssh.Password("h"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", "localhost:22", &config)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	fileDescriptor := int(os.Stdin.Fd())

	if term.IsTerminal(fileDescriptor) {

		envTerm := os.Getenv("TERM")

		originalState, err := term.MakeRaw(fileDescriptor)
		if err != nil {
			panic(err)
		}
		defer term.Restore(fileDescriptor, originalState)

		termWidth, termHeight, err := term.GetSize(fileDescriptor)
		if err != nil {
			panic(err)
		}

		err = session.RequestPty(envTerm, termHeight, termWidth, terminalModes)
		if err != nil {
			panic(err)
		}
	}

	err = session.Shell()
	if err != nil {
		panic(err)
	}

	err = session.Wait()
	if err != nil {
		panic(err)
	}
}
