package cmd

import (
	"fmt"
	"strconv"
	"willbehn/ht/internal"
	"willbehn/ht/models"

	"github.com/spf13/cobra"
)

var recentCmd = &cobra.Command{
	Use:   "recent",
	Short: "A brief description of your command",
	Long:  `Fiks senere`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var recentCount int

		if len(args) == 0 {
			recentCount = 10
		} else {
			rc, parseErr := strconv.Atoi(args[0])
			if parseErr != nil {
				return fmt.Errorf("invalid count: %w", parseErr)
			}
			recentCount = rc
		}

		fmt.Println(recentCount)

		db, err := internal.OpenDB()

		if err != nil {
			return err
		}

		defer db.Close()

		tx, err := db.BeginTx(cmd.Context(), nil)

		if err != nil {
			return err
		}

		query := `
			SELECT id, cmd, shell, dir, repo, branch, ts, exit_code, duration_ms  
			FROM commands
			ORDER BY ts DESC
			LIMIT ?`

		rows, err := tx.Query(query, recentCount)

		if err != nil {
			tx.Rollback()
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

		if err := tx.Commit(); err != nil {
			return err
		}

		isLong, _ := cmd.Flags().GetBool("long")
		fmt.Println("10 most recent commands")

		if isLong {
			internal.ResultOutputLong(results)

		} else {
			internal.ResultOutputShort(results)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(recentCmd)

}
