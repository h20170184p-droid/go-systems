package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"strings"
)

func main() {
	// Based on the number of bytes the number of rows and columns need to be allocated first.
	// fmt.Printf("\nEnter the number: > ")
	re := bufio.NewReader(os.Stdin)
BEGINNING:
	fmt.Printf("\nEnter the file name: > ")
	file_name, e1 := re.ReadString('\n')
	if e1 != nil {
		fmt.Println("Error: ", e1)
		goto BEGINNING
	}
	file_name = strings.TrimSpace(file_name)
	file_contents, e2 := os.ReadFile(file_name)
	if e2 != nil {
		fmt.Println("Error: ", e2)
		goto BEGINNING
	}

	var n, reqrows, reqcols int
	n = len(file_contents)
	// fmt.Scan(&n)
	rn := int(math.Round(math.Sqrt(float64(n)))) // square root of n
	diff := n - (rn * rn)                        // the difference
	// fmt.Println(rn, diff)
	if diff > 0 {
		reqrows = rn + 1
		reqcols = rn
	} else if diff <= 0 {
		reqrows = rn
		reqcols = rn
	}
	// fmt.Println(reqrows, reqcols)

	scram := make([][]byte, reqrows)
	for b := range scram {
		scram[b] = make([]byte, reqcols)
	}
	// fmt.Println(scram)

	point_sheet := make([][]int, reqrows)
	for c := range point_sheet {
		point_sheet[c] = make([]int, reqcols)
	}

	var k, m int = 0, 0
	for i := range reqrows {
		for j := range reqcols {
			if k < n {
				scram[i][j] = file_contents[k]
				k++
			}
			point_sheet[i][j] = m
			m += 1
		}
	}
	// fmt.Println(scram)
	// fmt.Println(point_sheet)

	// Scramble both scram and point sheet alike
	// row scramble with upper and lower bounds

	scramout := make([][]byte, reqrows)
	for b := range scramout {
		scramout[b] = make([]byte, reqcols)
	}
	// fmt.Println(scram)

	point_out := make([][]int, reqrows)
	for c := range point_out {
		point_out[c] = make([]int, reqcols)
	}

	ur, lr := 0, (reqrows - 1)
	var sr int = 0 // row swap variable
	for {
		if sr%2 == 0 {
			scramout[sr] = scram[ur]
			point_out[sr] = point_sheet[ur]
			ur += 1
		} else {
			scramout[sr] = scram[lr]
			point_out[sr] = point_sheet[lr]
			lr -= 1
		}
		sr += 1

		// Break condition 1 is when lr < ur
		// Break condition 2 is when lr = ur, then execute once and break
		if lr < ur {
			break
		}
	}

	// fmt.Printf("\n\n\n")
	// fmt.Println(point_out)

	// Transpose
	transposescram1 := make([][]byte, reqcols)
	for d := range transposescram1 {
		transposescram1[d] = make([]byte, reqrows)
	}
	for i := range reqrows {
		for j := range reqcols {
			transposescram1[j][i] = scramout[i][j]
		}
	}

	transposepoint1 := make([][]int, reqcols)
	for d := range transposepoint1 {
		transposepoint1[d] = make([]int, reqrows)
	}
	for i := range reqrows {
		for j := range reqcols {
			transposepoint1[j][i] = point_out[i][j]
		}
	}

	// column scramble
	scramout2 := make([][]byte, reqcols)
	for b := range scramout2 {
		scramout2[b] = make([]byte, reqrows)
	}
	// fmt.Println(scram)

	point_out2 := make([][]int, reqcols)
	for c := range point_out2 {
		point_out2[c] = make([]int, reqrows)
	}

	ur2, lr2 := 0, (reqcols - 1)
	var sr2 int = 0 // row swap variable
	for {
		if sr2%2 == 0 {
			scramout2[sr2] = transposescram1[ur2]
			point_out2[sr2] = transposepoint1[ur2]
			ur2 += 1
		} else {
			scramout2[sr2] = transposescram1[lr2]
			point_out2[sr2] = transposepoint1[lr2]
			lr2 -= 1
		}
		sr2 += 1

		// Break condition 1 is when lr < ur
		// Break condition 2 is when lr = ur, then execute once and break
		if lr2 < ur2 {
			break
		}
	}

	// Transpose
	scrambled := make([][]byte, reqrows)
	for d := range scrambled {
		scrambled[d] = make([]byte, reqcols)
	}
	for i := range reqcols {
		for j := range reqrows {
			scrambled[j][i] = scramout2[i][j]
		}
	}

	point := make([][]int, reqrows)
	for d := range point {
		point[d] = make([]int, reqcols)
	}
	for i := range reqcols {
		for j := range reqrows {
			point[j][i] = point_out2[i][j]
		}
	}

	// fmt.Println(scrambled)
	// fmt.Printf("\n\n")
	// fmt.Println(point)
	// fmt.Println(len(point), len(point[0]))

	scram_stream := make([]byte, (reqrows * reqcols))
	point_stream := make([]byte, (reqcols * reqrows * 4))
	var df int = 0
	var ps int = 0
	for i := range reqrows {
		for j := range reqcols {
			if df < (reqcols * reqrows) {
				scram_stream[df] = scrambled[i][j]
				v := point[i][j]
				point_stream[ps] = byte(v >> 24)
				point_stream[ps+1] = byte(v >> 16)
				point_stream[ps+2] = byte(v >> 8)
				point_stream[ps+3] = byte(v)
				df += 1
				ps += 4
			}
		}
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var sb strings.Builder
	sb.Grow(8)
	for l := 0; l < 8; l++ {
		randomIndex := rand.IntN(len(charset))
		sb.WriteByte(charset[randomIndex])
	}
	common := []byte(sb.String())
	dlen := []byte{
		byte(n >> 24),
		byte(n >> 16),
		byte(n >> 8),
		byte(n),
	}

	ss := append(common, scram_stream...)
	pt := append(common, dlen...)
	ptm := append(pt, point_stream...)
	os.WriteFile(file_name+"scrambled.txt", ss, 0644)
	os.WriteFile(file_name+"key.txt", ptm, 0644)
	// os.WriteFile(file_name+"scrambled", )
}
