package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(analyzeDataCMD)
}

var rootCmd = &cobra.Command{
	Use:   "insightful-job",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println("Error when execute command, detail: ", err)
		os.Exit(0)
	}
}
