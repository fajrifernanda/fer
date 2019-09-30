package cmd

import (
	"fmt"
	"github.com/kumparan/fer/generator"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init a microservice",
	Long:  `example 'trafo init content-service',and then input the proto path`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		if name != "" {
			generator.Generate(name)
		} else {
			fmt.Println("dont forget 'name' for microservice")
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().String("name", "", "name for new microservice")
}
