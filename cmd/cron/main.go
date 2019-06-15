package main

import (
	"context"
	"fmt"
	"os"

	"github.com/olcolabs/roomba/config"
	roomba "github.com/olcolabs/roomba/lib"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func main() {
	conf, err := config.Load(getEnv())
	if err != nil {
		panic(err)
	}

	if conf.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	// Github auth
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	ghClient := githubv4.NewClient(httpClient)

	checkOpenPRs(ghClient, conf)
	//addLabels()
	setupProjects(ghClient, conf)
}

// setupProjects will add the required labels to all repos if they dont yet exist
func setupProjects(ghClient *githubv4.Client, conf config.Config) {
	//---------------------- GET LABELS FROM REPO ----------------------
	repoName := "personal_capital_foreign_currency"
	// GraphQL query
	var q struct {
		//TODO: update with repo names
		Repository roomba.Repository `graphql:"repository(owner:amokhtar, name:personal_capital_foreign_currency)"`
	}
	// results gets mapped into `q`
	err := ghClient.Query(context.Background(), &q, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to reach github")
		return
	}
	log.Debug().Msgf("Response: %+v", q)
	labelSet := roomba.GetLabelTypeSetFromRepository(q.Repository)

	//---------------------- Get Missing labels (labelsToAdd) ----------------------
	//TODO: update with actual count
	if len(labelSet) == 4 {
		log.Debug().Msgf("All labels already present for repo +v%", "personal_capital_foreign_currency")
	} else {
		labelsToAdd := make([]githubv4.String, 0)
		if !labelSet[roomba.NeedsOneReview] {
			labelsToAdd = append(labelsToAdd, roomba.NeedsOneReview)
		}
		if !labelSet[roomba.NeedsTwoReviews] {
			labelsToAdd = append(labelsToAdd, roomba.NeedsTwoReviews)
		}
		if !labelSet[roomba.READY] {
			labelsToAdd = append(labelsToAdd, roomba.READY)
		}
		if !labelSet[roomba.WIP] {
			labelsToAdd = append(labelsToAdd, roomba.WIP)
		}
		//---------------------- Update labels to Repo ----------------------
		//TODO: update with repo name
		addLabels(labelsToAdd, repoName)
	}
}

func checkOpenPRs(ghClient *githubv4.Client, conf config.Config) {
	slackSvc, err := roomba.NewSlackSvc(conf)
	if err != nil {
		log.Fatal().Err(err).Msg("Can create slack service")
	}

	// GraphQL query
	var q struct {
		Search roomba.Search `graphql:"search(query:$query, type:ISSUE, first:30)"`
	}
	vars := map[string]interface{}{
		"query": githubv4.String(fmt.Sprintf("is:pr is:open user:%s", conf.Organization)),
	}

	// results gets mapped into `q`
	err = ghClient.Query(context.Background(), &q, vars)
	if err != nil {
		log.Error().Err(err).Msg("Failed to reach github")
		return
	}

	if len(q.Search.Edges) < 1 {
		return
	}

	for _, record := range q.Search.Edges {
		labelSet := roomba.GetLabelTypeSet(record.Node.PullRequest.Labels)
		if !roomba.HasValidLabel(labelSet) {
			log.Debug().Msgf("No Label: %+v", record)
			//addLabels()
		}
	}

	// parse results and report to slack
	err = slackSvc.Report(q.Search.Edges)
	if err != nil {
		log.Error().Err(err).Msg("Failed to issue PullRequest report")
		return
	}
}

func addLabels(labelsToAdd []githubv4.String, repoName string) {

	log.Info().Msgf("Adding Labels s% to Repo:s%", labelsToAdd, repoName)
}

// helper function to figure environment
func getEnv() string {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}
	return env
}
