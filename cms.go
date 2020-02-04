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

	"github.com/BurntSushi/toml"
)

//Auto-generated from https://xuri.me/toml-to-go/
type ConfigType struct {
	Notion struct {
		StartPage string `toml:"startPage"`
	} `toml:"notion"`
	Hugo struct {
		Repo string `toml:"repo"`
	} `toml:"hugo"`
	Cms struct {
		CacheDir string `toml:"cacheDir"`
		PostsDir string `toml:"postsDir"`
		ImgDir   string `toml:"imgDir"`
	} `toml:"cms"`
}

var (
	config ConfigType
)

func newNotionClient() *notionapi.Client {
	token := os.Getenv("NOTION_TOKEN")
	if token == "" {
		fmt.Println("must set NOTION_TOKEN env variable\n")
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
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		fmt.Println(err)
		return
	}
	
	client := newNotionClient()
	cache, err := caching_downloader.NewDirectoryCache(config.Cms.CacheDir)
	u.Must(err)
	d := caching_downloader.New(cache, client)
	_ = d //to bypass "declared and not used" error

	fmt.Printf("hello world")

	fmt.Printf("Title: %s\n", config.Hugo.Repo)
}
