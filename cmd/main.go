package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/pvoliveira/gofparse"
)

func main() {
	var err error

	var parser gofparse.FParser
	var configFile *os.File
	var inputFile *os.File
	var pathConfigFile *string
	var pathInputFile *string

	pathConfigFile = flag.String("config", "", "config file like doc (.json)")
	pathInputFile = flag.String("input", "", "text file (.txt)")
	flag.Parse()

	if configFile, err = os.Open(*pathConfigFile); err != nil {
		fmt.Printf("Error on trying open config file: \nFile %s\n Error: %s\n", *pathConfigFile, err.Error())
		os.Exit(1)
	}
	defer configFile.Close()

	if inputFile, err = os.Open(*pathInputFile); err != nil {
		fmt.Printf("Error on trying open input file: %s\n", err.Error())
		os.Exit(1)
	}
	inputFile.Close()

	err = json.NewDecoder(configFile).Decode(&parser)
	if err != nil {
		fmt.Printf("Configuration file is incorrect: %s\n", err.Error())
		os.Exit(2)
	}

	chSucess := make(chan gofparse.FParserLine, 10)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go (func() {
		defer wg.Done()
		for range chSucess {
			// to implement outputs
		}
	})()

	ctx := context.Background()
	ctx, fnCancel := context.WithCancel(ctx)

	go handlingInterrupt(ctx, fnCancel)

	<-time.After(time.Second * 5)

	err = parser.Analize(ctx, *pathInputFile, chSucess)
	if err != nil {
		close(chSucess)
		fmt.Printf("Error on trying parse file: %s\n", err.Error())
		os.Exit(3)
	}
	close(chSucess)

	wg.Wait()
}

// gracefully exit on ctrl+c / interrupt signal
func handlingInterrupt(ctx context.Context, cancel context.CancelFunc) {
	scInterrupt := make(chan os.Signal, 1)
	signal.Notify(scInterrupt, os.Interrupt)

	select {
	case <-scInterrupt:
		fmt.Println("User requested cancelation")
	}

	cancel()

	os.Exit(99)
}
