package cmd

import (
	"fmt"
	"strings"
	"willbehn/ht/internal"
	"willbehn/ht/models"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "A brief description of your command",
	Long:  `Fiks senere`,
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := internal.OpenDB()

		if err != nil {
			return err
		}
		defer db.Close()

		var conditions []string
		var parameters []any

		query := `SELECT id, cmd, shell, dir, repo, branch, ts, exit_code, duration_ms  
		FROM commands `

		for _, arg := range args {
			conditions = append(conditions, " cmd LIKE ?")
			parameters = append(parameters, "%"+arg+"%")
		}

		if len(conditions) > 0 {
			query += "WHERE" + strings.Join(conditions, "AND")
		}

		query += "ORDER BY ts"

		rows, err := db.Query(query, parameters...)

		if err != nil {
			return err
		}

		var results []models.CmdEvent

		for rows.Next() {
			var ev models.CmdEvent

			if err := rows.Scan(
				&ev.Id,
				&ev.Cmd,
				&ev.Shell,
				&ev.Dir,
				&ev.Repo,
				&ev.Branch,
				&ev.TS,
				&ev.Exit,
				&ev.Dur,
			); err != nil {
				return err
			}
			results = append(results, ev)
		}

		resultCount := len(results)
		fmt.Printf("%d results found - newest first\n", resultCount)

		isLong, _ := cmd.Flags().GetBool("long")

		if isLong {
			internal.ResultOutputLong(results)

		} else {
			internal.ResultOutputShort(results)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
