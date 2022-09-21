package helpers

import (
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
		token, _, err := GenerateToken(e.id, e.email)
		if err != nil {
			return
		}
		ValidateToken(token)

	}

}
