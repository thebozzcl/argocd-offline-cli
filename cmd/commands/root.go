package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/touchardv/argocd-offline-cli/preview"
)

func NewCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "argocd-offline-cli",
		Short: "An Argo CD CLI offline utility",
		Long: `A utility, based on Argo CD, that can be used "offline" (without requiring a running Argo CD server),
to preview the Kubernetes resource manifests being created and managed by Argo CD.`,
	}

	rootCmd.AddCommand(AppSetCommand())

	return rootCmd
}

func AppSetCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "appset",
		Short: "Preview ApplicationSets",
	}
	command.AddCommand(PreviewApplicationsCommand())
	command.AddCommand(PreviewApplicationResourcesCommand())
	return command
}

func PreviewApplicationsCommand() *cobra.Command {
	var name string
	var output string
	var localRepos []string
	command := &cobra.Command{
		Use:   "preview-apps APPSETMANIFEST",
		Short: "Preview Application(s) generated from an ApplicationSet",
		Run: func(c *cobra.Command, args []string) {
			if len(args) == 0 {
				c.HelpFunc()
				os.Exit(1)
			}
			filename := args[0]
			preview.SetLocalRepoMappings(parseLocalRepos(localRepos))
			preview.PreviewApplications(filename, name, output)
		},
	}
	command.Flags().StringVarP(&name, "name", "n", "", "Name of the Application to preview")
	command.Flags().StringVarP(&output, "output", "o", "name", "Output format. One of: name|json|yaml")
	command.Flags().StringArrayVarP(&localRepos, "local-repo", "l", nil, "Local repository mapping in the format 'repoURL=localPath' (can be specified multiple times)")
	return command
}

func PreviewApplicationResourcesCommand() *cobra.Command {
	var kind string
	var name string
	var output string
	var localRepos []string
	command := &cobra.Command{
		Use:   "preview-resources APPSETMANIFEST",
		Short: "Preview Kubernetes resource(s) generated from an ApplicationSet/Application",
		Run: func(c *cobra.Command, args []string) {
			if len(args) == 0 {
				c.HelpFunc()
				os.Exit(1)
			}
			filename := args[0]
			preview.SetLocalRepoMappings(parseLocalRepos(localRepos))
			preview.PreviewResources(filename, name, kind, output)
		},
	}
	command.Flags().StringVarP(&kind, "kind", "k", "", "Kind of resources to preview")
	command.Flags().StringVarP(&name, "name", "n", "", "Name of the Application to preview")
	command.Flags().StringVarP(&output, "output", "o", "name", "Output format. One of: name|json|yaml")
	command.Flags().StringArrayVarP(&localRepos, "local-repo", "l", nil, "Local repository mapping in the format 'repoURL=localPath' (can be specified multiple times)")
	return command
}

// parseLocalRepos parses the --local-repo flag values into a map.
// Each value should be in the format "repoURL=localPath".
func parseLocalRepos(localRepos []string) map[string]string {
	result := make(map[string]string)
	for _, mapping := range localRepos {
		parts := strings.SplitN(mapping, "=", 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result
}
