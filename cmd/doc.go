package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

const gendocFrontmatterTemplate = `---
date: %s
title: "%s"
slug: %s
url: %s
---
`

var genDocDir string
var docCmd = &cobra.Command{
	Use:   "doc",
	Short: "Generate Markdown documentation for the Dkron CLI.",
	Long: `Generate Markdown documentation for the Dkron CLI.
This command is, mostly, used to create up-to-date documentation
of Dkron's command-line interface for http://dkron.io/.
It creates one Markdown file per command with front matter suitable
for rendering in Hugo.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat(genDocDir); os.IsNotExist(err) {
			fmt.Println("Directory", genDocDir, "does not exist, creating...")
			if err := os.MkdirAll(genDocDir, 0777); err != nil {
				return err
			}
		}
		now := time.Now().Format("2006-01-02")
		prepender := func(filename string) string {
			name := filepath.Base(filename)
			base := strings.TrimSuffix(name, path.Ext(name))
			url := "/cli/" + strings.ToLower(base) + "/"
			return fmt.Sprintf(gendocFrontmatterTemplate, now, strings.Replace(base, "_", " ", -1), base, url)
		}

		linkHandler := func(name string) string {
			base := strings.TrimSuffix(name, path.Ext(name))
			return "/cli/" + strings.ToLower(base) + "/"
		}

		fmt.Println("Generating Dkron command-line documentation in", genDocDir, "...")
		doc.GenMarkdownTreeCustom(cmd.Root(), genDocDir, prepender, linkHandler)
		fmt.Println("Done.")

		return nil
	},
}

func init() {
	dkronCmd.AddCommand(docCmd)
	docCmd.PersistentFlags().StringVar(&genDocDir, "dir", "/tmp/dkrondoc/", "the directory to write the doc.")

	// For bash-completion
	docCmd.PersistentFlags().SetAnnotation("dir", cobra.BashCompSubdirsInDir, []string{})
}
