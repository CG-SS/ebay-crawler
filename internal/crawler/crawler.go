package crawler

import (
	"ebay-crawler/internal/config"
	"ebay-crawler/internal/itemcondition"
	"ebay-crawler/internal/model"
	"ebay-crawler/internal/persistence"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var app *config.AppConfig

func Init(config *config.AppConfig) {
	app = config
}

func Run(rootWebsite string) {

	infoMsg := "Starting application with " + strconv.Itoa(app.NumWorkers) + " workers and"
	if app.CrawlItemCondition == itemcondition.Unknown {
		app.InfoLog.Println(infoMsg, "no filter on item condition.")
	} else {
		app.InfoLog.Println(infoMsg, "crawling only", app.CrawlItemCondition.String(), "items.")
	}

	bufferSize := app.CrawlerChanBufferSize
	crawlerWorkerChan := make(chan string, bufferSize)
	remainingLinksChan := make(chan int)
	foundLinksChan := make(chan string)
	persistenceChan := make(chan model.ItemModel, app.PersistenceChanBufferSize)

	go func() {
		if app.IsDebugEnabled {
			app.DebugLog.Println("Sending the root website")
		}
		foundLinksChan <- rootWebsite
	}()

	Init(app)
	go Manager(crawlerWorkerChan, remainingLinksChan, foundLinksChan, bufferSize)

	persistence.Init(app)
	go persistence.Manager(persistenceChan)

	wg := sync.WaitGroup{}

	for i := 0; i < app.NumWorkers; i++ {
		app.InfoLog.Println("Creating worker ", i)
		wg.Add(1)
		go crawler(crawlerWorkerChan, remainingLinksChan, foundLinksChan, persistenceChan, &wg)
	}

	wg.Wait()

	close(persistenceChan)
}

func crawler(crawlerWorkerChan chan string, remainingLinksChan chan int, foundLinksChan chan string, persistenceChan chan model.ItemModel, wg *sync.WaitGroup) {
	for link := range crawlerWorkerChan {
		if app.IsDebugEnabled {
			app.DebugLog.Println("Crawling ", link)
		}
		exploreLink(link, remainingLinksChan, foundLinksChan, persistenceChan)
	}

	wg.Done()
}

func exploreLink(link string, remainingLinksChan chan int, foundLinksChan chan string, persistenceChan chan model.ItemModel) {
	res, err := http.Get(link)

	if err != nil {
		remainingLinksChan <- -1
		app.ErrorLog.Println("Failed to crawl ", link, " reason: ", err)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			remainingLinksChan <- -1
		}
	}(res.Body)

	if res.StatusCode != 200 {
		remainingLinksChan <- -1
		return
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		app.ErrorLog.Println("Failed to crawl ", link, " reason: ", err)
		remainingLinksChan <- -1
		return
	}

	doc.Find("a[href]").Each(func(index int, item *goquery.Selection) {
		href, exists := item.Attr("href")
		if exists && isLink(href) {
			if app.IsDebugEnabled {
				app.DebugLog.Println("Found link: ", href)
			}
			foundLinksChan <- href
		}
	})

	if isItemLink(link) {
		item, ok := extractItem(link, doc)

		if ok {
			if app.IsDebugEnabled {
				app.DebugLog.Println("Found item: ", item)
			}
			persistenceChan <- *item
		}
	}

	remainingLinksChan <- -1
}

func extractItem(link string, doc *goquery.Document) (*model.ItemModel, bool) {

	id, ok := extractIdFromLink(link)
	if !ok {
		return nil, false
	}

	var itemCondition = itemcondition.Unknown
	doc.Find(".item-highlight").Each(func(i int, selection *goquery.Selection) {
		currentItemCondition := itemcondition.ParseItemCondition(selection.Text())
		if currentItemCondition != itemcondition.Unknown {
			itemCondition = currentItemCondition
		}
	})

	if itemCondition == itemcondition.Unknown {
		return nil, false
	}
	if app.CrawlItemCondition != itemcondition.Unknown && itemCondition != itemCondition {
		return nil, false
	}

	title := doc.Find(".product-title").First().Text()
	if title == "" {
		return nil, false
	}

	price := doc.Find(".display-price").First().Text()
	if price == "" {
		return nil, false
	}

	return &model.ItemModel{
		Id:         id,
		Title:      title,
		Condition:  itemCondition.String(),
		Price:      price,
		ProductUrl: link,
	}, true

}

func extractIdFromLink(link string) (string, bool) {
	linkSplit := strings.Split(link, "/")
	if len(linkSplit) <= 0 {
		return "", false
	}

	urlParams := linkSplit[len(linkSplit)-1]
	urlParamsSplit := strings.Split(urlParams, "?")
	if len(urlParamsSplit) <= 0 {
		return "", false
	}

	return urlParamsSplit[0], true
}

func isItemLink(link string) bool {
	return strings.HasPrefix(link, "https://www.ebay.com/itm/")
}

func isLink(l string) bool {
	return strings.HasPrefix(l, "https://") || strings.HasPrefix(l, "http://")
}
