package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jung-kurt/gofpdf"
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

func generatePDF(username string, repos []Repo) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Arial", "B", 16)
	pdf.AddPage()

	pdf.Cell(0, 10, fmt.Sprintf("GitHub Repositories for User: %s", username))
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 10)
	pdf.SetFillColor(240, 240, 240)

	headers := []string{"Repository", "Stars", "Forks", "Issues", "Size (KB)", "Description"}
	columnWidths := []float64{50, 20, 20, 20, 20, 80}

	for i, header := range headers {
		pdf.CellFormat(columnWidths[i], 10, header, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	for _, repo := range repos {
		rowHeight := 10.0
		pdf.CellFormat(columnWidths[0], rowHeight, repo.Name, "1", 0, "L", false, 0, "")
		pdf.CellFormat(columnWidths[1], rowHeight, fmt.Sprintf("%d", repo.Stars), "1", 0, "C", false, 0, "")
		pdf.CellFormat(columnWidths[2], rowHeight, fmt.Sprintf("%d", repo.Forks), "1", 0, "C", false, 0, "")
		pdf.CellFormat(columnWidths[3], rowHeight, fmt.Sprintf("%d", repo.Issues), "1", 0, "C", false, 0, "")
		pdf.CellFormat(columnWidths[4], rowHeight, fmt.Sprintf("%d", repo.Size), "1", 0, "C", false, 0, "")
		pdf.Ln(rowHeight)
	}

	return pdf.OutputFileAndClose("Github_Repos.pdf")
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

	err = generatePDF(username, repos)
	if err != nil {
		log.Fatalf("Error generating PDF: %v", err)
	}
}
