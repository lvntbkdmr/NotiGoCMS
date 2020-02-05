package main

import (
	"flag"
	_ "net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/kjk/notionapi"
	"github.com/kjk/notionapi/caching_downloader"
	"github.com/kjk/u"

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
)

const (
	analyticsCode = "UA-194516-1"
)

var (
	flgVerbose bool
)

func rebuildAll(d *caching_downloader.Downloader) *Articles {
	//loadTemplates()
	articles := loadArticles(d)
	//readRedirects(articles)
	//netlifyBuild(articles)
	return articles
}

var (
	nDownloadedPage = 0
)

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

func cmdAddNetlifyToken(cmd *exec.Cmd) {
	token := os.Getenv("NETLIFY_TOKEN")
	if token == "" {
		logf("No NETLIFY_TOKEN\n")
		return
	}
	logf("Has NETLIFY_TOKEN\n")
	cmd.Args = append(cmd.Args, "--auth", token)
}

func main() {
		var (
		flgDeployDraft     bool
		flgDeployProd      bool
		flgNoCache         bool
	)

	{
		flag.BoolVar(&flgVerbose, "verbose", false, "if true, verbose logging")
		flag.BoolVar(&flgNoCache, "no-cache", false, "if true, disables cache for downloading notion pages")
		flag.BoolVar(&flgDeployDraft, "deploy-draft", false, "deploy to netlify as draft")
		flag.BoolVar(&flgDeployProd, "deploy-prod", false, "deploy to netlify production")
		flag.Parse()
	}

	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		must(err)
		return
	}

	openLog()
	defer closeLog()

	netlifyExe := filepath.Join("./node_modules", ".bin", "netlify")

	if flgDeployDraft || flgDeployProd {
		if !u.FileExists(netlifyExe) {
			cmd := exec.Command("yarn", "install")
			u.RunCmdMust(cmd)
		}
	}

	client := newNotionClient()
	cache, err := caching_downloader.NewDirectoryCache(config.Cms.CacheDir)
	must(err)
	d := caching_downloader.New(cache, client)
	d.EventObserver = eventObserver
	d.RedownloadNewerVersions = true
	d.NoReadCache = flgNoCache

	doOpen := runtime.GOOS == "darwin"
	//os.Setenv("PATH", )

	if flgDeployDraft {
		rebuildAll(d)
		cmd := exec.Command(netlifyExe, "deploy", "--dir=netlify_static", "--site=a1bb4018-531d-4de8-934d-8d5602bacbfb")
		cmdAddNetlifyToken(cmd)
		if doOpen {
			cmd.Args = append(cmd.Args, "--open")
		}
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		u.RunCmdMust(cmd)
		return
	}

	if flgDeployProd {
		rebuildAll(d)
		cmd := exec.Command(netlifyExe, "deploy", "--prod", "--dir=netlify_static", "--site=a1bb4018-531d-4de8-934d-8d5602bacbfb")
		cmdAddNetlifyToken(cmd)
		if doOpen {
			cmd.Args = append(cmd.Args, "--open")
		}
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		u.RunCmdMust(cmd)
		return
	}

	articles := rebuildAll(d)

	_ = articles
}
