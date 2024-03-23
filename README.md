# GitHub Repository Cloner
Matt Brocklehurst / www.xor.me.uk

This tool allows you to clone all GitHub repositories for a specified user. It utilizes a GitHub token for authentication, enabling access to both public and private repositories (depending on the permissions granted to the token).

## Prerequisites

- Go (1.15 or later recommended)
- Git installed and configured on your system
- A GitHub Personal Access Token with the appropriate permissions (at least `public_repo` for public repositories, `repo` for full access)

## Installation

1. Ensure you have Go installed on your system. You can download and install Go from [the official website](https://golang.org/dl/).

2. Clone this repository or download the Go script directly.

## Usage

1. Set your GitHub Personal Access Token as an environment variable:
   ```bash
   export GITHUB_TOKEN='your_github_token_here'
    ```

Run the script with the following command, replacing destination_dir with the path to the directory where the repositories should be cloned.

go run checkout_all_repo.go destination_dir
