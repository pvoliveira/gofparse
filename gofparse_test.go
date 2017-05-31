package gofparse

import "testing"
import "encoding/json"
import "os"

func TestFParserInitConfig(t *testing.T) {
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

func TestFParserAnalize(t *testing.T) {
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
	chError := make(chan *FParserLine, 10)

	go func() {
		for {
			select {
			case <-chSucess:
			case <-chError:
			}
		}
	}()

	if err := parser.Analize("./test.txt", chSucess, chError); err != nil {
		t.Error(err)
		return
	}

}
