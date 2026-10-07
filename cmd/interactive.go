/*
Copyright © 2026 jktjia <EMAIL ADDRESS>
*/
package cmd

import (
	"i-need-star-trek/internal/interactive"

	"github.com/spf13/cobra"
)

// interactiveCmd represents the interactive command
var interactiveCmd = &cobra.Command{
	Use:     "interactive",
	Aliases: []string{"i", "tui"},
	Short:   "Opens up interactive terminal UI",
	Long: `This application is a tool to select random Star Trek episodes to watch.

The interactive version of the app has allows users to view all possible
episodes that fit filters in addition to selecting random episodes.`,
	Run: func(cmd *cobra.Command, args []string) {
		interactive.Start()
	},
}

func init() {
	rootCmd.AddCommand(interactiveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// interactiveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// interactiveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
