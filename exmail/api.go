package exmail

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/wealthworks/go-tencent-api/client"
)

var (
	corpId string

	apiContact *API
	apiLogin   *API
	apiCheck   *API
)

func init() {
	corpId = os.Getenv("EXMAIL_CORP_ID")

}

type API struct {
	c *client.Client
}

// apps: Contact, Login, Check
func New(apiCat string) *API {
	if corpId == "" {
		log.Fatal("EXMAIL_CORP_ID is empty or not found")
	}
	if apiCat == "" {
		log.Print("empty apiCat")
	}
	k := fmt.Sprintf("EXMAIL_API_%s_SECRET", strings.ToUpper(apiCat))
	corpSecret := os.Getenv(k)
	if corpSecret == "" {
		log.Fatal(fmt.Sprintf("%s are empty or not found", k))
	}
	c := client.NewClient(urlToken)
	c.SetContentType("application/json")
	c.SetCorp(corpId, corpSecret)
	api := &API{
		c: c,
	}
	// log.Printf("api client: %s", api.c)
	return api
}

func ApiContact() *API {
	if apiContact == nil {
		apiContact = New("Contact")
	}
	return apiContact
}

func ApiLogin() *API {
	if apiLogin == nil {
		apiLogin = New("Login")
	}
	return apiLogin
}

func ApiCheck() *API {
	if apiCheck == nil {
		apiCheck = New("Check")
	}
	return apiCheck
}
