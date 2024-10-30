package cleaner

import (
	"context"

	"github.com/loureirovinicius/cleanup/config"
	"github.com/loureirovinicius/cleanup/helpers/logger"
	"github.com/loureirovinicius/cleanup/provider"
	"github.com/spf13/cobra"
)

var (
	services map[string]provider.Cleanable
	ctx      context.Context
	debug    bool
	rootCmd  = &cobra.Command{
		Use:   "cleanup",
		Short: "Cleanup - Cloud Provider Sanitization tool",
		Long:  "Cleanup is a tool designed to accomplish effective costs on Cloud Providers (AWS, GCP, etc...) without wasting money on unused resources - an empty Load Balancer, for example. Such tool was thought to be one of the greatest allies in a FinOps culture for its simplicity, efficiency and security.",
	}

	listCommand = &cobra.Command{
		Use:   "list",
		Short: "Lists all the created resources for a certain provider's service",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// args[0] = service name (like ebs, eni, etc...)
			var err error

			debug, err = cmd.Flags().GetBool("debug")
			if err != nil {
				logger.Log(ctx, "error", err.Error())
			}

			service := args[0]

			setup()
			err = list(ctx, services[service])
			if err != nil {
				logger.Log(ctx, "error", err.Error())
			}
		},
	}

	validateCommand = &cobra.Command{
		Use:   "validate",
		Short: "Validates if resources can be deleted or not",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// args[0] = service name (like ebs, eni, etc...)
			var err error

			debug, err = cmd.Flags().GetBool("debug")
			if err != nil {
				logger.Log(ctx, "error", err.Error())
			}

			service := args[0]

			setup()
			err = validate(ctx, services[service])
			if err != nil {
				logger.Log(ctx, "error", err.Error())
			}
		},
	}

	deleteCommand = &cobra.Command{
		Use:   "delete",
		Short: "Deletes the unused resource",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// args[0] = service name (like ebs, eni, etc...)
			var err error

			debug, err = cmd.Flags().GetBool("debug")
			if err != nil {
				logger.Log(ctx, "error", err.Error())
			}

			service := args[0]

			setup()
			err = delete(ctx, services[service])
			if err != nil {
				logger.Log(ctx, "error", err.Error())
			}
		},
	}
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enables debug mode")
	rootCmd.AddCommand(listCommand, validateCommand, deleteCommand)
}

func setup() {
	if debug {
		logger.SetLevel("debug")
	}

	// Start initialization of configuration
	logger.Log(ctx, "info", "Initializing configs...")
	err := config.Start()
	if err != nil {
		logger.Log(ctx, "error", err.Error())
		return
	}
	logger.Log(ctx, "info", "Configs were initialized successfully!")

	// Load the provider configuration
	logger.Log(ctx, "info", "Initializing provider's services...")
	services, err = provider.LoadProvider(ctx, "aws")
	if err != nil {
		logger.Log(ctx, "error", err.Error())
		return
	}
	logger.Log(ctx, "info", "Services were initialized successfully!")
}

func Run() error {
	ctx = context.Background()

	if err := rootCmd.Execute(); err != nil {
		return err
	}

	return nil
}
