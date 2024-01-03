package main

import (
	"os"
	"os/exec"

	"github.com/labstack/echo/v4"
	"github.com/shynome/err0/try"
	"github.com/shynome/fuck-qq/onebot"
)

func main() {
	if os.Getenv("RUN_COPYQ") != "" {
		go func() {
			cmd := exec.Command("copyq")
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout
			try.To(cmd.Start())
			try.To(cmd.Wait())
		}()
	}
	e := echo.New()
	onebot.Inject(e)
	e.Start(":5700")
}
