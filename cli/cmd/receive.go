/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/matt-wisdom/filesharego/cli/client"
	"github.com/spf13/cobra"
)

var fromUser, downloadFolder string

// receiveCmd represents the receive command
var receiveCmd = &cobra.Command{
	Use:   "receive",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fs := client.FileShareServer{ServerAddress: serverAddress, DownloadFolder: "downloads"}
		err := fs.ReceiveFiles(fromUser, toUser)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(receiveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// receiveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// receiveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	receiveCmd.PersistentFlags().StringVarP(&fromUser, "sender", "s", "", "Sender email/username")
	receiveCmd.PersistentFlags().StringVarP(&toUser, "receiver", "r", "", "Receiver username/email")
	receiveCmd.PersistentFlags().StringVarP(&toUser, "downloads", "d", "downloads", "Download folder")
}
