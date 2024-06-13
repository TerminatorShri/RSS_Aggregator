package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"desc"`
		Language    string    `xml:"language"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title        string `xml:"title"`
	Link         string `xml:"link"`
	Description  string `xml:"desc"`
	Publish_Date string `xml:"publish_date"`
}

func urlToFeed(url string) (RSSFeed, error) {
	httpClient := http.Client {
		Timeout: 10 * time.Second,
	}
	res, err := httpClient.Get(url)
	if err != nil {
		return RSSFeed{}, err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return RSSFeed{}, err
	}
	rss_feed := RSSFeed{}
	err = xml.Unmarshal(data, &rss_feed)
	if err != nil {
		return RSSFeed{}, err
	}
	return rss_feed, nil
}