package cmd

import "github.com/spf13/cobra"

func genCzCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cz",
		Short: "",
	}
	return cmd
}

func genReleaseCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "release",
	}
	return cmd
}
