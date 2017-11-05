package gofparse

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"sync"
	"testing"

	"gopkg.in/yaml.v2"
)

func handlingInterrupt(ctx context.Context, cancel context.CancelFunc) {
	scInterrupt := make(chan os.Signal, 1)
	signal.Notify(scInterrupt, os.Interrupt)

	select {
	case <-scInterrupt:
		fmt.Println("User requested cancelation")
	}

	cancel()

	os.Exit(1)
}

func TestFParser_InitConfig(t *testing.T) {
	var parser FParser

	rawConfig, err := os.Open("./config-test.yml")
	if err != nil {
		t.Error(err)
	}
	defer rawConfig.Close()

	dat, err := ioutil.ReadFile("./config-test.yml")

	err = yaml.Unmarshal(dat, &parser)
	if err != nil {
		fmt.Printf("Error %v\n", err)
		t.Error(err)
	}

	if len(parser.LinesConfig) != 3 {
		t.Fail()
	}
}

func TestFParser_CallAnalize(t *testing.T) {
	var parser FParser

	rawConfig, err := os.Open("./config-test.yml")
	if err != nil {
		t.Error(err)
	}
	defer rawConfig.Close()

	dat, err := ioutil.ReadFile("./config-test.yml")

	err = yaml.Unmarshal(dat, &parser)
	if err != nil {
		t.Error(err)
	}

	if len(parser.LinesConfig) != 3 {
		t.Fail()
	}

	chSucess := make(chan FParserLine, 1)

	go (func() {
		for range chSucess {

		}
	})()

	ctx := context.Background()
	ctx, fnCancel := context.WithCancel(ctx)

	go handlingInterrupt(ctx, fnCancel)

	if err := parser.Analize(ctx, "./test.txt", chSucess); err != nil {
		t.Error(err)
		return
	}
}

func TestFParser_ResultsOfAnalize(t *testing.T) {
	var parser FParser

	rawConfig, err := os.Open("./config-test.yml")
	if err != nil {
		t.Error(err)
	}
	defer rawConfig.Close()

	dat, err := ioutil.ReadFile("./config-test.yml")

	err = yaml.Unmarshal(dat, &parser)
	if err != nil {
		t.Error(err)
	}

	if len(parser.LinesConfig) != 3 {
		t.Fail()
	}

	chSucess := make(chan FParserLine, 10)

	wg := &sync.WaitGroup{}

	totalResults := 0

	wg.Add(1)
	go (func() {
		defer wg.Done()
		for range chSucess {
			//fmt.Printf("Success: %v\n", lnParsed.Fields)
			totalResults++
		}
	})()

	ctx := context.Background()
	ctx, fnCancel := context.WithCancel(ctx)

	go handlingInterrupt(ctx, fnCancel)

	err = parser.Analize(ctx, "./test.txt", chSucess)
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
