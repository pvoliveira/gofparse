package gofparse

import (
	"bufio"
	"context"
	"errors"
	"os"
	"sync"

	"runtime"

	"github.com/pvoliveira/gofparse/stringhandlers"
)

// FParser is the struct responsible by de process and contains of configuration
type FParser struct {
	FileDescription string        `json:"file-description" yaml:"file-description,omitempty"`
	Options         []string      `json:"options" yaml:"options,omitempty"`
	LinesConfig     []FParserLine `json:"lines-config" yaml:"lines-config"`
}

// FParserLine is the struct that hold the configuration to read the lines
type FParserLine struct {
	Description     string         `json:"description,omitempty" yaml:"description"`
	IdentifierField int            `json:"identifier-field-index" yaml:"identifier-field-index"`
	Fields          []FParserField `json:"fields" yaml:"fields"`
	Value           string         `json:"value" yaml:"value,omitempty"`
}

// FParserField is the struct that hold the configuration to identify a field in the line
type FParserField struct {
	Description string      `json:"description" yaml:"description,omitempty"`
	InitPos     int         `json:"init-pos" yaml:"init-pos"`
	Size        int         `json:"size" yaml:"size"`
	TypeData    string      `json:"type-data" yaml:"type-data"`
	Format      string      `json:"format" yaml:"format,omitempty"`
	Key         string      `json:"key" yaml:"key,omitempty"`
	Value       interface{} `json:"value" yaml:"value,omitempty"`
}

// Analize parse the file following the configs
func (parser *FParser) Analize(ctx context.Context, pathFile string, chParsedLine chan FParserLine) (err error) {

	worker := func(ctxW context.Context, queue chan string, out chan FParserLine, wg *sync.WaitGroup) {
		defer wg.Done()
		// for line := range queue {
		// 	if r, errWrk := breakLineToFields(line, parser.LinesConfig); errWrk != nil {
		// 		out <- &FParserLine{Value: errWrk.Error()}
		// 	} else {
		// 		out <- r
		// 	}
		// }
		for line := range queue {
			select {
			case <-ctx.Done(): // context canceled
				return
			default:
				if r, errWrk := breakLineToFields(line, parser.LinesConfig); errWrk != nil {
					out <- FParserLine{Value: errWrk.Error()}
				} else {
					out <- r
				}
			}
		}
	}

	reader := func(ctxR context.Context, pathFile string, lnQueue chan string) {
		var fileToParse *os.File

		fileToParse, errFile := os.Open(pathFile)
		if errFile != nil {
			return
		}
		defer fileToParse.Close()

		fScanner := bufio.NewScanner(fileToParse)
		for fScanner.Scan() {
			select {
			case <-ctxR.Done():
				return
			default:
				lnQueue <- fScanner.Text()
			}
		}

	}

	queue := make(chan string, 1)

	wg := &sync.WaitGroup{}

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go worker(ctx, queue, chParsedLine, wg)
	}

	reader(ctx, pathFile, queue)
	close(queue)

	wg.Wait()

	return nil
}

func breakLineToFields(strLine string, linesConfig []FParserLine) (line FParserLine, err error) {
	var cfg FParserLine
	configFounded := false
	// iterate between the lines config to get the right config to the line
	for _, lnCfg := range linesConfig {

		// substring
		if stringhandlers.Substr(strLine, lnCfg.Fields[lnCfg.IdentifierField-1].InitPos-1, lnCfg.Fields[lnCfg.IdentifierField-1].Size) == lnCfg.Fields[lnCfg.IdentifierField-1].Key {
			cfg = lnCfg
			configFounded = true
			break
		}
	}

	if !configFounded {
		return FParserLine{}, errors.New("Line configuration not found")
	}

	fields := make([]FParserField, len(cfg.Fields))

	// iterate between the fields to extract values from line
	for i, fieldCfg := range cfg.Fields {
		// substring
		rawField := stringhandlers.Substr(strLine, fieldCfg.InitPos-1, fieldCfg.Size)

		convertedValue, _ := stringhandlers.ConvertField(fieldCfg.TypeData, fieldCfg.Format, rawField)

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

	return FParserLine{
		Description:     cfg.Description,
		IdentifierField: cfg.IdentifierField,
		Value:           strLine,
		Fields:          fields,
	}, nil
}
