package seventv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	ApiUrl *url.URL
}

func New(apiUrl *url.URL) *Client {
	return &Client{ApiUrl: apiUrl}
}

func (c *Client) GetEmoteUrlById(id string, format FileFormat, scaler FileScaler) (string, error) {
	emote, err := c.getEmoteById(id)
	if err != nil {
		return "", err
	}

	return extractEmoteUrl(emote, format, scaler)
}

func (c *Client) GetEmoteURLByQuery(query string, format FileFormat, scaler FileScaler) (string, error) {
	emotes, err := c.getEmotesByQuery(query)
	if err != nil {
		return "", err
	}

	if len(emotes) > 0 {
		// TODO: Need smarter way to select emote for now I assume the first result is the best fit
		return c.GetEmoteUrlById(emotes[0].ID, format, scaler)
	}

	return "", fmt.Errorf("no emotes found")
}

func (c *Client) getEmoteById(id string) (*Emote, error) {
	resp, err := http.Get(c.ApiUrl.String() + "/emotes/" + id)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status returned by 7tv: %s", resp.Status)
	}

	var emote Emote
	if err := json.Unmarshal(bodyBytes, &emote); err != nil {
		return nil, err
	}

	return &emote, nil
}

func (c *Client) getEmotesByQuery(query string) ([]GQLEmote, error) {
	url := c.ApiUrl.String() + "/gql"
	queries := map[string]string{
		"all": `query SearchEmotes($query: String!, $page: Int, $sort: Sort, $limit: Int, $filter: EmoteSearchFilter) {
            emotes(query: $query, page: $page, sort: $sort, limit: $limit, filter: $filter) {
                items {
                    id
                    name
                    owner {
                        username
                    }
                    host {
                        url
                    }
                }
            }
        }`,
		"url": `query SearchEmotes($query: String!, $page: Int, $sort: Sort, $limit: Int, $filter: EmoteSearchFilter) {
            emotes(query: $query, page: $page, sort: $sort, limit: $limit, filter: $filter) {
                items {
                    host {
                        url
                    }
                }
            }
        }`,
	}

	payload := map[string]interface{}{
		"operationName": "SearchEmotes",
		"variables": map[string]interface{}{
			"query": query,
			"limit": 12,
			"page":  1,
			"sort": map[string]string{
				"value": "popularity",
				"order": "DESCENDING",
			},
			"filter": map[string]interface{}{
				"category":       "TOP",
				"exact_match":    false,
				"case_sensitive": false,
				"ignore_tags":    false,
				"zero_width":     false,
				"animated":       false,
				"aspect_ratio":   "",
			},
		},
		"query": queries["all"],
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var emotes GQLResponse
	err = json.Unmarshal(bodyBytes, &emotes)
	if err != nil {
		return nil, err
	}

	return emotes.Data.Emotes.Items, nil
}

func extractEmoteUrl(emote *Emote, format FileFormat, scaler FileScaler) (string, error) {
	emoteBaseUrl := fmt.Sprintf("https:%s", emote.Host.URL)
	for _, file := range emote.Host.Files {
		if file.Format == FileFormatMap[format] {
			if strings.HasPrefix(file.Name, FileScalerMap[scaler]) {
				return emoteBaseUrl + "/" + file.Name, nil
			}
		}
	}
	return "", fmt.Errorf("no file found for format %s and scaling %s", FileFormatMap[format], FileScalerMap[scaler])
}
