package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jivid/passman/passman"
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
			resp, err := http.Get(fmt.Sprintf("http://localhost:8080/passwords/%s", getOpts.site))
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return err
				}
				var r map[string]interface{}
				json.Unmarshal(body, &r)
				fmt.Println(r["message"])
			} else {
				body, err := ioutil.ReadAll(resp.Body)
				var entries []passman.PassmanEntry
				err = json.Unmarshal(body, &entries)
				if err != nil {
					return err
				}
				for _, entry := range entries {
					fmt.Println(fmt.Sprintf("Site: %s, Username: %s, Password: %s", entry.Site, entry.Username, entry.Password))
				}
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
			resp, err := http.Get("http://localhost:8080/passwords")
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			var entries []passman.PassmanEntry
			err = json.Unmarshal(body, &entries)
			if err != nil {
				return err
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
