package crawler

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/startup-job-board/crawler/internal/domain/entity"
)

// Paginator handles pagination through job listings
type Paginator struct {
	config entity.PaginationConfig
}

// NewPaginator creates a new Paginator
func NewPaginator(config entity.PaginationConfig) *Paginator {
	return &Paginator{config: config}
}

// GetPageURLs generates URLs for all pages to crawl
func (p *Paginator) GetPageURLs(baseURL string) ([]string, error) {
	var urls []string

	switch p.config.Type {
	case "query_param":
		urls = p.generateQueryParamURLs(baseURL)
	case "url_pattern":
		urls = p.generateURLPatternURLs(baseURL)
	case "link_follow":
		// This will be handled during crawling
		urls = []string{baseURL}
	case "api_pagination":
		urls = p.generateAPIURLs()
	default:
		urls = []string{baseURL}
	}

	return urls, nil
}

// HasNextPage checks if there's a next page based on the current document
func (p *Paginator) HasNextPage(doc *goquery.Document) bool {
	if p.config.Type != "link_follow" {
		return false
	}

	if p.config.NextPageSelector == "" {
		return false
	}

	return doc.Find(p.config.NextPageSelector).Length() > 0
}

// GetNextPageURL extracts the next page URL from the document
func (p *Paginator) GetNextPageURL(doc *goquery.Document, baseURL string) (string, error) {
	if p.config.Type != "link_follow" {
		return "", fmt.Errorf("next page URL extraction only supported for link_follow pagination")
	}

	link := doc.Find(p.config.NextPageSelector).First()
	if link.Length() == 0 {
		return "", fmt.Errorf("next page link not found")
	}

	href, exists := link.Attr("href")
	if !exists {
		return "", fmt.Errorf("next page link has no href attribute")
	}

	// Handle relative URLs
	if !strings.HasPrefix(href, "http") {
		base, err := url.Parse(baseURL)
		if err != nil {
			return "", err
		}
		nextURL, err := base.Parse(href)
		if err != nil {
			return "", err
		}
		return nextURL.String(), nil
	}

	return href, nil
}

// generateQueryParamURLs generates URLs with query parameters
func (p *Paginator) generateQueryParamURLs(baseURL string) []string {
	var urls []string
	paramName := p.config.ParamName
	if paramName == "" {
		paramName = "page"
	}

	startPage := p.config.StartPage
	if startPage == 0 {
		startPage = 1
	}

	increment := p.config.Increment
	if increment == 0 {
		increment = 1
	}

	maxPages := p.config.MaxPages
	if maxPages == 0 {
		maxPages = 100 // Default limit
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return []string{baseURL}
	}

	for i := 0; i < maxPages; i++ {
		page := startPage + (i * increment)
		query := base.Query()
		query.Set(paramName, strconv.Itoa(page))
		base.RawQuery = query.Encode()
		urls = append(urls, base.String())
	}

	return urls
}

// generateURLPatternURLs generates URLs using a pattern
func (p *Paginator) generateURLPatternURLs(baseURL string) []string {
	var urls []string

	startPage := p.config.StartPage
	if startPage == 0 {
		startPage = 1
	}

	increment := p.config.Increment
	if increment == 0 {
		increment = 1
	}

	maxPages := p.config.MaxPages
	if maxPages == 0 {
		maxPages = 100
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return []string{baseURL}
	}

	basePath := base.Path
	if !strings.HasSuffix(basePath, "/") {
		basePath += "/"
	}

	for i := 0; i < maxPages; i++ {
		page := startPage + (i * increment)
		pattern := strings.ReplaceAll(p.config.URLPattern, "{page}", strconv.Itoa(page))
		base.Path = basePath + strings.TrimPrefix(pattern, "/")
		urls = append(urls, base.String())
	}

	return urls
}

// generateAPIURLs generates API URLs for pagination
func (p *Paginator) generateAPIURLs() []string {
	if p.config.APIConfig == nil {
		return []string{}
	}

	var urls []string
	startPage := p.config.StartPage
	if startPage == 0 {
		startPage = 1
	}

	increment := p.config.Increment
	if increment == 0 {
		increment = 1
	}

	maxPages := p.config.APIConfig.MaxPages
	if maxPages == 0 {
		maxPages = 100
	}

	base, err := url.Parse(p.config.APIConfig.Endpoint)
	if err != nil {
		return []string{}
	}

	for i := 0; i < maxPages; i++ {
		page := startPage + (i * increment)
		query := base.Query()
		query.Set(p.config.APIConfig.PageParam, strconv.Itoa(page))
		if p.config.APIConfig.PageSize > 0 {
			query.Set("page_size", strconv.Itoa(p.config.APIConfig.PageSize))
		}
		base.RawQuery = query.Encode()
		urls = append(urls, base.String())
	}

	return urls
}

