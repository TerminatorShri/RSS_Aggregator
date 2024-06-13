package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/TerminatorShri/RSS_Aggregator/internal/database"
	"github.com/google/uuid"
)

func startScraping(db *database.Queries, concurrency int, timeBetweenReq time.Duration) {
	log.Printf("Scraping using %v goroutines very %s duration", concurrency, timeBetweenReq)
	ticker := time.NewTicker(timeBetweenReq)
	for ; ; <- ticker.C {
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)
		if err != nil {
			log.Printf("Error while fetching feeds: %v", err)
			continue
		}
		waitgrp := &sync.WaitGroup{}
		for _, feed := range feeds {
			waitgrp.Add(1)
			go scrapeFeed(db, waitgrp, feed)
		}
		waitgrp.Wait()
	}
}

func scrapeFeed(db *database.Queries, waitgrp *sync.WaitGroup, feed database.Feed) {
	defer waitgrp.Done()
	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Error marking as fetched: %v", err)
	}
	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Printf("Error fetchin feed: %v", err)
		return
	}
	for _, item := range rssFeed.Channel.Item {
		desc := sql.NullString {}
		if item.Description != "" {
			desc.String = item.Description
			desc.Valid = true		
		}
		publishTime, err := time.Parse(time.RFC1123Z, item.Publish_Date)
		if err != nil {
			log.Printf("Couldn't parse data %v Error: %v", item.Publish_Date, err)
		}
		_, err = db.CreatePost(context.Background(), 
			database.CreatePostParams{
				ID: uuid.New(),
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
				Title: item.Title,
				Description: desc,
				PublishedAt: publishTime,
				Url: item.Link,
				FeedID: feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Printf("Failed to Create Post: %v", err)
		}
	}
} 