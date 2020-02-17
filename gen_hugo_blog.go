package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"path/filepath"

	"github.com/kjk/u"
	"github.com/gosimple/slug"
)

func netlifyRequestGetFullHost() string {
	return "https://lvnt.be"
}

func netlifyPath(fileName string) string {
	fileName = strings.TrimLeft(fileName, "/")
	path := filepath.Clean(fileName)
	u.CreateDirForFileMust(path)
	return path
}

func copyImages(config *ConfigType) {
	srcDir := filepath.Join(config.Cms.CacheDir, "img")
	dstDir := filepath.Clean(config.Cms.ImgDir)
	u.DirCopyRecurMust(dstDir, srcDir, nil)
}

func genArticle(article *Article, w io.Writer) error {
	canonicalURL := netlifyRequestGetFullHost() + article.URL()
	model := struct {
		AnalyticsCode    string
		Article          *Article
		CanonicalURL     string
		CoverImage       string
		PageTitle        string
		TagsDisplay      string
		HeaderImageURL   string
		NotionEditURL    string
		Description      string
	}{
		AnalyticsCode:    analyticsCode,
		Article:          article,
		CanonicalURL:     canonicalURL,
		CoverImage:       article.HeaderImageURL,
		PageTitle:        article.Title,
		Description:      article.Description,
	}
	if article.page != nil {
		id := normalizeID(article.page.ID)
		model.NotionEditURL = "https://notion.so/" + id
	}
	path := fmt.Sprintf("%s/%s.md", config.Cms.PostsDir, slug.Make(article.Title))
	logVerbose("%s => %s, %s, %s\n", article.ID, path, article.URL(), article.Title)
	return execTemplate(path, "article.tmpl.md", model, w)
}

func hugoBuild(store *Articles, config *ConfigType) {
	//Recreate postsdir
	outDir := filepath.Join(config.Cms.PostsDir)
	err := os.RemoveAll(outDir)
	must(err)
	err = os.MkdirAll(outDir, 0755)
	must(err)

	//Recreate imgdir
	outDir = filepath.Join(config.Cms.ImgDir)
	err = os.RemoveAll(outDir)
	must(err)
	err = os.MkdirAll(outDir, 0755)
	must(err)

	copyImages(config)

	{
		logVerbose("%d articles\n", len(store.idToPage))
		for _, article := range store.articles {
			genArticle(article, nil)
		}
	}
}