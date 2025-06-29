package datamatrix

import (
	"bytes"
	"testing"
)

func codeFromStr(str string, size *dmCodeSize) *datamatrixCode {
	code := newDataMatrixCode(size)
	idx := 0
	for _, r := range str {
		x := idx % size.Columns
		y := idx / size.Columns

		switch r {
		case '#':
			code.set(x, y, true)
		case '.':
			code.set(x, y, false)
		default:
			continue
		}

		idx++
	}
	return code
}

func Test_Issue12(t *testing.T) {
	data := `{"po":12,"batchAction":"start_end"}`
	realData := addPadding(encodeText(data), 36)
	wantedData := []byte{124, 35, 113, 112, 35, 59, 142, 45, 35, 99, 98, 117, 100, 105, 66, 100, 117, 106, 112, 111, 35, 59, 35, 116, 117, 98, 115, 117, 96, 102, 111, 101, 35, 126, 129, 181}

	if bytes.Compare(realData, wantedData) != 0 {
		t.Error("Data Encoding failed")
		return
	}

	var codeSize *dmCodeSize
	for _, s := range codeSizes {
		if s.DataCodewords() >= len(wantedData) {
			codeSize = s
			break
		}
	}
	realECC := ec.calcECC(realData, codeSize)[len(realData):]
	wantedECC := []byte{196, 53, 147, 192, 151, 213, 107, 61, 98, 251, 50, 71, 186, 15, 43, 111, 165, 243, 209, 79, 128, 109, 251, 4}
	if bytes.Compare(realECC, wantedECC) != 0 {
		t.Errorf("Error correction calculation failed\nGot: %v", realECC)
		return
	}

	barcode := `
#.#.#.#.#.#.#.#.#.#.#.#.
#....###..#..#....#...##
##.......#...#.#.#....#.
#.###...##..#...##.##..#
##...####..##..#.#.#.##.
#.###.##.###..#######.##
#..###...##.##..#.##.##.
#.#.#.#.#.#.###....#.#.#
##.#...#.#.#..#...#####.
#...####..#...##..#.#..#
##...#...##.###.#.....#.
#.###.#.##.#.....###..##
##..#####...#..##...###.
###...#.####.##.#.#.#..#
#..###..#.#.####.#.###..
###.#.#..#..#.###.#.##.#
#####.##.###..#.####.#..
#.##.#......#.#..#.#.###
###.#....######.#...##..
##...#..##.###..#...####
#.######.###.##..#...##.
#..#..#.##.#..####...#.#
###.###..#..##.#.##...#.
########################`

	bc, err := Encode(data)

	if err != nil {
		t.Error(err)
		return
	}
	realResult := bc.(*datamatrixCode)
	if realResult.Columns != 24 || realResult.Rows != 24 {
		t.Errorf("Got wrong barcode size %dx%d", realResult.Columns, realResult.Rows)
		return
	}

	wantedResult := codeFromStr(barcode, realResult.dmCodeSize)

	for x := 0; x < wantedResult.Columns; x++ {
		for y := 0; y < wantedResult.Rows; y++ {
			r := realResult.get(x, y)
			w := wantedResult.get(x, y)
			if w != r {
				t.Errorf("Failed at: c%d/r%d", x, y)
			}
		}
	}
}


func Test_encode(t *testing.T) {
    data:="12"
    b := encodeText(data)
    if bytes.Compare(b, []byte{142}) != 0 {
		t.Errorf("Error encode 12 \nGot: %v", b)
		return
	}
	
	data ="123"
    b = encodeText(data)
    if bytes.Compare(b, []byte{142, 52}) != 0 {
		t.Errorf("Error encode 12 \nGot: %v", b)
		return
	}
	
	data ="abc"
    b = encodeText(data)
    if bytes.Compare(b, []byte{98,99,100}) != 0 {
		t.Errorf("Error encode abc \nGot: %v", b)
		return
	}
	
	data ="ABC"
    b = encodeText(data)
    if bytes.Compare(b, []byte{66,67,68}) != 0 {
		t.Errorf("Error encode abc \nGot: %v", b)
		return
	}
	
	data ="1ABC"
    b = encodeText(data)
    if bytes.Compare(b, []byte{50, 66,67,68}) != 0 {
		t.Errorf("Error encode abc \nGot: %v", b)
		return
	}
		
	data ="A\r\nA"
    b = encodeText(data)
    if bytes.Compare(b, []byte{66,14,11,66}) != 0 {
		t.Errorf("Error encode A\\r\\nA \nGot: %v", b)
		return
	}
	
	data ="\x2900"
    b = encodeText(data)
    if bytes.Compare(b, []byte{232,130}) != 0 {
		t.Errorf("Error encode GS1 \nGot: %v", b)
		return
	}
}
