package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func cobraAutoCompleteFileName(_ *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	cfg := GetConfig()
	return AutoCompleteFileName(cfg.SnippetsDir, cfg.Exclude, toComplete)
}

// TODO: Read config from this config file. default: ~/.snip/config.yaml.
var cfgFile string

var rootCmd = &cobra.Command{
	Use:               "snip",
	Short:             "snip is a snippet manager on the command line.",
	Long:              `snip is a snippet manager on the command line.`,
	Args:              cobra.ExactArgs(1),
	RunE:              CmdViewSnippet,
	ValidArgsFunction: cobraAutoCompleteFileName,
}

var dirCmd = &cobra.Command{
	Use:   "dir [subPath]",
	Short: "cd into the snippets directory",
	Args:  cobra.MaximumNArgs(1),
	RunE:  CmdGoToSnippetsDir,
}

var openCmd = &cobra.Command{
	Use:               "open",
	Short:             "open the snippet in the editor",
	Args:              cobra.ExactArgs(1),
	RunE:              CmdOpenSnippet,
	ValidArgsFunction: cobraAutoCompleteFileName,
}

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "pull the snippets from git",
	RunE:  CmdPullSnippets,
}

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "push the snippets into the git repository",
	Args:  cobra.MaximumNArgs(1),
	RunE:  CmdPushSnippets,
}

var editorCmd = &cobra.Command{
	Use:   "editor",
	Short: "Opens the snippets directory in your editor",
	RunE:  CmdOpenEditor,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.AddCommand(dirCmd, openCmd, pullCmd, pushCmd, editorCmd)
}

func run() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func CmdViewSnippet(_ *cobra.Command, args []string) error {
	cfg := GetConfig()
	return cfg.ViewerCmd(cfg.SnippetPath(args[0])).Run()
}

func CmdGoToSnippetsDir(_ *cobra.Command, args []string) error {
	subPath := ""
	if len(args) != 0 {
		subPath = filepath.Dir(args[0])
	}

	fmt.Print(filepath.Join(GetConfig().SnippetsDir, subPath))
	return nil
}

func CmdOpenSnippet(_ *cobra.Command, args []string) error {
	cfg := GetConfig()
	fpath := cfg.SnippetPath(args[0])

	// Make parent directories
	if err := os.MkdirAll(filepath.Dir(fpath), 0777); err != nil {
		log.Fatal(err)
	}

	return Command(cfg.Editor, fpath).Run()
}

func CmdPullSnippets(_ *cobra.Command, args []string) error {
	cfg := GetConfig()
	return Command(cfg.Git, "-C", cfg.SnippetsDir, "pull", "origin").Run()
}

func CmdPushSnippets(_ *cobra.Command, args []string) error {
	cfg := GetConfig()

	msg := "Update snippets"
	if len(args) > 0 {
		msg = args[0]
	}

	err := Command(cfg.Git, "-C", cfg.SnippetsDir, "commit", "-Am", msg).Run()
	if err != nil {
		return err
	}

	return Command(cfg.Git, "-C", cfg.SnippetsDir, "push", "origin").Run()
}

func CmdOpenEditor(_ *cobra.Command, args []string) error {
	cfg := GetConfig()
	return Command(cfg.Editor, cfg.SnippetsDir).Run()
}
