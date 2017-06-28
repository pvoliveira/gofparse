package gofparse

import (
	"bufio"
	"errors"
	"os"
	"sync"

	"runtime"

	utilgofparse "github.com/pvoliveira/gofparse/util"
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
	Value       interface{}
}

// Analize - responsible by the processing of a file
func (parser *FParser) Analize(pathFile string, chParsedLine chan *FParserLine) (err error) {

	worker := func(queue chan string, out chan *FParserLine, wg *sync.WaitGroup) {
		defer wg.Done()
		for line := range queue {
			if r, errWrk := breakLineToFields(line, parser.LinesConfig); errWrk != nil {
				out <- &FParserLine{Value: errWrk.Error()}
			} else {
				out <- r
			}
		}
	}

	reader := func(pathFile string, lnQueue chan string) {
		var fileToParse *os.File

		fileToParse, errFile := os.Open(pathFile)
		if errFile != nil {
			return
		}

		fScanner := bufio.NewScanner(fileToParse)
		for fScanner.Scan() {
			lnQueue <- fScanner.Text()
		}
		fileToParse.Close()
	}

	queue := make(chan string, 1)
	//done := make(chan bool, 1)
	wg := &sync.WaitGroup{}

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go worker(queue, chParsedLine, wg)
	}

	reader(pathFile, queue)
	close(queue)

	wg.Wait()

	return nil
}

func breakLineToFields(strLine string, linesConfig []FParserLine) (line *FParserLine, err error) {
	var cfg FParserLine
	configFounded := false
	// iterate between the lines config to get the right config to the line
	for _, lnCfg := range linesConfig {
		// substring
		if utilgofparse.Substr(strLine, lnCfg.IdentifierField.InitPos-1, lnCfg.IdentifierField.Size) == lnCfg.IdentifierField.Key {
			cfg = lnCfg
			configFounded = true
			break
		}
	}

	if !configFounded {
		return nil, errors.New("Line configuration not found")
	}

	fields := make([]FParserField, len(cfg.Fields))

	// iterate between the fields to extract values from line
	for i, fieldCfg := range cfg.Fields {
		// substring
		rawField := utilgofparse.Substr(strLine, fieldCfg.InitPos-1, fieldCfg.Size)

		convertedValue, _ := utilgofparse.ConvertField(fieldCfg.TypeData, rawField)

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

	return &FParserLine{
		Description:     cfg.Description,
		IdentifierField: cfg.IdentifierField,
		Value:           strLine,
		Fields:          fields,
	}, nil
}
