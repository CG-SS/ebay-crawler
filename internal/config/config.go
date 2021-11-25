package config

import (
	"ebay-crawler/internal/itemcondition"
	"log"
)

type AppConfig struct {
	InfoLog                   *log.Logger
	ErrorLog                  *log.Logger
	DebugLog                  *log.Logger
	CrawlItemCondition        itemcondition.ItemCondition
	IsDebugEnabled            bool
	NumWorkers                int
	CrawlerChanBufferSize     int
	PersistenceChanBufferSize int
}
