package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	FeedURL     = "https://ansuman-satapathy.github.io/index.xml"
	ReadmePath  = "README.md"
	StartMarker = "<!-- BLOG-START -->"
	EndMarker   = "<!-- BLOG-END -->"
)

type RSS struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Items []Item `xml:"item"`
}

type Item struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
}

func main() {
	fmt.Println("1. Fetching RSS Feed...")
	items, err := fetchFeed(FeedURL)
	if err != nil {
		fmt.Printf("Error fetching feed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("2. Generating Markdown...")
	newContent := generateMarkdown(items)

	fmt.Println("3. Updating README...")
	if err := updateReadme(ReadmePath, newContent); err != nil {
		fmt.Printf("Error updating README: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Success! README updated.")
}

func fetchFeed(url string) ([]Item, error) {
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rss RSS
	if err := xml.Unmarshal(data, &rss); err != nil {
		return nil, err
	}

	limit := 5
	if len(rss.Channel.Items) < limit {
		limit = len(rss.Channel.Items)
	}
	return rss.Channel.Items[:limit], nil
}

func generateMarkdown(items []Item) string {
	var sb strings.Builder

	sb.WriteString(StartMarker + "\n")

	for _, item := range items {
		sb.WriteString(fmt.Sprintf("- [%s](%s)\n", item.Title, item.Link))
	}

	sb.WriteString(EndMarker)
	return sb.String()
}

func updateReadme(path, newContent string) error {
	contentBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	content := string(contentBytes)

	startIndex := strings.Index(content, StartMarker)
	endIndex := strings.Index(content, EndMarker)

	if startIndex == -1 || endIndex == -1 {
		return fmt.Errorf("markers not found in README.md")
	}

	finalContent := content[:startIndex] + newContent + content[endIndex+len(EndMarker):]

	return os.WriteFile(path, []byte(finalContent), 0644)
}
