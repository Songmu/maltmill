package maltmill

import (
	"net/http"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func newGithubClient(token string) *github.Client {
	var oauthCli *http.Client
	if token != "" {
		oauthCli = oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: token,
		}))
	}
	return github.NewClient(oauthCli)
}
