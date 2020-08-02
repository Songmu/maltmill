package maltmill

import (
	"context"
	"flag"
	"io"
	"log"
)

const (
	exitCodeOK = iota
	exitCodeErr
)

// Run the maltmill
func Run(ctx context.Context, args []string, outStream, errStream io.Writer) int {
	err := (&cli{outStream: outStream, errStream: errStream}).run(ctx, args)
	if err != nil {
		if err == flag.ErrHelp || err == errOpt {
			return exitCodeOK
		}
		log.Printf("[!!ERROR!!] %s\n", err)
		return exitCodeErr
	}
	return exitCodeOK
}
