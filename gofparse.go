package gofparse

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// FParser - Entity responsible by de process and container of configuration
type FParser struct {
	FileDescription string
	Options         []string
	LinesConfig     []FParserLine
}

// FParserLine - Struct have the configuration to read the lines
type FParserLine struct {
	Description     string
	IdentifierField FParserField
	Fields          []FParserField
	Value           *string
}

// FParserField - Struct have the configuration to identify a field in the line
type FParserField struct {
	Description string
	InitPos     int
	Size        int
	TypeData    string
	Key         string
	Value       interface{}
}

type fnCallBackLine func(ln *string)
type fnCallBackAnalize func(lnParsed *FParserLine)

// Analize - responsible by the processing of a file
func (parser *FParser) Analize(pathFile string, chParsedLine chan<- *FParserLine) (err error) {

	// channel which receive the lines
	chLine := make(chan *string, 100)

	wg := &sync.WaitGroup{}

	// goroutine to process lines
	// wg.Add(1)
	// go callBreakLine(wg, parser.LinesConfig, chLine, chSucesses, chErrors)

	wg.Add(1)
	go (func() {
		defer wg.Done()
		defer close(chParsedLine)

		for ln := range chLine {
			breakLineToFields(ln, parser.LinesConfig, chParsedLine)
		}
	})()

	// goroutine to read de file
	wg.Add(1)
	go (func() {
		defer close(chLine)
		defer wg.Done()
		fmt.Println("Reading file...")
		readFile(pathFile, chLine)
		fmt.Println("Reading file... Done!")
	})()

	wg.Wait()

	fmt.Println("Analize done!")

	return
}

func readFile(pathFile string, chLines chan<- *string) (err error) {
	var fileToParse *os.File

	fileToParse, err = os.Open(pathFile)
	if err != nil {
		return err
	}
	defer fileToParse.Close()

	fScanner := bufio.NewScanner(fileToParse)
	for fScanner.Scan() {
		ln := fScanner.Text()
		chLines <- &ln
	}
	return
}

func breakLineToFields(strLine *string, linesConfig []FParserLine, chParsedLine chan<- *FParserLine) {
	var cfg FParserLine
	configFounded := false
	// iterate between the lines config to get the right config to the line
	for _, lnCfg := range linesConfig {
		// substring
		if substr(strLine, lnCfg.IdentifierField.InitPos-1, lnCfg.IdentifierField.Size) == lnCfg.IdentifierField.Key {
			cfg = lnCfg
			configFounded = true
			break
		}
	}

	if !configFounded {
		chParsedLine <- &FParserLine{Value: strLine}
		return
	}

	fields := make([]FParserField, len(cfg.Fields))

	// iterate between the fields to extract values from line
	for i, fieldCfg := range cfg.Fields {
		// substring
		rawField := substr(strLine, fieldCfg.InitPos-1, fieldCfg.Size)

		convertedValue, _ := convertField(fieldCfg.TypeData, rawField)

		// instance FParseField with the value extracted
		fields[i] = FParserField{
			Description: fieldCfg.Description,
			InitPos:     fieldCfg.InitPos,
			Size:        fieldCfg.Size,
			TypeData:    fieldCfg.TypeData,
			Key:         fieldCfg.Key,
			Value:       convertedValue,
		}
	}

	chParsedLine <- &FParserLine{
		Description:     cfg.Description,
		IdentifierField: cfg.IdentifierField,
		Value:           strLine,
		Fields:          fields,
	}
}

// extract chars from string using runes
func substr(s *string, pos, length int) string {
	value := *s
	runes := []rune(value)
	l := pos + length

	if l > len(runes) {
		l = len(runes)
	}

	return string(runes[pos:l])
}

// convert values from string to the type configured
func convertField(typeData, value string) (newValue interface{}, err error) {

	switch typeData {
	case "date":
		newValue, err = time.Parse(time.RFC3339, strings.TrimSpace(value))
		break
	case "number":
		newValue, err = strconv.ParseFloat(strings.TrimSpace(value), 64)
		break
	default:
		newValue = value
	}
	return
}
