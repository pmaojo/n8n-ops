package cmd

import (
        "fmt"
        "os"

        "github.com/spf13/cobra"
        "github.com/spf13/viper"
        "github.com/n8n-workflows/cli/internal/config"
        "github.com/n8n-workflows/cli/internal/utils"
)

var (
        cfgFile     string
        environment string
        verbose     bool
        logger      = utils.NewLogger()
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
        Use:   "n8n-cli",
        Short: "A collaborative n8n workflow management CLI tool",
        Long: `n8n-cli is a command-line tool for managing n8n workflows across multiple environments.
It supports syncing workflows from n8n instances, deploying local changes, validating workflow files,
and integrating with GitLab CI/CD pipelines for collaborative development.

Examples:
  n8n-cli sync --env development    # Sync workflows from development environment
  n8n-cli deploy --env staging      # Deploy workflows to staging environment  
  n8n-cli validate ./workflows/     # Validate workflow files
  n8n-cli rollback --env production # Rollback to previous deployment`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
        return rootCmd.Execute()
}

func init() {
        cobra.OnInitialize(initConfig)

        // Global flags
        rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.n8n-cli.yaml)")
        rootCmd.PersistentFlags().StringVarP(&environment, "env", "e", "development", "target environment (development, staging, production)")
        rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

        // Bind flags to viper
        viper.BindPFlag("environment", rootCmd.PersistentFlags().Lookup("env"))
        viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
        if cfgFile != "" {
                // Use config file from the flag.
                viper.SetConfigFile(cfgFile)
        } else {
                // Find home directory.
                home, err := os.UserHomeDir()
                cobra.CheckErr(err)

                // Search config in home directory with name ".n8n-cli" (without extension).
                viper.AddConfigPath(home)
                viper.AddConfigPath(".")
                viper.SetConfigType("yaml")
                viper.SetConfigName(".n8n-cli")
        }

        viper.AutomaticEnv() // read in environment variables that match

        // If a config file is found, read it in.
        if err := viper.ReadInConfig(); err == nil {
                if verbose {
                        fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
                }
        }

        // Initialize configuration
        config.InitConfig()
        
        // Set logger verbosity
        if viper.GetBool("verbose") {
                utils.SetLogLevel("debug")
        }
}
