package github_api

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func CheckIfRepoExist(owner string, repo string) bool {

	// Create an HTTP client to handle the API rate limit
	client := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://github.com/%s/%s", owner, repo), nil)

	if err != nil {
		panic(err)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	for {
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		if resp.StatusCode == 403 {
			resetTime := time.Unix(0, 0)
			for k, v := range resp.Header {
				if k == "X-Ratelimit-Reset" {
					i, _ := strconv.ParseInt(v[0], 10, 64)
					resetTimeUnix := time.Unix(i, 0)
					resetTime = resetTimeUnix
				}
			}
			fmt.Printf("API rate limit exceeded. Waiting until %s...\n", resetTime)
			time.Sleep(resetTime.Sub(time.Now()))
			fmt.Println("Wait 20 seconds...")
			time.Sleep(time.Second * 20)
		} else if resp.StatusCode == 200 {
			return true

		} else if resp.StatusCode == 404 {
			return false

		} else {
			return false

		}

	}
}

func GetLatestRelease(owner string, repo string) (string, error) {

	// Send a GET request to the GitHub releases page
	resp, err := http.Get(fmt.Sprintf("https://github.com/%s/%s/releases", owner, repo))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read the HTML response into a string
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	bodyString := string(bodyBytes)

	// Find the first non-pre-release tag in the HTML string
	startString := fmt.Sprintf("<a href=\"/%s/%s/releases/tag/", owner, repo)
	startIndex := strings.Index(bodyString, startString)
	if startIndex == -1 {
		// fmt.Printf("No releases found for '%s/%s'.\n", owner, repo)
		return "0.0.0", nil
	}
	endIndex := strings.Index(bodyString[startIndex+len(startString):], "\">")
	if endIndex == -1 {
		// fmt.Printf("No releases found for '%s/%s'.\n", owner, repo)
		return "0.0.0", errors.New("No releases found for '" + owner + "/" + repo + "'")
	}
	tagName := strings.Replace(strings.Split(bodyString[startIndex+len(startString):startIndex+len(startString)+endIndex], " ")[0], "\"", "", 1)

	return tagName, nil
}
