package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	// "time"
)

func worker(linkChan chan string, wg *sync.WaitGroup) {
	// Decreasing internal counter for wait-group as soon as goroutine finishes
	defer wg.Done()

	for url := range linkChan {
		urlobj, e := http.Get("http://" + url + ".com")
		if e == nil {
			b, _ := json.Marshal(urlobj.Header)
			fmt.Println("%s", string(b))
		}
	}

}

func main() {
	b, e := ioutil.ReadFile("./list.txt")
	if e != nil {
		panic(e)
	}
	File := strings.Split(string(b), "\n")
	// yourLinksSlice := make([]string, 50)
	// for i := 0; i < 50; i++ {
	// 	yourLinksSlice[i] = fmt.Sprintf("%d", i+1)
	// }

	lCh := make(chan string)
	wg := new(sync.WaitGroup)

	// Adding routines to workgroup and running then
	for i := 0; i < 30; i++ {
		wg.Add(1)
		go worker(lCh, wg)
	}

	// Processing all links by spreading them to `free` goroutines
	for _, link := range File {
		lCh <- link
	}

	// Closing channel (waiting in goroutines won't continue any more)
	close(lCh)

	// Waiting for all goroutines to finish (otherwise they die as main routine dies)
	wg.Wait()
}

