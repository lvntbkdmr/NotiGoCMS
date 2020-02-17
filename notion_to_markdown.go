package main

import (
	"html"
	"strconv"
	"strings"
	"path/filepath"

	"github.com/kjk/notionapi"
	"github.com/kjk/notionapi/tomarkdown"
	"github.com/gosimple/slug"
)

// Converter renders article as html
type Converter struct {
	article      *Article
	page         *notionapi.Page
	notionClient *notionapi.Client
	idToArticle  func(string) *Article
	galleries    [][]string

	r *tomarkdown.Converter
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

func genGalleryMainHTML(galleryID int, imageURL string) string {
	s := `
  <div class="img-wrapper-wrapper">
    <div class="img-wrapper">
      <img id="id-gallery-{galleryID}" src="{imageURL}" />
      <a class="for-nav-icon nav-icon-left" href="#" onclick="imgPrev("{galleryID}"); return false;">
        <svg viewBox="0 0 24 24" preserveAspectRatio="xMidYMid meet" focusable="false" class="nav-icon">
          <g>
            <path d="M15.41 7.41L14 6l-6 6 6 6 1.41-1.41L10.83 12z" class="style-scope yt-icon">
            </path>
          </g>
        </svg>
      </a>
      <a class="for-nav-icon nav-icon-right" href="#" onclick="imgNext({galleryID}); return false;">
        <svg viewBox="0 0 24 24" preserveAspectRatio="xMidYMid meet" focusable="false" class="nav-icon" style="">
          <g>
            <path d="M10 6L8.59 7.41 13.17 12l-4.58 4.59L10 18l6-6z" class="yt-icon"></path>
          </g>
        </svg>
      </a>
    </div>
  </div>
`
	s = strings.Replace(s, "{galleryID}", strconv.Itoa(galleryID), -1)
	s = strings.Replace(s, "{imageURL}", imageURL, -1)
	return s
}

func genGalleryThumbHTML(galleryID int, n int, im *ImageMapping) string {
	s := `
    <div id="id-thumb-{galleryID}-{imageNo}" class="pa1 ib">
      <a href="#" onclick="changeShot({galleryID}, {imageNo}); return false;">
        <img id="id-thumb-img-{galleryID}-{imageNo}" src="{imageURL}" width="80" height="60" />
      </a>
	</div>
`
	s = strings.Replace(s, "{galleryID}", strconv.Itoa(galleryID), -1)
	ns := strconv.Itoa(n)
	s = strings.Replace(s, "{imageNo}", ns, -1)
	s = strings.Replace(s, "{imageURL}", im.relativeURL, -1)
	return s
}

func (c *Converter) renderGallery(block *notionapi.Block) bool {
	imageURLS := c.article.getGalleryImages(block)
	if len(imageURLS) == 0 {
		return false
	}
	panicIf(len(imageURLS) < 2, "expected gallery to have at least 2 images, got %d", len(imageURLS))
	galleryID := len(c.galleries)
	c.galleries = append(c.galleries, imageURLS)
	var images []*ImageMapping
	for _, link := range imageURLS {
		im := findImageMapping(c.article.Images, link)
		panicIf(im == nil, "didn't find ImageMapping for %s", link)
		images = append(images, im)
	}
	firstImage := images[0]
	s := genGalleryMainHTML(galleryID, firstImage.relativeURL)
	c.r.Printf(s)

	c.r.Printf(`<div class="center mt3 mb6">`)
	for i, im := range images {
		s := genGalleryThumbHTML(galleryID, i, im)
		c.r.Printf(s)
	}
	c.r.Printf(`</div>`)
	return true
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

// RenderPage renders BlockPage
func (c *Converter) RenderPage(block *notionapi.Block) bool {
	if c.r.Page.IsRoot(block) {
		c.r.RenderChildren(block)
		return true
	}

	url, title := c.getURLAndTitleForBlock(block)
	title = html.EscapeString(title)
	c.r.Printf("[%s](%s)\n", url, title)
	
	return false
}

// RenderCode renders BlockCode
func (c *Converter) RenderCode(block *notionapi.Block) bool {
	// code := html.EscapeString(block.Code)
	// fmt.Fprintf(g.f, `<div class="%s">Lang for code: %s</div>
	// <pre class="%s">
	// %s
	// </pre>`, levelCls, block.CodeLanguage, levelCls, code)
	//err := htmlHighlight(c.r.Buf, string(block.Code), block.CodeLanguage, "")
	//must(err)
	return true
}

// if returns false, the block will be rendered with default
func (c *Converter) blockRenderOverride(block *notionapi.Block) bool {
	if c.article.shouldSkipBlock(block) {
		return true
	}
	if c.renderGallery(block) {
		return true
	}
	switch block.Type {
	case notionapi.BlockPage:
		return c.RenderPage(block)
	case notionapi.BlockCode:
		return c.RenderCode(block)
	case notionapi.BlockImage:
		return c.RenderImage(block)
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