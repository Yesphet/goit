package cmd

import (
	"github.com/spf13/cobra"
	"github.com/Yesphet/goit/commit"
	"fmt"
	"os"
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
		Use: "release",
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
