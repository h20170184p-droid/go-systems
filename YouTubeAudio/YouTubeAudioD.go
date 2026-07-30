package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

type VideoData struct {
	Title string
	URL   string
}

func main() {
	fmt.Println("\n=== YouTube Audio Extractor ===")
	fmt.Println("Note: For direct extraction, ensure 'yt-dlp' or 'ffmpeg' is installed on your system.")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\nTo exit, type '!Exi!T'")
		fmt.Print("Type the YouTube Video Name / Search Query: > ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		query := strings.TrimSpace(input)
		if query == "!Exi!T" {
			fmt.Println("Exiting...")
			os.Exit(0)
		}
		if query == "" {
			continue
		}

		videoLink, videoName := YouTubeSearch(reader, query)
		if videoLink == "" {
			continue
		}

		fmt.Printf("\nSelected: %s\n", videoName)
		fmt.Println("Choose Action:")
		fmt.Println("1. Download Audio directly (via yt-dlp - Recommended)")
		fmt.Println("2. Open in Browser & Record System Audio (Legacy Mode)")
		fmt.Print("Select option (1/2): ")

		optStr, _ := reader.ReadString('\n')
		optStr = strings.TrimSpace(optStr)

		saveFileName := sanitizeFilename(videoName) + ".flac"

		if optStr == "1" {
			downloadDirectAudio(videoLink, saveFileName)
		} else {
			launchAndRecord(videoLink, saveFileName)
		}
	}
}

// YouTubeSearch scrapes/searches and presents user selection cleanly
func YouTubeSearch(reader *bufio.Reader, query string) (string, string) {
	searchURL := "https://www.youtube.com/results?search_query=" + url.QueryEscape(query)

	req, _ := http.NewRequest("GET", searchURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Network error:", err)
		return "", ""
	}
	defer resp.Body.Close()

	// Use safely parsed JSON or yt-dlp search fallback
	results := searchYtDlp(query)
	if len(results) == 0 {
		fmt.Println("No videos found or search blocked.")
		return "", ""
	}

	fmt.Println("\nSearch Results:")
	for idx, item := range results {
		fmt.Printf("%2d. %s (%s)\n", idx+1, item.Title, item.URL)
	}
	fmt.Println(" 0. Cancel search")

	for {
		fmt.Print("\nSelect index (0 to cancel): ")
		input, _ := reader.ReadString('\n')
		idx, err := strconv.Atoi(strings.TrimSpace(input))

		if err != nil || idx < 0 || idx > len(results) {
			fmt.Println("Invalid selection. Try again.")
			continue
		}
		if idx == 0 {
			return "", ""
		}
		return results[idx-1].URL, results[idx-1].Title
	}
}

// Uses yt-dlp for search JSON metadata safely (no manual fragile scraping needed)
func searchYtDlp(query string) []VideoData {
	cmd := exec.Command("yt-dlp", fmt.Sprintf("ytsearch10:%s", query), "--dump-json", "--default-search", "ytsearch")
	out, err := cmd.Output()
	if err != nil {
		// If yt-dlp isn't installed, fallback notice
		fmt.Println("Notice: 'yt-dlp' CLI not found. Falling back to basic browser mode...")
		return []VideoData{{Title: query, URL: "https://www.youtube.com/results?search_query=" + url.QueryEscape(query)}}
	}

	var results []VideoData
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var entry struct {
			Title    string `json:"title"`
			WebpageURL string `json:"webpage_url"`
		}
		if err := json.Unmarshal([]byte(line), &entry); err == nil {
			results = append(results, VideoData{Title: entry.Title, URL: entry.WebpageURL})
		}
	}
	return results
}

// Directly extracts audio stream in high quality without recording desktop noise
func downloadDirectAudio(videoURL, outputFile string) {
	fmt.Printf("Downloading high-quality FLAC audio directly to '%s'...\n", outputFile)
	cmd := exec.Command("yt-dlp", "-x", "--audio-format", "flac", "-o", outputFile, videoURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Download failed: %v\n", err)
	} else {
		fmt.Println("Download complete successfully!")
	}
}

// Legacy record option
func launchAndRecord(videoURL, outputFile string) {
	fmt.Println("Launching default browser...")
	if err := openBrowser(videoURL); err != nil {
		fmt.Println("Error opening browser:", err)
		return
	}

	fmt.Println("Recording system audio via FFmpeg... Press Ctrl+C or stop FFmpeg when done.")
	var cmd *exec.Cmd

	if runtime.GOOS == "linux" {
		cmd = exec.Command("ffmpeg", "-f", "pulse", "-i", "@DEFAULT_SINK@.monitor", "-c:a", "flac", outputFile)
	} else {
		fmt.Println("System audio recording in this script is configured for Linux PulseAudio only.")
		return
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		fmt.Printf("FFmpeg finished/stopped: %v\n", err)
	}
}

func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
	return cmd.Start()
}

func sanitizeFilename(name string) string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9_\-]+`)
	return reg.ReplaceAllString(name, "_")
}
