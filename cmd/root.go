package cmd

import (
	"fmt"
	"os"

	"github.com/opsline/swizzle"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var cfgFile string
var cfg = swizzle.Config{}

// RootCmd sets up cobra cli app
var RootCmd = &cobra.Command{
	Use:   "swizzle",
	Short: "A simple webapp for testing deployments",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		injectViper(viper.GetViper(), cmd)
		os.Setenv("AWS_REGION", cfg.S3Config.Region)
		if !cmd.Flag("s3_file").Changed {
			// not sure why this was getting set to ""
			cfg.S3Config.File = cmd.Flag("s3_file").DefValue
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// run by default
		swizzle.RunWeb(&cfg)
	},
}

// Execute starts our app
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.swizzle.yaml)")
	RootCmd.PersistentFlags().IntVarP(&cfg.Port, "port", "p", 8080, "port number to listen on")

	RootCmd.PersistentFlags().StringVar(&cfg.S3Config.Bucket, "s3_bucket", "com.opsline.swizzle", "S3 bucket to read from")
	RootCmd.PersistentFlags().StringVar(&cfg.S3Config.File, "s3_file", "hello.txt", "S3 file name to read")

	RootCmd.PersistentFlags().StringVar(&cfg.S3Config.Region, "aws_region", "us-east-1", "AWS region")
	RootCmd.PersistentFlags().StringVar(&cfg.S3Config.AccessKey, "aws_access_key", "", "AWS access key (default: read ~/.aws/credentials)")
	RootCmd.PersistentFlags().StringVar(&cfg.S3Config.SecretKey, "aws_secret_key", "", "AWS secret key (default: read ~/.aws/credentials)")
	RootCmd.PersistentFlags().BoolVar(&cfg.S3Config.UseMetadata, "aws_use_metadata", false, "Load AWS credentials from instance metadata endpoint (default: disabled)")

	RootCmd.PersistentFlags().StringVar(&cfg.Redis.Host, "redis_host", "localhost", "Redis hostname")
	RootCmd.PersistentFlags().IntVar(&cfg.Redis.Port, "redis_port", 6379, "Redis port number")

	RootCmd.PersistentFlags().StringVar(&cfg.Pg.URI, "pg_uri", "postgres://postgres:postgres@localhost/postgres?sslmode=disable", "pgsql connection uri")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".swizzle")
	viper.AddConfigPath(os.Getenv("HOME"))
	viper.SetEnvPrefix("SWIZZLE")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// Read explicitly set values from viper and override Flags
// values with the same long-name if they were not explicitly set
// via cmd line
//
// via https://github.com/spf13/cobra/issues/367
func injectViper(cmdViper *viper.Viper, cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if !f.Changed {
			if cmdViper.IsSet(f.Name) {
				cmd.Flags().Set(f.Name, cmdViper.GetString(f.Name))
			}
		}
	})
}
