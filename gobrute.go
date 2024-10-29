package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/schollz/progressbar/v3"
)

var (
	url           string
	wordlist      string
	login         string
	failMsg       string
	payload       string
	numThreads    int
	timeout       int
	totalAttempts int32
)

func main() {
	// Parse command-line arguments
	flag.StringVar(&url, "u", "", "Target API URL")
	flag.StringVar(&wordlist, "w", "", "Path to the wordlist")
	flag.StringVar(&login, "l", "", "Username to test with")
	flag.StringVar(&failMsg, "f", "", "Failure message to identify unsuccessful attempts")
	flag.StringVar(&payload, "p", "", "Payload format with ^USER^ for username and ^PASS^ for password")
	flag.IntVar(&numThreads, "t", 10, "Number of threads to use (default 10)")
	flag.IntVar(&timeout, "o", 2000, "Timeout for each request in milliseconds (default 2000)")
	flag.Parse()

	if url == "" || wordlist == "" || login == "" || failMsg == "" || payload == "" {
		fmt.Println("All parameters (-u, -w, -l, -f, -p) are required.")
		flag.Usage()
		os.Exit(1)
	}

	// Read the wordlist file
	passwords, err := readWordlist(wordlist)
	if err != nil {
		fmt.Printf("Error reading wordlist: %v\n", err)
		os.Exit(1)
	}
	totalAttempts = int32(len(passwords))

	// Initialize progress bar
	bar := progressbar.NewOptions64(int64(totalAttempts),
		progressbar.OptionSetDescription("Attempting passwords"),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(20),
		progressbar.OptionClearOnFinish(),
	)

	// Set up concurrency control and shared state
	var wg sync.WaitGroup
	var attemptCount int32
	success := false
	passwordsChan := make(chan string, numThreads)

	// Launch worker goroutines
	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := &http.Client{
				Timeout: time.Duration(timeout) * time.Millisecond,
			}
			for password := range passwordsChan {
				if success {
					break
				}

				// Prepare payload
				data := strings.ReplaceAll(payload, "^USER^", login)
				data = strings.ReplaceAll(data, "^PASS^", password)

				// Send POST request
				resp, err := sendRequest(client, url, data)
				if err != nil {
					fmt.Printf("Error sending request: %v\n", err)
					continue
				}
				defer resp.Body.Close()

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Printf("Error reading response: %v\n", err)
					continue
				}

				// Check for failure message
				if !strings.Contains(string(body), failMsg) {
					fmt.Printf("\nSuccess! Username: %s Password: %s\n", login, password)
					success = true
					break
				}

				// Update progress bar
				atomic.AddInt32(&attemptCount, 1)
				bar.Add(1)
			}
		}()
	}

	// Distribute passwords to workers
	for _, password := range passwords {
		if success {
			break
		}
		passwordsChan <- password
	}
	close(passwordsChan)

	// Wait for all workers to finish
	wg.Wait()

	if !success {
		fmt.Println("\nFailed to find a valid password in the wordlist.")
	}
}

// Reads wordlist file and returns a slice of passwords
func readWordlist(filepath string) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var passwords []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		passwords = append(passwords, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return passwords, nil
}

// Sends a POST request with the specified payload data
func sendRequest(client *http.Client, url, data string) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0")
	req.Header.Set("Content-Type", "application/json")
	return client.Do(req)
}
