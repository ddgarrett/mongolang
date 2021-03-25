package mongolang

import (
	"fmt"
	_ "time" // test import of unused package
)

// Hello returns a greeting for the named person.
func Hello(name string) string {
	// Return a greeting that embeds the name in a message.
	message := fmt.Sprintf("Hi, %v. Welcome to version 0.2.2 of MonGolang!", name)
	return message
}
