package cmd

import (
        "fmt"
        "os"

        "github.com/spf13/cobra"
        "github.com/spf13/viper"
        "github.com/n8n-workflows/n8n-ops/internal/ascii"
        "github.com/n8n-workflows/n8n-ops/internal/config"
        "github.com/n8n-workflows/n8n-ops/internal/i18n"
        "github.com/n8n-workflows/n8n-ops/internal/utils"
)

var (
        cfgFile     string
        environment string
        verbose     bool
        language    string
        showVersion bool
        demoMode    bool
        logger      = utils.NewLogger()
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
        Use:   "n8n-ops",
        Short: "A collaborative n8n workflow operations tool",
        Long: ascii.N8nLogo() + `

n8n-ops is a command-line tool for managing n8n workflows across multiple environments.
It supports syncing workflows from n8n instances, deploying local changes, validating workflow files,
and integrating with GitLab CI/CD pipelines for collaborative development.

Examples:
  n8n-ops sync --env development    # Sync workflows from development environment
  n8n-ops deploy --env staging      # Deploy workflows to staging environment  
  n8n-ops validate ./workflows/     # Validate workflow files
  n8n-ops rollback --env production # Rollback to previous deployment`,
        Run: func(cmd *cobra.Command, args []string) {
                if showVersion {
                        fmt.Printf("n8n-ops version %s\n", Version)
                        fmt.Printf("Git commit: %s\n", GitCommit)
                        fmt.Printf("Build date: %s\n", BuildDate)
                        return
                }
                cmd.Help()
        },
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
        return rootCmd.Execute()
}

// Version information - will be set at build time
var (
        Version   = "1.0.0"
        GitCommit = "unknown"
        BuildDate = "unknown"
)

func init() {
        cobra.OnInitialize(initConfig)

        // Global flags
        rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.n8n-ops.yaml)")
        rootCmd.PersistentFlags().StringVarP(&environment, "env", "e", "development", "target environment (development, staging, production)")
        rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
        rootCmd.PersistentFlags().StringVarP(&language, "lang", "l", "en", "language (en, es)")
        rootCmd.PersistentFlags().BoolVar(&demoMode, "demo", false, "use demo mode with mock n8n server")
        rootCmd.Flags().BoolVar(&showVersion, "version", false, "show version information")

        // Bind flags to viper
        viper.BindPFlag("environment", rootCmd.PersistentFlags().Lookup("env"))
        viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
        viper.BindPFlag("language", rootCmd.PersistentFlags().Lookup("lang"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
        // Set language first
        if language != "" {
                i18n.SetLanguage(language)
        } else if viper.IsSet("language") {
                i18n.SetLanguage(viper.GetString("language"))
        }
        if cfgFile != "" {
                // Use config file from the flag.
                viper.SetConfigFile(cfgFile)
        } else {
                // Find home directory.
                home, err := os.UserHomeDir()
                cobra.CheckErr(err)

                // Search config in home directory with name ".n8n-ops" (without extension).
                viper.AddConfigPath(home)
                viper.AddConfigPath(".")
                viper.SetConfigType("yaml")
                viper.SetConfigName(".n8n-ops")
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
