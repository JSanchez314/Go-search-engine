package search

import (
	"fmt"
	"pro/jet/db"
	"time"
)

func RunEngine() {
	fmt.Println("started search engine crawl...")
	defer fmt.Println("search engine crawl has finished")

	settings := &db.SearchSetting{}
	err := settings.Get()
	if err != nil {
		fmt.Println("something went wrong getting the seattings")
		return
	}

	if !settings.SearchOn {
		fmt.Println("search is turned off")
		return
	}
	crawl := &db.CrawledUrl{}
	nextUrls, err := crawl.GetNextCrawlUrls(int(settings.Amount))
	if err != nil {
		fmt.Println("somethin went wrong getting next urls")
		return
	}

	newUrls := []db.CrawledUrl{}
	testedTime := time.Now()
	for _, next := range nextUrls {
		result := runCrawl(next.Url)
		if !result.Success {
			err := next.UpdateUrl(db.CrawledUrl{
				ID:              next.ID,
				Url:             next.Url,
				Success:         false,
				CrawlDuration:   result.CrawlData.CrawlTime,
				ResponseCode:    result.ResponseCode,
				PageTitle:       result.CrawlData.PageDescription,
				PageDescription: result.CrawlData.PageDescription,
				Heading:         result.CrawlData.Headings,
				LastTested:      &testedTime,
			})
			if err != nil {
				fmt.Println("something went wrong updating a failed url")
			}
			continue
		}
		// Success
		err := next.UpdateUrl(db.CrawledUrl{
			ID:              next.ID,
			Url:             next.Url,
			Success:         result.Success,
			CrawlDuration:   result.CrawlData.CrawlTime,
			ResponseCode:    result.ResponseCode,
			PageTitle:       result.CrawlData.PageDescription,
			PageDescription: result.CrawlData.PageDescription,
			Heading:         result.CrawlData.Headings,
			LastTested:      &testedTime,
		})
		if err != nil {
			fmt.Println("something went wrong updating a success url")
			fmt.Println(next.Url)
		}
		for _, newUrl := range result.CrawlData.Links.External {
			newUrls = append(newUrls, db.CrawledUrl{Url: newUrl})
		}
	} // End of range
	if !settings.AddNew {
		return
	}
	// Insert new urls
	for _, newUrl := range newUrls {
		err := newUrl.Save()
		if err != nil {
			fmt.Println("something went wrong adding the new url to the database")
		}
	}
	fmt.Printf("\n Added %d added new urls to the database", len(newUrls))
}

func RunIndex() {
	fmt.Println("started search indexing...")
	defer fmt.Println("search indexing has finished")
	crawled := &db.CrawledUrl{}
	notIndexed, err := crawled.GetNotIndex()
	if err != nil {
		return
	}

	idx := make(Index)
	idx.Add(notIndexed)
	searchIndex := &db.SearchIndex{}
	err = searchIndex.Save(idx, notIndexed)
	if err != nil {
		return
	}
	err = crawled.SetIndexTrue(notIndexed)
	if err != nil {
		return
	}
}
