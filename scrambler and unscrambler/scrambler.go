package main

import (
	"bufio"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter the file name to scramble: ")
		fileName, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		fileName = strings.TrimSpace(fileName)
		fileContents, err := os.ReadFile(fileName)
		if err != nil {
			fmt.Println("Error reading file:", err)
			continue
		}

		n := len(fileContents)
		rn := int(math.Round(math.Sqrt(float64(n))))

		reqRows, reqCols := rn, rn
		if n-(rn*rn) > 0 {
			reqRows = rn + 1
		}

		// Allocate Matrix
		scram := make([][]byte, reqRows)
		pointSheet := make([][]int, reqRows)
		for i := range scram {
			scram[i] = make([]byte, reqCols)
			pointSheet[i] = make([]int, reqCols)
		}

		// Populate Initial Data
		k, m := 0, 0
		for i := 0; i < reqRows; i++ {
			for j := 0; j < reqCols; j++ {
				if k < n {
					scram[i][j] = fileContents[k]
					k++
				}
				pointSheet[i][j] = m
				m++
			}
		}

		// 1. Interleave Rows (Top/Bottom Swap)
		scramOut := make([][]byte, reqRows)
		pointOut := make([][]int, reqRows)

		ur, lr := 0, reqRows-1
		for sr := 0; sr < reqRows; sr++ {
			if sr%2 == 0 {
				scramOut[sr] = scram[ur]
				pointOut[sr] = pointSheet[ur]
				ur++
			} else {
				scramOut[sr] = scram[lr]
				pointOut[sr] = pointSheet[lr]
				lr--
			}
		}

		// 2. Transpose 1
		transScram1 := make([][]byte, reqCols)
		transPoint1 := make([][]int, reqCols)
		for d := range transScram1 {
			transScram1[d] = make([]byte, reqRows)
			transPoint1[d] = make([]int, reqRows)
		}
		for i := 0; i < reqRows; i++ {
			for j := 0; j < reqCols; j++ {
				transScram1[j][i] = scramOut[i][j]
				transPoint1[j][i] = pointOut[i][j]
			}
		}

		// 3. Interleave Columns
		scramOut2 := make([][]byte, reqCols)
		pointOut2 := make([][]int, reqCols)

		ur2, lr2 := 0, reqCols-1
		for sr2 := 0; sr2 < reqCols; sr2++ {
			if sr2%2 == 0 {
				scramOut2[sr2] = transScram1[ur2]
				pointOut2[sr2] = transPoint1[ur2]
				ur2++
			} else {
				scramOut2[sr2] = transScram1[lr2]
				pointOut2[sr2] = transPoint1[lr2]
				lr2--
			}
		}

		// 4. Final Transpose
		scrambled := make([][]byte, reqRows)
		point := make([][]int, reqRows)
		for d := range scrambled {
			scrambled[d] = make([]byte, reqCols)
			point[d] = make([]int, reqCols)
		}
		for i := 0; i < reqCols; i++ {
			for j := 0; j < reqRows; j++ {
				scrambled[j][i] = scramOut2[i][j]
				point[j][i] = pointOut2[i][j]
			}
		}

		// Flatten Streams
		scramStream := make([]byte, reqRows*reqCols)
		pointStream := make([]byte, reqRows*reqCols*4)

		df, ps := 0, 0
		for i := 0; i < reqRows; i++ {
			for j := 0; j < reqCols; j++ {
				scramStream[df] = scrambled[i][j]
				binary.BigEndian.PutUint32(pointStream[ps:ps+4], uint32(point[i][j]))
				df++
				ps += 4
			}
		}

		// Generate 8-byte Random Header
		common := make([]byte, 8)
		_, _ = rand.Read(common)

		dlen := make([]byte, 4)
		binary.BigEndian.PutUint32(dlen, uint32(n))

		// Write Final Files
		scramFinal := append(common, scramStream...)
		keyFinal := append(append(common, dlen...), pointStream...)

		_ = os.WriteFile(fileName+"scrambled.txt", scramFinal, 0644)
		_ = os.WriteFile(fileName+"key.txt", keyFinal, 0644)

		fmt.Println("File scrambled successfully!")
		fmt.Println("Saved:", fileName+"scrambled.txt")
		fmt.Println("Saved:", fileName+"key.txt")
		break
	}
}
