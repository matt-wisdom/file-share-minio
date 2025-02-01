/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/matt-wisdom/filesharego/cli/client"
)

// shareCmd represents the share command
var (
	fromUsername  string
	fromUserEmail string
	toUser        string
	shareCmd      = &cobra.Command{
		Use:   "share -s [receiver] -r [sender] [files]",
		Args:  cobra.MinimumNArgs(1),
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fs := client.FileShareServer{ServerAddress: serverAddress}
			err := fs.ShareFiles(args, fromUsername, fromUserEmail, toUser)
			if err != nil {
				panic(err)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(shareCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	shareCmd.PersistentFlags().StringVarP(&fromUsername, "sender", "s", "", "Sender username")
	shareCmd.PersistentFlags().StringVarP(&fromUserEmail, "email", "m", "", "Sender email")
	shareCmd.PersistentFlags().StringVarP(&toUser, "receiver", "r", "", "Receiver username/email")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// shareCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
