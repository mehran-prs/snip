package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"
)

func cobraAutoCompleteFileName(_ *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return AutoCompleteFileName(Cfg.SnippetsDir, Cfg.Exclude, toComplete)
}

var flagConfigFile string
var flagLogLevel string

var completionCmd = &cobra.Command{
	Use:                   "completion [bash|zsh|fish|powershell]",
	Short:                 "Generate completion script",
	Long:                  completionDocs(rootCmd.Root().Name()),
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE:                  CmdCompletionGenerator,
}

var rootCmd = &cobra.Command{
	Use:               "snip",
	Short:             "snip is a snippet manager on the command line.",
	Long:              `snip is a snippet manager on the command line.`,
	Args:              cobra.ExactArgs(1),
	RunE:              CmdViewSnippet,
	PersistentPreRunE: boot,
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
	rootCmd.PersistentFlags().StringVarP(&flagConfigFile, "config", "c", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.PersistentFlags().StringVarP(&flagLogLevel, "log", "l", "", "Set log level. default is warning. values: debug,info,warn,error")
	rootCmd.AddCommand(completionCmd, dirCmd, openCmd, pullCmd, pushCmd, editorCmd)
}

func run() {
	if err := rootCmd.Execute(); err != nil {
		Error("app error", err)
		os.Exit(1)
	}
}

// boot boots the app. it loads config for us.
func boot(_ *cobra.Command, _ []string) error {
	// Set log level if it's given as a flag
	if flagLogLevel != "" {
		SetLogLevel(LogLevelFromString(flagLogLevel))
	}

	if flagConfigFile == "" {
		flagConfigFile = path.Join(userHomeDir(), ".snip/config.yaml")
	}

	if err := loadConfig(flagConfigFile, flagLogLevel); err != nil {
		return err
	}

	SetLogLevel(LogLevelFromString(Cfg.LogLevel))
	return nil
}

func CmdCompletionGenerator(cmd *cobra.Command, args []string) error {
	switch args[0] {
	case "bash":
		if err := cmd.Root().GenBashCompletionV2(os.Stdout, true); err != nil {
			return err
		}
	case "zsh":
		if err := cmd.Root().GenZshCompletion(os.Stdout); err != nil {
			return err
		}
		fmt.Print("\n\n")
		// In zsh we support fzf too:
		if err := WriteAutocompleteScript(os.Stdout, "fzf.zsh"); err != nil {
			return err
		}
	case "fish":
		if err := cmd.Root().GenFishCompletion(os.Stdout, true); err != nil {
			return err
		}
	case "powershell":
		if err := cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout); err != nil {
			return err
		}
	}
	return nil
}

func CmdViewSnippet(_ *cobra.Command, args []string) error {
	return Cfg.ViewerCmd(Cfg.SnippetPath(args[0])).Run()
}

func CmdGoToSnippetsDir(_ *cobra.Command, args []string) error {
	subPath := ""
	if len(args) != 0 {
		subPath = filepath.Dir(args[0])
	}

	fmt.Print(filepath.Join(Cfg.SnippetsDir, subPath))
	return nil
}

func CmdOpenSnippet(_ *cobra.Command, args []string) error {
	fpath := Cfg.SnippetPath(args[0])

	// Make parent directories
	if err := os.MkdirAll(filepath.Dir(fpath), 0777); err != nil {
		return fmt.Errorf("can not create snippet directory: %w", err)
	}

	return Command(Cfg.Editor, fpath).Run()
}

func CmdPullSnippets(_ *cobra.Command, args []string) error {
	return Command(Cfg.Git, "-C", Cfg.SnippetsDir, "pull", "origin").Run()
}

func CmdPushSnippets(_ *cobra.Command, args []string) error {
	msg := "Update snippets"
	if len(args) > 0 {
		msg = args[0]
	}

	err := Command(Cfg.Git, "-C", Cfg.SnippetsDir, "commit", "-Am", msg).Run()
	if err != nil {
		return err
	}

	return Command(Cfg.Git, "-C", Cfg.SnippetsDir, "push", "origin").Run()
}

func CmdOpenEditor(_ *cobra.Command, args []string) error {
	return Command(Cfg.Editor, Cfg.SnippetsDir).Run()
}
