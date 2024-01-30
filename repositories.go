package main

import (
	"sort"

	"github.com/google/go-github/v58/github"
)

type Totals struct {
	Critical int `json:"critical"`
	High     int `json:"high"`
	Moderate int `json:"moderate"`
	Low      int `json:"low"`
	Total    int `json:"total"`
}

type SeverityLevel string

const (
	Critical SeverityLevel = "critical"
	High     SeverityLevel = "high"
	Moderate SeverityLevel = "medium"
	Low      SeverityLevel = "low"
)

type Repository struct {
	Name       string `json:"name"`
	Archived   bool   `json:"archived"`
	Topics     []string
	Compliance string `json:"compliance"`
	Alerts     []*github.DependabotAlert
	Totals     Totals `json:"totals"`
}

type OrderedMapRepositories struct {
	RepoMap      map[string]*Repository
	OrderedNames []string
}

func (o *OrderedMapRepositories) Init(ghRepos []*github.Repository) {
	o.RepoMap = make(map[string]*Repository)

	for _, repo := range ghRepos {
		o.OrderedNames = append(o.OrderedNames, repo.GetName())
		o.RepoMap[repo.GetName()] = &Repository{
			Name:       repo.GetName(),
			Archived:   repo.GetArchived(),
			Compliance: "unknown",
		}
	}
}

func (o *OrderedMapRepositories) order(sortBy string) {
	switch sortBy {
	case "compliance":
		sort.Slice(o.OrderedNames, func(i, j int) bool {
			return o.RepoMap[o.OrderedNames[i]].Compliance < o.RepoMap[o.OrderedNames[j]].Compliance
		})
	case "archived":
		sort.Slice(o.OrderedNames, func(i, j int) bool {
			return !o.RepoMap[o.OrderedNames[i]].Archived
		})
	case "name":
		sort.Slice(o.OrderedNames, func(i, j int) bool {
			return o.RepoMap[o.OrderedNames[i]].Name < o.RepoMap[o.OrderedNames[j]].Name
		})
	case "total":
		sort.Slice(o.OrderedNames, func(i, j int) bool {
			return o.RepoMap[o.OrderedNames[i]].Totals.Total > o.RepoMap[o.OrderedNames[j]].Totals.Total
		})
	case "critical":
		sort.Slice(o.OrderedNames, func(i, j int) bool {
			return o.RepoMap[o.OrderedNames[i]].Totals.Critical > o.RepoMap[o.OrderedNames[j]].Totals.Critical
		})
	case "high":
		sort.Slice(o.OrderedNames, func(i, j int) bool {
			return o.RepoMap[o.OrderedNames[i]].Totals.High > o.RepoMap[o.OrderedNames[j]].Totals.High
		})
	case "medium":
		sort.Slice(o.OrderedNames, func(i, j int) bool {
			return o.RepoMap[o.OrderedNames[i]].Totals.Moderate > o.RepoMap[o.OrderedNames[j]].Totals.Moderate
		})
	case "low":
		sort.Slice(o.OrderedNames, func(i, j int) bool {
			return o.RepoMap[o.OrderedNames[i]].Totals.Low > o.RepoMap[o.OrderedNames[j]].Totals.Low
		})
	}

}

func (o *OrderedMapRepositories) CountAlerts() {
	for _, repo := range o.RepoMap {
		for _, alert := range repo.Alerts {
			switch SeverityLevel(alert.SecurityVulnerability.GetSeverity()) {
			case Critical:
				repo.Totals.Critical++
			case High:
				repo.Totals.High++
			case Moderate:
				repo.Totals.Moderate++
			case Low:
				repo.Totals.Low++
			}
		}
		repo.Totals.Total = repo.Totals.Critical + repo.Totals.High + repo.Totals.Moderate + repo.Totals.Low
	}
}

func (o *OrderedMapRepositories) Compliance() {
	for _, repo := range o.RepoMap {
		for _, topic := range repo.Topics {
			if topic == "glcp-production" || topic == "glcp-not-production" {
				repo.Compliance = topic
				break
			}
			repo.Compliance = "unknown"
		}
	}
}
