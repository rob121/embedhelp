package main

import (
 "embed"
 _ "embed"
 "github.com/rob121/embedhelp"
 "log"
)

//go:embed testdir
var sys embed.FS


func main(){ 

    embedhelp.Register("example",sys,"/tmp/embedtest",false)


    err := embedhelp.DumpAll()

    if(err!=nil){

      log.Println(err)

    }

}
