package main

import (
	"github.com/labstack/echo/v4"
	"github.com/shynome/fuck-qq/onebot"
)

func main() {
	e := echo.New()
	onebot.Inject(e)
	e.Start(":5700")
}
