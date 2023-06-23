package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"github.com/gliderlabs/ssh"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/mustache/v2"
	gossh "golang.org/x/xrypto/ssh"
)

type Pipe struct {
	wch    chan io.Writer
	donech chan struct{}
}

func newPipe() Pipe {
	return Pipe{
		wch:    make(chan io.Writer),
		donech: make(chan struct{}),
	}
}

var pipes = map[string]Pipe{}

func main() {
	engine := mustache.New("www", ".html")
	engine.Reload(true)
	app := fiber.New(fiber.Config{Views: engine})
	app.Get("/", handleHome)
	app.Get("/:link", handleLink)
	app.Get("/d/:link", handleDownload)

	go func() {
		fmt.Println("HTTP server running")
		app.Listen(":3000")
	}()

	b, err := ioutil.ReadFile("privatekey")
	if err != nil {
		log.Fatal(err)
	}
	privatekey, err := gossh.ParsePrivateKey(b)
	if err != nil {
		log.Fatal("Failed to parse private key: ", err)
	}

	sshPort := getenv("SENDIT_SSH_PORT", ":2222")

	server := ssh.Server{
		Addr:    sshPort,
		Handler: handleSshSession,
		PublicKeyHandler: func(ctx ssh.Context, key ssh.PublicKey) bool {
			fmt.Println("--->", key)
			return true
		},
		ServerConfigCallback: func(c ssh.Context) *gossh.ServerConfig {
			conf := &gossh.ServerConfig{}
			conf.AddHostKey(privatekey)
			return conf
		},
	}

	log.Fatal(server.ListenAndServe())
}
