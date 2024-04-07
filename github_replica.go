package main

import (
	"encoding/json"
	"flag"
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

var updateExisting = flag.Bool("update-existing", false, "If set, will update existing repositories by pulling or cloning as necessary.")

func checkoutOrUpdateRepository(repo Repository, destinationDir string) {
	repoPath := filepath.Join(destinationDir, repo.Name)
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		fmt.Printf("Cloning repository: %s\n", repo.Name)
		cmd := exec.Command("git", "clone", repo.SshURL, repoPath)
		if err := cmd.Run(); err != nil {
			log.Printf("Failed to clone repository %s: %v", repo.Name, err)
		} else {
			fmt.Printf("Successfully cloned %s into %s\n", repo.Name, repoPath)
		}
	} else {
		if *updateExisting {
			fmt.Printf("Updating repository: %s\n", repo.Name)
			cmd := exec.Command("git", "-C", repoPath, "fetch", "--all")
			cmd.Run() // Fetch changes
			cmd = exec.Command("git", "-C", repoPath, "reset", "--hard", "origin/master")
			cmd.Run() // Reset any local changes
			cmd = exec.Command("git", "-C", repoPath, "pull")
			if err := cmd.Run(); err != nil {
				log.Printf("Failed to update repository %s: %v", repo.Name, err)
			} else {
				fmt.Printf("Successfully updated %s\n", repo.Name)
			}
		} else {
			fmt.Printf("Warning: Repository '%s' already exists and '-update-existing' not set\n", repo.Name)
		}
	}
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
		checkoutOrUpdateRepository(repo, destinationDir)
	}
}

func main() {
	flag.Parse()
	fmt.Println("GitHub Replica\nMatt Brocklehurst / www.xor.me.uk")
	if GitCommitHash != "" {
		fmt.Printf("Git Commit Hash: %s\n", GitCommitHash)
	}

	if len(flag.Args()) != 1 {
		fmt.Println("Usage: go run github_replica.go [-update-existing] <destination_dir>")
		os.Exit(1)
	}

	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		log.Fatal("GITHUB_TOKEN environment variable is not set")
	}

	destinationDir := flag.Arg(0) // Assumes the non-flag argument is the destination directory
	checkoutRepositories(githubToken, destinationDir)
}
