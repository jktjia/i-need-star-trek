/*
Copyright © 2026 jktjia <EMAIL ADDRESS>
*/
package cmd

import (
	"i-need-star-trek/internal"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	dateFormats = []string{"2006", "Jan 2006", "January 2006", "1 Jan 2006", "Jan 1 2006", "1 January 2006", "January 1 2006", time.DateOnly}

	shows   []string
	species []string
	other   []string

	after  time.Time
	before time.Time
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "i-need-star-trek",
	Short: "random Star Trek episode selector",
	Long: `This application is a tool to select random Star Trek episodes to watch.
Use flags to limit the random selection to specific shows and categories`,
	Args: cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return internal.RandomEp(shows, species, other, after, before)
	},
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

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.i-need-star-trek.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringSliceVarP(&shows, "shows", "s", []string{}, "show to filter on")
	rootCmd.Flags().StringSliceVarP(&species, "species", "p", []string{}, "major species in the episode")
	rootCmd.Flags().StringSliceVarP(&other, "tag", "t", []string{}, "episode tags to filter on")

	rootCmd.Flags().TimeVar(&after, "start-date", time.Date(1960, time.January, 1, 0, 0, 0, 0, time.UTC), dateFormats, "earliest date to include")
	rootCmd.Flags().TimeVar(&before, "end-date", time.Now(), dateFormats, "latest date to include")
}
