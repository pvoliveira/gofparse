package gofparse

import (
	"bufio"
	"os"
	"sync"
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
	Value           string
}

// FParserField - Struct have the configuration to identify a field in the line
type FParserField struct {
	Description string
	InitPos     int
	Size        int
	TypeData    string
	Key         string
	Value       string
}

// Analize - responsible by the processing of a file
func (parser *FParser) Analize(pathFile string, chSucesses, chErrors chan<- *FParserLine) (err error) {

	// channel which receive the lines
	chLine := make(chan string, 10)

	wg := &sync.WaitGroup{}

	// goroutine to process lines
	wg.Add(1)
	go func() {
		for lineStr := range chLine {
			breakLineToFields(lineStr, parser.LinesConfig, chSucesses, chErrors)
		}
		wg.Done()
	}()

	// goroutine to read de file
	wg.Add(1)
	go func() {
		var fileToParse *os.File

		fileToParse, err = os.Open(pathFile)
		if err != nil {
			wg.Done()
			panic(err)
		}
		defer fileToParse.Close()

		fScanner := bufio.NewScanner(fileToParse)
		for fScanner.Scan() {
			chLine <- fScanner.Text()
		}
		close(chLine)
		wg.Done()
	}()

	wg.Wait()

	return err
}

func breakLineToFields(strLine string, linesConfig []FParserLine, chOk, chErr chan<- *FParserLine) {

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
		chErr <- &FParserLine{Value: strLine}
		return
	}

	fields := make([]FParserField, len(cfg.Fields))

	// iterate between the fields to extract values from line
	for i, fieldCfg := range cfg.Fields {
		// substring
		rawField := substr(strLine, fieldCfg.InitPos-1, fieldCfg.Size)
		// instance FParseField with the value extracted
		fields[i] = FParserField{
			Description: fieldCfg.Description,
			InitPos:     fieldCfg.InitPos,
			Size:        fieldCfg.Size,
			TypeData:    fieldCfg.TypeData,
			Key:         fieldCfg.Key,
			Value:       rawField,
		}
	}

	chOk <- &FParserLine{
		Description:     cfg.Description,
		IdentifierField: cfg.IdentifierField,
		Value:           strLine,
		Fields:          fields,
	}
}

// extract chars from string using runes
func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length

	if l > len(runes) {
		l = len(runes)
	}

	return string(runes[pos:l])
}
