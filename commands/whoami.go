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

	fmt.Println("BusinessUnitId:", whoamiResponse.BusinessUnitId)
	fmt.Println("UserId:", whoamiResponse.UserId)
	fmt.Println("OrganizationId:", whoamiResponse.OrganizationId)
	return nil
}
