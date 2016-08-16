package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/brimstone/go-pushbullet"
)

func main() {
	myID := os.Getenv("ID")
	if myID == "" {
		fmt.Println("ID must be specified in environment")
		os.Exit(1)
	}
	key := os.Getenv("KEY")
	if key == "" {
		fmt.Println("KEY must be specified in environment")
		os.Exit(1)
	}
	client := pushbullet.ClientWithKey(key, myID)
	err := client.ListenForPushes(time.Unix(0, 0), func(push pushbullet.PushMessage) {
		var env []string
		env = append(env, "ID="+push.ID)
		if push.Type == "link" {
			log.Println(push.URL)
			env = append(env, "TYPE=url")
			env = append(env, "URL="+push.URL)
		} else if push.Type == "note" {
			env = append(env, "TYPE=note")
			env = append(env, "TITLE="+push.Title)
		} else {
			log.Println("I can't handle type", push.Type)
		}

		cmd := exec.Command(os.Args[1])
		cmd.Env = env

		// Pass through output
		output, _ := cmd.StdoutPipe()
		go io.Copy(os.Stdout, output)
		// Pass through error
		stderr, _ := cmd.StderrPipe()
		go io.Copy(os.Stderr, stderr)

		// handle cmd.StdinPipe
		if push.Body != "" {
			input, _ := cmd.StdinPipe()
			go func() {
				input.Write([]byte(push.Body))
				input.Close()
			}()
		}

		err := cmd.Run()
		if err != nil {
			log.Println("Error running command:", err)
			return
		}
		client.DeletePush(push.ID)
	})

	log.Println("Error:", err)
}
