package cmd

import (
	"context"

	"github.com/shurcooL/githubv4"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"

	"github.com/skarlso/caretaker/pkg/client"
	"github.com/skarlso/caretaker/pkg/logger"
	"github.com/skarlso/caretaker/pkg/moveissue"
)

// CreateMoveIssueCommand gets the issue and updates its status to the desired Status.
func CreateMoveIssueCommand(rootArgs *rootArgsStruct) *cobra.Command {
	moveIssueCmd := &cobra.Command{
		Use:   "move-issue",
		Short: "Moves an issue into a specific column location on a Project Board.",
	}

	moveIssueCmd.RunE = moveIssueRunE(rootArgs)

	return moveIssueCmd
}

func moveIssueRunE(rootArgs *rootArgsStruct) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: rootArgs.token},
		)
		tc := oauth2.NewClient(ctx, ts)
		gclient := githubv4.NewClient(tc)

		// setup logger
		var log logger.Logger = &logger.QuiteLogger{}
		if rootArgs.verbose {
			log = &logger.VerboseLogger{}
		}

		log.Log("running move command")

		client := client.NewCaretaker(log, gclient, client.Options{
			Repo:       rootArgs.repo,
			Owner:      rootArgs.owner,
			StatusName: rootArgs.statusOption,
			Interval:   rootArgs.staleInterval,
			StaleLabel: rootArgs.pullRequestProcessedLabel,
		})
		mover := moveissue.NewMoveIssueAction(log, client, moveissue.Options{
			PullRequestNumber: rootArgs.pullRequestNumber,
			StatusName:        rootArgs.statusOption,
			StaleLabel:        rootArgs.pullRequestProcessedLabel,
		})

		return mover.Move(ctx)
	}
}
