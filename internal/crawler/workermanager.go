package crawler

func countMonitor(crawlerWorkerChan chan string, remainingLinksChan chan int, foundLinksChan chan string) {
	remainingLinks := 0

	for n := range remainingLinksChan {
		remainingLinks += n

		if remainingLinks <= 0 {
			close(crawlerWorkerChan)
			close(foundLinksChan)
			close(remainingLinksChan)
		}
	}
}

func workerMonitor(crawlerWorkerChan chan string, foundLinksChan chan string, remainingLinksChan chan int, bufferSize int) {
	visitedLinksMap := make(map[string]int)

	for link := range foundLinksChan {
		if len(crawlerWorkerChan) != bufferSize && visitedLinksMap[link] == 0 {
			visitedLinksMap[link] = 1
			crawlerWorkerChan <- link
			remainingLinksChan <- 1

		}
	}
}

func Manager(crawlerWorkerChan chan string, remainingLinksChan chan int, foundLinksChan chan string, bufferSize int) {

	go countMonitor(crawlerWorkerChan, remainingLinksChan, foundLinksChan)
	go workerMonitor(crawlerWorkerChan, foundLinksChan, remainingLinksChan, bufferSize)

}
