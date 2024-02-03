/*
Exercise: Web Crawler
In this exercise you'll use Go's concurrency features to parallelize a web crawler.

Modify the Crawl function to fetch URLs in parallel without fetching the same URL twice.

Hint: you can keep a cache of the URLs that have been fetched on a map, but maps alone are not safe for concurrent use!
*/

package main

import (
	"fmt"
	"sync"
)

// fetcher is an instance of fakeFetcher populated with urls and results for url's, sumulating some get response.
// fetcher represents a list of pages, maximum depth(length) of a list is four.
var fetcher = fakeFetcher{
	// key : value(body, []urls)
	"https://golang.org/": &fakeResult{ //element 0
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},

	"https://golang.org/pkg/": &fakeResult{ //element 1
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},

	"https://golang.org/pkg/fmt/": &fakeResult{ //element 2
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},

	"https://golang.org/pkg/os/": &fakeResult{ //element 3
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}

// Fetcher holds methods signatures for handling fetched data.
type Fetcher interface {
	//method Fetch() recieves a URL and returns the body of URL and a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

// fakeResult is a canned result for Fetcher.
type fakeResult struct {
	body string
	urls []string
}

// fakeFetcher is Fetcher that returs canned results.
type fakeFetcher map[string]*fakeResult

func main() {
	//Crawl() call will print all the URLs inside fetcher's first 4 pages which starts with "https://golang.org/"
	Crawl("https://golang.org/", 4, fetcher)
	fmt.Println()

	//Crawl("https://golang.org/pkg/", 4, fetcher)
}

// function Crawl() uses Fetcher to recursively crawl pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
	// TODO: Fetch URLs in parallel.
	// TODO: Don't fetch the same URL twice.

	//safeCache is used for safe access to cached values for sub-goroutines.
	type safeCache struct {
		cachedValue map[string]bool //value is cached if map[value] == true
		mu          sync.Mutex      //lock for cached value
	}

	cache := &safeCache{cachedValue: make(map[string]bool)} //cache initialized by zero-value map, lock is open
	var wg sync.WaitGroup                                   //wg is a wait group counter; used for synchronizing goroutines
	var asyncCrawl func(string, int)                        //asyncCrawl is a recursive crawl function

	//recursive call
	asyncCrawl = func(url string, depth int) {
		defer wg.Done() //deferred wg counter decrement

		//fetcher(our DATA) length check
		if depth <= 0 { //if there is nowhere to crawl to
			return //finish this call
		}

		if !cache.cachedValue[url] { //if a URL value found crawling is not in the cache
			cache.mu.Lock()               //lock cache from other goroutines
			cache.cachedValue[url] = true //add this value to cache
			cache.mu.Unlock()             //unlock cache
		} else { //if a URL value found is already in the cache
			return //finish this call, no need to fetch results
		}

		//fetching results from DATA by URL key
		body, urls, err := fetcher.Fetch(url)
		if err != nil { //if URL is NOT present in DATA
			fmt.Println(err) //print error message(URL not found)
			return           //finish this call
		}

		//if URL is in the DATA
		fmt.Printf("found: %s %q\n", url, body) //print results for URL found

		//then for all URL results fetched from DATA
		for _, u := range urls {
			wg.Add(1)                 //increment wg counter
			go asyncCrawl(u, depth-1) //do next x*URLs async calls with lower length(depth of search)
		}
	}

	//start crawling pages
	wg.Add(1)                 //wg counter increment
	go asyncCrawl(url, depth) //first async call
	wg.Wait()                 //blocking this(main) goroutine until all other sub-goroutines are finished(until wg counter becomes zero)
	fmt.Printf("Cache: %v", cache)
}

// method Fetch() is implemented by fakeFetcher: Fetch() tests reciever for containing a value by url key given in argument.
func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	//if fakeFetcher fetches a fakeResult for url in argument, return fakeResult + success error message
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	//operator if would trigger if res != nil and ok == true;
	//res != nil if f[url]'s value is not nil, ok == true if f[url] contains a value

	//if fakeFetcher cannot fetch a fakeResult for url in argument, return zero-value fakeResult + error message of not fetched url
	return "", nil, fmt.Errorf("not found - url: %s", url)
}
