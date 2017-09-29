package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

func readFile(path string) [9][9]byte {
	fp, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()
	scn := bufio.NewScanner(fp)
	var ret [9][9]byte
	row := 0
	for scn.Scan() {
		txt := scn.Text()

		for i := 0; i < 9; i++ {
			if txt[i] == '_' {
				ret[row][i] = 0
			} else {
				ret[row][i] = txt[i] - '0'
			}
		}
		row++
	}
	if err := scn.Err(); err != nil {
		log.Fatal(err)
	}

	return ret
}

func writeBoard(board [9][9]byte) {
	for _, row := range board {
		for _, num := range row {
			if num == 0 {
				fmt.Print("_")
			} else {
				fmt.Printf("%c", num+'0')
			}
		}
		fmt.Println()
	}
}

func next(row int, col int) (int, int) {
	if col == 8 {
		return row + 1, 0
	}
	return row, col + 1
}

func invalidSet(row, col int, board *[9][9]byte) (trueRet [9]bool) {
	// ret[k] == true means the value k is our row on the board
	// for k = 1..9, that number is now invalid
	// for k = 0, that's just a blank
	var ret [10]bool
	for i := 0; i < 9; i++ {
		ret[board[row][i]] = true
		ret[board[i][col]] = true
	}

	sqRow := row - row%3
	sqCol := col - col%3
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			ret[board[sqRow+i][sqCol+j]] = true
		}
	}

	copy(trueRet[:], ret[1:10])
	return trueRet
}

func solve(row int, col int, board *[9][9]byte) bool {
	if row == 9 {
		return true
	}

	nextRow, nextCol := next(row, col)
	if board[row][col] != 0 {
		return solve(nextRow, nextCol, board)
	}

	invalids := invalidSet(row, col, board)
	for i, invalid := range invalids {
		if !invalid {
			number := byte(i + 1)
			board[row][col] = number
			if solve(nextRow, nextCol, board) {
				return true
			}
		}
	}

	board[row][col] = 0
	return false
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("Need filename")
		os.Exit(1)
	}
	filename := args[0]

	board := readFile(filename)
	success := solve(0, 0, &board)
	if success {
		writeBoard(board)
	} else {
		fmt.Println("Unsolvable")
	}
}
