package cleaner

import (
	"context"
	"fmt"
	"os"

	"github.com/loureirovinicius/cleanup/config"
	"github.com/loureirovinicius/cleanup/helpers/logger"
	"github.com/loureirovinicius/cleanup/providers"
	"github.com/spf13/cobra"
)

var (
	ctx      context.Context
	debug    bool
	output   string
	provider string
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
			serviceName := args[0]

			// Load cloud provider resource that is being verified
			service, err := providers.LoadProvider(ctx, provider, serviceName)
			if err != nil {
				logger.Log(ctx, "error", err.Error())
				return
			}

			// List instances of a determined cloud provider resource
			err = list(ctx, service, serviceName)
			if err != nil {
				logger.Log(ctx, "error", err.Error())
				return
			}
		},
	}

	validateCommand = &cobra.Command{
		Use:   "validate",
		Short: "Validates if resources can be deleted or not",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// args[0] = service name (like ebs, eni, etc...)
			serviceName := args[0]

			// Load cloud provider resource that is being verified
			service, err := providers.LoadProvider(ctx, provider, serviceName)
			if err != nil {
				logger.Log(ctx, "error", err.Error())
				return
			}

			// Validate resources checking if they're unused
			err = validate(ctx, service, serviceName)
			if err != nil {
				logger.Log(ctx, "error", err.Error())
				return
			}
		},
	}

	deleteCommand = &cobra.Command{
		Use:   "delete",
		Short: "Deletes the unused resource",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// args[0] = service name (like ebs, eni, etc...)
			serviceName := args[0]

			// Load cloud provider resource that is being verified
			service, err := providers.LoadProvider(ctx, provider, serviceName)
			if err != nil {
				logger.Log(ctx, "error", err.Error())
				return
			}

			// Delete unused resources found by the execution
			err = delete(ctx, service, serviceName)
			if err != nil {
				logger.Log(ctx, "error", err.Error())
				return
			}
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&provider, "provider", "p", "aws", "Cloud Provider being used during execution")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enables debug mode")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "text", "Chooses between output format (text or JSON)")
	rootCmd.PersistentFlags().BoolP("help", "h", false, "Display help information")
	rootCmd.AddCommand(listCommand, validateCommand, deleteCommand)
}

// Start the cleaner
func Run() error {
	ctx = context.Background()

	// Explicitly parse flags early
	err := rootCmd.PersistentFlags().Parse(os.Args[1:])
	if err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	// Access the parsed flags and get the cloud provider being used based on their values
	provider, err := rootCmd.PersistentFlags().GetString("provider")
	if err != nil {
		return fmt.Errorf("could not get 'provider' flag: %w", err)
	}

	// Access the parsed flags and set the logger based on their values
	debug, err = rootCmd.PersistentFlags().GetBool("debug")
	if err != nil {
		return fmt.Errorf("could not get 'debug' flag: %w", err)
	}

	// Access the parsed flags and set the log format based on their values
	output, err = rootCmd.PersistentFlags().GetString("output")
	if err != nil {
		return fmt.Errorf("could not get 'output' flag: %w", err)
	}

	level := "info"
	// Enable debug logs
	if debug {
		level = "debug"
	}

	logger.InitializeLogger(level, output, os.Stdout)

	// Start initialization of configuration
	logger.Log(ctx, "debug", "Initializing configs...")
	err = config.Start(provider)
	if err != nil {
		return fmt.Errorf("could not initialize configs: %w", err)
	}
	logger.Log(ctx, "debug", "Configs were initialized successfully!")

	return rootCmd.Execute()
}
