# GitHub Replica

Matt Brocklehurst / www.xor.me.uk

This tool allows you to clone all GitHub repositories for a specified user. It utilizes a GitHub token for authentication, enabling access to both public and private repositories (depending on the permissions granted to the token).

# Prerequisites

Go (1.15 or later recommended)
Git installed and configured on your system
A GitHub Personal Access Token with the appropriate permissions (at least public_repo for public repositories, repo for full access)

# Installation

## Binary Download:
Users can download pre-compiled binaries for their specific architecture (amd64 or arm64) from the Releases section on GitHub. Debian packages (.deb) are also available for both architectures.

## From Source:

Ensure you have Go installed on your system. You can download and install Go from the official Go website.

Clone this repository or download the Go script directly.

### Usage

First, set your GitHub Personal Access Token as an environment variable:

```bash
export GITHUB_TOKEN='your_github_token_here'
```

### Running the Pre-compiled Binary

After downloading the binary for your platform, you can run the application directly from the terminal. Replace ./github-replica with the path to the downloaded binary and destination_dir with the path to the directory where the repositories should be cloned:

```bash
./github-replica destination_dir
```

Ensure the binary has execute permission:

```bash
chmod +x ./github-replica
```

### Advanced

If you prefer to compile and run from source, navigate to the directory containing the Go script and execute:

```bash
go run github_replica.go destination_dir
```

Note: Running from source is primarily recommended for development or troubleshooting purposes.
