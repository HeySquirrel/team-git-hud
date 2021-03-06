package cmd

import (
	"fmt"
	"os"

	"github.com/HeySquirrel/tribe/blame"
	"github.com/HeySquirrel/tribe/blame/model"
	"github.com/HeySquirrel/tribe/work"
	"github.com/spf13/cobra"
)

var endpoints []int

var blameCmd = &cobra.Command{
	Use:   "blame",
	Short: "Why the @*$% does this code exist?",
	Long:  `Access historical work items or issues, frequent contributors and your entire git history with one simple command so that you quickly determine why a line of code exists.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]
		var start, end int

		switch len(endpoints) {
		case 1:
			start = endpoints[0]
			end = endpoints[0]
		case 2:
			start = endpoints[0]
			end = endpoints[1]
		default:
			cmd.Help()
			os.Exit(1)
		}

		api, err := work.NewItemServer()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		annotate := model.NewCachingAnnotate(model.NewAnnotate(api))

		blame := blame.NewApp(annotate)
		defer blame.Close()

		blame.SetFile(filename, start, end)
		blame.Loop()
	},
}

func init() {
	RootCmd.AddCommand(blameCmd)
	blameCmd.Flags().IntSliceVarP(&endpoints, "lines", "L", []int{1, 1}, "line numbers to blame")
}
