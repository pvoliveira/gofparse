package main

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

func (parser *FParser) Analize(pathFile string) (opOk chan *FParserLine, opErr chan *FParserLine, err error) {
	var fileToParse *os.File

	fileToParse, err = os.Open(pathFile)
	if err != nil {
		return
	}
	defer fileToParse.Close()

	opOk = make(chan *FParserLine)
	opErr = make(chan *FParserLine)

	fileReader := bufio.NewReader(fileToParse)

	var ln []byte
	for {
		ln, _, err = fileReader.ReadLine()
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			return
		}

		breakLineToFields(ln, parser.LinesConfig, opOk, opErr)
	}

	return
}

func breakLineToFields(lineRaw []byte, linesConfig []FParserLine, successParse <-chan *FParserLine, errorParse <-chan *FParserLine) {
	fmt.Printf("breakLine called to:\n%s\n", string(lineRaw))
	//
	return
}

func main() {

}
