package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
)

type createOptions struct {
	Site     string
	Username string
	Password string
}

var (
	createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new password",
		Long:  "Create a new password",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(fmt.Sprintf(
				"Creating password for %s. "+
					"Username %s, Password %s",
				createOpts.Site, createOpts.Username, createOpts.Password,
			))

			body, err := json.Marshal(&createOpts)
			if err != nil {
				return err
			}
			buf := bytes.NewBuffer(body)
			resp, err := http.Post("http://127.0.0.1:8080/passwords", "application/json", buf)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusCreated {
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return err
				}
				var r map[string]interface{}
				json.Unmarshal(body, &r)
				fmt.Println(r["message"])
				return errors.New("Couldn't create password")
			}
			return nil
		},
	}

	createOpts = createOptions{}
)

func init() {
	createCmd.Flags().StringVarP(&createOpts.Site, "site", "s", "", "Site to create password for")
	createCmd.Flags().StringVarP(&createOpts.Username, "username", "u", "", "Username for the password")
	createCmd.Flags().StringVarP(&createOpts.Password, "password", "p", "", "Password")
	createCmd.MarkFlagsRequiredTogether("site", "username", "password")
}
