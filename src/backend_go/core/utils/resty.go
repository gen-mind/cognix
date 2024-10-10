package utils

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"golang.org/x/oauth2"
	"net/http"
)

// WrapRestyError wraps a Resty response error into a standard Go error.
// The function takes a pointer to a Resty response and an error as input.
// If the error is not nil, it returns the error.
// If the response is not an error, it returns nil.
// If the response contains an error message, it creates a new error using fmt.Errorf and returns it.
// Otherwise, it returns nil.
func WrapRestyError(resp *resty.Response, err error) error {
	if err != nil {
		return err
	}
	if !resp.IsError() {
		return nil
	}
	errMsg := string(resp.Body())
	if errMsg != "" {
		return fmt.Errorf(errMsg)
	}
	return nil
}

const AuthorizationHeader = "Authorization"

type Transport struct {
	token *oauth2.Token
	http.Transport
}

func NewTransport(token *oauth2.Token) http.RoundTripper {
	return &Transport{token: token}
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add(AuthorizationHeader, fmt.Sprintf("%s %s", t.token.TokenType, t.token.AccessToken))
	return t.Transport.RoundTrip(req)
}
