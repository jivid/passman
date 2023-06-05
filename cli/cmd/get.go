package cmd

import (
	"fmt"

	"github.com/jivid/passman/passman/client"
	"github.com/spf13/cobra"
)

type getOptions struct {
	site string
}

type getAllOptions struct{}

var (
	getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get the password for a website",
		Long:  "Get the password for a website",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(fmt.Sprintf("Getting passwords for %s", getOpts.site))
			passman := client.PassmanClient{ServerHost: "localhost", ServerPort: "8080"}
			entries, err := passman.Get(getOpts.site)
			if err != nil {
				fmt.Println(err)
				return nil
			}
			for _, entry := range entries {
				fmt.Println(fmt.Sprintf("Site: %s, Username: %s, Password: %s", entry.Site, entry.Username, entry.Password))
			}
			return nil

		},
	}

	getAllCmd = &cobra.Command{
		Use:   "get-all",
		Short: "Get all passwords",
		Long:  "Get all passwords currently stored in passman",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Getting all passwords")
			passman := client.PassmanClient{ServerHost: "localhost", ServerPort: "8080"}
			entries, err := passman.GetAll()
			if err != nil {
				fmt.Println(err)
				return nil
			}
			for _, entry := range entries {
				fmt.Println(fmt.Sprintf("Site: %s, Username: %s, Password: %s", entry.Site, entry.Username, entry.Password))
			}
			return nil
		},
	}

	getOpts    = getOptions{}
	getAllOpts = getAllOptions{}
)

func init() {
	getCmd.Flags().StringVarP(&getOpts.site, "site", "s", "", "Site to get the password for")
	getCmd.MarkFlagRequired("site")

}
