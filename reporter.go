package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/tealeg/xlsx"
)

func (app *application) writeMarkdown(repositories OrderedMapRepositories) {
	app.logger.Info(fmt.Sprintf("Writing report to %s", app.config.outputFile))

	// Open a new file for writing only
	file, err := os.OpenFile(
		app.config.outputFile,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666, // For read and write permissions
	)

	if err != nil {
		log.Fatalf("Failed creating file: %s", err)
	}

	defer file.Close()

	// Create a new writer instance
	writer := bufio.NewWriter(file)

	// Write the table header
	_, err = writer.WriteString("| Repository | Compliance | Archived | Critical | High | Moderate | Low | Total |\n")
	if err != nil {
		log.Fatalf("Failed writing to file: %s", err)

	}
	_, err = writer.WriteString("| --- | --- | --- | --- | --- | --- | --- | --- |\n")
	if err != nil {
		log.Fatalf("Failed writing to file: %s", err)
	}

	// Initialize the totals
	totalCritical, totalHigh, totalModerate, totalLow, totalAll := 0, 0, 0, 0, 0

	// Iterate over the repositories in order
	for _, repoName := range repositories.OrderedNames {
		repo := repositories.RepoMap[repoName]

		if app.config.excludeNonProd && repo.Compliance != "glcp-production" {
			continue
		}
		if app.config.excludeArchived && repo.Archived {
			continue
		}
		if app.config.excludeZeroAlerts && repo.Totals.Total == 0 {
			continue
		}
		repoName := strings.Replace(repo.Name, "glcp/", "", 1)
		row := fmt.Sprintf("| [%s](https://github.com/glcp/%s/security/dependabot) | %s | %v | [%d](https://github.com/glcp/%s/security/dependabot?q=is%%3Aopen+severity%%3Acritical) | [%d](https://github.com/glcp/%s/security/dependabot?q=is%%3Aopen+severity%%3Ahigh) | [%d](https://github.com/glcp/%s/security/dependabot?q=is%%3Aopen+severity%%3Amoderate) | [%d](https://github.com/glcp/%s/security/dependabot?q=is%%3Aopen+severity%%3Alow) | [%d](https://github.com/glcp/%s/security/dependabot) |\n",
			repo.Name, repoName,
			repo.Compliance,
			repo.Archived,
			repo.Totals.Critical, repoName,
			repo.Totals.High, repoName,
			repo.Totals.Moderate, repoName,
			repo.Totals.Low, repoName,
			repo.Totals.Total, repoName)
		_, err = writer.WriteString(row)
		if err != nil {
			log.Fatalf("Failed writing to file: %s", err)
		}

		// Add the counts to the totals
		totalCritical += repo.Totals.Critical
		totalHigh += repo.Totals.High
		totalModerate += repo.Totals.Moderate
		totalLow += repo.Totals.Low
		totalAll += repo.Totals.Total
	}

	// Write the totals row
	totalsRow := fmt.Sprintf("| **Total** |  |  | **%d** | **%d** | **%d** | **%d** | **%d** |\n",
		totalCritical,
		totalHigh,
		totalModerate,
		totalLow,
		totalAll)
	_, err = writer.WriteString(totalsRow)
	if err != nil {
		log.Fatalf("Failed writing to file: %s", err)
	}

	// Save the changes
	err = writer.Flush()
	if err != nil {
		log.Fatalf("Failed flushing writer: %s", err)
	}

}

