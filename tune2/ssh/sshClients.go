package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

var (
	ip, port,username,password string
)

var pubKey ssh.PublicKey

func init(){
	flag.StringVar(&ip, "ip", "linuxzoo.net", "ip")
	flag.StringVar(&port, "port", "22", "port")
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
	// hostkey :=  ssh.FixedHostKey()

	return file
}


func main(){
	flag.Parse()

	// os.Exit(1)

	addr := net.JoinHostPort(ip, port)

	// fmt.Println(addr, "port")

	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.Password("secure"),
			publicKeyAuth("./../../../../privateSecret"),
		},
		// HostKeyCallback:hostKeyCallback("~/.ssh/known_hosts"),
		HostKeyCallback:ssh.InsecureIgnoreHostKey(),
		// HostKeyCallback: ssh.FixedHostKey(pubKey),
	}

	// if len(username) < 1 {
	// 	fmt.Println("-username is empty or unused")
	// 	os.Exit(2)
	// }

	client, err := ssh.Dial("tcp",addr,config)
	if err != nil {
		log.Fatalf("There's something wrong with connection, %v\n", err)
	}
	defer client.Close()

	

	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("There's something wrong with client session %v\n", err)
	}

	//////////////////////
	reader, writer, _:= os.Pipe()

	// writer = os.Stdout
	// oldStd := 
	
	// session.Stdout = os.Stdout

	out, err  := session.StdoutPipe()
	if err != nil {
		log.Fatalf("asdf")
	}
		
	 outStdout := make([]byte , 2048)  
	 go func(){
		outScan := bufio.NewScanner(out)
		// outStdout = []byte{}

		for{
			 if outScan.Scan() {
				
				 fmt.Printf("out %s\n", outStdout)		
				 os.Stdout.Write([]byte(fmt.Sprintf("%v\n",  outScan.Text())))
				//  fmt.Printf("teststset, %v\n",)
			 }
		}
	}()

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
	// readsTo := bufio.NewReader(os.Stdin)
	
	// writer := bufio.NewWriter(input)

	
	writer.Write([]byte("test"))

	
	
	scanner := bufio.NewScanner(os.Stdin)
	outScanner := bufio.NewScanner(reader)
	for {
		fmt.Printf(" writer %v \n ", outScanner.Text())
		// fmt.Printf("test %v \n ", outStdout)
		// fmt.Printf("\ntestststst %v \n ", outScanner.Text())
		// fmt.Printf("%s <--- \n", buff.String())

		
		
		if scanner.Scan(){
			
			// s,err := os.Stdout.Write(str)
			
			// fmt.Printf("\n%v\n", outScanner.Text())
			if err != nil {
				log.Fatalf("there's something wrong with the stdout, %v\n", err)
			}
			input.Write([]byte(fmt.Sprintf("%v\n", scanner.Text())))		
			// outStdout= []byte("")
		}
	}
	
	// for {
	// 	fmt.Println("success")
	// 	io.Copy(input, os.Stdin)
	// }

}