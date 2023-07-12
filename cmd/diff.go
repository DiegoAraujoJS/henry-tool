/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/DiegoAraujoJS/henry-tool/pkg"
	"github.com/DiegoAraujoJS/henry-tool/utils"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
)

func listAllCommitsDiff(c *object.Commit) error {
    c.Parents().ForEach(func(parent *object.Commit) error {
        diff, err := parent.Patch(c)

        if err != nil {
            fmt.Printf("Failed to get diff: %v\n", err)
        }

        stats := diff.Stats()

        fmt.Printf("Hash: %v, Author: %v, Date: %v\n",c.Hash, c.Author.Name, c.Author.When)
        for _, s := range stats {
            fmt.Printf("\tFile: %v, Additions: %v, Deletions: %v\n", s.Name, s.Addition, s.Deletion)
        }

        return nil
    })
    return nil
}

// diffCmd represents the diff command
var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Este comando muestra los avances de cada persona en términos de líneas de código y files editados",
	Long: `Este comando muestra los avances de cada persona en términos de líneas de código y files editados`,
	Run: func(cmd *cobra.Command, args []string) {
        repo, err := pkg.OpenRepositoryAtRoot()

        if err != nil {
            fmt.Printf("Failed to open repository: %v\n", err)
            return
        }

        yesterday := time.Now().Add(-24 * time.Hour)

        if wide, _ := cmd.Flags().GetBool("wide"); wide {
            utils.LogTimeline(repo, &yesterday, listAllCommitsDiff)
        } else {
            user_additions_deletions := map[string][2]int{}

            utils.LogTimeline(repo, &yesterday, func (c *object.Commit) error {
                c.Parents().ForEach(func(parent *object.Commit) error {
                    diff, err := parent.Patch(c)

                    if err != nil {
                        fmt.Printf("Failed to get diff: %v\n", err)
                    }

                    stats := diff.Stats()

                    for _, s := range stats {
                        if s.Name == "package-lock.json" {
                            continue
                        }
                        if value, ok := user_additions_deletions[c.Author.Name]; !ok {
                            user_additions_deletions[c.Author.Name] = [2]int{0, 0}
                        } else {
                            user_additions_deletions[c.Author.Name] = [2]int{value[0] + s.Addition, value[1] + s.Deletion}
                        }
                    }

                    return nil
                })
                return nil
            })

            display(user_additions_deletions, "Cambios", func(count [2]int) string {
                return fmt.Sprintf("Additions %d\nDeletions %d", count[0], count[1])
            })

        }

        if err != nil {
            fmt.Printf("Failed to iterate over commits: %v\n", err)
            return
        }
	},
}


func init() {
	rootCmd.AddCommand(diffCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// diffCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	diffCmd.Flags().BoolP("wide", "w", false, "Mostrar por cada commit los files modificados y las líneas agregadas y eliminadas")
}
