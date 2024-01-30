package main

import "regexp"

func (app *application) Run() {
	allRepos := app.orgRepos(&app.config.org)

	var repositories OrderedMapRepositories

	repositories.Init(allRepos)

	for _, repo := range allRepos {

		repositories.RepoMap[repo.GetName()].Topics = app.topics(repo)
		repositories.RepoMap[repo.GetName()].Alerts = app.dependabotAlerts(repo)

	}

	repositories.Compliance()
	repositories.CountAlerts()
	repositories.order(app.config.sortBy)

	caseRegex := regexp.MustCompile(`\.md$`)
	if caseRegex.MatchString(app.config.outputFile) {
		app.writeMarkdown(repositories)
	}
	caseRegex = regexp.MustCompile(`\.html$`)
	if caseRegex.MatchString(app.config.outputFile) {
		app.writeHTML(repositories)
	}
	caseRegex = regexp.MustCompile(`\.xlsx$`)
	if caseRegex.MatchString(app.config.outputFile) {
		app.writeExcel(repositories)
	}

	caseRegex = regexp.MustCompile(`\.json$`)
	if caseRegex.MatchString(app.config.outputFile) {
		app.writeJSON(repositories)
	}

}
