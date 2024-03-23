package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// only used by CI
var GitCommitHash string

// Repository represents the minimal data for a GitHub repository we need.
type Repository struct {
	Name   string `json:"name"`
	SshURL string `json:"ssh_url"`
}

func checkoutRepositories(githubToken, destinationDir string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/user/repos", nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Add("Authorization", "token "+githubToken)
	req.Header.Add("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to fetch repositories: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("Failed to fetch repositories: %s", resp.Status)
	}

	var repos []Repository
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		log.Fatalf("Failed to decode response: %v", err)
	}

	for _, repo := range repos {
		fmt.Printf("Cloning repository: %s\n", repo.Name)
		repoPath := filepath.Join(destinationDir, repo.Name)

		// Clone repo into destinationDir/repo.Name
		cmd := exec.Command("git", "clone", repo.SshURL, repoPath)
		err := cmd.Run()
		if err != nil {
			log.Printf("Failed to clone repository %s: %v", repo.Name, err)
			continue
		}
		fmt.Printf("Successfully cloned %s into %s\n", repo.Name, repoPath)
	}
}

func main() {
	fmt.Println("GitHub Replica\nMatt Brocklehurst / www.xor.me.uk")
	// Print the GitCommitHash only if it's set
	if GitCommitHash != "" {
		fmt.Printf("Git Commit Hash: %s\n", GitCommitHash)
	}

	if len(os.Args) != 2 {
		fmt.Println("Usage: go run github_replica.go <destination_dir>")
		os.Exit(1)
	}

	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		log.Fatal("GITHUB_TOKEN environment variable is not set")
	}

	destinationDir := os.Args[1]
	checkoutRepositories(githubToken, destinationDir)
}
