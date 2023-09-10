/*
Copyright © 2023 khml 43922054+khml@users.noreply.github.com
*/

package cmd

import (
	"csv2sqlite/libc2s"
	"github.com/spf13/cobra"
)

var pathToCsv, tableName, pathToDB string

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "read CSV and convert to SQLite DB",
	Long: `read CSV and convert to SQLite DB
Usage:
    ./csv2sqlite read <csv_file_path> <table_name> <database_file_path>
`,

	Run: func(cmd *cobra.Command, args []string) {
		libc2s.Csv2sqlite(pathToCsv, pathToDB, tableName)
	},
}

func init() {
	rootCmd.AddCommand(readCmd)
	readCmd.Flags().StringVarP(&pathToCsv, "csv", "c", "", "path to csv")
	readCmd.Flags().StringVarP(&pathToDB, "db", "d", "", "path to sqlite db file")
	readCmd.Flags().StringVarP(&tableName, "table", "t", "", "table name")

	readCmd.MarkFlagRequired("csv")
	readCmd.MarkFlagRequired("db")
	readCmd.MarkFlagRequired("table")
}
