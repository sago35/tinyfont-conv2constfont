package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
	"tinygo.org/x/tinyfont/freesans"
	"tinygo.org/x/tinyfont/freeserif"
	"tinygo.org/x/tinyfont/gophers"
	"tinygo.org/x/tinyfont/notoemoji"
	"tinygo.org/x/tinyfont/notosans"
	"tinygo.org/x/tinyfont/proggy"
)

type target struct {
	font     tinyfont.Font
	path     string
	fontname string
	pkgname  string
}

var targets = []target{
	{path: "org_01.go", font: tinyfont.Org01, fontname: "Org01", pkgname: "tinyfont"},
	{path: "picopixel.go", font: tinyfont.Picopixel, fontname: "Picopixel", pkgname: "tinyfont"},
	{path: "tiny3x3a2pt7b.go", font: tinyfont.Tiny3x3a2pt7b, fontname: "Tiny3x3a2pt7b", pkgname: "tinyfont"},
	{path: "tomthumb.go", font: tinyfont.TomThumb, fontname: "TomThumb", pkgname: "tinyfont"},
	{path: "freemono/freemono12pt7b.go", font: freemono.Regular12pt7b, fontname: "Regular12pt7b", pkgname: "freemono"},
	{path: "freemono/freemono18pt7b.go", font: freemono.Regular18pt7b, fontname: "Regular18pt7b", pkgname: "freemono"},
	{path: "freemono/freemono24pt7b.go", font: freemono.Regular24pt7b, fontname: "Regular24pt7b", pkgname: "freemono"},
	{path: "freemono/freemono9pt7b.go", font: freemono.Regular9pt7b, fontname: "Regular9pt7b", pkgname: "freemono"},
	{path: "freemono/freemonobold12pt7b.go", font: freemono.Bold12pt7b, fontname: "Bold12pt7b", pkgname: "freemono"},
	{path: "freemono/freemonobold18pt7b.go", font: freemono.Bold18pt7b, fontname: "Bold18pt7b", pkgname: "freemono"},
	{path: "freemono/freemonobold24pt7b.go", font: freemono.Bold24pt7b, fontname: "Bold24pt7b", pkgname: "freemono"},
	{path: "freemono/freemonobold9pt7b.go", font: freemono.Bold9pt7b, fontname: "Bold9pt7b", pkgname: "freemono"},
	{path: "freemono/freemonoboldoblique12pt7b.go", font: freemono.BoldOblique12pt7b, fontname: "BoldOblique12pt7b", pkgname: "freemono"},
	{path: "freemono/freemonoboldoblique18pt7b.go", font: freemono.BoldOblique18pt7b, fontname: "BoldOblique18pt7b", pkgname: "freemono"},
	{path: "freemono/freemonoboldoblique24pt7b.go", font: freemono.BoldOblique24pt7b, fontname: "BoldOblique24pt7b", pkgname: "freemono"},
	{path: "freemono/freemonoboldoblique9pt7b.go", font: freemono.BoldOblique9pt7b, fontname: "BoldOblique9pt7b", pkgname: "freemono"},
	{path: "freemono/freemonooblique12pt7b.go", font: freemono.Oblique12pt7b, fontname: "Oblique12pt7b", pkgname: "freemono"},
	{path: "freemono/freemonooblique18pt7b.go", font: freemono.Oblique18pt7b, fontname: "Oblique18pt7b", pkgname: "freemono"},
	{path: "freemono/freemonooblique24pt7b.go", font: freemono.Oblique24pt7b, fontname: "Oblique24pt7b", pkgname: "freemono"},
	{path: "freemono/freemonooblique9pt7b.go", font: freemono.Oblique9pt7b, fontname: "Oblique9pt7b", pkgname: "freemono"},
	{path: "freesans/freesans12pt7b.go", font: freesans.Regular12pt7b, fontname: "Regular12pt7b", pkgname: "freesans"},
	{path: "freesans/freesans18pt7b.go", font: freesans.Regular18pt7b, fontname: "Regular18pt7b", pkgname: "freesans"},
	{path: "freesans/freesans24pt7b.go", font: freesans.Regular24pt7b, fontname: "Regular24pt7b", pkgname: "freesans"},
	{path: "freesans/freesans9pt7b.go", font: freesans.Regular9pt7b, fontname: "Regular9pt7b", pkgname: "freesans"},
	{path: "freesans/freesansbold12pt7b.go", font: freesans.Bold12pt7b, fontname: "Bold12pt7b", pkgname: "freesans"},
	{path: "freesans/freesansbold18pt7b.go", font: freesans.Bold18pt7b, fontname: "Bold18pt7b", pkgname: "freesans"},
	{path: "freesans/freesansbold24pt7b.go", font: freesans.Bold24pt7b, fontname: "Bold24pt7b", pkgname: "freesans"},
	{path: "freesans/freesansbold9pt7b.go", font: freesans.Bold9pt7b, fontname: "Bold9pt7b", pkgname: "freesans"},
	{path: "freesans/freesansboldoblique12pt7b.go", font: freesans.BoldOblique12pt7b, fontname: "BoldOblique12pt7b", pkgname: "freesans"},
	{path: "freesans/freesansboldoblique18pt7b.go", font: freesans.BoldOblique18pt7b, fontname: "BoldOblique18pt7b", pkgname: "freesans"},
	{path: "freesans/freesansboldoblique24pt7b.go", font: freesans.BoldOblique24pt7b, fontname: "BoldOblique24pt7b", pkgname: "freesans"},
	{path: "freesans/freesansboldoblique9pt7b.go", font: freesans.BoldOblique9pt7b, fontname: "BoldOblique9pt7b", pkgname: "freesans"},
	{path: "freesans/freesansoblique12pt7b.go", font: freesans.Oblique12pt7b, fontname: "Oblique12pt7b", pkgname: "freesans"},
	{path: "freesans/freesansoblique18pt7b.go", font: freesans.Oblique18pt7b, fontname: "Oblique18pt7b", pkgname: "freesans"},
	{path: "freesans/freesansoblique24pt7b.go", font: freesans.Oblique24pt7b, fontname: "Oblique24pt7b", pkgname: "freesans"},
	{path: "freesans/freesansoblique9pt7b.go", font: freesans.Oblique9pt7b, fontname: "Oblique9pt7b", pkgname: "freesans"},
	{path: "freeserif/freeserif12pt7b.go", font: freeserif.Regular12pt7b, fontname: "Regular12pt7b", pkgname: "freeserif"},
	{path: "freeserif/freeserif18pt7b.go", font: freeserif.Regular18pt7b, fontname: "Regular18pt7b", pkgname: "freeserif"},
	{path: "freeserif/freeserif24pt7b.go", font: freeserif.Regular24pt7b, fontname: "Regular24pt7b", pkgname: "freeserif"},
	{path: "freeserif/freeserif9pt7b.go", font: freeserif.Regular9pt7b, fontname: "Regular9pt7b", pkgname: "freeserif"},
	{path: "freeserif/freeserifbold12pt7b.go", font: freeserif.Bold12pt7b, fontname: "Bold12pt7b", pkgname: "freeserif"},
	{path: "freeserif/freeserifbold18pt7b.go", font: freeserif.Bold18pt7b, fontname: "Bold18pt7b", pkgname: "freeserif"},
	{path: "freeserif/freeserifbold24pt7b.go", font: freeserif.Bold24pt7b, fontname: "Bold24pt7b", pkgname: "freeserif"},
	{path: "freeserif/freeserifbold9pt7b.go", font: freeserif.Bold9pt7b, fontname: "Bold9pt7b", pkgname: "freeserif"},
	{path: "freeserif/freeserifbolditalic12pt7b.go", font: freeserif.BoldItalic12pt7b, fontname: "BoldItalic12pt7b", pkgname: "freeserif"},
	{path: "freeserif/freeserifbolditalic18pt7b.go", font: freeserif.BoldItalic18pt7b, fontname: "BoldItalic18pt7b", pkgname: "freeserif"},
	{path: "freeserif/freeserifbolditalic24pt7b.go", font: freeserif.BoldItalic24pt7b, fontname: "BoldItalic24pt7b", pkgname: "freeserif"},
	{path: "freeserif/freeserifbolditalic9pt7b.go", font: freeserif.BoldItalic9pt7b, fontname: "BoldItalic9pt7b", pkgname: "freeserif"},
	{path: "freeserif/freeserifitalic12pt7b.go", font: freeserif.Italic12pt7b, fontname: "Italic12pt7b", pkgname: "freeserif"},
	{path: "freeserif/freeserifitalic18pt7b.go", font: freeserif.Italic18pt7b, fontname: "Italic18pt7b", pkgname: "freeserif"},
	{path: "freeserif/freeserifitalic24pt7b.go", font: freeserif.Italic24pt7b, fontname: "Italic24pt7b", pkgname: "freeserif"},
	{path: "freeserif/freeserifitalic9pt7b.go", font: freeserif.Italic9pt7b, fontname: "Italic9pt7b", pkgname: "freeserif"},
	{path: "gophers/gophers121pt.go", font: gophers.Regular121pt, fontname: "Regular121pt", pkgname: "gophers"},
	{path: "gophers/gophers14pt.go", font: gophers.Regular14pt, fontname: "Regular14pt", pkgname: "gophers"},
	{path: "gophers/gophers18pt.go", font: gophers.Regular18pt, fontname: "Regular18pt", pkgname: "gophers"},
	{path: "gophers/gophers22pt.go", font: gophers.Regular22pt, fontname: "Regular22pt", pkgname: "gophers"},
	{path: "gophers/gophers32pt.go", font: gophers.Regular32pt, fontname: "Regular32pt", pkgname: "gophers"},
	{path: "gophers/gophers58pt.go", font: gophers.Regular58pt, fontname: "Regular58pt", pkgname: "gophers"},
	{path: "notoemoji/NotoEmoji-Regular-12pt.go", font: notoemoji.NotoEmojiRegular12pt, fontname: "NotoEmojiRegular12pt", pkgname: "notoemoji"},
	{path: "notoemoji/NotoEmoji-Regular-16pt.go", font: notoemoji.NotoEmojiRegular16pt, fontname: "NotoEmojiRegular16pt", pkgname: "notoemoji"},
	{path: "notoemoji/NotoEmoji-Regular-20pt.go", font: notoemoji.NotoEmojiRegular20pt, fontname: "NotoEmojiRegular20pt", pkgname: "notoemoji"},
	{path: "notosans/notosans12pt.go", font: notosans.Notosans12pt, fontname: "Notosans12pt", pkgname: "notosans"},
	{path: "proggy/tinysz8pt7b.go", font: proggy.TinySZ8pt7b, fontname: "TinySZ8pt7b", pkgname: "proggy"},
}

