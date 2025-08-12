package handlers

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Article represents a blog article
type Article struct {
	Title       string
	Slug        string
	Content     template.HTML
	Excerpt     string
	Author      string
	Date        time.Time
	Tags        []string
	ReadingTime int
}

// ArticleList represents a list of articles
type ArticleList struct {
	Articles []Article
	Total    int
}

// ArticlesHandler handles the articles listing page
func ArticlesHandler(w http.ResponseWriter, r *http.Request) {
	articles, err := loadArticles()
	if err != nil {
		log.Printf("Error loading articles: %v", err)
		http.Error(w, "Error loading articles", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"PageTitle":    "Blog Olho Urbano",
		"PageSubtitle": "Artigos, insights e histórias sobre cidades inteligentes",
		"Content":      "articles_content",
		"Articles":     articles,
		"Total":        len(articles),
	}

	if err := renderTemplate(w, "05_footer_pages.html", data); err != nil {
		log.Printf("Error rendering articles template: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// ArticleHandler handles individual article pages
func ArticleHandler(w http.ResponseWriter, r *http.Request) {
	// Extract slug from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}
	slug := pathParts[2]

	article, err := loadArticle(slug)
	if err != nil {
		log.Printf("Error loading article %s: %v", slug, err)
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}

	// Construct full URL for sharing
	scheme := "https"
	if r.TLS == nil {
		scheme = "http"
	}
	fullURL := fmt.Sprintf("%s://%s%s", scheme, r.Host, r.URL.Path)

	data := map[string]interface{}{
		"PageTitle":    article.Title,
		"PageSubtitle": "Artigo do Blog Olho Urbano",
		"Content":      "article_content",
		"Article":      article,
		"URL":          fullURL,
	}

	if err := renderTemplate(w, "05_footer_pages.html", data); err != nil {
		log.Printf("Error rendering article template: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// loadArticles loads all articles from the articles directory
func loadArticles() ([]Article, error) {
	articlesDir := "./articles"
	files, err := ioutil.ReadDir(articlesDir)
	if err != nil {
		if os.IsNotExist(err) {
			// Articles directory doesn't exist yet, return empty list
			return []Article{}, nil
		}
		return nil, err
	}

	var articles []Article
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".md" {
			article, err := loadArticleFromFile(filepath.Join(articlesDir, file.Name()))
			if err != nil {
				log.Printf("Error loading article %s: %v", file.Name(), err)
				continue
			}
			articles = append(articles, article)
		}
	}

	// Sort articles by date (newest first)
	// This is a simple implementation - in production you might want more sophisticated sorting
	return articles, nil
}

// loadArticle loads a specific article by slug
func loadArticle(slug string) (Article, error) {
	articlesDir := "./articles"
	filePath := filepath.Join(articlesDir, slug+".md")

	article, err := loadArticleFromFile(filePath)
	if err != nil {
		return Article{}, err
	}

	return article, nil
}

// loadArticleFromFile loads an article from a markdown file
func loadArticleFromFile(filePath string) (Article, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return Article{}, err
	}

	// Parse front matter and content
	// This is a simple implementation - in production you might want to use a proper front matter parser
	lines := strings.Split(string(content), "\n")

	article := Article{
		Slug: strings.TrimSuffix(filepath.Base(filePath), ".md"),
	}

	var inFrontMatter bool
	var contentLines []string

	for _, line := range lines {
		if strings.TrimSpace(line) == "---" {
			if !inFrontMatter {
				inFrontMatter = true
				continue
			} else {
				inFrontMatter = false
				continue
			}
		}

		if inFrontMatter {
			// Parse front matter
			if strings.HasPrefix(line, "title:") {
				article.Title = strings.TrimSpace(strings.TrimPrefix(line, "title:"))
			} else if strings.HasPrefix(line, "author:") {
				article.Author = strings.TrimSpace(strings.TrimPrefix(line, "author:"))
			} else if strings.HasPrefix(line, "date:") {
				dateStr := strings.TrimSpace(strings.TrimPrefix(line, "date:"))
				if date, err := time.Parse("2006-01-02", dateStr); err == nil {
					article.Date = date
				}
			} else if strings.HasPrefix(line, "tags:") {
				tagsStr := strings.TrimSpace(strings.TrimPrefix(line, "tags:"))
				article.Tags = strings.Split(tagsStr, ",")
				for i, tag := range article.Tags {
					article.Tags[i] = strings.TrimSpace(tag)
				}
			} else if strings.HasPrefix(line, "excerpt:") {
				article.Excerpt = strings.TrimSpace(strings.TrimPrefix(line, "excerpt:"))
			}
		} else {
			// This is content
			contentLines = append(contentLines, line)
		}
	}

	// Convert markdown to HTML (simple implementation)
	htmlContent := convertMarkdownToHTML(strings.Join(contentLines, "\n"))
	article.Content = template.HTML(htmlContent)

	// Calculate reading time (rough estimate: 200 words per minute)
	wordCount := len(strings.Fields(string(article.Content)))
	article.ReadingTime = (wordCount + 199) / 200 // Round up

	// If no excerpt was provided, generate one from content
	if article.Excerpt == "" {
		article.Excerpt = generateExcerpt(string(article.Content), 150)
	}

	return article, nil
}

// convertMarkdownToHTML converts markdown to HTML
// This is a simple implementation - in production you might want to use a proper markdown parser
func convertMarkdownToHTML(markdown string) string {
	// Split into lines for better processing
	lines := strings.Split(markdown, "\n")
	var htmlLines []string

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Skip empty lines
		if trimmedLine == "" {
			htmlLines = append(htmlLines, "")
			continue
		}

		// Process headers
		if strings.HasPrefix(trimmedLine, "# ") {
			htmlLines = append(htmlLines, "<h1>"+strings.TrimPrefix(trimmedLine, "# ")+"</h1>")
		} else if strings.HasPrefix(trimmedLine, "## ") {
			htmlLines = append(htmlLines, "<h2>"+strings.TrimPrefix(trimmedLine, "## ")+"</h2>")
		} else if strings.HasPrefix(trimmedLine, "### ") {
			htmlLines = append(htmlLines, "<h3>"+strings.TrimPrefix(trimmedLine, "### ")+"</h3>")
		} else if strings.HasPrefix(trimmedLine, "#### ") {
			htmlLines = append(htmlLines, "<h4>"+strings.TrimPrefix(trimmedLine, "#### ")+"</h4>")
		} else if strings.HasPrefix(trimmedLine, "- ") {
			// Process list items
			content := strings.TrimPrefix(trimmedLine, "- ")
			// Process bold and italic within list items
			content = processInlineMarkdown(content)
			htmlLines = append(htmlLines, "<li>"+content+"</li>")
		} else {
			// Regular paragraph content
			content := processInlineMarkdown(trimmedLine)
			htmlLines = append(htmlLines, "<p>"+content+"</p>")
		}
	}

	// Join lines and handle list wrapping
	html := strings.Join(htmlLines, "\n")

	// Wrap consecutive list items in <ul> tags
	html = wrapListItems(html)

	return html
}

// processInlineMarkdown processes inline markdown elements like bold, italic, and links
func processInlineMarkdown(content string) string {
	// Bold
	content = strings.ReplaceAll(content, "**", "<strong>")
	content = strings.ReplaceAll(content, "__", "<strong>")

	// Italic
	content = strings.ReplaceAll(content, "*", "<em>")
	content = strings.ReplaceAll(content, "_", "<em>")

	// Links (basic implementation)
	// This is a very basic implementation - you might want to use regex for better parsing
	content = strings.ReplaceAll(content, "[", "<a href=\"")
	content = strings.ReplaceAll(content, "](", "\">")
	content = strings.ReplaceAll(content, ")", "</a>")

	return content
}

// wrapListItems wraps consecutive <li> elements in <ul> tags
func wrapListItems(html string) string {
	lines := strings.Split(html, "\n")
	var result []string
	var inList bool

	for _, line := range lines {
		if strings.Contains(line, "<li>") {
			if !inList {
				result = append(result, "<ul>")
				inList = true
			}
			result = append(result, line)
		} else {
			if inList {
				result = append(result, "</ul>")
				inList = false
			}
			result = append(result, line)
		}
	}

	// Close any open list
	if inList {
		result = append(result, "</ul>")
	}

	return strings.Join(result, "\n")
}

// generateExcerpt generates an excerpt from content
func generateExcerpt(content string, maxLength int) string {
	// Remove HTML tags for excerpt
	plainText := content
	// Simple HTML tag removal - in production you might want to use a proper HTML parser
	plainText = strings.ReplaceAll(plainText, "<p>", "")
	plainText = strings.ReplaceAll(plainText, "</p>", " ")
	plainText = strings.ReplaceAll(plainText, "<h1>", "")
	plainText = strings.ReplaceAll(plainText, "</h1>", " ")
	plainText = strings.ReplaceAll(plainText, "<h2>", "")
	plainText = strings.ReplaceAll(plainText, "</h2>", " ")
	plainText = strings.ReplaceAll(plainText, "<h3>", "")
	plainText = strings.ReplaceAll(plainText, "</h3>", " ")
	plainText = strings.ReplaceAll(plainText, "<strong>", "")
	plainText = strings.ReplaceAll(plainText, "</strong>", "")
	plainText = strings.ReplaceAll(plainText, "<em>", "")
	plainText = strings.ReplaceAll(plainText, "</em>", "")

	if len(plainText) <= maxLength {
		return plainText
	}

	// Truncate and add ellipsis
	truncated := plainText[:maxLength]
	lastSpace := strings.LastIndex(truncated, " ")
	if lastSpace > 0 {
		truncated = truncated[:lastSpace]
	}
	return truncated + "..."
}
