package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("\nPokedex > ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		parts := cleanInput(line)
		fmt.Printf("Your command was: %s", parts[0])
	}


	// fmt.Printf("%q",cleanInput("   hello  world "))
}
