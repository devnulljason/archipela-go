package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/devnulljason/archipelago"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	host := "archipelago.gg"
	port := 55965

	if err := archipelago.Connect(host, port); err != nil {
		fmt.Printf("error: %s", err)
	}
}
