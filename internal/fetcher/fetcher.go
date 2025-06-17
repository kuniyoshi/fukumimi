package fetcher

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/kuniyoshi/fukumimi/internal/auth"
	"github.com/kuniyoshi/fukumimi/internal/models"
)

const (
	episodesURL = "https://kitoakari-fc.com/special_contents/?category_id=4&page=%d"
)

type Fetcher struct {
	client *auth.Client
}

func New() *Fetcher {
	return &Fetcher{
		client: auth.NewClient(),
	}
}

func (f *Fetcher) FetchEpisodes() ([]models.Episode, error) {
	// First, ensure we're logged in
	if err := f.client.Login(); err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

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
	
	resp, err := f.client.Get(url)
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
	var currentYear int

	// Silent operation - no debug output

	// First, try to find the current year from any date on the page
	doc.Find("*").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if yearMatch := regexp.MustCompile(`(\d{4})\.\d{2}\.\d{2}`).FindStringSubmatch(text); len(yearMatch) > 1 {
			if year, err := strconv.Atoi(yearMatch[1]); err == nil && currentYear == 0 {
				currentYear = year
			}
		}
	})

	// If no year found, use current year
	if currentYear == 0 {
		currentYear = time.Now().Year()
	}


	// Look for all text containing STREAMING
	doc.Find("*").Each(func(i int, s *goquery.Selection) {
		// Get direct text content (not including children)
		text := s.Text()
		
		// Skip if this doesn't contain STREAMING
		if !strings.Contains(text, "STREAMING") {
			return
		}


		episode := models.Episode{
			Year: currentYear,
		}

		// Extract episode number (e.g., "#037")
		if numberMatch := regexp.MustCompile(`#(\d+)`).FindStringSubmatch(text); len(numberMatch) > 0 {
			episode.Number = numberMatch[0]
		}

		// Extract date (e.g., "6/11(水)")
		if dateMatch := regexp.MustCompile(`(\d{1,2})/(\d{1,2})\([^)]+\)`).FindStringSubmatch(text); len(dateMatch) > 2 {
			month, _ := strconv.Atoi(dateMatch[1])
			day, _ := strconv.Atoi(dateMatch[2])
			
			if month > 0 && day > 0 {
				episode.Date = time.Date(currentYear, time.Month(month), day, 0, 0, 0, 0, time.Local)
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

func (f *Fetcher) GenerateMarkdown(episodes []models.Episode) string {
	// Sort episodes by date (newest first)
	sortedEpisodes := make([]models.Episode, len(episodes))
	copy(sortedEpisodes, episodes)
	
	// Simple bubble sort for now
	for i := 0; i < len(sortedEpisodes)-1; i++ {
		for j := 0; j < len(sortedEpisodes)-i-1; j++ {
			if sortedEpisodes[j].Date.Before(sortedEpisodes[j+1].Date) {
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