package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

func downloadImage(url string, imageName string) {
	req, e := http.NewRequest("GET", url, nil)
	req.Header.Set("Referer", "http://www.minnano-av.com/actress_list.php?page=1")
	client := &http.Client{}

	// don't worry about errors
	response, e := client.Do(req)

	if e != nil {
		log.Fatal(e)
	}
	defer response.Body.Close()

	//open a file for writing
	file, err := os.Create(fmt.Sprintf("%s/%s.jpg", SAVE_PATH, imageName))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Success!")
}

const BASE_URL, SAVE_PATH = "http://www.minnano-av.com", "./images"

var wg sync.WaitGroup

func main() {
	doc, err := htmlquery.LoadURL("http://www.minnano-av.com/actress_list.php?page=1")
	nodes, err := htmlquery.QueryAll(doc, "//*[@id=\"main-area\"]/section/table/tbody/tr[*]/td[1]/a/img")
	if err != nil {
		panic(`not a valid XPath expression.`)
	}

	wg.Add(len(nodes))

	for index, value := range nodes {
		fmt.Println(index, value.Attr[0].Val, value.Attr[1].Val)
		go func(value *html.Node) {
			defer wg.Done()
			downloadImage(fmt.Sprintf("%s/%s", BASE_URL, value.Attr[0].Val), value.Attr[1].Val)
		}(value)
	}
	wg.Wait()
}
