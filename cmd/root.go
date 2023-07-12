/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/DiegoAraujoJS/henry-tool/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// Function display is a generic function to display a map of committer to a value. It will display the committers in alphabetical order and will write on each cell the value of the passed function.
func display [T any] (sum_map map[string]T, row_name string, cell_format func(T) string) {
    var committers []string
    for committer := range sum_map {
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

        sliced_committers := committers[start:end]
        table := tablewriter.NewWriter(os.Stdout)
        // Set the header to the sliced_committers.
        table.SetHeader(append([]string{"Row"}, sliced_committers...))

        // Get the commit counts for each committer in the same order as the header.
        var counts []string
        for _, committer := range sliced_committers {
            count := sum_map[committer]
            counts = append(counts, cell_format(count))
        }

        // Add a single row with the commit counts.
        var rows = [][]string{
            append([]string{row_name}, counts...),
        }
        for _, row := range rows {
            table.Append(row)
        }
        table.Render() // Send

    }
}
// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "henry-tool",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
    err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.henry-tool.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


