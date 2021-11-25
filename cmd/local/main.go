package main

import (
	"ebay-crawler/internal/config"
	"ebay-crawler/internal/crawler"
	"ebay-crawler/internal/itemcondition"
	"fmt"
	"github.com/akamensky/argparse"
	"log"
	"os"
)

func main() {
	parser := argparse.NewParser(
		"ebay-crawler",
		"Crawls ebay and saves found items using JSON format.",
	)
	allConditions := itemcondition.New.String() + "," + itemcondition.PreOwned.String() + "," + itemcondition.Refurbished.String() + "," + itemcondition.Used.String()
	s := parser.String("c", "condition", &argparse.Options{Required: false, Help: "Item condition (" + allConditions + ")"})
	err := parser.Parse(os.Args)

	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}
	itemCondition := itemcondition.ParseItemCondition(*s)

	crawler.Init(&config.AppConfig{
		InfoLog:                   log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog:                  log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		DebugLog:                  log.New(os.Stdout, "DEBUG\t", log.Ldate|log.Ltime),
		IsDebugEnabled:            true,
		CrawlItemCondition:        itemCondition,
		NumWorkers:                10,
		CrawlerChanBufferSize:     1024,
		PersistenceChanBufferSize: 1024,
	})
	crawler.Run(
		"https://www.ebay.com/",
	)

}