func conv2constfont() error {

	for _, target := range targets {
		font := target.font
		//fmt.Printf("%#v\n", font)
		//fmt.Printf("YAdvance : %d\n", font.YAdvance)
		//for _, g := range font.Glyphs {
		//	fmt.Printf("%#v\n", g)
		//}

		w, err := os.Create(filepath.Join("../../tinygo-org/tinyfont", target.path))
		if err != nil {
			return err
		}
		defer w.Close()

		err = conv(w, target.pkgname, target.fontname, font)
		if err != nil {
			return err
		}
	}
	return nil
}

func conv(w io.Writer, pkgname, fontname string, ufont tinyfont.Font) error {
	ufont.Glyphs = sortGlyphs(ufont.Glyphs)

	//tmp, err := ioutil.TempFile(``, `tinyfontgen`)
	//if err != nil {
	//	return err
	//}
	//defer os.Remove(tmp.Name())
	tmp, err := os.Create("out.font")
	if err != nil {
		return err
	}
	defer tmp.Close()

	//fmt.Fprintln(tmp, `//`, filepath.Base(os.Args[0]), strings.Join(os.Args[1:], ` `))
	//fmt.Fprintln(tmp)
	fmt.Fprintln(tmp, `package `+pkgname)
	fmt.Fprintln(tmp)

	if pkgname != `tinyfont` {
		fmt.Fprintln(tmp, `import (`)
		fmt.Fprintln(tmp, `	"tinygo.org/x/tinyfont"`)
		fmt.Fprintln(tmp, `)`)
		fmt.Fprintln(tmp)
	}

	fmt.Fprintf(tmp, "\n")
	fmt.Fprintf(tmp, "var %s = s%s{}\n", fontname, fontname)
	fmt.Fprintf(tmp, "\n")
	fmt.Fprintf(tmp, "type s%s struct {\n", fontname)
	fmt.Fprintf(tmp, "}\n")
	fmt.Fprintf(tmp, "\n")

	fmt.Fprintf(tmp, `const c%s = "" +`, fontname)
	fmt.Fprintln(tmp)

	idxMap := map[rune]int{}
	current := 0
	for i, g := range ufont.Glyphs {
		idxMap[g.Rune] = current
		c := fmt.Sprintf("%c", ufont.Glyphs[i].Rune)
		if ufont.Glyphs[i].Rune == 0 {
			c = ""
		}
		length := 4 + 1 + 1 + 1 + 1 + 1 + len(g.Bitmaps)
		fmt.Fprintf(tmp, `	/* %s */ `, c)
		fmt.Fprintf(tmp, `"\x%02X\x%02X" +`, (uint16(length) >> 8), uint8(length))
		fmt.Fprintf(tmp, `"\x%02X\x%02X\x%02X\x%02X" +`, ((uint32(g.Rune) & 0xFF000000) >> 24), ((uint32(g.Rune) & 0x00FF0000) >> 16), ((uint32(g.Rune) & 0x0000FF00) >> 8), (uint32(g.Rune) & 0x000000FF))
		fmt.Fprintf(tmp, `"\x%02X\x%02X" +`, uint8(g.Width), uint8(g.Height))
		fmt.Fprintf(tmp, `"\x%02X\x%02X" +`, uint8(g.XAdvance), uint8(g.XOffset))
		fmt.Fprintf(tmp, `"\x%02X" +`, uint8(g.YOffset))
		fmt.Fprintf(tmp, `"`)
		for _, b := range g.Bitmaps {
			fmt.Fprintf(tmp, `\x%02X`, uint8(b))
		}
		fmt.Fprintf(tmp, `" +`)
		fmt.Fprintln(tmp)
		current += 2 + length
	}
	fmt.Fprintln(tmp, `""`)
	fmt.Fprintln(tmp)

	fmt.Fprintf(tmp, "func (f *s%s) GetGlyph(r rune) tinyfont.Glyph {\n", fontname)
	fmt.Fprintf(tmp, "	idx := -1\n")
	fmt.Fprintf(tmp, "\n")

	// switch
	fmt.Fprintf(tmp, "switch r {\n")
	keys := []int{}
	for k := range idxMap {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)

	for _, k := range keys {
		v := idxMap[rune(k)]
		fmt.Fprintf(tmp, "case 0x%04X: idx = %d\n", k, v)
	}
	fmt.Fprintf(tmp, "}\n")
	fmt.Fprintf(tmp, "\n")

	fmt.Fprintf(tmp, "if idx == -1 {\n")
	fmt.Fprintf(tmp, "	return tinyfont.Glyph{\n")
	fmt.Fprintf(tmp, "		Rune:     r,\n")
	fmt.Fprintf(tmp, "		Width:    0,\n")
	fmt.Fprintf(tmp, "		Height:   0,\n")
	fmt.Fprintf(tmp, "		XAdvance: uint8(c%s[6]),\n", fontname)
	fmt.Fprintf(tmp, "		XOffset:  int8(c%s[7]),\n", fontname)
	fmt.Fprintf(tmp, "		YOffset:  int8(c%s[8]),\n", fontname)
	fmt.Fprintf(tmp, "		Bitmaps:  []uint8{},\n")
	fmt.Fprintf(tmp, "	}\n")
	fmt.Fprintf(tmp, "}\n")
	fmt.Fprintf(tmp, "\n")

	fmt.Fprintf(tmp, "	length := int((uint16(c%s[idx+0]) << 8) + uint16(c%s[idx+1]))\n", fontname, fontname)
	fmt.Fprintf(tmp, "	idx += 2\n")
	fmt.Fprintf(tmp, "	ret := tinyfont.Glyph{\n")
	fmt.Fprintf(tmp, "		Rune:     rune((uint32(c%s[idx+0]) << 24) + (uint32(c%s[idx+1]) << 16) + (uint32(c%s[idx+2]) << 8) + uint32(c%s[idx+3])),\n", fontname, fontname, fontname, fontname)
	fmt.Fprintf(tmp, "		Width:    uint8(c%s[idx+4]),\n", fontname)
	fmt.Fprintf(tmp, "		Height:   uint8(c%s[idx+5]),\n", fontname)
	fmt.Fprintf(tmp, "		XAdvance: uint8(c%s[idx+6]),\n", fontname)
	fmt.Fprintf(tmp, "		XOffset:  int8(c%s[idx+7]),\n", fontname)
	fmt.Fprintf(tmp, "		YOffset:  int8(c%s[idx+8]),\n", fontname)
	fmt.Fprintf(tmp, "		Bitmaps:  []uint8(c%s[idx+9 : idx+length]),\n", fontname)
	fmt.Fprintf(tmp, "	}\n")
	fmt.Fprintf(tmp, "\n")
	fmt.Fprintf(tmp, "	return ret\n")
	fmt.Fprintf(tmp, "}\n")
	fmt.Fprintf(tmp, "\n")
	fmt.Fprintf(tmp, "func (f *s%s) GetYAdvance() uint8 {\n", fontname)
	fmt.Fprintf(tmp, "	return 0x%02X\n", ufont.YAdvance)
	fmt.Fprintf(tmp, "}\n")

	// gofmt
	buf := bytes.Buffer{}
	cmd := exec.Command(`gofmt`, tmp.Name())
	cmd.Stdout = w
	cmd.Stderr = &buf
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("%s : %s", err.Error(), strings.TrimSpace(buf.String()))
	}

	return nil
}

func sortGlyphs(glyphs []tinyfont.Glyph) []tinyfont.Glyph {
	sort.Slice(glyphs, func(i, j int) bool { return glyphs[i].Rune < glyphs[j].Rune })
	return glyphs
}
