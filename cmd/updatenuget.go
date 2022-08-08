/*
Copyright © 2022 Raa Labs <post@raalabs.com>

*/
package cmd

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
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
		fmt.Println("updatenuget called")
		updateNugetPackage()
	},
}

func init() {
	rootCmd.AddCommand(updatenugetCmd)
	updatenugetCmd.PersistentFlags().String("package", "", "The package to update")
}

func updateNugetPackage() {
	xmlFile, err := os.Open(".csproj")
	if err != nil {
	    fmt.Println(err)
    }

	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)
	var propertyGroup PropertyGroup
	xml.Unmarshal(byteValue, &propertyGroup)

	var itemGroup ItemGroup
	xml.unmarshal(byteValue, &itemGroup)

	for i := 0; i < len(itemGroup.Items); i++ {
		fmt.Println(itemGroup.Items[i].InnerXML)
	}
}

type CsProjFile struct {
	XMLName xml.Name `xml:"Project"`
	PropertyGroup []PropertyGroup `xml:"PropertyGroup"`
	ItemGroup []ItemGroup `xml:"ItemGroup"`
}

type PropertyGroup struct {
	OutputType string `xml:"OutputType"`
	TargetFramework string `xml:"TargetFramework"`
	AssemblyName string `xml:"AssemblyName"`
}

type ItemGroup struct {
	PackageReferences []PackageReference `xml:"PackageReference"` 
}
