package cmd

import (
	"os"
	"regexp"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	"k8s.io/kubectl/pkg/util/term"

	"github.com/timebertt/kubectl-history/pkg/cmd/diff"
	"github.com/timebertt/kubectl-history/pkg/cmd/get"
	"github.com/timebertt/kubectl-history/pkg/cmd/util"
	"github.com/timebertt/kubectl-history/pkg/cmd/version"
)

type Options struct {
	genericclioptions.IOStreams

	ConfigFlags *genericclioptions.ConfigFlags
}

func NewOptions() *Options {
	return &Options{
		IOStreams:   genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr},
		ConfigFlags: genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag().WithDiscoveryBurst(300).WithDiscoveryQPS(50.0),
	}
}

func NewCommand() *cobra.Command {
	o := NewOptions()

	cmd := &cobra.Command{
		Use:   "history",
		Short: "Time-travel through your cluster",

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			warningHandler := rest.NewWarningWriter(o.IOStreams.ErrOut, rest.WarningWriterOptions{Deduplicate: true, Color: term.AllowsColorOutput(o.IOStreams.ErrOut)})
			rest.SetDefaultWarningHandler(warningHandler)
			return nil
		},
	}

	flags := cmd.PersistentFlags()

	o.ConfigFlags.AddFlags(flags)
	f := util.NewFactory(o.ConfigFlags)

	defaultGroup := &cobra.Group{
		ID:    "default",
		Title: "Available Commands:",
	}
	cmd.AddGroup(defaultGroup)

	for _, subcommand := range []*cobra.Command{
		diff.NewCommand(f, o.IOStreams),
		get.NewCommand(f, o.IOStreams),
	} {
		subcommand.GroupID = defaultGroup.ID
		cmd.AddCommand(subcommand)
	}

	otherGroup := &cobra.Group{
		ID:    "other",
		Title: "Other Commands:",
	}
	cmd.AddGroup(otherGroup)

	cmd.SetCompletionCommandGroupID(otherGroup.ID)
	cmd.SetHelpCommandGroupID(otherGroup.ID)

	versionCmd := version.NewCommand(o.IOStreams)
	versionCmd.GroupID = otherGroup.ID
	cmd.AddCommand(versionCmd)

	customizeUsageTemplate(cmd)

	return cmd
}

// customizeUsageTemplate makes sure the plugin's help output always has the `kubectl ` command prefix before the
// plugin's command path to match the expected command path when it is executed via kubectl.
// This implements https://krew.sigs.k8s.io/docs/developer-guide/develop/best-practices/#help-messages.
// I.e., the default template would output:
//
//	Usage:
//	  history [command]
//
// The modified template outputs:
//
//	Usage:
//	  kubectl history [command]
//
// Changing cmd.Use to `kubectl history` makes cobra remove `history` from all command paths and use lines.
func customizeUsageTemplate(cmd *cobra.Command) {
	defaultTmpl := cmd.UsageTemplate()

	r := regexp.MustCompile(`([{ ])(.CommandPath|.UseLine)([} ])`)
	tmpl := r.ReplaceAllString(defaultTmpl, `$1(printf "kubectl %s" $2)$3`)

	cmd.SetUsageTemplate(tmpl)
}
