package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

var signupTest = []struct {
	firstname    string
	lastname     string
	email        string
	password     string
	statuscode   int
	response     string
	httplocation string
	httpredirect string
}{
	{
		"Patrick",
		"Dennis",
		"pat@gmail.com",
		"passlic",
		http.StatusOK,
		"valid Request",
		"/user/signup",
		"/",
	},

	{
		"Elias",
		"Oli",
		"eli.com",
		"13777",
		http.StatusBadRequest,
		"invalid credentials",
		"/user/signup",
		"/user/signup",
	},
	{
		"Dennis",
		"styles",
		"denis@style.com",
		"traeytt",
		http.StatusOK,
		"Valid-Request",
		"user/signup",
		"/",
	},
}

// loginTests is the data for the Login handler tests
var loginTests = []struct {
	name               string
	email              string
	expectedStatusCode int
	expectedHTML       string
	expectedLocation   string
}{
	{
		"valid-credentials",
		"me@here.ca",
		http.StatusSeeOther,
		"",
		"/",
	},
	{
		"invalid-credentials",
		"jack@nimble.com",
		http.StatusSeeOther,
		"",
		"/user/login",
	},
	{
		"invalid-data",
		"j",
		http.StatusOK,
		`action="/user/login"`,
		"",
	},
}

func TestSignup(t *testing.T) {
	for _, i := range signupTest {
		//Values maps a string key to a list of values.
		//It is typically used for query parameters and form values.
		postdata := url.Values{}
		postdata.Add("firstname", i.firstname)
		postdata.Add("lastname", i.lastname)
		postdata.Add("email", i.email)
		postdata.Add("password", i.password)

		//create request
		req, err := http.NewRequest("POST", "/user/login", strings.NewReader(postdata.Encode())) ///Encode encodes the values into “URL encoded” form
		if err != nil {
			t.Errorf("error didnt fill the complete details %s,%s,%s,%s", i.firstname, i.email, i.lastname, i.password)
		}
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		// set the header
		//Header contains the request header fields either received
		//by the server or to be sent by the client.
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		//call the header
		handle := http.HandlerFunc(Repo.SignUp)
		handle.ServeHTTP(rec, req)

	}
}

func TestLogin(t *testing.T) {
	// range through all tests
	for _, e := range loginTests {
		postedData := url.Values{}
		postedData.Add("email", e.email)
		postedData.Add("password", "password")

		// create request
		req, _ := http.NewRequest("POST", "/user/login", strings.NewReader(postedData.Encode()))
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		// set the header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(Repo.PostShowLogin)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

		// if e.expectedLocation != "" {
		// 	// get the URL from test
		// 	actualLoc, _ := rr.Result().Location()
		// 	if actualLoc.String() != e.expectedLocation {
		// 		t.Errorf("failed %s: expected location %s, but got location %s", e.name, e.expectedLocation, actualLoc.String())
		// 	}
		// }

		// // checking for expected values in HTML
		// if e.expectedHTML != "" {
		// 	// read the response body into a string
		// 	html := rr.Body.String()
		// 	if !strings.Contains(html, e.expectedHTML) {
		// 		t.Errorf("failed %s: expected to find %s but did not", e.name, e.expectedHTML)
		// 	}
		// }
	}
}
