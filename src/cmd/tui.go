package cmd

import (
	"csv2sqlite/libc2s/tui"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// tuiCmd represents the tui command
var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Interactive TUI mode for csv2sqlite",
	Long:  `Launch an interactive Terminal User Interface (TUI) to convert CSV files to SQLite databases.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := tui.RunTUI(); err != nil {
			fmt.Printf("Error running TUI: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
