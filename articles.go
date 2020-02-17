package main

import (
	"html/template"
	"sort"

	"github.com/kjk/notionapi"
	"github.com/kjk/notionapi/caching_downloader"
	"github.com/kjk/u"
)

var (
	notionBlogsStartPage      = "300db9dc27c84958a08b8d0c37f4cfe5"
	notionWebsiteStartPage    = "1087c264ef6b450ca5c7d3c034b399d7"
	notionGoCookbookStartPage = "7495260a1daa46118858ad2e049e77e6"
)

// Articles has info about all articles downloaded from notion
type Articles struct {
	idToArticle map[string]*Article
	idToPage    map[string]*notionapi.Page
	// all downloaded articles
	articles []*Article
	// articles that are not hidden
	articlesNotHidden []*Article
	// articles that belong to a blog
	blog []*Article
	// blog articles that are not hidden
	blogNotHidden []*Article
}

func (a *Articles) getNotHidden() []*Article {
	if a.articlesNotHidden == nil {
		var arr []*Article
		for _, article := range a.articles {
			if !article.IsHidden() {
				arr = append(arr, article)
			}
		}
		a.articlesNotHidden = arr
	}
	return a.articlesNotHidden
}

func (a *Articles) getBlogNotHidden() []*Article {
	if a.blogNotHidden == nil {
		var arr []*Article
		for _, article := range a.blog {
			if !article.IsHidden() {
				arr = append(arr, article)
			}
		}
		a.blogNotHidden = arr
	}
	return a.blogNotHidden
}

func buildArticleNavigation(article *Article, isRootPage func(string) bool, idToBlock map[string]*notionapi.Block) {
	// some already have path (e.g. those that belong to a collection)
	if len(article.Paths) > 0 {
		return
	}

	page := article.page.Root()
	currID := normalizeID(page.ParentID)

	var paths []URLPath
	for !isRootPage(currID) {
		block := idToBlock[currID]
		if block == nil {
			break
		}
		// parent could be a column
		if block.Type != notionapi.BlockPage {
			currID = normalizeID(block.ParentID)
			continue
		}
		title := block.Title
		uri := "/article/" + normalizeID(block.ID) + "/" + urlify(title)
		path := URLPath{
			Name: title,
			URL:  uri,
		}
		paths = append(paths, path)
		currID = normalizeID(block.ParentID)
	}

	// set in reverse order
	n := len(paths)
	for i := 1; i <= n; i++ {
		path := paths[n-i]
		article.Paths = append(article.Paths, path)
	}
}

func normalizeID(id string) string {
	return notionapi.ToNoDashID(id)
}

func addIDToBlock(block *notionapi.Block, idToBlock map[string]*notionapi.Block) {
	id := normalizeID(block.ID)
	idToBlock[id] = block
	for _, block := range block.Content {
		if block == nil {
			continue
		}
		addIDToBlock(block, idToBlock)
	}
}

// build navigation bread-crumbs for articles
func buildArticlesNavigation(articles *Articles) {
	idToBlock := map[string]*notionapi.Block{}
	for _, a := range articles.articles {
		page := a.page
		if page == nil {
			continue
		}
		addIDToBlock(page.Root(), idToBlock)
	}

	isRoot := func(id string) bool {
		id = notionapi.ToNoDashID(id)
		switch id {
		case notionBlogsStartPage, notionWebsiteStartPage, notionGoCookbookStartPage:
			return true
		}
		return false
	}

	for _, article := range articles.articles {
		buildArticleNavigation(article, isRoot, idToBlock)
	}
}

func loadArticles(d *caching_downloader.Downloader, config *ConfigType) *Articles {
	res := &Articles{}
	_, err := d.DownloadPagesRecursively(config.Notion.StartPage, nil)
	must(err)
	res.idToPage = d.IdToPage

	c := d.GetClientCopy()
	res.idToArticle = map[string]*Article{}
	for id, page := range res.idToPage {
		u.PanicIf(id != notionapi.ToNoDashID(id), "bad id '%s' sneaked in", id)
		article := notionPageToArticle(c, page)
		if article.urlOverride != "" {
			verbose("url override: %s => %s\n", article.urlOverride, article.ID)
		}
		res.idToArticle[id] = article
		// this might be legacy, short id. If not, we just set the same value twice
		articleID := article.ID
		res.idToArticle[articleID] = article
		if article.IsBlog() {
			res.blog = append(res.blog, article)
		}
		res.articles = append(res.articles, article)
	}

	for _, article := range res.articles {
		html, images := notionToMarkdown(c, article, res)
		article.BodyHTML = string(html)
		article.HTMLBody = template.HTML(article.BodyHTML)
		article.Images = append(article.Images, images...)
	}

	buildArticlesNavigation(res)

	sort.Slice(res.blog, func(i, j int) bool {
		return res.blog[i].PublishedOn.After(res.blog[j].PublishedOn)
	})

	return res
}