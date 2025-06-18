package merger

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/kuniyoshi/fukumimi/internal/models"
)

type Merger struct{}

func New() *Merger {
	return &Merger{}
}

// ParseEpisodeLine parses a line in the format:
// [ ] 06/11 [#38](https://kitoakari-fc.com/special_contents/?contents_id=1&id=55)
// [x] 05/28 [#37](https://kitoakari-fc.com/special_contents/?contents_id=1&id=54)
func (m *Merger) ParseEpisodeLine(line string) (*models.Episode, error) {
	// Skip empty lines
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, nil
	}

	// Check if listened
	isListened := strings.HasPrefix(line, "[x]")
	
	// Extract date (MM/DD)
	dateMatch := regexp.MustCompile(`\[.\]\s+(\d{2})/(\d{2})`).FindStringSubmatch(line)
	if len(dateMatch) < 3 {
		return nil, fmt.Errorf("invalid episode format: %s", line)
	}
	
	month, _ := strconv.Atoi(dateMatch[1])
	day, _ := strconv.Atoi(dateMatch[2])
	
	// Extract episode number
	numberMatch := regexp.MustCompile(`\[(#\d+)\]`).FindStringSubmatch(line)
	if len(numberMatch) < 2 {
		return nil, fmt.Errorf("episode number not found: %s", line)
	}
	
	// Extract URL
	urlMatch := regexp.MustCompile(`\((https?://[^)]+)\)`).FindStringSubmatch(line)
	url := ""
	if len(urlMatch) >= 2 {
		url = urlMatch[1]
	}
	
	return &models.Episode{
		Number:     numberMatch[1],
		Date:       time.Date(0, time.Month(month), day, 0, 0, 0, 0, time.UTC),
		URL:        url,
		IsListened: isListened,
	}, nil
}

func (m *Merger) ReadEpisodesFromReader(r io.Reader) ([]models.Episode, error) {
	var episodes []models.Episode
	scanner := bufio.NewScanner(r)
	
	for scanner.Scan() {
		episode, err := m.ParseEpisodeLine(scanner.Text())
		if err != nil {
			// Skip invalid lines
			continue
		}
		if episode != nil {
			episodes = append(episodes, *episode)
		}
	}
	
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	
	return episodes, nil
}

func (m *Merger) ReadEpisodesFromStdin() ([]models.Episode, error) {
	return m.ReadEpisodesFromReader(os.Stdin)
}

func (m *Merger) ReadEpisodesFromFile(filename string) ([]models.Episode, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	return m.ReadEpisodesFromReader(file)
}

func (m *Merger) MergeEpisodes(newEpisodes, localEpisodes []models.Episode) []models.Episode {
	// Create a map of local episodes by number for quick lookup
	listenedMap := make(map[string]bool)
	for _, ep := range localEpisodes {
		if ep.IsListened {
			listenedMap[ep.Number] = true
		}
	}
	
	// Update new episodes with listened status
	merged := make([]models.Episode, len(newEpisodes))
	for i, ep := range newEpisodes {
		merged[i] = ep
		if listenedMap[ep.Number] {
			merged[i].IsListened = true
		}
	}
	
	return merged
}

func (m *Merger) GenerateOutput(episodes []models.Episode) string {
	var output strings.Builder
	for _, ep := range episodes {
		output.WriteString(ep.String() + "\n")
	}
	return output.String()
}