func (app *application) writeHTML(repositories OrderedMapRepositories) {
	app.logger.Info(fmt.Sprintf("Writing report to %s", app.config.outputFile))

	// Open a new file for writing only
	file, err := os.OpenFile(
		app.config.outputFile,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666, // For read and write permissions
	)

	if err != nil {
		log.Fatalf("Failed creating file: %s", err)
	}

	defer file.Close()

	// Create a new writer instance
	writer := bufio.NewWriter(file)

	// Write the HTML header and table header
	_, err = writer.WriteString("<!DOCTYPE html>\n<html>\n<head>\n<title>Report</title><link rel='stylesheet' href='https://cdn.simplecss.org/simple.min.css'>\n</head>\n<body>\n<table>\n<tr><th>Repository</th><th>Compliance</th><th>Archived</th><th>Critical</th><th>High</th><th>Moderate</th><th>Low</th><th>Total</th></tr>\n")
	if err != nil {
		log.Fatalf("Failed writing to file: %s", err)
	}

	// Initialize the totals
	totalCritical, totalHigh, totalModerate, totalLow, totalAll := 0, 0, 0, 0, 0

	// Iterate over the repositories in order
	for _, repoName := range repositories.OrderedNames {
		repo := repositories.RepoMap[repoName]

		if app.config.excludeNonProd && repo.Compliance != "glcp-production" {
			continue
		}
		if app.config.excludeArchived && repo.Archived {
			continue
		}
		if app.config.excludeZeroAlerts && repo.Totals.Total == 0 {
			continue
		}
		repoName := strings.Replace(repo.Name, "glcp/", "", 1)
		row := fmt.Sprintf("<tr><td><a href='https://github.com/glcp/%s/security/dependabot'>%s</a></td><td>%s</td><td>%v</td><td><a href='https://github.com/glcp/%s/security/dependabot?q=is%%3Aopen+severity%%3Acritical'>%d</a></td><td><a href='https://github.com/glcp/%s/security/dependabot?q=is%%3Aopen+severity%%3Ahigh'>%d</a></td><td><a href='https://github.com/glcp/%s/security/dependabot?q=is%%3Aopen+severity%%3Amoderate'>%d</a></td><td><a href='https://github.com/glcp/%s/security/dependabot?q=is%%3Aopen+severity%%3Alow'>%d</a></td><td><a href='https://github.com/glcp/%s/security/dependabot'>%d</a></td></tr>\n",
			repoName, repo.Name,
			repo.Compliance,
			repo.Archived,
			repoName, repo.Totals.Critical,
			repoName, repo.Totals.High,
			repoName, repo.Totals.Moderate,
			repoName, repo.Totals.Low,
			repoName, repo.Totals.Total)
		_, err = writer.WriteString(row)
		if err != nil {
			log.Fatalf("Failed writing to file: %s", err)
		}

		// Add the counts to the totals
		totalCritical += repo.Totals.Critical
		totalHigh += repo.Totals.High
		totalModerate += repo.Totals.Moderate
		totalLow += repo.Totals.Low
		totalAll += repo.Totals.Total
	}

	// Write the totals row
	totalsRow := fmt.Sprintf("<tr><td><strong>Total</strong></td><td></td><td></td><td><strong>%d</strong></td><td><strong>%d</strong></td><td><strong>%d</strong></td><td><strong>%d</strong></td><td><strong>%d</strong></td></tr>\n",
		totalCritical,
		totalHigh,
		totalModerate,
		totalLow,
		totalAll)
	_, err = writer.WriteString(totalsRow)
	if err != nil {
		log.Fatalf("Failed writing to file: %s", err)
	}

	// Write the HTML footer
	_, err = writer.WriteString("</table>\n</body>\n</html>")
	if err != nil {
		log.Fatalf("Failed writing to file: %s", err)
	}

	// Save the changes
	err = writer.Flush()
	if err != nil {
		log.Fatalf("Failed flushing writer: %s", err)
	}
}

func (app *application) writeExcel(repositories OrderedMapRepositories) {
	app.logger.Info(fmt.Sprintf("Writing report to %s", app.config.outputFile))

	// Create a new Excel file
	file := xlsx.NewFile()

	// Add a new sheet
	sheet, err := file.AddSheet("Report")
	if err != nil {
		log.Fatalf("Failed adding sheet: %s", err)
	}

	// Add the header row
	row := sheet.AddRow()
	row.AddCell().Value = "Repository"
	row.AddCell().Value = "Compliance"
	row.AddCell().Value = "Archived"
	row.AddCell().Value = "Critical"
	row.AddCell().Value = "High"
	row.AddCell().Value = "Moderate"
	row.AddCell().Value = "Low"
	row.AddCell().Value = "Total"

	// Iterate over the repositories in order
	for _, repoName := range repositories.OrderedNames {
		repo := repositories.RepoMap[repoName]

		if app.config.excludeNonProd && repo.Compliance != "glcp-production" {
			continue
		}
		if app.config.excludeArchived && repo.Archived {
			continue
		}
		if app.config.excludeZeroAlerts && repo.Totals.Total == 0 {
			continue
		}

		// Add a new row for each repository
		row := sheet.AddRow()
		row.AddCell().Value = repo.Name
		row.AddCell().Value = repo.Compliance
		row.AddCell().Value = strconv.FormatBool(repo.Archived)
		row.AddCell().Value = strconv.Itoa(repo.Totals.Critical)
		row.AddCell().Value = strconv.Itoa(repo.Totals.High)
		row.AddCell().Value = strconv.Itoa(repo.Totals.Moderate)
		row.AddCell().Value = strconv.Itoa(repo.Totals.Low)
		row.AddCell().Value = strconv.Itoa(repo.Totals.Total)
	}

	// Save the Excel file
	err = file.Save(app.config.outputFile)
	if err != nil {
		log.Fatalf("Failed saving file: %s", err)
	}
}

func (app *application) writeJSON(repositories OrderedMapRepositories) {
	app.logger.Info(fmt.Sprintf("Writing report to %s", app.config.outputFile))

	// Create a new file
	file, err := os.Create(app.config.outputFile)
	if err != nil {
		log.Fatalf("Failed creating file: %s", err)
	}
	defer file.Close()

	// Encode the repositories to JSON
	encoder := json.NewEncoder(file)
	err = encoder.Encode(repositories)
	if err != nil {
		log.Fatalf("Failed encoding to JSON: %s", err)
	}
}
