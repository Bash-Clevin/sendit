package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gliderlabs/ssh"
	"github.com/sirupsen/logrus"
)

const maxFileSize = 1 << 30

func handleSshSession(s ssh.Session) {
	start := time.Now()
	if len(s.Command() == 0) {
		return
	}

	userInput := s.Command()

	var (
		targetEmail string
		isDirectory bool
	)

	if len(userInput) == 3 {
		isDirectory = false
		targetEmail = userInput[2]
	}

	if len(userInput) == 4 {
		if userInput[1] != "-r" {
			logrus.Errorf("Invalid user input: %s", s.RawCommand())
			return
		}
		isDirectory = true
		targetEmail = userInput[3]
	}

	if validateEmail(targetEmail) {
		logrus.Errorf("Invaidemail:%s", targetEmail)
		return
	}

	if isDirectory {
		handleFolderTransfers(s)
		return
	}

	s.Write([]byte{0x00})
	reader := bufio.NewReader(s)
	header, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 3 {
		logrus.Errorf("Invalid header")
		return
	}

	filesizeStr := headerParts[1]
	filesize, _ := strconv.Atoi(filesizeStr)
	if filesize > maxFileSize {
		logrus.Errorf("max filesize exceeded: %d > %d", filesize, maxFileSize)
		return
	}
	s.Write([]byte{0x00})
	link := generateRandomString(12)
	fmt.Println("---", link)
	pipes[link] = newPipe()
	w := <-pipes[link].wch
	_, err = io.CopyN(w, s, int64(filesize))
	if err != nil {
		logrus.Errorf("io copy error: %s", err)
	}
	close(pipes[link].donech)
	delete(pipes, link)
	logrus.WithFields(logrus.Fields{
		"to":         targetEmail,
		"took":       time.Since(start),
		"user":       s.User(),
		"filesize":   filesize,
		"remoteAddr": s.RemoteAddr().String(),
	}).Info("transfer complete")
}

func handleFolderTransfers(s ssh.Session) {
	fmt.Println("Handling a folder")
	s.Write([]byte{0x00})

	reader := bufio.NewReader(s)
	raw, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	header, err := parseFileHeader(raw)
	if err != nil {
		logrus.Error(err)
		return
	}

	fmt.Println(header)
	s.Write([]byte{0x00})
	header, err = parseFileHeader(raw)

	if err != nil {
		log.Fatal(err)
	}

	// buf := make([]byte, 1000)
	// s.Read(buf)
	// fmt.Println(string(buf))
	// s.Write([]byte{0x00})
	// s.Read(buf)
	// fmt.Println(string(buf))
}
