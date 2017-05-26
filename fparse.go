package fparse

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type FParser struct {
	FileDescription string
	Options         []string
	LinesConfig     []FParserLine
}

type FParserLine struct {
	Description     string
	IdentifierField FParserField
	Fields          []FParserField
	Value           string
}

type FParserField struct {
	Description string
	InitPos     int
	Size        int
	TypeData    string
	Key         string
}

func (parser *FParser) Analize(pathFile string) (lnOk chan *FParserLine, lnErr chan *FParserLine, err error) {
	var fileToParse *os.File

	fileToParse, err = os.Open(pathFile)
	if err != nil {
		return
	}
	defer fileToParse.Close()

	lnOk = make(chan *FParserLine)
	lnErr = make(chan *FParserLine)

	fileReader := bufio.NewReader(fileToParse)

	for {
		var ln []byte
		ln, _, err = fileReader.ReadLine()
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			return
		}

		breakLineToFields(ln, parser.LinesConfig, lnOk, lnErr)
	}
	return
}

func breakLineToFields(lineRaw []byte, linesConfig []FParserLine, successLine <-chan *FParserLine, errorLine <-chan *FParserLine) {
	fmt.Printf("breakLine called to:\n%s\n", string(lineRaw))
	//
	return
}
