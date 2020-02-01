package main

import (
	"flag"
	_ "net/url"
	"os"
	_ "os/exec"
	_ "path/filepath"
	_ "runtime"
	_ "time"
	"fmt"

	"github.com/kjk/notionapi"
	"github.com/kjk/notionapi/caching_downloader"
	"github.com/kjk/u"
)

var (
	cacheDir = "cache"
)

func newNotionClient() *notionapi.Client {
	token := os.Getenv("NOTION_TOKEN")
	if token == "" {
		fmt.Printf("must set NOTION_TOKEN env variable\n")
		flag.Usage()
		os.Exit(1)
	}
	// TODO: verify token still valid, somehow
	client := &notionapi.Client{
		AuthToken: token,
	}

	return client
}

func recreateDir(dir string) {
	err := os.RemoveAll(dir)
	u.Must(err)
	err = os.MkdirAll(dir, 0755)
	u.Must(err)
}

func main() {

	recreateDir("static")

	client := newNotionClient()
	cache, err := caching_downloader.NewDirectoryCache(cacheDir)
	u.Must(err)
	d := caching_downloader.New(cache, client)
	_ = d //to bypass "declared and not used" error

	fmt.Printf("hello world")
}
