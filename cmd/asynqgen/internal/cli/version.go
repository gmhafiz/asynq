package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().StringVarP(&Domain, "version", "v", "", "asynqgen version")
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of asynq_gen",
	Long:  `All software has versions. This is asynqgen's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("asynq task scaffolder v0.1.0")
	},
}
