package gofparse

import "testing"
import "encoding/json"
import "os"
import "fmt"
import "sync"

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

	chSucess := make(chan *FParserLine, 100)

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

	chSucess := make(chan *FParserLine, 100)

	wg := &sync.WaitGroup{}
	//var fileMutex sync.Mutex

	totalResults := 0

	wg.Add(1)
	go (func() {
		defer wg.Done()

		//fileMutex.Lock()
		//defer fileMutex.Unlock()

		for lnParsed := range chSucess {
			fmt.Printf("Success: %v\n", lnParsed.Fields)
			totalResults++
		}
	})()

	err = parser.Analize("./test.txt", chSucess)
	if err != nil {
		t.Error(err)
	}

	wg.Wait()
	//time.Sleep(time.Second * 1)

	if totalResults != 3 {
		fmt.Printf("Error after read results (waiting 3): %d\n", totalResults)
		t.Error()
		return
	}

}
