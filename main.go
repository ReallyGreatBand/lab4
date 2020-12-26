package main

import (
	"bufio"
	"flag"
	"os"

	"main/engine"

	"main/parse"
)

var (
	fStringPtr = flag.String("f", "", "enter the filename with commands")
)

func main() {
	eventLoop := new(engine.EventLoop)
	eventLoop.Start()
	flag.Parse()
	fString := *fStringPtr

	if fString == "" {
		cmd := &engine.PrintCmd{Message: "ERROR: Enter the filename with commands via -f key"}
		eventLoop.Post(cmd)
	} else if input, err := os.Open(fString); err == nil {
		defer input.Close()
		scanner := bufio.NewScanner(input)
		isEmpty := true
		for scanner.Scan() {
			isEmpty = false
			commandLine := scanner.Text()
			cmd := parse.Parse(commandLine)
			eventLoop.Post(cmd)
		}
		if isEmpty {
			cmd := &engine.PrintCmd{Message: "WARNING: File is empty"}
			eventLoop.Post(cmd)
		}
	} else {
		cmd := &engine.PrintCmd{Message: "SYNTAX ERROR: Cannot open file or it does not exist"}
		eventLoop.Post(cmd)
	}
	eventLoop.AwaitFinish()
}
