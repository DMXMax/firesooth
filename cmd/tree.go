/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"strings"

	"cloud.google.com/go/firestore"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/api/iterator"
)

var (
	project string
	// treeCmd represents the tree command
	treeCmd = &cobra.Command{
		Use:   "tree",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: tree,
	}
)

func tree(cmd *cobra.Command, args []string) {
	log.Trace().Msg("Tree Called")
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, viper.GetString("project"))

	if err != nil {
		log.Err(err).Send()
	}
	defer client.Close()

	log.Trace().Str("project", viper.GetString("project")).Send()
	colItr := client.Collections(ctx)
	for {
		col, err := colItr.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Err(err).Send()
			return
		}

		walk(ctx, client, col, 0)
	}

}

func walk(ctx context.Context, client *firestore.Client, ifc interface{}, level int) {
	switch ref := ifc.(type){
	case *firestore.CollectionRef:
		fmt.Printf("%s col: %s", strings.Repeat("\t", level), ref.ID)
		fmt.Printf("collection: %s\n", ref.ID)
		Itr := ref.Documents(ctx)
		for {
			ref, err := Itr.Next()

			if err == iterator.Done {
			break
			}
			if err != nil {
				log.Err(err).Send()
				return
			}
			walk(ctx, client, ref, level+1)
		}
	case *firestore.DocumentSnapshot:
		fmt.Printf("%s doc: %s - %v\n", strings.Repeat("\t", level), ref.Ref.ID, ref.Data())
		Itr :=  ref.Ref.Collections(ctx)
		for {
			ref, err := Itr.Next()

			if err == iterator.Done {
			break
			}
			if err != nil {
				log.Err(err).Send()
				return
			}
			walk(ctx, client, ref, level+1)
		}
	}

}
func init() {

	rootCmd.AddCommand(treeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// treeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// treeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
