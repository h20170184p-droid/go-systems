package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	// This is the unscrambler
BEGINNING:
	fmt.Println("Enter the name of scrambled file:")
	i := bufio.NewReader(os.Stdin)
	scramFile, e1 := i.ReadString('\n')
	if e1 != nil {
		fmt.Println("Error: ", e1)
		goto BEGINNING
	}
	scramFile = strings.TrimSpace(scramFile)
	scramContents, e2 := os.ReadFile(scramFile)
	if e2 != nil {
		fmt.Println("Error: ", e2)
		goto BEGINNING
	}

	// .......................................

	fmt.Println("Enter the name of key file:")
	j := bufio.NewReader(os.Stdin)
	keyFile, e3 := j.ReadString('\n')
	if e3 != nil {
		fmt.Println("Error: ", e3)
		goto BEGINNING
	}
	keyFile = strings.TrimSpace(keyFile)
	keyContents, e4 := os.ReadFile(keyFile)
	if e4 != nil {
		fmt.Println("Error: ", e4)
		goto BEGINNING
	}

	for ini := 0; ini < 8; ini++ {
		if scramContents[ini] != keyContents[ini] {
			fmt.Println("Wrong pair!")
			goto BEGINNING
		}
	}

	fmt.Println("The scram file can be decrypted with the key file!")
	n := int(keyContents[8])<<24 | int(keyContents[9])<<16 | int(keyContents[10])<<8 | int(keyContents[11])
	fmt.Printf("%d data bytes can be unscrambled\n", n)

	scramData := scramContents[8:]
	pointData := keyContents[12:]

	// fmt.Println(len(scramData))
	total := len(scramData)
	unscrambledOutput := make([]byte, total)

	var wg sync.WaitGroup
	cap := make(chan struct{}, 2000)
	for k := range total {
		wg.Add(1)
		cap <- struct{}{}
		go func(k int) {
			defer wg.Done()
			defer func() { <-cap }()
			position := (int(pointData[k*4]) << 24) | (int(pointData[((k*4)+1)]) << 16) | (int(pointData[((k*4)+2)]) << 8) | int(pointData[((k*4)+3)])
			unscrambledOutput[position] = scramData[k]
		}(k)
	}

	wg.Wait()
	file_n := filepath.Base(scramFile)
	outputFileName := strings.TrimSuffix("unscrambled"+file_n, "scrambled.txt")
	os.WriteFile(outputFileName, unscrambledOutput[:n], 0644)
}
