package main

import (
	proto "Examdisys/grpc"
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


type Client struct {
	id         int
	portNumber int
	proto.UnimplementedDictionaryServiceServer
	servers []Server
	amount  int32
}

type Server struct {
	server   proto.DictionaryServiceClient
	port     int32
}

var (
	clientPort = flag.Int("cPort", 0, "client port number")
)


func main() {

	flag.Parse()
	

	client := &Client{
		servers:    make([]Server, 0),
		portNumber: *clientPort,
	}

	go client.connectToServer(5001)
	go client.connectToServer(5002)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
	
		if input == "add" {
			scanner.Scan() //The next scan is the word to add to the dictionary
			toAddWord := scanner.Text()

			scanner.Scan() //The next scan is the definition to add to the dictionary
			toAddDef := scanner.Text()

			for i := 0; i < len(client.servers); i++ {
				reponse, err := client.servers[i].server.Add(context.Background(), &proto.WordDef{Word: toAddWord, Definition: toAddDef})

				if(err != nil){
					log.Println("The client found out that a server is crashed")
					client.servers = removeServer(client.servers,i)
					fmt.Println("Something went wrong in adding your word to the dictionary. Try again")
				} else {
					if(reponse.Response == true){
						fmt.Println("You added a new word to the dictionary")
					} else {
						//if we recieve a false response something went wrong when the leader was updating the replica(s)
						fmt.Println("Something went wrong in adding your word to the dictionary. Try again")
					}
				}
			}
		} else if (input == "read"){
			scanner.Scan() //The next scan is the word the client wants to read in the dictionary
			wordToRead := scanner.Text()

			var definition string
			for i := 0; i < len(client.servers); i++ {
				def, err := client.servers[i].server.Read(context.Background(), &proto.Word{Word: wordToRead})

				if(err != nil){
					log.Println("The client found out that a server is crashed")
					client.servers = removeServer(client.servers,i)
				} else {
					if(def.Definition == ""){
						fmt.Println("The word did not exist in the dictionary")
					} else {
						definition = def.Definition
					}
				}
			}
			fmt.Printf("The definition for the word was: %v\n", definition)
		}
	}
	for {
	}
}

func (c *Client) connectToServer(portNumber int32) {
	//dialing the server
	conn, err := grpc.Dial("localhost:"+strconv.Itoa(int(portNumber)), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Could not connect: %s\n", err)
	}

	log.Printf("Client connected to server at port: %v\n", portNumber)

	newServerToAdd := proto.NewDictionaryServiceClient(conn)
	
	c.servers = append(c.servers, Server{
		server:   newServerToAdd,
		port:     portNumber,
	})

	wait := make(chan bool)
	<-wait
}

func removeServer(s []Server, i int) []Server {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]

}

// go run client/client.go -cPort 4041