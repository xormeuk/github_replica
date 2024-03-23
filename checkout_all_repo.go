package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// Repository structure to match JSON data from GitHub API
type Repository struct {
	Name   string `json:"name"`
	SshURL string `json:"ssh_url"`
}

func checkoutRepositories(username, token, destinationDir string) {
	url := fmt.Sprintf("https://api.github.com/users/%s/repos", username)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return
	}

	req.Header.Add("Authorization", "token "+token)
	req.Header.Add("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to fetch repositories: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Failed to fetch repositories: received status code %d\n%s\n", resp.StatusCode, string(body))
		return
	}

	var repositories []Repository
	if err := json.NewDecoder(resp.Body).Decode(&repositories); err != nil {
		fmt.Printf("Failed to decode response: %v\n", err)
		return
	}

	for _, repo := range repositories {
		fmt.Printf("Cloning repository: %s\n", repo.Name)
		cmd := exec.Command("git", "clone", repo.SshURL, filepath.Join(destinationDir, repo.Name))
		if err := cmd.Run(); err != nil {
			fmt.Printf("Failed to clone repository: %v\n", err)
		}
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run checkout_all_repo.go GITHUB_USERNAME destination_dir")
		os.Exit(1)
	}

	username := os.Args[1]
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Println("GITHUB_TOKEN environment variable is not set.")
		os.Exit(1)
	}
	destinationDir := os.Args[2]
	checkoutRepositories(username, token, destinationDir)
}
