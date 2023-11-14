package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"insightful/src/worker/services"
	"log"
	"strconv"
	"time"
)

var analyzeDataCMD = &cobra.Command{
	Use:   "analyze-data",
	Short: "analyze-data",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// validate args
		if len(args) > 0 {
			if len(args) != 2 {
				log.Println(ctx, "Args is invalid: %v", args)
				return
			}

			numDistribute, err0 := strconv.ParseInt(args[0], 10, 64)
			order, err1 := strconv.ParseInt(args[1], 10, 64)
			if err0 != nil || err1 != nil {
				log.Println(ctx, "Args params is invalid: ", numDistribute, order)
				return
			}
		}

		var (
			startTime     = time.Now()
			rewardService services.AnalyzeDataService
		)

		err := services.GetServiceContainer().Invoke(func(s services.AnalyzeDataService) {
			rewardService = s
		})
		if err != nil {
			log.Println(ctx, "Error when get service container: %v", err)
			return
		}

		if err := rewardService.TriggerAnalyzeData(ctx, args); err != nil {
			log.Println(ctx, "Failed TriggerAnalyzeData, error: %v", err)
		}

		log.Println(ctx, "analyze-data done, took: %v", time.Now().Sub(startTime).String())
	},
}
