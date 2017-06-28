package gofparse

import (
	"encoding/json"
	"os"
	"sync"
	"testing"
)

func TestFParser_InitConfig(t *testing.T) {
	var parser FParser

	rawConfig, err := os.Open("./config-test.json")
	if err != nil {
		t.Error(err)
	}
	defer rawConfig.Close()

	err = json.NewDecoder(rawConfig).Decode(&parser)
	if err != nil {
		t.Error(err)
	}

	if len(parser.LinesConfig) != 3 {
		t.Fail()
	}
}

func TestFParser_CallAnalize(t *testing.T) {
	var parser FParser

	rawConfig, err := os.Open("./config-test.json")
	if err != nil {
		t.Error(err)
	}
	defer rawConfig.Close()

	err = json.NewDecoder(rawConfig).Decode(&parser)
	if err != nil {
		t.Error(err)
	}

	if len(parser.LinesConfig) != 3 {
		t.Fail()
	}

	chSucess := make(chan *FParserLine, 1)

	go (func() {
		for _ = range chSucess {

		}
	})()

	if err := parser.Analize("./test.txt", chSucess); err != nil {
		t.Error(err)
		return
	}

}

func TestFParser_ResultsOfAnalize(t *testing.T) {
	var parser FParser

	rawConfig, err := os.Open("./config-test.json")
	if err != nil {
		t.Error(err)
	}
	defer rawConfig.Close()

	err = json.NewDecoder(rawConfig).Decode(&parser)
	if err != nil {
		t.Error(err)
	}

	if len(parser.LinesConfig) != 3 {
		t.Fail()
	}

	chSucess := make(chan *FParserLine, 10)

	wg := &sync.WaitGroup{}
	//var fileMutex sync.Mutex

	totalResults := 0

	wg.Add(1)
	go (func() {
		defer wg.Done()
		for _ = range chSucess {
			//fmt.Printf("Success: %v\n", lnParsed.Fields)
			totalResults++
		}
	})()

	err = parser.Analize("./test.txt", chSucess)
	if err != nil {
		t.Error(err)
	}
	close(chSucess)

	wg.Wait()

	if totalResults != 3 {
		//fmt.Printf("Error after read results (waiting 3): %d\n", totalResults)
		t.Error()
		return
	}

}
