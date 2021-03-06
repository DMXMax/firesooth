/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// fillCmd represents the fill command
var (
	docPath string

	fillCmd = &cobra.Command{
		Use:   "create",
		Short: "create a firestore document",
		Long: `
Provide a path to a Document, and a json object. 

EXAMPLE: fsquery --docpath users/user1 "{\"id\":\"0001\", \"name\":{\"first\":\"Ada\", \"last\":\"Lovelace\"}}"
`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Trace().Interface("args", args).Send()
			fmt.Println("addDoc called", viper.GetString("docPath"))

			ctx := context.Background()
			client, err := firestore.NewClient(ctx, viper.GetString("project"))

			if err != nil {
				log.Err(err).Send()
			}
			defer client.Close()
			if len(args) < 1 {
				log.Err(errors.New("Not enough Args"))
				return
			}
			data := make(map[string]interface{})
			err = json.Unmarshal([]byte(args[0]), &data)
			if err != nil {
				log.Err(err).Send()
				return
			}
			_, err = client.Doc(viper.GetString("docpath")).Set(ctx, data)
			if err != nil {
				log.Err(err).Send()
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(fillCmd)
	fillCmd.Flags().StringVar(&docPath, "docpath", "", "Document Path")
	fillCmd.MarkFlagRequired("docpath")
	viper.BindPFlag("docpath", fillCmd.Flags().Lookup("docpath"))

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fillCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fillCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
