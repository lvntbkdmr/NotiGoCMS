package main

import (
	"html"
	"path/filepath"
	"strings"

	"github.com/kjk/notionapi"
	"github.com/kjk/notionapi/tomarkdown"
	"github.com/kjk/notionapi/tohtml"
	"github.com/gosimple/slug"
)

//Chroma highlighting option map for Notion codeblock type
var codeLanguageMap = map[string]string {
	"C++": "cpp",
	"Bash": "bash",
	"Shell": "shell",
	"Python": "python",
}

// Converter renders article as html
type Converter struct {
	article      *Article
	page         *notionapi.Page
	notionClient *notionapi.Client
	idToArticle  func(string) *Article
	galleries    [][]string

	r *tomarkdown.Converter
	h *tohtml.Converter
}

func (c *Converter) maybeGetID(block *notionapi.Block) string {
	return notionapi.ToNoDashID(block.ID)
}

// change https://www.notion.so/Advanced-web-spidering-with-Puppeteer-ea07db1b9bff415ab180b0525f3898f6
// =>
// /article/${id}
func (c *Converter) rewriteURL(uri string) string {
	id := notionapi.ExtractNoDashIDFromNotionURL(uri)
	if id == "" {
		return uri
	}
	article := c.idToArticle(id)
	// this might happen when I link to some-one else's public notion pages
	if article == nil {
		return uri
	}
	return article.URL()
}

func (c *Converter) getURLAndTitleForBlock(block *notionapi.Block) (string, string) {
	id := notionapi.ToNoDashID(block.ID)
	article := c.idToArticle(id)
	if article == nil {
		title := block.Title
		logf("No article for id %s %s\n", id, title)
		pageURL := "https://notion.so/" + notionapi.ToNoDashID(c.page.ID)
		logf("Link from page: %s\n", pageURL)
		url := "/posts/" + slug.Make(title)
		return url, title
	}

	return article.URL(), article.Title
}

// RenderImage renders BlockImage
func (c *Converter) RenderImage(block *notionapi.Block) bool {
	link := block.Source
	im := findImageMapping(c.article.Images, link)
	imName := filepath.Base(im.relativeURL)
	relURL := im.relativeURL
	imgURL := c.article.getImageBlockURL(block)
	if imgURL != "" {
		c.r.Printf(`![%s](%s)`, imName, relURL)
	} else {
		c.r.Printf(`![%s](%s)`, imName, relURL)
	}
	return true
}

// RenderBlockCode renders BlockCode
func (c *Converter) RenderBlockCode(block *notionapi.Block) bool {

	code := block.Code
	language := block.CodeLanguage
	start := "{{< highlight " + codeLanguageMap[language] + " >}}"
	end := "{{< / highlight >}}"

	c.r.Printf(start + "\n")

	parts := strings.Split(code, "\n")
	for _, part := range parts {
		c.r.Printf(part + "\n")
	}

	c.r.Printf(end + "\n")

	return true
}

// RenderBlockCallout renders BlockCallout
func (c *Converter) RenderBlockCallout(block *notionapi.Block) bool {

	c.r.Printf("> ")
	c.r.RenderInlines(block.InlineContent, false)

	return true
}

// RenderPage renders BlockPage
func (c *Converter) RenderPage(block *notionapi.Block) bool {
	if c.r.Page.IsRoot(block) {
		c.r.RenderChildren(block)
		return true
	}

	url, title := c.getURLAndTitleForBlock(block)
	title = html.EscapeString(title)
	c.r.Printf("[%s](%s)\n", url, title)
	
	return true
}

// if returns false, the block will be rendered with default
func (c *Converter) blockRenderOverride(block *notionapi.Block) bool {
	if c.article.shouldSkipBlock(block) {
		return true
	}
	switch block.Type {
	case notionapi.BlockPage:
		return c.RenderPage(block)
	case notionapi.BlockImage:
		return c.RenderImage(block)
	case notionapi.BlockCode:
		return c.RenderBlockCode(block)
	case notionapi.BlockCallout:
		return c.RenderBlockCallout(block)
	}
	return false
}

// NewMarkdownConverter returns new HTMLGenerator
func NewMarkdownConverter(c *notionapi.Client, article *Article) *Converter {
	res := &Converter{
		notionClient: c,
		article:      article,
		page:         article.page,
	}

	r := tomarkdown.NewConverter(article.page)
	notionapi.PanicOnFailures = true
	r.RenderBlockOverride = res.blockRenderOverride
	r.RewriteURL = res.rewriteURL
	res.r = r

	return res
}

// Gen returns generated Markdown
func (c *Converter) GenereateMarkdown() []byte {
	inner := c.r.ToMarkdown()
	s := string(inner)
	return []byte(s)
}

func notionToMarkdown(client *notionapi.Client, article *Article, articles *Articles) ([]byte, []*ImageMapping) {
	//logf("notionToMarkdown: %s\n", notionapi.ToNoDashID(article.ID))
	c := NewMarkdownConverter(client, article)
	if articles != nil {
		c.idToArticle = func(id string) *Article {
			return articles.idToArticle[id]
		}
	}
	return c.GenereateMarkdown(), c.article.Images
}