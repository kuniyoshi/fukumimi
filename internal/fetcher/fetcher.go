package fetcher

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/kuniyoshi/fukumimi/internal/models"
)

const (
	episodesURL = "https://kitoakari-fc.com/special_contents/?category_id=4&page=%d"
	// UserAgent identifies this client
	UserAgent = "Fukumimi/0.1.0 (https://github.com/kuniyoshi/fukumimi)"
	// Cache directory for URL following results
	cacheDir = ".fukumimi-cache"
)

type cacheEntry struct {
	URL string `json:"url"`
}

type Fetcher struct {
	client      *http.Client
	IgnoreCache bool
}

func New() *Fetcher {
	return &Fetcher{
		client: &http.Client{},
	}
}

func (f *Fetcher) SetIgnoreCache(ignore bool) {
	f.IgnoreCache = ignore
}

func (f *Fetcher) FetchEpisodes() ([]models.Episode, error) {
	var allEpisodes []models.Episode
	page := 1
	maxPages := 100 // Safety limit to prevent infinite loops
	consecutiveEmptyPages := 0

	for page <= maxPages {
		episodes, hasMore, err := f.fetchPage(page)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch page %d: %w", page, err)
		}

		if len(episodes) == 0 {
			consecutiveEmptyPages++
			// If we get 3 consecutive empty pages, assume we've reached the end
			if consecutiveEmptyPages >= 3 {
				break
			}
		} else {
			consecutiveEmptyPages = 0
			allEpisodes = append(allEpisodes, episodes...)
		}

		// Always try the next page unless we explicitly know there's no more
		// This handles cases where pagination detection might fail
		if !hasMore && len(episodes) == 0 {
			break
		}

		page++
	}

	return allEpisodes, nil
}

func (f *Fetcher) fetchPage(page int) ([]models.Episode, bool, error) {
	url := fmt.Sprintf(episodesURL, page)

	// Create request with User-Agent
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, false, err
	}
	req.Header.Set("User-Agent", UserAgent)

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, false, err
	}

	var episodes []models.Episode

	// Look for all text containing STREAMING
	doc.Find("*").Each(func(i int, s *goquery.Selection) {
		// Get direct text content (not including children)
		text := s.Text()

		// Skip if this doesn't contain STREAMING
		if !strings.Contains(text, "STREAMING") {
			return
		}

		episode := models.Episode{}

		// Extract episode number (e.g., "#037")
		if numberMatch := regexp.MustCompile(`#(\d+)`).FindStringSubmatch(text); len(numberMatch) > 0 {
			episode.Number = numberMatch[0]
		}

		// Extract date (e.g., "6/11(水)")
		// Since we can't determine the year reliably, we'll use year 0 as a placeholder
		// The actual year doesn't matter for sorting by month/day
		if dateMatch := regexp.MustCompile(`(\d{1,2})/(\d{1,2})\([^)]+\)`).FindStringSubmatch(text); len(dateMatch) > 2 {
			month, _ := strconv.Atoi(dateMatch[1])
			day, _ := strconv.Atoi(dateMatch[2])

			if month > 0 && day > 0 {
				// Use year 0 as placeholder - we'll format without year in output
				episode.Date = time.Date(0, time.Month(month), day, 0, 0, 0, 0, time.UTC)
			}
		}

		// Try to find a link in this element or its children
		if link, exists := s.Find("a").Attr("href"); exists {
			// Extract contents_id and id from the URL
			if matches := regexp.MustCompile(`contents_id=(\d+)&id=(\d+)`).FindStringSubmatch(link); len(matches) > 2 {
				// Construct the proper URL format
				episode.URL = fmt.Sprintf("https://kitoakari-fc.com/special_contents/?contents_id=%s&id=%s", matches[1], matches[2])
			}
		} else if link, exists := s.Attr("href"); exists {
			// Extract contents_id and id from the URL
			if matches := regexp.MustCompile(`contents_id=(\d+)&id=(\d+)`).FindStringSubmatch(link); len(matches) > 2 {
				// Construct the proper URL format
				episode.URL = fmt.Sprintf("https://kitoakari-fc.com/special_contents/?contents_id=%s&id=%s", matches[1], matches[2])
			}
		}

		// Only add if we found an episode number and it's not a duplicate
		if episode.Number != "" && episode.URL != "" {
			// Check for duplicates
			isDuplicate := false
			for _, existing := range episodes {
				if existing.Number == episode.Number && existing.Date.Equal(episode.Date) {
					isDuplicate = true
					break
				}
			}
			if !isDuplicate {
				episodes = append(episodes, episode)
			}
		}
	})

	// Follow URLs concurrently to find streaming pages
	f.followURLsConcurrently(episodes)

	// Check if there's a next page - look for various pagination patterns
	hasMore := false
	nextPatterns := []string{
		".pagination .next",
		"a[rel='next']",
		".pager .next",
		"a:contains('次へ')",
		"a:contains('Next')",
		".page-link:contains('>')",
	}

	for _, pattern := range nextPatterns {
		if doc.Find(pattern).Length() > 0 {
			hasMore = true
			break
		}
	}

	// Also check if current page number exists in pagination
	if !hasMore && page == 1 {
		// Check if there's a page 2 link
		pageLinks := doc.Find("a[href*='page=2']")
		hasMore = pageLinks.Length() > 0
	}

	return episodes, hasMore, nil
}

