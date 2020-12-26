package parse

import (
	"strings"

	"main/engine"
)

func Parse(str string) engine.Command {
	array := strings.Fields(str)
	if len(array) == 0 {
		return &engine.PrintCmd{Message: "WARNING: Empty line was found"}
	}
	if array[0] == "delete" {
		if len(array) == 3 {
			byteArray := []byte(array[2])
			if len(byteArray) == 1 {
				return &engine.DeleteCmd{Str: array[1], Symbol: byteArray[0]}
			}
			return &engine.PrintCmd{Message: "SYNTAX ERROR: Not a single symbol"}
		}
		return &engine.PrintCmd{Message: "SYNTAX ERROR: Wrong amount of arguments"}
	} else if array[0] == "print" {
		if len(array) == 2 {
			return &engine.PrintCmd{Message: array[1]}
		}
		return &engine.PrintCmd{Message: "SYNTAX ERROR: Wrong amount of arguments"}
	}
	return &engine.PrintCmd{Message: "SYNTAX ERROR: Wrong command name"}
}
