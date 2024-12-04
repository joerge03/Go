package main

import (
	"os"

	"github.com/joho/godotenv"
)
const SHODAN_URL = "https://api.shodan.io"


type Client struct {
	API_KEY string
}



func init(){
	godotenv.Load("./../../.env")
}


func main(){
	// client := &http.Client{}
	// fmt.Println()

	c := NewClient(os.Getenv("SHODAN_API_KEY"))

	c.APIInfo()
	
	// s,err := shodan.GetClient(API_KEY_SHODAN,client,false)

	// if err != nil {
	// 	log.Panic(err)
	// }

	// s.
	
}