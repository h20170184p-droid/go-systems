package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Let's create a map of key value pairs
	kv := make(map[string]([2]int))
	fmt.Println(len(kv))
	var count int = 0

RESET1:
	// Take input from the user, a string
	fmt.Printf("\nEnter a string: ")
	ur := bufio.NewReader(os.Stdin)
	ui, e1 := ur.ReadString('\n')
	if e1 != nil {
		fmt.Println("Error reading the input", e1)
		goto RESET1
	}
	ui = strings.TrimSpace(ui)

	if ui == "!Exit!" {
		os.Exit(0)
	}

	if strings.HasPrefix(ui, "!Del! ") {
		// Remove the particular key from the map
		key := strings.TrimPrefix(ui, "!Del! ")
		a := kv[key]
		fmt.Println(a)
		delete(kv, key)
		// Restructure the values in the map
		for key, val := range kv {
			if val[0] == a[1] {
				kv[key] = [2]int{(a[1] - 1), a[1]}
				a[0], a[1] = (a[0] + 1), (a[1] + 1)
			}
		}
		fmt.Println(kv)
		count -= 1
		goto RESET1
	}

	if len(kv) == 0 {
		kv[ui] = [2]int{0, 1}
		fmt.Println(kv)
		count += 1
		goto RESET1
	} else {
		for _, val := range kv {
			if val[1] == count {
				kv[ui] = [2]int{val[1], (val[1] + 1)}
				fmt.Println(kv)
				count += 1
				goto RESET1
			}
		}
	}

}
