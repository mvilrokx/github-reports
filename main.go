package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v58/github"
)

type config struct {
	org               string
	appID             int64
	appInstallID      int64
	privateKeyFile    string
	excludeNonProd    bool
	excludeArchived   bool
	excludeZeroAlerts bool
	sortBy            string
	outputFile        string
}

type application struct {
	config         config
	logger         *slog.Logger
	ghClient       *github.Client
	RepoLsByOrgOpt *github.RepositoryListByOrgOptions
}

func main() {
	var cfg config

	flag.StringVar(&cfg.org, "org", "glcp", "The Github Organization to report on")

	flag.Int64Var(&cfg.appID, "appID", 0, "The Github Application ID")
	flag.Int64Var(&cfg.appInstallID, "appInstallID", 0, "The Github Application ID")
	flag.StringVar(&cfg.privateKeyFile, "privateKeyFile", "", "The Github Application Private Key File")

	flag.BoolVar(&cfg.excludeNonProd, "excludeNonProd", false, "Exclude repositories flagged as glcp-not-production")
	flag.BoolVar(&cfg.excludeArchived, "excludeArchived", false, "Exclude archived repositories")
	flag.BoolVar(&cfg.excludeZeroAlerts, "excludeZeroAlerts", false, "Exclude repositories that have no Dependabot alerts")
	flag.StringVar(&cfg.sortBy, "sortBy", "total", "Which column to sort by (name, compliance, archived, total, critical, high, moderate, low)")
	flag.StringVar(&cfg.outputFile, "outputFile", "output.md", "The name of the output file")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	logger.Info("----------------------------------------")
	logger.Info("Running with the following options:")
	logger.Info(fmt.Sprintf("   excludeNonProd    = %v", cfg.excludeNonProd))
	logger.Info(fmt.Sprintf("   excludeArchived   = %v", cfg.excludeArchived))
	logger.Info(fmt.Sprintf("   excludeZeroAlerts = %v", cfg.excludeZeroAlerts))
	logger.Info(fmt.Sprintf("   sortBy            = %s", cfg.sortBy))
	logger.Info(fmt.Sprintf("   outputFile        = %s", cfg.outputFile))
	logger.Info("----------------------------------------")

	itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, cfg.appID, cfg.appInstallID, cfg.privateKeyFile)

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	ghClient := github.NewClient(&http.Client{Transport: itr})

	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	app := &application{
		config:         cfg,
		logger:         logger,
		ghClient:       ghClient,
		RepoLsByOrgOpt: opt,
	}

	app.Run()

}
