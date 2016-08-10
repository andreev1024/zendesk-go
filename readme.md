This package implements a minimal wrapper for Zendesk API for Go (Golang).

##Authentication
In current version implemented API token auth.

##Errors
Zendesk returns many useful information about error. Hence all package endpoints
 return `body`, `response` and `err`. 
This is because:
*   Error format not consistent. We can get errors in `json`, `plain/text` formats.
*   Sometimes we should check headers (rate limit error).

Handle error example:
```
	user, body, response, err := api.CreateUser(req)
	if err != nil {
		if err.Error() == zd.ERR { //Zendesk API error (for e.g validation)
			if response != nil {
				log.Printf("STATUS %v\n", response.Status) //check status
				log.Printf("BODY %v\n", string(body)) //check body
				//  do something
			}
			return
		}
		//  Not API Zendesk error. Handle error as usual (don't check body/status)
		return
	}
```

##Example
```
package main

import (
	"fmt"
	zd "github.com/andreev1024/zendesk-go"
	"log"
)

func main() {
	email := "email@gmail.com"
	token := "myToken1234567890"
	host := "https://testcompany.zendesk.com"
	errorHandler := func(e error) { //not required
		log.Println(e)
	}
	api := zd.NewAPI(email, token, host, errorHandler)

	//Create user
	req := map[string]interface{}{
		"name":  "Vasya Pupkin",
		"email": "iam@example.com",
	}
	user, body, res, err := api.CreateUser(req)

	//Show user
	userId := 3643536945
	user, body, res, err := api.ShowUser(userId)

	//Update user
	req := map[string]interface{}{
		"name":      "Petya Vasilkow",
		"locale_id": 1,
	}
	u, body, res, e := api.UpdateUser(userId, req)

	//Update user profile image
	userId := 3643536945
	body, res, err := api.UpdateUserProfileImage(userId, "./avatar.jpg", "")
	//or
	body, res, err := api.UpdateUserProfileImage(userId, "", "https://site.com/picture.jpg")

	//Create ticket
	req := map[string]interface{}{
		"subject": "My printer is on fire!",
		"comment": map[string]interface{}{
			"body": "The smoke is very colorful.",
		},
	}
	t, body, res, err := api.CreateTicket(req)
	//or
	t, body, res, err := api.CreateTicketAsync(req)
	
	//Show ticket
	ticketId := 1
	ticket, body, res, e := api.ShowTicket(ticketId)

	//ListTickets
	url := "users/3618912245/tickets/assigned.json"
	sortBy := "status"
	sortOrder := "desc"
	//with sort
	tickets, body, res, e := api.ListTickets(url, sortBy, sortOrder)
	//or without (default sort)
	tickets, body, res, e := api.ListTickets(url, "", "")
}

```