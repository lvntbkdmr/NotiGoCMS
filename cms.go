package main

import (
	"flag"
	_ "net/url"
	"os"

	"github.com/kjk/notionapi"
	"github.com/kjk/notionapi/caching_downloader"

	"github.com/BurntSushi/toml"
)

//Auto-generated from https://xuri.me/toml-to-go/
type ConfigType struct {
	Notion struct {
		StartPage string `toml:"startPage"`
	} `toml:"notion"`
	Cms struct {
		CacheDir string `toml:"cacheDir"`
		PostsDir string `toml:"postsDir"`
		ImgDir   string `toml:"imgDir"`
	} `toml:"cms"`
}

var (
	config ConfigType
	flgVerbose bool
	nDownloadedPage = 0
)

const (
	analyticsCode = "UA-194516-1"
)

func rebuildAll(d *caching_downloader.Downloader, config *ConfigType) *Articles {
	loadTemplates()
	articles := loadArticles(d, config)
	//readRedirects(articles)
	hugoBuild(articles, config)
	return articles
}

func eventObserver(ev interface{}) {
	switch v := ev.(type) {
	case *caching_downloader.EventError:
		logf(v.Error)
	case *caching_downloader.EventDidDownload:
		nDownloadedPage++
		logf("%03d '%s' : downloaded in %s\n", nDownloadedPage, v.PageID, v.Duration)
	case *caching_downloader.EventDidReadFromCache:
		// TODO: only verbose
		nDownloadedPage++
		logf("%03d '%s' : read from cache in %s\n", nDownloadedPage, v.PageID, v.Duration)
	case *caching_downloader.EventGotVersions:
		logf("downloaded info about %d versions in %s\n", v.Count, v.Duration)
	}
}

func newNotionClient() *notionapi.Client {
	token := os.Getenv("NOTION_TOKEN")
	if token == "" {
		logf("must set NOTION_TOKEN env variable\n")
		flag.Usage()
		os.Exit(1)
	}
	// TODO: verify token still valid, somehow
	client := &notionapi.Client{
		AuthToken: token,
	}
	if flgVerbose {
		client.Logger = os.Stdout
	}
	return client
}

func main() {
		var (
		flgNoCache         bool
	)

	{
		flag.BoolVar(&flgVerbose, "verbose", false, "if true, verbose logging")
		flag.BoolVar(&flgNoCache, "no-cache", false, "if true, disables cache for downloading notion pages")
		flag.Parse()
	}

	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		must(err)
		return
	}

	openLog()
	defer closeLog()

	client := newNotionClient()
	cache, err := caching_downloader.NewDirectoryCache(config.Cms.CacheDir)
	must(err)
	d := caching_downloader.New(cache, client)
	d.EventObserver = eventObserver
	d.RedownloadNewerVersions = true
	d.NoReadCache = flgNoCache

	articles := rebuildAll(d, &config)

	_ = articles
}
