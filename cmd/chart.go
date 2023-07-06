/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/DiegoAraujoJS/henry-tool/pkg"
	"github.com/DiegoAraujoJS/henry-tool/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type CommitterInfo struct {
	CommitCount int
	FilesChanged map[string]int
}

// chartCmd represents the chart command
var chartCmd = &cobra.Command{
	Use:   "chart",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
    Run: func(cmd *cobra.Command, args []string) {
        // Open the repository.
        repo, err := pkg.OpenRepositoryAtRoot()

        if err != nil {
            fmt.Printf("Failed to open repository: %v\n", err)
            return
        }

        // Get the HEAD reference.
        ref, err := repo.Head()
        if err != nil {
            fmt.Printf("Failed to get HEAD: %v\n", err)
            return
        }

        // Get an iterator over the repository's commit history.
        iter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
        if err != nil {
            fmt.Printf("Failed to get log: %v\n", err)
            return
        }


        

        

        
        
        
        
        
        
        
        
        

        

        
        
        
        
        
        
        
        
        
        
        
        
        
        
        

        
        
        
        
        
        
        
        
        if err != nil {
            fmt.Printf("Failed to iterate over commits: %v\n", err)
            return
        }

        // Map to store the committer and their corresponding commit count.
        commitsByCommitter := make(map[string]int)

        // Iterate over each commit.
        err = iter.ForEach(func(c *object.Commit) error {
            committer := c.Committer.String()
            commitsByCommitter[committer]++
            return nil
        })
        if err != nil {
            fmt.Printf("Failed to iterate over commits: %v\n", err)
            return
        }

        // Get a sorted slice of the committers.
        var committers []string
        for committer := range commitsByCommitter {
            committers = append(committers, committer)
        }
        committers = utils.MergeSort(committers, func(a string, b string) bool {
            return a < b
        })

        // Maximum number of committers to show in each row
        const maxCommittersPerRow = 5

        for start := 0; start < len(committers); start += maxCommittersPerRow {
            end := start + maxCommittersPerRow
            if end > len(committers) {
                end = len(committers)
            }
            displayCommitters(commitsByCommitter, committers[start:end])
        }

    },
}


func displayCommitters(commitsByCommitter map[string]int, committers []string) {
	table := tablewriter.NewWriter(os.Stdout)
	// Set the header to the committers.
	table.SetHeader(append([]string{"Row"}, committers...))

	// Get the commit counts for each committer in the same order as the header.
	var counts []string
	for _, committer := range committers {
		count := commitsByCommitter[committer]
		counts = append(counts, strconv.Itoa(count))
	}

	// Add a single row with the commit counts.
	table.Append(append([]string{"Commits"}, counts...))
	table.Render() // Send

}




func init() {
	rootCmd.AddCommand(chartCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// chartCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// chartCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
