/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/DiegoAraujoJS/henry-tool/pkg"
	"github.com/DiegoAraujoJS/henry-tool/utils"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
)

// chartCmd represents the chart command
var commitsCmd = &cobra.Command{
	Use:   "commits",
	Short: "Este comando muestra cuantos commits tiene cada persona",
	Long: `Este comando muestra cuantos commits tiene cada persona`,
    Run: func(cmd *cobra.Command, args []string) {
        // Open the repository.
        repo, err := pkg.OpenRepositoryAtRoot()

        if err != nil {
            fmt.Printf("Failed to open repository: %v\n", err)
            return
        }

        since, err := cmd.Flags().GetDuration("since")
        var yesterday *time.Time
        if err != nil {
            fmt.Printf("Failed to get duration: %v\n", err)
            return
        }
        if since != 0 {
            yst := time.Now().Add(-since)
            yesterday = &yst
        }

        commitsByCommitter := make(map[string]int)
        err = utils.LogTimeline(repo, yesterday, func(c *object.Commit) error {
            committer := c.Committer.Name
            commitsByCommitter[committer]++
            return nil
        })
        if err != nil {
            fmt.Printf("Failed to iterate over commits: %v\n", err)
            return
        }

        display(commitsByCommitter, "Commits", strconv.Itoa)
    },
}

func init() {
	rootCmd.AddCommand(commitsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// chartCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// chartCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
