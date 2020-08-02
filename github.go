package maltmill

import (
	"context"
	"net/http"

	"github.com/Songmu/gitconfig"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func newGithubClient(ctx context.Context, token string) *github.Client {
	if token == "" {
		token, _ = gitconfig.GitHubToken("")
	}
	var oauthCli *http.Client
	if token != "" {
		oauthCli = oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: token,
		}))
	}
	return github.NewClient(oauthCli)
}
