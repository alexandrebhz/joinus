package crawler

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/startup-job-board/crawler/internal/domain/entity"
)

// Extractor handles extracting job data from HTML
type Extractor struct {
	rules entity.ExtractionRules
}

// NewExtractor creates a new Extractor
func NewExtractor(rules entity.ExtractionRules) *Extractor {
	return &Extractor{rules: rules}
}

// ExtractJobs extracts job listings from HTML document
func (e *Extractor) ExtractJobs(doc *goquery.Document, baseURL string) ([]*entity.CrawledJob, error) {
	var jobs []*entity.CrawledJob

	doc.Find(e.rules.JobListSelector).Each(func(i int, s *goquery.Selection) {
		job := e.extractJob(s, baseURL)
		if job != nil {
			jobs = append(jobs, job)
		}
	})

	return jobs, nil
}

// extractJob extracts a single job from a selection
func (e *Extractor) extractJob(selection *goquery.Selection, baseURL string) *entity.CrawledJob {
	job := &entity.CrawledJob{}

	// Extract detail URL
	detailURL := e.extractDetailURL(selection, baseURL)
	if detailURL == "" {
		return nil // Skip if no detail URL found
	}
	job.DetailURL = detailURL

	// Extract fields
	for fieldName, rule := range e.rules.Fields {
		value := e.extractField(selection, rule)
		value = e.applyTransformations(value, rule.Transformations)

		switch fieldName {
		case "title":
			job.Title = value
		case "description":
			job.Description = value
		case "requirements":
			job.Requirements = value
		case "company":
			job.Company = value
		case "location":
			job.Location = value
		case "city":
			job.City = value
		case "country":
			job.Country = value
		case "job_type":
			job.JobType = value
		case "location_type":
			job.LocationType = value
		case "salary_min":
			if val := e.parseInt(value); val != nil {
				job.SalaryMin = val
			}
		case "salary_max":
			if val := e.parseInt(value); val != nil {
				job.SalaryMax = val
			}
		case "currency":
			job.Currency = value
		case "application_url":
			if value != "" {
				job.ApplicationURL = &value
			}
		case "application_email":
			if value != "" {
				job.ApplicationEmail = &value
			}
		case "expires_at":
			if val := e.parseDate(value); val != nil {
				job.ExpiresAt = val
			}
		}
	}

	// Store raw HTML for debugging
	html, _ := selection.Html()
	job.RawHTML = html

	return job
}

// extractDetailURL extracts the job detail page URL
func (e *Extractor) extractDetailURL(selection *goquery.Selection, baseURL string) string {
	rule := e.rules.JobDetailURL
	link := selection.Find(rule.Selector).First()

	if link.Length() == 0 {
		return ""
	}

	var url string
	switch rule.Type {
	case "attribute":
		url, _ = link.Attr(rule.Attribute)
	case "relative":
		url, _ = link.Attr(rule.Attribute)
		if url != "" && !strings.HasPrefix(url, "http") {
			if rule.BaseURL != "" {
				url = rule.BaseURL + url
			} else {
				url = baseURL + url
			}
		}
	case "absolute":
		url, _ = link.Attr(rule.Attribute)
	default:
		url, _ = link.Attr("href")
	}

	return url
}

// extractField extracts a field value based on the rule
func (e *Extractor) extractField(selection *goquery.Selection, rule entity.FieldRule) string {
	elem := selection.Find(rule.Selector).First()
	if elem.Length() == 0 {
		return rule.DefaultValue
	}

	var value string
	switch rule.Type {
	case "text":
		value = strings.TrimSpace(elem.Text())
	case "html":
		html, _ := elem.Html()
		value = strings.TrimSpace(html)
	case "attribute":
		value, _ = elem.Attr(rule.Attribute)
	case "regex":
		text := elem.Text()
		if rule.RegexPattern != "" {
			re := regexp.MustCompile(rule.RegexPattern)
			matches := re.FindStringSubmatch(text)
			if len(matches) > 1 {
				value = matches[1]
			}
		}
	default:
		value = strings.TrimSpace(elem.Text())
	}

	if value == "" {
		return rule.DefaultValue
	}

	return value
}

// applyTransformations applies transformations to a value
func (e *Extractor) applyTransformations(value string, transformations []string) string {
	result := value
	for _, transform := range transformations {
		switch transform {
		case "trim":
			result = strings.TrimSpace(result)
		case "lowercase":
			result = strings.ToLower(result)
		case "uppercase":
			result = strings.ToUpper(result)
		case "strip_html":
			result = e.stripHTML(result)
		case "remove_commas":
			result = strings.ReplaceAll(result, ",", "")
		case "parse_int":
			// Already handled in parseInt
		case "parse_date":
			// Already handled in parseDate
		}
	}
	return result
}

// stripHTML removes HTML tags
func (e *Extractor) stripHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return html
	}
	return strings.TrimSpace(doc.Text())
}

// parseInt parses an integer from a string
func (e *Extractor) parseInt(s string) *int {
	// Remove commas and spaces
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, " ", "")
	s = strings.TrimSpace(s)

	// Extract numbers
	re := regexp.MustCompile(`\d+`)
	matches := re.FindString(s)
	if matches == "" {
		return nil
	}

	var val int
	fmt.Sscanf(matches, "%d", &val)
	return &val
}

// parseDate parses a date string (simplified - should use proper date parsing)
func (e *Extractor) parseDate(s string) *time.Time {
	// This is a simplified implementation
	// In production, use a proper date parser like time.Parse
	return nil
}

