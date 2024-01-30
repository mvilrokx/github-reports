package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/v58/github"
)

func (app *application) orgRepos(org *string) []*github.Repository {
	var allRepos []*github.Repository

	app.logger.Info(fmt.Sprintf("Getting repositories for org '%s', using App Id '%d' and App Install Id '%d' ", app.config.org, app.config.appID, app.config.appInstallID))

	for {
		repos, resp, err := app.ghClient.Repositories.ListByOrg(context.Background(), *org, app.RepoLsByOrgOpt)
		if err != nil {
			app.logger.Error(err.Error())
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		app.RepoLsByOrgOpt.Page = resp.NextPage
	}

	return allRepos
}

func (app *application) dependabotAlerts(repo *github.Repository) []*github.DependabotAlert {
	var allAlerts []*github.DependabotAlert

	app.logger.Info("Getting dependabot alerts")

	listAlertOpt := &github.ListAlertsOptions{
		State:       github.String("open"),
		ListOptions: github.ListOptions{PerPage: 100},
	}

	if !repo.GetArchived() { // ListRepoAlerts fails for archived repos

		for {
			alerts, resp, err := app.ghClient.Dependabot.ListRepoAlerts(context.Background(), repo.GetOwner().GetLogin(), repo.GetName(), listAlertOpt)
			if err != nil {
				app.logger.Error(err.Error())
			}
			allAlerts = append(allAlerts, alerts...)
			if resp.NextPage == 0 {
				break
			}
			listAlertOpt.ListOptions.Page = resp.NextPage
		}
	}
	return allAlerts
}

func (app *application) topics(repo *github.Repository) []string {
	app.logger.Info("Getting topics")

	topics, _, err := app.ghClient.Repositories.ListAllTopics(context.Background(), repo.GetOwner().GetLogin(), repo.GetName())
	if err != nil {
		app.logger.Error(err.Error())
	}
	return topics
}
