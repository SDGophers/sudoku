package sudoku

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestSolver(t *testing.T) {
	const (
		testdir = "sample-puzzles"
	)
	// test files are in the sample_puzzles directory.
	// Each file has solution along with it in the form of
	// the puzzle name _solution.
	// First thing we need to do is grab all the solution filenames.
	files, err := ioutil.ReadDir(testdir)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		fn := f.Name()
		fnlen := len(fn)
		if fnlen > 4 && fn[fnlen-4:] != ".txt" {
			continue
		}
		if len(fn) > 13 && fn[fnlen-13:] == "_solution.txt" {
			continue
		}
		b, err := LoadFromFile(filepath.Join(testdir, fn))
		if err != nil {
			t.Fatal(err)
		}
		bs, err := LoadFromFile(filepath.Join(testdir, fn[:fnlen-4]+"_solution.txt"))
		if err != nil {
			t.Fatal(err)
		}
		t.Run(fn, func(t *testing.T) {
			if err = b.Solve(); err != nil {

				t.Fatal(err)
			}
			for r := range bs {
				for c := range bs[r] {
					if bs[r][c] != b[r][c] {
						t.Errorf("For (%v, %v) Got %v Expected %v:\nGot Puzzle:\n%v\nExpected Puzzle:\n%v\n", r, c, b[r][c], bs[r][c], b, bs)
					}
				}
			}
		})
	}
}

func BenchmarkSolver(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b := Board{
			{0, 0, 2, 0, 0, 4, 0, 0, 1},
			{0, 0, 0, 1, 2, 0, 0, 0, 4},
			{3, 0, 4, 0, 0, 0, 8, 2, 0},
			{6, 0, 0, 8, 4, 0, 0, 5, 0},
			{2, 0, 7, 0, 3, 0, 4, 0, 8},
			{0, 4, 0, 0, 5, 2, 0, 0, 7},
			{0, 7, 8, 0, 0, 0, 0, 0, 5},
			{4, 0, 0, 0, 7, 9, 0, 0, 0},
			{5, 0, 0, 4, 0, 0, 9, 0, 0},
		}
		b.Solve()
	}
}
