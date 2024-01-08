package main

import (
	"flag"
	"log"
	"os"
	"os/exec"

	"github.com/labstack/echo/v4"
	"github.com/shynome/err0/try"
	"github.com/shynome/fuck-qq/onebot"
)

var args struct {
	addr string
}

func init() {
	flag.StringVar(&args.addr, "addr", ":5700", "http server listen addr")
}

func main() {
	flag.Parse()
	if os.Getenv("RUN_COPYQ") != "" {
		go func() {
			cmd := exec.Command("copyq")
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout
			try.To(cmd.Start())
			log.Println("copyq is running")
			try.To(cmd.Wait())
		}()
	}
	e := echo.New()
	onebot.Inject(e)
	try.To(e.Start(args.addr))
}
