package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

    "github.com/lonegunmanb/genv/pkg"
	"github.com/spf13/cobra"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	downloadInstaller, _ := pkg.NewDownloadInstaller("https://github.com/terraform-linters/tflint/releases/download/{{ .Version }}/tflint_{{ .Os }}_{{ .Arch }}.zip", ctx)
	goBuildInstaller := pkg.NewGoBuildInstaller("https://github.com/terraform-linters/tflint.git", "tflint", "", ctx)
	fallbackInstaller := pkg.NewFallbackInstaller(downloadInstaller, goBuildInstaller)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err.Error())
	}
	env := pkg.NewEnv(homeDir, "tflintenv", "tflint", fallbackInstaller)

	// Listen for interrupt signal (Ctrl + C) and cancel the context when received
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			cancel()
		}
	}()

	var rootCmd = &cobra.Command{Use: "tflintenv"}

	var cmdInstall = &cobra.Command{
		Use:   "install [version]",
		Short: "Install a specific version",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			version := args[0]
			fmt.Printf("Installing version: %s\n", version)
			return env.Install(version)
		},
	}

	var cmdUse = &cobra.Command{
		Use:   "use [version]",
		Short: "Use a specific version",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			version := args[0]
			fmt.Printf("Using version: %s\n", version)
			return env.Use(version)
		},
	}

    var cmdBinaryPath = &cobra.Command{
		Use:   "path",
		Short: "Get the full path to current binary",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := env.CurrentBinaryPath()
			if err != nil {
				return err
			}
			if path == nil {
				return fmt.Errorf("no version selected, please run use first")
			}
			fmt.Print(*path)
			return nil
		},
	}

	var cmdUninstall = &cobra.Command{
		Use:   "uninstall [version]",
		Short: "Uninstall a specific version",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			version := args[0]
			fmt.Printf("Uninstalling version: %s\n", version)
			return env.Uninstall(version)
		},
	}

	var cmdList = &cobra.Command{
		Use:   "list",
		Short: "List all installed versions",
		RunE: func(cmd *cobra.Command, args []string) error {
			installed, err := env.ListInstalled()
			if err != nil {
				return err
			}
			for _, i := range installed {
				fmt.Println(i)
			}
			return nil
		},
	}

	rootCmd.AddCommand(cmdInstall, cmdUse, cmdUninstall, cmdList, cmdBinaryPath)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error executing command:", err)
	}
}
