# Embed help

This package supports taking embedded items in a go binary and dumping them to the os.
Perhaps you want to pack up a service file or dump the html and serve locally, this will facilitate that


## Install

```
go get github.com/rob121/embedhelp
```

## Usage

```
// Setup embed items
//go:embed example.service
var svcvar string
//go:embed example
var sys embed.FS
//go:embed test.txt
var test string

//Mark which items are to be dumped
embedhelp.Register("service",svcvar,"/var/lib/systemd/example.service",false)
embedhelp.Register("example",sys,"/tmp/embedtest",true)
embedhelp.Register("example2",test,"/tmp/test.txt",true)

// Dump All
embedhelp.DumpAll()
// Dump Specific item 
embedhelp.Dump("example2")
```

## Caveats
This package supports byte,embed.FS and string at the moment