func (f *Fetcher) followURL(url string) (string, error) {
	// Check cache first (unless ignoring cache)
	if !f.IgnoreCache {
		if cachedURL, found := f.getCachedURL(url); found {
			return cachedURL, nil
		}
	}

	// Create request with User-Agent
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", UserAgent)

	resp, err := f.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	// Look for links with text '配信ページはこちら'
	var streamingURL string
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		linkText := strings.TrimSpace(s.Text())
		if linkText == "配信ページはこちら" {
			if href, exists := s.Attr("href"); exists {
				// Make URL absolute if it's relative
				if strings.HasPrefix(href, "/") {
					streamingURL = "https://kitoakari-fc.com" + href
				} else if strings.HasPrefix(href, "http") {
					streamingURL = href
				} else {
					streamingURL = "https://kitoakari-fc.com/" + href
				}
			}
		}
	})

	// Cache the result
	if streamingURL != "" {
		f.cacheURL(url, streamingURL)
	}

	return streamingURL, nil
}

func (f *Fetcher) followURLsConcurrently(episodes []models.Episode) {
	const maxConcurrency = 10 // Limit concurrent requests to avoid overwhelming the server

	// Create a channel to limit concurrency
	semaphore := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup

	for i := range episodes {
		if episodes[i].URL != "" {
			wg.Add(1)
			go func(episodeIndex int) {
				defer wg.Done()

				// Acquire semaphore
				semaphore <- struct{}{}
				defer func() { <-semaphore }()

				streamingURL, err := f.followURL(episodes[episodeIndex].URL)
				if err == nil && streamingURL != "" {
					episodes[episodeIndex].URL = streamingURL
				}
			}(i)
		}
	}

	wg.Wait()
}

func (f *Fetcher) getCacheFilePath(url string) string {
	// Create a hash of the URL to use as filename
	hasher := sha256.New()
	hasher.Write([]byte(url))
	hash := fmt.Sprintf("%x", hasher.Sum(nil))
	return filepath.Join(cacheDir, hash+".json")
}

func (f *Fetcher) getCachedURL(url string) (string, bool) {
	cacheFile := f.getCacheFilePath(url)

	// Check if cache file exists
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return "", false
	}

	var entry cacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return "", false
	}

	return entry.URL, true
}

func (f *Fetcher) cacheURL(originalURL, streamingURL string) {
	// Ensure cache directory exists
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return // Silently fail if we can't create cache directory
	}

	entry := cacheEntry{
		URL: streamingURL,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return // Silently fail if we can't marshal
	}

	cacheFile := f.getCacheFilePath(originalURL)
	os.WriteFile(cacheFile, data, 0644) // Silently fail if we can't write
}

func (f *Fetcher) GenerateMarkdown(episodes []models.Episode) string {
	// Sort episodes by episode number (newest/highest first)
	sortedEpisodes := make([]models.Episode, len(episodes))
	copy(sortedEpisodes, episodes)

	// Extract number from episode number string (e.g., "#038" -> 38)
	getEpisodeNum := func(ep models.Episode) int {
		if match := regexp.MustCompile(`#(\d+)`).FindStringSubmatch(ep.Number); len(match) > 1 {
			num, _ := strconv.Atoi(match[1])
			return num
		}
		return 0
	}

	// Simple bubble sort by episode number (descending)
	for i := 0; i < len(sortedEpisodes)-1; i++ {
		for j := 0; j < len(sortedEpisodes)-i-1; j++ {
			if getEpisodeNum(sortedEpisodes[j]) < getEpisodeNum(sortedEpisodes[j+1]) {
				sortedEpisodes[j], sortedEpisodes[j+1] = sortedEpisodes[j+1], sortedEpisodes[j]
			}
		}
	}

	// Generate simple list format as per README
	var content strings.Builder
	for _, episode := range sortedEpisodes {
		content.WriteString(fmt.Sprintf("%s\n", episode.String()))
	}

	return content.String()
}
