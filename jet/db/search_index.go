package db

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

// SearchIndex represents a searchable keyword and its associated URLs.
type SearchIndex struct {
	ID        string         `gorm:"type:uuid;default:uuid_generate_v4()"` // Unique ID generated using uuid v4.
	Value     string         // The keyword or term to be indexed.
	Urls      []CrawledUrl   `gorm:"many2many:token_urls;"` // Many-to-many relationship with CrawledUrl via join table "token_urls".
	CreatedAt time.Time      `gorm:"autoCreateTime"`        // Timestamp for when the record was created.
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`        // Timestamp for when the record was last updated.
	DeletedAt gorm.DeletedAt `gorm:"index"`                 // Soft delete field with index.
}

// TableName overrides the default table name for SearchIndex.
func (s *SearchIndex) TableName() string {
	return "search_index"
}

// Save stores or updates search terms and associates them with crawled URLs.
func (s *SearchIndex) Save(index map[string][]string, crawledUrls []CrawledUrl) error {
	for value, ids := range index {
		// Create or retrieve the existing SearchIndex for each term (value).
		newIndex := &SearchIndex{
			Value: value,
		}
		if err := DBConn.Where(SearchIndex{Value: value}).FirstOrCreate(newIndex).Error; err != nil {
			return err
		}

		// Find the CrawledUrl objects that match the given IDs.
		var urlsToAppend []CrawledUrl
		for _, id := range ids {
			for _, url := range crawledUrls {
				if url.ID == id {
					urlsToAppend = append(urlsToAppend, url)
					break
				}
			}
		}

		// Associate the found URLs with the current search index entry.
		if err := DBConn.Model(&newIndex).Association("Urls").Append(&urlsToAppend); err != nil {
			return err
		}
	}
	return nil
}

// FullTextSearch searches for terms and returns associated crawled URLs.
func (s *SearchIndex) FullTextSearch(value string) ([]CrawledUrl, error) {
	// Split the search input into individual terms.
	terms := strings.Fields(value)
	var urls []CrawledUrl

	for _, term := range terms {
		var searchIndexes []SearchIndex
		// Look for search index entries where the value contains the term.
		if err := DBConn.Preload("Urls").Where("value LIKE ?", "%"+term+"%").Find(&searchIndexes).Error; err != nil {
			return nil, err
		}

		// Collect all associated URLs from the matching search indexes.
		for _, searchIndex := range searchIndexes {
			urls = append(urls, searchIndex.Urls...)
		}
	}
	return urls, nil
}
