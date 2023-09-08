package version

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// version can be set via:
// -ldflags="-X 'github.com/timebertt/kubectl-history/pkg/cmd/version.version=$TAG'"
// If not overwritten via ldflags, it defaults to the go module's version if installed via `go install`.
var version string

func init() {
	if version == "" {
		i, ok := debug.ReadBuildInfo()
		if !ok {
			return
		}
		version = i.Main.Version
	}
}

type Options struct {
	genericclioptions.IOStreams
}

func NewOptions(streams genericclioptions.IOStreams) *Options {
	return &Options{
		IOStreams: streams,
	}
}

func NewCommand(streams genericclioptions.IOStreams) *cobra.Command {
	o := NewOptions(streams)

	return &cobra.Command{
		Use:                   "version",
		DisableFlagsInUseLine: true,
		Short:                 "Print the version of kubectl-history",
		Long: `The version command prints the source version that was used to build the binary.
Note that the version string's format can be different depending on how the binary was built.
E.g, release builds inject the version via -ldflags, while installing with 'go install' injects
the go module's version (which can also be "(devel)").`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if version == "" {
				fmt.Fprintln(o.Out, "could not determine build information")
			} else {
				fmt.Fprintf(o.Out, "kubectl-history %s\n", version)
			}

			return nil
		},
	}
}
