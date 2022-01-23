package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	twitterscraper "github.com/n0madic/twitter-scraper"
)

var wg sync.WaitGroup
var mwg sync.WaitGroup
var usr string
var update bool
var onlyrtw bool
var vidz bool
var imgs bool

func download(url string, filetype string, output string, dwn_type string) {
	defer wg.Done()
	segments := strings.Split(url, "/")
	name := segments[len(segments)-1]
	resp, _ := http.Get(url)
	if resp.StatusCode != 200 {
		return
	}
	var f *os.File
	defer f.Close()
	if dwn_type == "user" {
		if update {
			if _, err := os.Stat(output + "/" + filetype + "/" + name); !errors.Is(err, os.ErrNotExist) {
				fmt.Println(name + ": alrady exist")
				return
			}
		}
		if filetype == "rtimg" {
			f, _ = os.Create(output + "/img/RE-" + name)
		} else if filetype == "rtvideo" {
			f, _ = os.Create(output + "/video/RE-" + name)
		} else {
			f, _ = os.Create(output + "/" + filetype + "/" + name)
		}
	} else {
		if update {
			if _, err := os.Stat(output + "/" + name); !errors.Is(err, os.ErrNotExist) {
				fmt.Println("File exist")
				return
			}
		}
		f, _ = os.Create(output + "/" + name)
	}
	defer resp.Body.Close()
	io.Copy(f, resp.Body)
	fmt.Println("Downloaded " + name)
}

func vidUrl(video string) string {
	vid := strings.Split(string(video), " ")
	v := vid[len(vid)-1]
	v = strings.TrimSuffix(v, "}")
	vid = strings.Split(v, "?")
	return vid[0]
}

func videoSingle(tweet *twitterscraper.Tweet, output string, dwn_tweet string) {
	if len(tweet.Videos) > 0 {
		for _, i := range tweet.Videos {
			j := fmt.Sprintf("%s", i)
			v := vidUrl(j)
			wg.Add(1)
			if usr != "" {
				go download(v, "rtvideo", output, "user")
			} else {
				go download(v, "tweet", output, "tweet")
			}
		}
		wg.Wait()
	}
}

func photoSingle(tweet *twitterscraper.Tweet, output string, dwn_type string) {
	if len(tweet.Photos) > 0 {
		for _, i := range tweet.Photos {
			if !strings.Contains(i, "video_thumb/") {
				wg.Add(1)
				if usr != "" {
					go download(i, "rtimg", output, "user")
				} else {
					go download(i, "tweet", output, "tweet")
				}
			}
		}
		wg.Wait()
	}
}

func singleTweet(output string, id string) {
	scraper := twitterscraper.New()
	tweet, err := scraper.GetTweet(id)
	if err != nil {
		fmt.Println(err)
	}
	if usr != "" {
		if vidz {
			videoSingle(tweet, output, "video")
		}
		if imgs {
			photoSingle(tweet, output, "img")
		}
	} else {
		videoSingle(tweet, output, "tweet")
		photoSingle(tweet, output, "tweet")
	}
}

func main() {
	var single, output string
	if single != "" {
		if output == "" {
			output = "./"
			fmt.Println("2")
		} else {
			os.MkdirAll(output, os.ModePerm)
			fmt.Println("3")
		}
		singleTweet(output, single)
		fmt.Println("4")
	}
	mwg.Wait()
}
