package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		// 1. Read Scrambled File
		scramFile, err := promptInput(reader, "Enter the name of scrambled file: ")
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		scramContents, err := os.ReadFile(scramFile)
		if err != nil {
			fmt.Println("Error reading scrambled file:", err)
			continue
		}

		// 2. Read Key File
		keyFile, err := promptInput(reader, "Enter the name of key file: ")
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		keyContents, err := os.ReadFile(keyFile)
		if err != nil {
			fmt.Println("Error reading key file:", err)
			continue
		}

		// 3. Validation Check (First 8 bytes header match)
		if len(scramContents) < 8 || len(keyContents) < 12 || !bytes.Equal(scramContents[:8], keyContents[:8]) {
			fmt.Println("Error: Wrong file pair or invalid header!")
			continue
		}

		fmt.Println("Header match verified! Decrypting...")

		// Decode original byte length from big-endian uint32
		dataLen := binary.BigEndian.Uint32(keyContents[8:12])
		fmt.Printf("%d data bytes will be unscrambled\n", dataLen)

		scramData := scramContents[8:]
		pointData := keyContents[12:]

		total := len(scramData)
		if len(pointData) < total*4 {
			fmt.Println("Error: Key file is corrupted or truncated!")
			continue
		}

		// 4. Unscramble Direct Memory Mapping (Fast & Safe)
		unscrambledOutput := make([]byte, total)
		for k := 0; k < total; k++ {
			pos := binary.BigEndian.Uint32(pointData[k*4 : (k+1)*4])
			if int(pos) < len(unscrambledOutput) {
				unscrambledOutput[pos] = scramData[k]
			}
		}

		// 5. Output File Creation
		baseName := filepath.Base(scramFile)
		outFileName := "unscrambled_" + strings.TrimSuffix(baseName, "scrambled.txt")

		err = os.WriteFile(outFileName, unscrambledOutput[:dataLen], 0644)
		if err != nil {
			fmt.Println("Error writing unscrambled file:", err)
			continue
		}

		fmt.Println("Successfully restored file:", outFileName)
		break
	}
}

func promptInput(r *bufio.Reader, prompt string) (string, error) {
	fmt.Print(prompt)
	input, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}
