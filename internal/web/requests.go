package web

import (
	"encoding/xml"
	"net/http"
	"context"
	"fmt"
	"html"
	"io"

)

type RSSFeed struct{
	Channel struct{
		Title        string    `xml:"title"`
		Link         string    `xml:"link"`
		Description  string    `xml:"description"`
		Item         []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title        string `xml:"title"`
	Link         string `xml:"link"`
	Description  string `xml:"description"`
	PubDate      string `xml:"item"`
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx,http.MethodGet,feedURL,nil)
	if err != nil{
		err = fmt.Errorf("Error creating new HTTP request: %v\n", err)
		return nil, err
	}

	req.Header.Add("User-Agent", "gator")

	resp, err := client.Do(req)

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)

	var newFeed RSSFeed
	err = xml.Unmarshal(data, &newFeed)
	if err !=nil{
		err = fmt.Errorf("Error unmarshalling response: %v", err)
		return nil, err
	}

	newFeed.Channel.Title = html.UnescapeString(newFeed.Channel.Title)
	newFeed.Channel.Description = html.UnescapeString(newFeed.Channel.Description)
	for i, item := range newFeed.Channel.Item{
		newFeed.Channel.Item[i].Title = html.UnescapeString(item.Title)
		newFeed.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}

	return &newFeed, nil
}

