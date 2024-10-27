package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Repo struct {
	Name        string `json:"name"`
	Stars       int    `json:"stargazers_count"`
	Forks       int    `json:"forks_count"`
	Issues      int    `json:"open_issues_count"`
	Size        int    `json:"size"`
	Description string `json:"description"`
}

func fetchRepos(username string) ([]Repo, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/repos", username)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch repositories: %s", resp.Status)
	}

	var repos []Repo
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, err
	}
	return repos, nil
}

func main() {
	fmt.Print("Enter GitHub username: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	username := scanner.Text()

	repos, err := fetchRepos(username)
	if err != nil {
		log.Fatalf("Error fetching repositories: %v", err)
	}

	fmt.Printf("\nRepositories for user '%s':\n", username)
	for _, repo := range repos {
		fmt.Printf("Repo: %s\n", repo.Name)
		fmt.Printf("  Stars: %d, Forks: %d, Issues: %d, Size: %dKB\n", repo.Stars, repo.Forks, repo.Issues, repo.Size)
	}
}
