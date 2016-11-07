package main

import (
	"fmt"
	"bytes"
	"golang.org/x/crypto/ssh"
	"os"
	"unsafe"
	cfg "cssh/config"
)

func main() {
	conf := cfg.Init()

	config := &ssh.ClientConfig{
		User: conf.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(conf.Passwd),
		},
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", conf.Host, conf.Port), config)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}


	b := make([]byte, 100)
	f := os.Stdin
	w := os.Stdout
	defer f.Close()
	defer w.Close()
	for {
		w.WriteString("input:")
		c, _ := f.Read(b)
		bb := b[:c-1]
		str := *(*string)(unsafe.Pointer(&bb))


		session, err := client.NewSession()
		if err != nil {
			panic("Failed to create session: " + err.Error())
		}
		defer session.Close()

		var ob bytes.Buffer
		session.Stdout = &ob
		if err := session.Run(str); err != nil {
			//panic("Failed to run: " + err.Error())
			fmt.Println("Input error")
		}
		fmt.Println(ob.String())
		if str == "exit" {
			break
		}
	}
}
