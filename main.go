package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const Version = "" // fill-in at compile time.

var completionCmd = &cobra.Command{
	Use:                   "completion [bash|zsh|fish|powershell]",
	Short:                 "Generate completion script",
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE:                  CmdCompletionGenerator,
}

var rootCmd = &cobra.Command{
	Use:                "snip [command]",
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

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "sync the snippets changes with your git repository",
	Long:  "Sync command first pull changes from yourb git repository and then commit and pushes your changes.",
	Args:  cobra.MaximumNArgs(1),
	RunE:  CmdSync,
}

var editorCmd = &cobra.Command{
	Use:   "editor",
	Short: "Opens the snippets directory in your editor",
	RunE:  CmdOpenEditor,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Run:   func(*cobra.Command, []string) { fmt.Println(DefaultStr(Version, "unknown")) },
}

func init() {
	rootCmd.AddCommand(completionCmd, dirCmd, openCmd, syncCmd, editorCmd, versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		Error("app error: ", err)
		os.Exit(1)
	}
}

// boot boots the app. it loads config for us.
func boot(cmd *cobra.Command, _ []string) error {
	envPrefix := strings.ToUpper(cmd.Root().Name()) + "_" // e.g., SNIP_
	if err := loadConfig(envPrefix); err != nil {
		return err
	}

	// Init Logger
	if Cfg.LogTmpFileName != "" {
		if err := SetLoggerFile(Cfg.LogTmpFileName); err != nil {
			return nil
		}
	}

	return nil
}

func shutdown(_ *cobra.Command, _ []string) error {
	return CloseLoggerFile(os.Stderr)
}

func cobraAutoCompleteFileName(_ *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	exclude := Cfg.Exclude
	searchDir := filepath.Dir(toComplete)
	root := path.Join(Cfg.Dir, searchDir)

	if searchDir != "." && searchDir != "/" { // Currently we support exclude only on the root dir.
		exclude = nil
	}

	res, err := findFiles(root, baseName(toComplete), exclude, searchDir)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	return res, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
}

func CmdCompletionGenerator(cmd *cobra.Command, args []string) error {
	switch args[0] {
	case "bash":
		if err := cmd.Root().GenBashCompletionV2(os.Stdout, true); err != nil {
			return err
		}
		// In bash, we support fzf too:
		fmt.Println(genFzfBashCompletion(cmd.Root().Name()))
	case "zsh":
		if err := cmd.Root().GenZshCompletion(os.Stdout); err != nil {
			return err
		}
		fmt.Print("\n\n")

		// In zsh, we support fzf too:
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

func CmdSync(_ *cobra.Command, args []string) error {
	msg := "snip: Update snippets"
	if len(args) > 0 {
		msg = args[0]
	}

	fmt.Println("pull new changes")
	if err := Command(Cfg.Git, "-C", Cfg.Dir, "pull", "origin").Run(); err != nil {
		return err
	}

	fmt.Println("Add new changes to the git index")
	if err := Command(Cfg.Git, "-C", Cfg.Dir, "add", "-A").Run(); err != nil {
		return err
	}

	var exitErr *exec.ExitError
	if err := Command(Cfg.Git, "-C", Cfg.Dir, "diff", "HEAD", "--quiet").Run(); err == nil {
		fmt.Println("You don't have any changes since last push")
		return nil
	} else if !errors.As(err, &exitErr) || exitErr.ExitCode() != 1 {
		return err
	}

	fmt.Println("commit new changes")
	err := Command(Cfg.Git, "-C", Cfg.Dir, "commit", "-m", msg).Run()
	if err != nil {
		return err
	}

	fmt.Println("Push changes")
	return Command(Cfg.Git, "-C", Cfg.Dir, "push", "origin").Run()
}

func CmdOpenEditor(_ *cobra.Command, _ []string) error {
	return Command(Cfg.Editor, Cfg.Dir).Run()
}
