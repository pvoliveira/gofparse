package gofparse

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/pvoliveira/gofparse/stringhandlers"
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

func TestStringHandlers_Substr(t *testing.T) {
	input := "test"
	result := "es"

	if stringhandlers.Substr(input, 1, 2) != result {
		t.Error("Substring are wrong")
	}
}

func TestStringHandlers_ConvertField(t *testing.T) {
	strDate := "20170923"
	strDateFormat := "20060102"

	strNumber := "120120"

	retDate, err := stringhandlers.ConvertField("date", strDateFormat, strDate)
	if err != nil {
		t.Error("Fail to parse date value (err): " + err.Error())
	}

	if time.Date(2017, time.September, 23, 0, 0, 0, 0, time.UTC) != retDate {
		t.Error("Fail to parse date value (values)")
	}

	retInt, err := stringhandlers.ConvertField("number", "", strNumber)
	if err != nil {
		t.Error("Fail to parse number value (err): " + err.Error())
	}

	if reflect.ValueOf(retInt).Int() != 120120 {
		t.Error("Fail to parse number value (values)")
	}

	retDefault, err := stringhandlers.ConvertField("not_exists", "", "abc")
	if err != nil {
		t.Error("Fail default case")
	}

	if retDefault != "abc" {
		t.Error("Fail default case")
	}
}
