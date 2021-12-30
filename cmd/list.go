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
	project     string
	entityLimit int
	// listCmd represents the tree command
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "list everything from the root",
		Long: `
		Provides the entire collection and documents. Use with caution on real firestore
`,
		Run: list,
	}
)

func list(cmd *cobra.Command, args []string) {
	log.Info().Int("limit", viper.GetInt("limit")).Msg("list")
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, viper.GetString("project"))

	if err != nil {
		log.Err(err).Send()
	}
	defer client.Close()

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

		walk(ctx, client, col, 0, viper.GetUint("limit"))
	}

}

// The walk fu
func walk(ctx context.Context, client *firestore.Client, ifc interface{}, level int, entityLimit uint) bool {
	if entityLimit == 0 {
		fmt.Println("entity limit reached")
		log.Trace().Msg("Entity Limit")
		return false
	}

	switch ref := ifc.(type) {
	case *firestore.CollectionRef:
		fmt.Printf("%s col: %s", strings.Repeat("\t", level), ref.ID)
		fmt.Printf("collection: %s\n", ref.ID)
		Itr := ref.Documents(ctx)
		for {
			ref, err := Itr.Next()

			if err == iterator.Done {
				return true
			}
			if err != nil {
				log.Err(err).Send()
				return false
			}
			// Returning false breaks the for loop for this loop every parent loop.
			if walk(ctx, client, ref, level+1, entityLimit-1) == false {
				return false
			}
			entityLimit = entityLimit - 1
		}
	case *firestore.DocumentSnapshot:
		fmt.Printf("%s doc: %s - %v\n", strings.Repeat("\t", level), ref.Ref.ID, ref.Data())
		Itr := ref.Ref.Collections(ctx)
		for {
			ref, err := Itr.Next()

			if err == iterator.Done {
				return true
			}
			if err != nil {
				log.Err(err).Send()
				return false
			}
			// Returning false breaks the for loop for this loop every parent loop.
			if walk(ctx, client, ref, level+1, entityLimit-1) == false {
				return false
			}
			entityLimit = entityLimit - 1
		}
	}
	return true
}
func init() {

	rootCmd.AddCommand(listCmd)
	listCmd.Flags().Uint("limit", 50, "maximum number of entities returned")
	viper.BindPFlag("limit", listCmd.Flags().Lookup("limit"))

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
