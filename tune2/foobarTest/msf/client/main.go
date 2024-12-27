package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type SessionListRes struct {
	ID uint32 `json:",omitempty"`
	Type string `json:"type"`
	TunnelLocal string `json:"tunnel_local"`
	TunnelPeer string `json:"tunnel_peer"`
	ViaPayload string `json:"via_payload"`
	Description string `json:"desc"`
	Info string `json:"info"`
	Workspace string `json:"workspace"`
	TargetHost string `json:"target_host"`
	Username string `json:"username"`
	UUID string `json:"uuid"`
	ExploitUUID string `json:"exploit_uuid"`
	Routes any `json:"routes,omitempty"`	
}

type sessionListReq struct{
	// _json struct{} `json:",asArray"`
	Method string 
	Token string
}

type loginRes struct {	
	Result string `json:"result"`
	Token string `json:"token"`
	IsError bool `json:"error"`
	ErrorClass string `json:"error_class"`
	ErrorMessage string `json:"error_message"`
}
type loginReq struct {
	// _json struct{} `json:",asArray"`
	Method string `json:"method"`
	Username string `json:"username"`
	Password string `json:"password"`
}


type logoutRes struct {
	Result string `json:"result"`
}

type logoutReq struct{
	// _json struct{} `json:",asArray"`
	Method string 
	Token string 
	LogoutToken string 
}

type Metasploit struct{
	host,
	user,
	pass,
	token string
}

func New(host, user, pass string) (*Metasploit, error){
	m := &Metasploit{
		host: host,
		user: user,
		pass: pass,
	}
	if err := m.Login(); err != nil {
		return nil, err
	}
	return m, nil 
}


func (m *Metasploit) Login() error {
	req := &loginReq{
		Method: "host.login",
		Username: m.user,
		Password: m.pass,
	}
	
	var res loginRes
	
	err := m.Send(req,res)
	if err != nil {
		return err
	}
	fmt.Println(res, "login")
	
	m.token = res.Token
	return nil
}

func (m *Metasploit) Logout() error {	
	
	req := &logoutReq{
		Method: "auth.logout",
		Token: m.token,
		LogoutToken: m.token,
	}
	var res logoutRes
	
	if err := m.Send(req, res); err != nil {
		return err 
	}
	fmt.Println(res, "logout")
	
	fmt.Printf("Logged out, %+v", res.Result)
	return nil
}
func (m *Metasploit) SessionList() (map[uint32]SessionListRes,error){
	sessionReq := &sessionListReq{
		Method: "session.list",
		Token: m.token,
	}
	sessionRes := make(map[uint32]SessionListRes)
	
	if err := m.Send( sessionReq,sessionRes); err != nil {
		return nil,err
	}
	
	fmt.Println(sessionRes)
	
	for id, session := range sessionRes {
		fmt.Println(session, "session")
		session.ID = id
		sessionRes[id] = session
	}
	return sessionRes,nil
}
func (m *Metasploit) Send(req, res any) error{
	buf := new(bytes.Buffer)

	json.NewEncoder(buf).Encode(req)

	// certData, err := os.ReadFile("./rpcCert.pem")
	certData, err:=os.ReadFile("./localhost.pem")
	if err != nil {
		return err
	}
	// cert, _ := tls.LoadX509KeyPair("metasploit.crt", "metasploit.key")
	//       if err != nil {
	//               log.Fatal(err, "err")
	//       }
	
	certPool := x509.NewCertPool()
	// fmt.Println(cert)
	
	if ok := certPool.AppendCertsFromPEM(certData); !ok {
		log.Fatal("Failed to append cert")
    }
	tlsConfig := &tls.Config{
		// ClientAuth: tls.RequireAndVerifyClientCert,
		RootCAs: certPool,
		// Certificates: []tls.Certificate{certPool},	
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,		
		},
	}

	des := fmt.Sprintf("http://%s/api/v1/json-rpc",m.host)

	r, err := client.Post(des,"binary/message-pack", buf)

	// data := make([]byte, 2049)

	// b,err :=r.Body.Read(data)
	
	// // test,err:=bufio.NewReader(r.Body).Read(data)
	// if err!= nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(data[:b])

	

	
	
	if err != nil{
		fmt.Println("error post nga")
		return err
	}

	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(res)
	if err != nil {
		return err
	}
	return nil
}

func main(){

	
	host := "localhost:8081"
	pass:= "123"
	// test, err := http.Get(fmt.Sprintf("http://%v/api",host)) 
	// if err != nil {
	// 	log.Fatal("err", err)
	// }
	
	
	
	
	buf := make([]byte, 2048)
	
	// test.Body.Read(buf)
	fmt.Printf("%+v", string(buf))
	user := "msf"
	if len(host) == 0 || len(pass) == 0 {
		log.Panic("empty host/pass")	
	}
	
	m, err := New(host,user,pass)
	if err != nil {
		log.Panic(err, "new")
	}
	
	defer m.Logout()
	
	sessions, err := m.SessionList()
	
	if err != nil {
		log.Panic(err, "sess")
	}

	
	for id, sessions := range sessions{
		fmt.Printf("id %v sessions %+v\n", id, sessions)
	}


	
}