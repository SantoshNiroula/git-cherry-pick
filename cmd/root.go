/*
Copyright © 2024 Santosh Niroula <niroulasantosh624@gmail.com>
*/
package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/SantoshNiroula/git-cherry-pick/helper"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "git-cherry-pick",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: startCherryPick,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringP("file-name", "f", "pr.txt", "Read PR number from given file with new line")
}

func startCherryPick(cmd *cobra.Command, args []string) {
	var uniquePr []string

	if len(args) == 0 {
		fileName, err := cmd.Flags().GetString("file-name")
		if err != nil {
			log.Fatal(err)
		}
		prList := helper.ReadPRNumberFromFile(fileName)
		uniquePr = removeDuplicate(prList)
	} else {
		uniquePr = removeDuplicate(args)
	}

	formattedPr := formatPRNumberForGrep(uniquePr)
	runCommand(formattedPr)
}

func formatPRNumberForGrep(args []string) string {
	formattedGrepCmd := ""
	for _, arg := range args {
		formattedGrepCmd += fmt.Sprintf("-e %s ", arg)
	}

	return formattedGrepCmd
}

func removeDuplicate(args []string) []string {
	m := make(map[string]bool)
	var result []string
	for _, str := range args {
		if _, notUnique := m[str]; !notUnique {
			m[str] = true
			result = append(result, str)
		}
	}
	return result
}

func runCommand(formattedPr string) {
	gitCmd := exec.Command("git", "log", `--format=pretty:%ah %h %S`)
	gitOut, err := gitCmd.Output()
	if err != nil {
		fmt.Println("Git cmd", err)
		return
	}

	grepCmd := exec.Command("grep", formattedPr)
	grepCmd.Stdin = bytes.NewReader(gitOut)
	grepOut, err := grepCmd.Output()
	if err != nil {
		fmt.Println("grep cmd", err)
		return
	}

	sortCmd := exec.Command("sort", "-n")
	sortCmd.Stdin = bytes.NewReader(grepOut)
	sortOut, err := sortCmd.Output()
	if err != nil {
		fmt.Println("Sort cmd", err)
		return
	}

	if err := os.WriteFile("gi-cp.txt", sortOut, 0644); err != nil {
		fmt.Println("Unable to write", err)
		fmt.Println("\n", string(sortOut))
		return
	}
}
