package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

type SSHServer struct {
	Address   string
	Host      string
	Port      int
	IsSSH     bool
	Banner    string
	Cert      ssh.Certificate
	HostName  string
	PublicKey ssh.PublicKey
}

type strList []string

func (s *strList) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *strList) Set(str string) error {
	if str == "" {
		return fmt.Errorf("the arguments are empty")
	}

	argumentsData := strings.Split(str, ",")
	*s = append(*s, argumentsData...)
	return nil
}

var (
	in, out, address string
	logSSH           log.Logger
	targets          strList
	sshWG 			sync.WaitGroup
	password 			string
	username 			string
)

func errorExit(s string, err error) {
	if err != nil {
		logSSH.Panicf("%v, test %v \n", s, err)
	}
	log.Panicf("%v\n", s)
}

func init() {
	flag.StringVar(&in, "i", "", "-i for input file")
	flag.StringVar(&in, "input", "", "-input for input file")

	flag.StringVar(&out, "o", "", "-o for output file")
	flag.StringVar(&out, "out", "", "-out for output file")

	flag.StringVar(&username, "u", "root", "username")
	flag.StringVar(&username, "username", "root", "username")

	flag.StringVar(&password, "p", "secure", "password")
	flag.StringVar(&password, "password", "secure", "password")

	flag.StringVar(&address, "a", "", "-a address with a format of ip:port or [ipv6:4f::]:23")
	flag.StringVar(&address, "address", "-address", "-out address with a format of ip:port or [ipv6:4f::]:23")

	flag.Var(&targets, "t", "-t for targets separated with a comma (,)")
	flag.Var(&targets, "target", "-target for targets separated with a comma (,)")

	flag.Parse()

	logSSH = *log.New(os.Stdout, "[]", log.Ltime)

	if in == "" && len(targets) == 0 {
		errorExit("-i (input) and -t (targets) needs to be populated", nil)
	}
}

func newSSHServer(address string) (*SSHServer, error) {
	server := new(SSHServer)
	i, p, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}
	server.Address = address
	server.HostName = i
	server.Port, err = strconv.Atoi(p)
	if err != nil {
		return nil, err
	}

	if 0 > server.Port || server.Port > 65535 {
		return nil, errors.New(p + " Invalid port")
	}

	return server, nil
}

func ToJSON(s any, pretty bool )(string, error){
	 jsonString :=  new([]byte)
	 var err error
	if pretty {
		*jsonString, err = json.MarshalIndent(s, "", "\t")
		if err!= nil{
			return "", err
		}
		return string(*jsonString), nil
	}

	*jsonString, err = json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(*jsonString), nil
}

type SSHServers []*SSHServer

func (servers *SSHServers) String() string{
	serverString := new(string)
	var err error
	*serverString, err = ToJSON(*servers, true)
	if err != nil{		
		for _, server := range *servers {
			fmt.Println("updateing servers")
			*serverString += fmt.Sprintf("%+v\n%s\n", server, strings.Repeat("-",30))
		}
	}
	return *serverString	
}

type IsHostAuthorityCallback func(auth ssh.PublicKey, addr string) bool 

func IsHostAuthority() IsHostAuthorityCallback{
	
	return func (auth ssh.PublicKey, addr string) bool {		
			// s.Address = addr
			// s.PublicKey = auth 
		return true
	}
}

type IsRevokedCallback func(cert *ssh.Certificate) bool 

func IsRevokedFunc(s *SSHServer)IsRevokedCallback{
	return func(cert *ssh.Certificate) bool{
		(*s).Cert = *cert
		(*s).IsSSH = true
		return true
	}
}

type IsUserAuthorityCallback func(p ssh.PublicKey) bool

func IsUserAuthority(s *SSHServer) IsUserAuthorityCallback{
	return func(p ssh.PublicKey) bool {
		return true
	}
}

// type HostKeyCallbackType func(hostName string, remote net.Addr, key ssh.PublicKey) error 


func HostKeyCallback(s *SSHServer) ssh.HostKeyCallback {
	return func(hostName string, remote net.Addr, key ssh.PublicKey) error{

		host, port, _ := net.SplitHostPort(remote.String())
		portInt,_  := strconv.Atoi(port)
		
		(*s).Address = remote.String()
		(*s).Host = host
		(*s).Port = portInt
		(*s).PublicKey = key
		
		return nil
	}
} 

// type BannerCallBackType func(message string) error 


func BannerCallBack(s *SSHServer) ssh.BannerCallback{
	return func(message string) error {
		fmt.Println(message)
		s.Banner = message
		return nil
	}
}


func readFile(fileLoc string)([]string, error) {

	add := new([]string)

	file, err := os.Open(fileLoc)

	if err != nil {
		return nil, err
	}

	fileScanner := bufio.NewScanner(file)

	for fileScanner.Scan() {
		*add = append(*add, fileScanner.Text())
	}
	return *add, nil
	
}


func (s *SSHServer) discover(){
	defer sshWG.Done()

	certC := &ssh.CertChecker{
		IsHostAuthority: IsHostAuthority(),
		IsUserAuthority: IsUserAuthority(s),
		IsRevoked: IsRevokedFunc(s),
	}
	
	sshC := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: certC.CheckHostKey,
		BannerCallback: BannerCallBack(s),
		Timeout: time.Duration(time.Second *5),
	}

	c ,err := ssh.Dial("tcp",s.Address,sshC)

	if err != nil {
		logSSH.Println("error ", err)
		return
	}
	c.Close()	
}

func main3() {

	// certC := ssh.ClientConfig{
	// 	HostKeyCallback: ,
	// }

	allTargets := new([]string)
	servers :=  new(SSHServers)
	var err error 
	if in != ""{
		*allTargets, err = readFile(in)
		if err != nil {
			errorExit("test1", err)
		}
	}

	if len(targets) != 0 {
		*allTargets = append(*allTargets, targets...)		
	}

	for _,allT :=  range *allTargets {
		server, err:= newSSHServer(allT)
		if err != nil {
			errorExit("test", err)
		}
		*servers = append(*servers, server)		
	}

	if err != nil {
		errorExit("there's something wrong with the server", err)
	}

	for _ , server := range *servers {
		sshWG.Add(1)
		fmt.Println(server)
		go server.discover()
	}

	
	sshWG.Wait()
	fmt.Printf("servers %+v", servers)

	fmt.Println("123", username, password)
	
	if out != "" {
		file, err := os.OpenFile(out,os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

		if err != nil {
			log.Panicf("can't open file, %v", err)
		}
		defer file.Close()
		// writer := bufio.NewWriter(file)
		formattedServers := servers.String()
		
		_,err = file.WriteString(formattedServers)

		fmt.Printf("err write %v\n", err)
	}	
}