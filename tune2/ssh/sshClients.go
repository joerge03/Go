package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

var (
	ip, port,username,password string
)

func init(){
	flag.StringVar(&ip, "ip", "127.0.0.1", "ip")
	flag.StringVar(&port, "port", "127.0.0.1", "port")
	flag.StringVar(&username, "username", "", "username")
	flag.StringVar(&password, "password", "", "password")
}



func publicKeyAuth(path string) ssh.AuthMethod {
	key, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("there's something wrong reading the file path %v\n", err)
	}

	signedPrivateKey, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("there's something wrong with key %v\n", err)
	}

	return ssh.PublicKeys(signedPrivateKey)
}

func hostKeyCallback(path string) ssh.HostKeyCallback{
	file, err := knownhosts.New(path)
	if err != nil {
		log.Fatalf("There's something wrong to known hosts %v\n", err)
	}

	return file
}


func main(){
	flag.Parse()

	os.Exit(1)

	addr := net.JoinHostPort(ip, port)

	config := ssh.ClientConfig{
		User: "proxy",
		Auth: []ssh.AuthMethod{
			publicKeyAuth("~/Desktop/privateSecret"),
		},
		// HostKeyCallback:hostKeyCallback("~/.ssh/known_hosts"),
		HostKeyCallback:ssh.InsecureIgnoreHostKey(),
	}

	if len(username) < 1 {
		fmt.Println("-username is empty or unused")
		os.Exit(2)
	}

	client, err := ssh.Dial("tcp",addr,&config)

	if err != nil {
		log.Fatalf("There's something wrong with connection %v\n", err)
	}

	defer client.Close()

	

	session, err := client.NewSession()

	if err != nil {
		log.Fatalf("There's something wrong with client session %v\n", err)
	}

	// session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	input, err := session.StdinPipe()


	if err != nil {
		log.Fatalf("there's something wrong with the stdin pipe %v\n", err)
	}

	termModes := ssh.TerminalModes{
		ssh.ECHO:0 ,
	}


	err = session.RequestPty("vt220", 40,80,termModes)

	if err != nil {
		log.Fatalf("can't proceed due to error, %v\n", err)
	}


	err = session.Shell()
	if err != nil {
		log.Fatalf("there' something wrong with the session shell %v\n", err)
	}

	for {
		io.Copy(input,os.Stdin)
	}

		
}