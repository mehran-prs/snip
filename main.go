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

var appName string

var (
	Version = "" // fill-in at compile time.
	Commit  = "" // fill-in at compile time.
	Date    = "" // fill-in at compile time.
)

var FlagRecursiveRemove = false

func init() {
	appName = DefaultStr(baseName(os.Args[0]), "snip")
	appName = strings.TrimSuffix(appName, filepath.Ext(appName))
}

func main() {
	if err := run(); err != nil {
		Error("app error: ", err)
		os.Exit(1)
	}
}

func run() error {

	var completionCmd = &cobra.Command{
		Use:                   "completion [bash|zsh|fish|powershell]",
		Short:                 "Generate completion script",
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE:                  CmdCompletionGenerator,
	}

	var rootCmd = &cobra.Command{
		Use:                fmt.Sprintf("%s [command]", appName),
		Short:              fmt.Sprintf("%s is a snippet manager on the command line.", appName),
		Long:               "The snip tool is a snippet manager on the command line.",
		Args:               cobra.ExactArgs(1),
		RunE:               CmdViewSnippet,
		PersistentPreRunE:  boot,
		PersistentPostRunE: shutdown,
		ValidArgsFunction:  cobraAutoCompleteFileName,
	}

	var dirCmd = &cobra.Command{
		Use:   "dir [subPath]",
		Short: "prints the snippets directory",
		Long:  "you can run 'cd $(snip dir)' to cd into your snippets directory.",
		Args:  cobra.MaximumNArgs(1),
		RunE:  CmdSnippetsDir,
	}

	var editCmd = &cobra.Command{
		Use:               "edit",
		Short:             "Create|Edit the snippet in the editor",
		Long:              "If do not provide any snippet name, it'll open the snippets directory in the editor",
		Args:              cobra.MaximumNArgs(1),
		RunE:              CmdEditSnippet,
		ValidArgsFunction: cobraAutoCompleteFileName,
	}

	var RemoveCmd = &cobra.Command{
		Use:               "rm [-r] [file|dir(append a slash to it)]",
		Short:             "Remove a snippet or directory",
		Long:              `Removes a snippet file or a snippet directory. To specify a directory, append a '/' to it.`,
		Args:              cobra.MaximumNArgs(1),
		RunE:              CmdRemoveSnippet,
		ValidArgsFunction: cobraAutoCompleteFileName,
	}

	var syncCmd = &cobra.Command{
		Use:   "sync",
		Short: "sync the snippets changes with your remote git repository",
		Long:  "Sync command first pull changes from yourb git repository and then commit and pushes your changes.",
		Args:  cobra.MaximumNArgs(1),
		RunE:  CmdSync,
	}

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version and build information",
		Run:   CmdPrintVersion,
	}

	RemoveCmd.Flags().BoolVarP(&FlagRecursiveRemove, "recursive", "r", false, "Remove recursively")
	rootCmd.AddCommand(completionCmd, dirCmd, editCmd, RemoveCmd, syncCmd, versionCmd)

	return rootCmd.Execute()
}

// boot boots the app. it loads config for us.
func boot(cmd *cobra.Command, _ []string) error {
	if err := loadConfig(prefix, strings.ToUpper(cmd.Root().Name())); err != nil {
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
		// Currently in bash we do not support fzf, because when we define the fzf function and enable completion
		// for it, it'll override the original completion function, this is while we need to just call it when
		// user's input ends with "**" to the fzf trigger value.
		// TODO: enable completion for bash, but keep the original completion too.
		//fmt.Println(genFzfBashCompletion(cmd.Root().Name()))
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

func CmdViewSnippet(c *cobra.Command, args []string) error {
	fpath := Cfg.SnippetPath(args[0])

	// If the file doesn't exist and is not a directory path, ask for creating it.
	if _, err := os.Stat(fpath); errors.Is(err, os.ErrNotExist) && !EndsWithDirectoryPath(fpath) {
		msg := fmt.Sprintf("File %s doesn't exist, create it? (y/n) [y] ", strings.TrimPrefix(fpath, Cfg.Dir))
		edit, err := BoolPrompt(c.InOrStdin(), c.OutOrStdout(), msg)
		if err != nil || !edit {
			return err
		}

		return CmdEditSnippet(c, args)
	}
	return Cfg.ViewerCmd(fpath).Run()
}

func CmdSnippetsDir(_ *cobra.Command, args []string) error {
	subPath := ""
	if len(args) != 0 {
		subPath = filepath.Dir(args[0])
	}

	fmt.Println(filepath.Join(Cfg.Dir, subPath))
	return nil
}

func CmdEditSnippet(_ *cobra.Command, args []string) error {
	fpath := Cfg.Dir
	if len(args) != 0 {
		fpath = Cfg.SnippetPath(args[0])

		// Make parent directories
		if err := os.MkdirAll(filepath.Dir(fpath), 0777); err != nil {
			return fmt.Errorf("can not create snippet directory: %w", err)
		}
	}

	return Command(Cfg.Editor, fpath).Run()
}

func CmdRemoveSnippet(_ *cobra.Command, args []string) error {
	fpath := Cfg.SnippetPath(args[0])

	f := os.Remove
	if FlagRecursiveRemove {
		f = os.RemoveAll
	}
	if err := f(fpath); err != nil {
		return err
	}

	fmt.Printf("Removed: %s\n", fpath)
	return nil
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

func CmdPrintVersion(*cobra.Command, []string) {
	fmt.Println("Version: ", DefaultStr(Version, "unknown"))
	fmt.Println("Commit: ", DefaultStr(Commit, "unknown"))
	fmt.Println("Build Date: ", DefaultStr(Date, "unknown"))
}
