package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
)

type board [9][9]uint8

var filename = "sample-puzzles/hard_1_puzzle.txt"

func (b *board) LoadFromFile(filename string) (err error) {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(file)
	row := 0
	for scanner.Scan() {
		line := scanner.Text()
		for col, ch := range line {
			if ch == '_' {
				continue
			}
			b[row][col] = uint8((ch - '1') + 1)
		}
		row++
	}
	if err = scanner.Err(); err != nil {
		return err
	}
	return nil
}

func (b *board) fillMaskForRow(mask *[10]bool, row int) {
	for _, val := range b[row] {
		if val != 0 {
			mask[int(val)] = true
		}
	}
}
func (b *board) fillMaskForCol(mask *[10]bool, col int) {
	for _, row := range b {
		val := row[col]
		if val != 0 {
			mask[int(val)] = true
		}
	}
}

func (b *board) fillMaskForQuad(mask *[10]bool, row, col int) {
	// Get upper left coord for quad.
	wrow := (row / 3) * 3
	wcol := (col / 3) * 3
	for r := wrow; r < wrow+3; r++ {
		for c := wcol; c < wcol+3; c++ {
			val := b[r][c]
			if val != 0 {
				mask[int(val)] = true
			}
		}
	}
}

func (b *board) MaskForPos(row, col int) (mask [10]bool) {
	b.fillMaskForRow(&mask, row)
	b.fillMaskForCol(&mask, col)
	b.fillMaskForQuad(&mask, row, col)
	return
}

func (b board) String() (out string) {
	template := []byte(`
┏━━━┯━━━┯━━━┓
┃WWW│WWW│WWW┃
┃WWW│WWW│WWW┃
┃WWW│WWW│WWW┃
┠───┼───┼───┨
┃WWW│WWW│WWW┃
┃WWW│WWW│WWW┃
┃WWW│WWW│WWW┃
┠───┼───┼───┨
┃WWW│WWW│WWW┃
┃WWW│WWW│WWW┃
┃WWW│WWW│WWW┃
┗━━━┷━━━┷━━━┛
	`)

	for r := range b {
		for c := range b[r] {
			val := b[r][c]
			replacement := byte(' ')
			if val != 0 {
				replacement = byte('1' + (val - 1))
			}
			idx := bytes.IndexRune(template, 'W')
			if idx == -1 {
				// More values then slots in templates?
				log.Println("More values then slots in template?", template, r, c)
				break
			}
			template[idx] = replacement
		}
	}
	return string(template)
}

func main() {

	var b board
	b.LoadFromFile(filename)
	fmt.Println(&b)

	var backtrack [][3]int
	var skip int
	for r := 0; r < len(b); r++ {
	LoopCol:
		for c := 0; c < len(b[r]); {
			if b[r][c] != 0 {
				c++
				continue
			}

			mask := b.MaskForPos(r, c)
			// Number of positions to skip
			skipping := skip
			for val, unaval := range mask[1:] {
				if unaval {
					continue
				}
				if skipping > 0 {
					skipping--
					continue
				}
				// we filled out the position
				b[r][c] = uint8(val) + 1
				backtrack = append(backtrack, [3]int{skip + 1, r, c})
				skip = 0
				c++
				continue LoopCol
			}
			// Need to backtrack.
			if len(backtrack) == 0 {
				panic("Unsolvable.")
			}
			bkr := backtrack[len(backtrack)-1]
			skip, r, c = bkr[0], bkr[1], bkr[2]
			// Let it know we need to recalculate this value.
			b[r][c] = 0
			// Remove the last value from the stack.
			backtrack = backtrack[0 : len(backtrack)-1]
		}
	}
	log.Println(b)
}
