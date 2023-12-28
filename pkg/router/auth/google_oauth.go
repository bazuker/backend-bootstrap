package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

const (
	redirectURLEnvVar = "BACKEND_BOOSTRAP_GOOGLE_REDIRECT_URL"
	oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
)

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv(redirectURLEnvVar),
		ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"openid",
		},
		Endpoint: google.Endpoint,
	}
	httpClient = &http.Client{
		Timeout: time.Second * 15,
	}
)

func oauthGoogleLogin(w http.ResponseWriter, r *http.Request) (state string) {
	// Create oauthState cookie.
	oauthState := generateStateOauthCookie(w)
	// AuthCodeURL receive state that is a token to protect the user from CSRF attacks. You must always provide a non-empty string and
	// validate that it matches the state query parameter on your redirect callback.
	u := googleOauthConfig.AuthCodeURL(oauthState)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
	return oauthState
}

func oauthGoogleCallback(_ http.ResponseWriter, r *http.Request) (GoogleUserInfo, string, error) {
	// Read oauthState from Cookie
	oauthState, err := r.Cookie("oauthstate")
	if err != nil {
		return GoogleUserInfo{}, "", fmt.Errorf("user is missing oauthState cookie: %w", err)
	}

	if r.FormValue("state") != oauthState.Value {
		return GoogleUserInfo{}, oauthState.Value, errors.New("invalid oauth google state")
	}

	data, err := getUserDataFromGoogle(r.FormValue("code"))
	if err != nil {
		return GoogleUserInfo{}, oauthState.Value, err
	}

	var userInfo GoogleUserInfo
	err = json.Unmarshal(data, &userInfo)
	return userInfo, oauthState.Value, err
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)
	return state
}

func getUserDataFromGoogle(code string) ([]byte, error) {
	// Use code to get token and get user info from Google.
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := httpClient.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, nil
}
