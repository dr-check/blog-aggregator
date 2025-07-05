package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	req.Header.Set("User-Agent", "gator")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading response body: %v\n", err)
		return nil, err
	}

	var rFeed RSSFeed
	err = xml.Unmarshal(data, &rFeed)
	if err != nil {
		return nil, err
	}

	rFeed.Channel.Title = html.UnescapeString(rFeed.Channel.Title)
	rFeed.Channel.Description = html.UnescapeString(rFeed.Channel.Description)

	for i := 0; i < len(rFeed.Channel.Item); i++ {
		rFeed.Channel.Item[i].Title = html.UnescapeString(rFeed.Channel.Item[i].Title)
		rFeed.Channel.Item[i].Description = html.UnescapeString(rFeed.Channel.Item[i].Description)
	}
	return &rFeed, nil
}
