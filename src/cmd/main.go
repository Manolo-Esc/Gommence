package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Manolo-Esc/gommence/src/internal/server"
)

func main() {
	ctx := context.Background()
	if err := server.Run(ctx, os.Args, nil, os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
