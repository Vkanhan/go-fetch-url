package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	fmt.Println("fetch url concurrently")

	//validate user input
	if len(os.Args) < 2 {
		fmt.Println("Error: Please provide atleast one url to fetch")
		return
	}

	//Initialize time
	start := time.Now()

	//Initialize the channel to recive the result from fetch
	ch := make(chan string)

	//take user input and iterate over it and concurrently fetch url
	for _, url := range os.Args[1:] {
		go fetch(url, ch)
	}

	//recieve and print the result
	for i := 0; i < len(os.Args)-1; i++ {
		result := <-ch
		fmt.Println(result)
	}

	//print the elapsed time
	elapsedTime := time.Since(start).Seconds()
	fmt.Printf("Total time elapsed: %.2fs\n", elapsedTime)
}

func fetch(url string, ch chan<- string) {
	//initialize time
	start := time.Now()

	//send GET request to url
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}

	//get the bytes by reading response body and discarding it for measurement
	nbytes, err := io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
	if err != nil {
		ch <- fmt.Sprintf("error writing %s: %s", url, err)
		return
	}

	//calculate elapsed time
	timeElapsed := time.Since(start).Seconds()

	//format and send result to channel
	ch <- fmt.Sprintf("%.2fs elapsed from writing %7d bytes from %s", timeElapsed, nbytes, url)

}
