/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"
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

        since, err := cmd.Flags().GetDuration("since")

        if err != nil {
            fmt.Printf("Failed to get duration: %v\n", err)
            return
        }
        yesterday := time.Now().Add(-since)


        if wide, _ := cmd.Flags().GetBool("wide"); wide {
            utils.LogTimeline(repo, &yesterday, listAllCommitsDiff)
            return
        }

        user_additions_deletions := map[string]map[string][2]int{}

        utils.LogTimeline(repo, &yesterday, func (c *object.Commit) error {
            c.Parents().ForEach(func(parent *object.Commit) error {
                diff, err := parent.Patch(c)

                if err != nil {
                    fmt.Printf("Failed to get diff: %v\n", err)
                }

                stats := diff.Stats()

                for _, s := range stats {
                    if strings.Contains(s.Name, "package-lock.json") {
                        continue
                    }
                    if _, ok := user_additions_deletions[c.Author.Name]; !ok {
                        user_additions_deletions[c.Author.Name] = map[string][2]int{}
                    }
                    if value, ok := user_additions_deletions[c.Author.Name][s.Name]; !ok {
                        user_additions_deletions[c.Author.Name][s.Name] = [2]int{0, 0}
                    } else {
                        user_additions_deletions[c.Author.Name][s.Name] = [2]int{value[0] + s.Addition, value[1] + s.Deletion}
                    }

                }

                return nil
            })
            return nil
        })

        if short, _ := cmd.Flags().GetBool("short"); short {
            display(user_additions_deletions, "Cambios", func(changes map[string][2]int) string {
                additions := 0
                deletions := 0
                for _, change := range changes {
                    additions += change[0]
                    deletions += change[1]
                }
                return fmt.Sprintf("Additions: %d\nDeletions: %d\n", additions, deletions)
            })
            return
        } else {
            display(user_additions_deletions, "Cambios", func(changes map[string][2]int) string {
                list := ""
                for file, change := range changes {
                    list += fmt.Sprintf("%v %v %v\n\n", getFilename(file), change[0], change[1])
                }
                return list
            })
        }




        if err != nil {
            fmt.Printf("Failed to iterate over commits: %v\n", err)
            return
        }
    },
}

func getFilename(path string) string {
    last_instance := 0
    for i, char := range path {
        if char == '/' {
            last_instance = i
        }
    }
    return path[last_instance:]
}

func init() {
	diffCmd.Flags().BoolP("wide", "w", false, "Mostrar por cada commit los files modificados y las líneas agregadas y eliminadas")
    diffCmd.Flags().BoolP("short", "s", false, "Mostrar por cada persona el total de líneas agregadas y eliminadas")
    diffCmd.Flags().Duration("since", 24 * time.Hour, "Mostrar los commits desde hace un día")

	rootCmd.AddCommand(diffCmd)
}
