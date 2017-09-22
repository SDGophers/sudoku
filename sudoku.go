package sudoku

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
)

type Board [9][9]uint8

var ErrUnsolvable = errors.New("Unsolvable")
var ErrTemplateSolts = errors.New("More values then slots in template.")
var ErrInvalidInput = errors.New("Got an invalid input file.")

func Load(src io.Reader) (b *Board, err error) {
	b = new(Board)
	scanner := bufio.NewScanner(src)
	row := 0
	for scanner.Scan() {
		line := scanner.Text()
		for col, ch := range line {
			if ch != '_' && (ch < '1' || ch > '9') {
				return nil, ErrInvalidInput
			}
			if ch == '_' {
				continue
			}
			b[row][col] = uint8((ch - '1') + 1)
		}
		row++
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}
	return b, nil
}

func LoadFromFile(filename string) (b *Board, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return b, err
	}
	defer file.Close()
	return Load(file)
}

func (b *Board) fillMaskForRow(mask *[10]bool, row int) {
	for _, val := range b[row] {
		if val != 0 {
			mask[int(val)] = true
		}
	}
}
func (b *Board) fillMaskForCol(mask *[10]bool, col int) {
	for _, row := range b {
		val := row[col]
		if val != 0 {
			mask[int(val)] = true
		}
	}
}

func (b *Board) fillMaskForQuad(mask *[10]bool, row, col int) {
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

func (b *Board) MaskForPos(row, col int) (mask [10]bool) {
	b.fillMaskForRow(&mask, row)
	b.fillMaskForCol(&mask, col)
	b.fillMaskForQuad(&mask, row, col)
	return
}

func (b Board) FillTemplate(template []byte, repl rune, empty rune) error {
	if empty == 0 {
		empty = ' '
	}
	for r := range b {
		for c := range b[r] {
			val := b[r][c]
			replacement := byte(empty)
			if val != 0 {
				replacement = byte('1' + (val - 1))
			}
			idx := bytes.IndexRune(template, repl)
			if idx == -1 {
				// More values then slots in templates?
				return ErrTemplateSolts
			}
			template[idx] = replacement
		}
	}
	return nil
}

func (b Board) String() (out string) {
	template := []byte(`┏━━━┯━━━┯━━━┓
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
┗━━━┷━━━┷━━━┛`)

	b.FillTemplate(template, 'W', ' ')
	return string(template)
}

func (b *Board) Solve() error {
	// Stack of backtracking points. (value, row, col)
	var backtrack [][3]int
	startVal := 1
	for r := 0; r < len(b); r++ {
	LoopCol:
		for c := 0; c < len(b[r]); {
			if b[r][c] != 0 {
				c++
				startVal = 1
				continue
			}

			mask := b.MaskForPos(r, c)
			for offset, unaval := range mask[startVal:] {
				if unaval {
					continue
				}
				val := startVal + offset
				// we filled out the position
				b[r][c] = uint8(val)
				// Store the next value position, and the current row and column.
				backtrack = append(backtrack, [3]int{val + 1, r, c})
				c++
				startVal = 1
				continue LoopCol
			}

		BackTrack:
			// Need to backtrack.
			if len(backtrack) == 0 {
				return ErrUnsolvable
			}
			bkr := backtrack[len(backtrack)-1]
			startVal, r, c = bkr[0], bkr[1], bkr[2]
			// Let it know we need to recalculate this value.
			b[r][c] = 0
			// Remove the last value from the stack.
			backtrack = backtrack[0 : len(backtrack)-1]
			if startVal >= len(mask) {
				goto BackTrack
			}
		}
	}
	return nil
}
