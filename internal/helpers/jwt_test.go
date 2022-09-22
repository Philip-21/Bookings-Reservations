package helpers

import (
	"net/http"
	"testing"
)

var LoginTest = []struct {
	id    int
	email string
}{
	{
		id:    1,
		email: "harry@gmail.com",
	},
	{
		id:    2,
		email: "lexis@gmail.com",
	},
	{
		id:    3,
		email: "his@gmail.com",
	},
}

func TestLogin(t *testing.T) {
	for _, e := range LoginTest {
		_, err := GenerateJWT(e.email)
		if err != nil {
			return
		}
		ValidateToken(&http.Request{})

	}

}
