/*
Copyright © 2022 Raa Labs <post@raalabs.com>

*/
package cmd

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var updatenugetCmd = &cobra.Command{
	Use:   "updatenuget",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		nugetPackage, _ := cmd.Flags().GetString("nugetPackage")
		repos, _ := cmd.Flags().GetString("repos")
		fmt.Println("updatenuget called")
		updateNugetPackage(nugetPackage, repos)
	},
}

func init() {
	rootCmd.AddCommand(updatenugetCmd)
	updatenugetCmd.PersistentFlags().String("nugetPackage", "", "The nuget package to update")
	updatenugetCmd.PersistentFlags().String("repos", "", "Comma separated list of git repos to update")
}

func updateNugetPackage(nugetPackage string, repos string) {
	tmpReposDir := createTmpFolder()
	fmt.Println(tmpReposDir)

	for _, repo := range strings.Split(strings.ReplaceAll(repos, " ", ""), ",") {
		repoDir := cloneGitRepo(tmpReposDir, repo)

		fmt.Println(repoDir)

		projectFiles := getCsprojFiles(repoDir)
		for _, projectFile := range projectFiles {
			fmt.Println(projectFile)
			updateCsProjFile(projectFile)
		}
	}
}

func createTmpFolder() string {
	err := os.Mkdir("tmp", 0755)
	if err != nil {
		fmt.Println(err)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	return filepath.Join(currentDir, "tmp")
}

func cloneGitRepo(targetDir string, repoName string) string {
	repoDir := filepath.Join(targetDir, repoName)
	_ = os.Mkdir(repoDir, 0755)

	return repoDir
}

func getCsprojFiles(dir string) []string {
	files := []string{}
	err := filepath.Walk(
		dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(path, ".testcsproj") {
				files = append(files, path)
			}

			return nil
		})
	if err != nil {
		log.Println(err)
	}

	return files
}

func updateCsProjFile(file string) {
	xmlFile, err := os.Open(file)
	if err != nil {
	    fmt.Println(err)
    }
	defer xmlFile.Close()
	byteValue, _ := ioutil.ReadAll(xmlFile)
	var project Project
	xml.Unmarshal(byteValue, &project)

	for i := 0; i < len(project.ItemGroup.PackageReferences); i++ {
		fmt.Println(project.ItemGroup.PackageReferences[i].Include)
		fmt.Println(project.ItemGroup.PackageReferences[i].Version)
	}
}

type Project struct {
	XMLName xml.Name `xml:"Project"`
	Sdk string `xml:"Sdk,attr"`
	PropertyGroup PropertyGroup `xml:"PropertyGroup"`
	ItemGroup ItemGroup `xml:"ItemGroup"`
}

type PropertyGroup struct {
	XMLName xml.Name `xml:"PropertyGroup"`
	OutputType string `xml:"OutputType"`
	TargetFramework string `xml:"TargetFramework"`
	AssemblyName string `xml:"AssemblyName"`
}

type ItemGroup struct {
	XMLName xml.Name `xml:"ItemGroup"`
	PackageReferences []PackageReference `xml:"PackageReference"`
}

type PackageReference struct {
	XMLName xml.Name `xml:"PackageReference"`
	Include string `xml:"Include,attr"`
	Version string `xml:"Version,attr"`
}
