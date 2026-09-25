package main

import (
	"fmt"
	"log"

	"github.com/devnulljason/archipelago"
)

func main() {
	port := 55965

	client, err := archipelago.NewClient(port)
	if err != nil {
		fmt.Printf("error creating client: %s", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			fmt.Printf("error closing client: %s", err)
		}
	}()

	d := archipelago.GetDataPackage{
		Games: []string{"Archipelago"},
	}
	if err = client.Send(d); err != nil {
		fmt.Printf("error sending: %s", err)
	}

	packets, err := client.Receive()
	if err != nil {
		fmt.Println(err)
	}

	for _, p := range packets {
		log.Printf("%+v", p)
	}
}
