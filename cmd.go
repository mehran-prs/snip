package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func cobraAutoCompleteFileName(cmd *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	Error("cfg Dir: ", Cfg.Dir, "os Args: ", os.Args)
	return AutoCompleteFileName(Cfg.Dir, Cfg.Exclude, toComplete)
}

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
	Use:                "snip",
	Short:              "snip is a snippet manager on the command line.",
	Long:               `snip is a snippet manager on the command line.`,
	Args:               cobra.ExactArgs(1),
	RunE:               CmdViewSnippet,
	PersistentPreRunE:  boot,
	PersistentPostRunE: shutdown,
	ValidArgsFunction:  cobraAutoCompleteFileName,
}

var dirCmd = &cobra.Command{
	Use:   "dir [subPath]",
	Short: "cd into the snippets directory",
	Args:  cobra.MaximumNArgs(1),
	RunE:  CmdSnippetsDir,
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
	if err := loadConfig(); err != nil {
		return err
	}

	// Init Logger
	l, err := NewLogger(Cfg.LogTmpFileName, LogLevelFromString(Cfg.LogLevel))
	if err != nil {
		return err
	}

	SetGlobalLogger(l)
	return nil
}
func shutdown(_ *cobra.Command, _ []string) error {
	return GlobalLogger().Shutdown()
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
		fmt.Println(genFzfZshCompletion(cmd.Root().Name()))
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

func CmdSnippetsDir(_ *cobra.Command, args []string) error {
	subPath := ""
	if len(args) != 0 {
		subPath = filepath.Dir(args[0])
	}

	fmt.Println(filepath.Join(Cfg.Dir, subPath))
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
	return Command(Cfg.Git, "-C", Cfg.Dir, "pull", "origin").Run()
}

func CmdPushSnippets(_ *cobra.Command, args []string) error {
	msg := "Update snippets"
	if len(args) > 0 {
		msg = args[0]
	}

	err := Command(Cfg.Git, "-C", Cfg.Dir, "commit", "-Am", msg).Run()
	if err != nil {
		return err
	}

	return Command(Cfg.Git, "-C", Cfg.Dir, "push", "origin").Run()
}

func CmdOpenEditor(_ *cobra.Command, args []string) error {
	return Command(Cfg.Editor, Cfg.Dir).Run()
}
