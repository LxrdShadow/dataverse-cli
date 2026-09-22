package commands

import (
	"dvc/client"
	"fmt"
)

var WhoAmICommand = &Command{
	Name:        "whoami",
	Description: "Retrieves the current user's informations. (Business unit ID, User ID, Organization ID)",
	Run:         whoAmI,
}

func whoAmI(client *client.DataverseClient, args []string) error {
	whoamiResponse, err := client.WhoAmI()
	if err != nil {
		return err
	}

	fmt.Println("BusinessUnitID:", whoamiResponse.BusinessUnitID)
	fmt.Println("UserID:", whoamiResponse.UserID)
	fmt.Println("OrganizationID:", whoamiResponse.OrganizationID)
	return nil
}
