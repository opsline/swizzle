package cmd

import (
	"github.com/opsline/swizzle"
	"github.com/spf13/cobra"
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Run web server",
	Run: func(cmd *cobra.Command, args []string) {
		swizzle.RunWeb(&cfg)
	},
}

func init() {
	RootCmd.AddCommand(webCmd)
}
