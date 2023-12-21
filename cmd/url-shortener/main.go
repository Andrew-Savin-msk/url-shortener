package main

import (
	"fmt"
	"url-shorter/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)

	// TODO: init logger: slog

	// TODO: init storage: sqlite

	// TODO: init router: chi, chi render

	// TODO: init server:
}
