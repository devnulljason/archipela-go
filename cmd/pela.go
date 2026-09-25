// pela is a command-line program to interact with an Archipelago server.
package main

import (
	"fmt"
	"log"

	"github.com/devnulljason/archipelago"
)

func main() {
	port := 55965

	client := archipelago.NewClient(port)
	defer func() {
		if err := client.Close(); err != nil {
			fmt.Printf("error closing client: %s", err)
		}
	}()

	if err := client.Connect(); err != nil {
		log.Fatalf("websocket connection error: %s", err)
	}

	d := archipelago.GetDataPackage{
		Games: []string{"Archipelago"},
	}
	if err := client.Send(d); err != nil {
		fmt.Printf("error sending: %s", err)
	}

	packets, err := client.ReadPackets()
	if err != nil {
		fmt.Println(err)
	}

	for _, p := range packets {
		log.Printf("%+v", p)
	}
}
