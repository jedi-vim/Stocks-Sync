package main

import (
	"log"
        "github.com/alecthomas/kong"
)

func main() {
    cli := CLI{}
    ctx := kong.Parse(&cli)
    err := ctx.Run(&kong.Context{})
    if err != nil{
        log.Fatal(err)
    }
}
