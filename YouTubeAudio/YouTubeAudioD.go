package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

func main() {
	// Take the input here but run the execution in a separate function.
	// That might preserve the repeated runs
	// Avoid gotos?
	// Perhaps even preserve memory foot print and avoid unnecessary loops?

	fmt.Printf("\nWelcome to Youtube Audio Downloader (sort of.)\nNote: There are issues where advertisement audio will also be recorded. I am working on it.")
	for {
		fmt.Println("To exit the program, type !Exi!T and press enter!")
		fmt.Printf("\nType the name of the YouTube Video to extract the audio from: > ")
		r := bufio.NewReader(os.Stdin)
		yvname, err := r.ReadString('\n')
		if err == nil {
			yvname = strings.TrimSpace(yvname)
			if yvname == "!Exi!T" {
				os.Exit(0)
			}
			// Take the name and go into a function
			videoLink, videoname := YouTubeSearch(yvname)
			if videoLink == "" && videoname == "" {
				continue
			}
			// Once the video link is available, launch default web-browser of the linux system and use the provided video-link
			fmt.Println("Launching the default browser! Please handle the ad skips on your own. It is still under development.")
			fmt.Println("Once the video of choice is open, come back to the terminal to approve the audio download!")

			// fmt.Println(videoLink)

			err = openBrowser(videoLink)
			if err != nil {
				fmt.Println("Error occurred: ", err)
			}

			fmt.Println("Browser is open. The video is in paused state by default ideally which is better.")
			saveFileName := strings.Join(strings.Fields(videoname), "") + ".flac"
			fmt.Println("GO PLAY the video or SKIP the ads on the video and let the video play!")
			command := exec.Command("ffmpeg", "-f", "pulse", "-i", "@DEFAULT_SINK@.monitor", "-c:a", "flac", saveFileName)
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr
			command.Stdin = os.Stdin

			if err := command.Run(); err != nil {
				fmt.Printf("FFmpeg stopped: %v\n", err)
			}
			fmt.Println("Recording done. Go program continues!")
			// command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

			// // Intercepting the signal to handle the ffmpeg commands.
			// sigChan := make(chan os.Signal, 1)
			// signal.Notify(sigChan, os.Interrupt, syscall.SIGINT)
			// if err := command.Start(); err != nil {
			// 	fmt.Printf("Failed to start FFmpeg: %v\n", err)
			// 	return
			// }

			// doneChan := make(chan error, 1)
			// go func() {
			// 	doneChan <- command.Wait()
			// }()

			// select {
			// case <-sigChan:
			// 	fmt.Println("\nCtrl+C detected. Safely stopping FFmpeg...")
			// 	if err := command.Process.Signal(syscall.SIGTERM); err != nil {
			// 		fmt.Printf("Failed to signal FFmpeg: %v\n", err)
			// 	}
			// 	<-doneChan

			// case err := <-doneChan:
			// 	if err != nil {
			// 		fmt.Printf("FFmpeg closed with an error: %v\n", err)
			// 	}
			// }
			// fmt.Println("FFmpeg closed safely. Go program continues executing!")
			// fmt.Println("Doing next tasks...")

		} else {
			fmt.Println("Error occurred: ", err)
		}
	}

}

func YouTubeSearch(videoName string) (string, string) {
	searchquery := "https://www.youtube.com/results?search_query=" + strings.ReplaceAll(videoName, " ", "+")

	// Creating a client? Or emulating a fake client to YouTube
	client := &http.Client{}
	req, _ := http.NewRequest("GET", searchquery, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36") // Copy pasted this line from internet. Sets request header I guess.
	response, err := client.Do(req)                                                              // I generally like to use my own variations to describe errors like e1, e2. But fine I will buckle to industry standard.
	if err != nil {
		fmt.Println("Error: ", err)
		return "", "" // This return hands back the control to main.
		// For now the choice is to launch the browser from this function and nothing to be passed to main as a return type from this function.
	}

	// If the error is not an issue then we dig into finding the YouTube search results.
	// Display first 10 results and then let the user select which video to launch on default browser.

	defer response.Body.Close()
	searchResult_HTML_Bytes, _ := io.ReadAll(response.Body) // Assuming no error occurs reading the contents.
	searchContents := string(searchResult_HTML_Bytes)       // Now we have bunch of data saved as strings that can be opened up using json marshalling I guess.

	// Looking for initial Data, ytInitialData
	re := regexp.MustCompile(`ytInitialData\s*=\s*({.+?});`) // The contents inside the regular expression was copy pasted from website. Not my own brain
	matches := re.FindStringSubmatch(searchContents)
	if len(matches) < 2 {
		fmt.Println("Could not extract initial search data.")
		return "", "" // Also hands back the control to main
	}

	var data map[string]any
	err = json.Unmarshal([]byte(matches[1]), &data) // resort to err2, perhaps I could also have used err = json.Unmarshal....
	if err != nil {
		// fmt.Println("Error: ", err)
		return "", ""
	}
	var items []any
	// Digging to find the video ID and title
	if contents, ok := data["contents"].(map[string]any); ok {
		if searchResults, ok := contents["twoColumnSearchResultsRenderer"].(map[string]any); ok {
			if primaryContents, ok := searchResults["primaryContents"].(map[string]any); ok {
				if selectionList, ok := primaryContents["sectionListRenderer"].(map[string]any); ok {
					if contentList, ok := selectionList["contents"].([]any)[0].(map[string]any); ok {
						if itemSection, ok := contentList["itemSectionRenderer"].(map[string]any); ok {
							items = itemSection["contents"].([]any)
						}
					}
				}
			}
		}
	}

	type videoData struct {
		title string
		URL   string
	}

	var matchMap []videoData
	for _, item := range items {
		vR, ok := item.(map[string]any)["videoRenderer"].(map[string]any)
		if ok {
			title := vR["title"].(map[string]any)["runs"].([]any)[0].(map[string]any)["text"].(string)
			vID := vR["videoId"].(string)
			vURL := "https://youtu.be/" + vID
			matchMap = append(matchMap, videoData{title: title, URL: vURL})
		}
	}

	for index, value := range matchMap {
		fmt.Printf("\n%d. %s (%s)\n", index+1, value.title, value.URL)
		if index == 9 {
			break
		}
	}

	// Now let the user choose the url / video of choice
	fmt.Println("If the video of your choice is not in the list type 0")
	fmt.Printf("\nSelect the video of your choice using the index: ")
	var selection int
	var selReader = bufio.NewReader(os.Stdin)
	selInput, _ := selReader.ReadString('\n')
	selInput = strings.TrimSpace(selInput)
	fmt.Sscan(selInput, &selection)
	if selection == 0 {
		return "", ""
	}
	fmt.Println("Selected video: ", matchMap[selection-1].title)
	targetURL := matchMap[selection-1].URL
	targetname := matchMap[selection-1].title

	// Return this string back
	return targetURL, targetname
}

// Browser open function

func openBrowser(url string) error {
	var command *exec.Cmd

	switch runtime.GOOS {
	case "Windows":
		command = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		command = exec.Command("open", url)
	case "linux":
		command = exec.Command("xdg-open", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
	return command.Start()
}
