package cmd

import (
	"fmt"
	"os"

	"github.com/Yesphet/goit/commit"
	"github.com/spf13/cobra"
)

func genCzCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cz",
		Short: "",
	}

	cmd.Run = func(cmd *cobra.Command, args []string) {
		err := commit.Do()
		checkIfErr(err)

	}
	return cmd
}

func genReleaseCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "release",
		Short: "release <patch|minor|major>",
	}

	cmd.Run = func(cmd *cobra.Command, args []string) {

	}
	return cmd
}

func checkIfErr(err error) {
	if err != nil {
		fatal(err.Error())
	}
}

func fatalf(format string, args ...interface{}) {
	fmt.Printf("Failed: \n\t"+format+"\n", args...)
	os.Exit(1)
}

func fatal(s string) {
	fatalf("%s", s)
}
