package dto

// OAuthProvider represents an external authentication provider.
type OAuthProvider string

const (
	// OAuthProviderGitHub represents authentication via GitHub OAuth.
	OAuthProviderGitHub OAuthProvider = "github"
	// OAuthProviderGoogle represents authentication via Google OAuth.
	OAuthProviderGoogle OAuthProvider = "google"
)

// OAuthUser represents oauth user object.
type OAuthUser struct {
	ID             string        `db:"id"`
	UserID         string        `db:"user_id"`
	ProviderUserID string        `db:"provider_user_id"`
	Provider       OAuthProvider `db:"provider"`
}

// GithubUser represents payload received from github when fetching its user.
type GithubUser struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// GithubTokenResponse represents payload received from github when fetching access_token.
type GithubTokenResponse struct {
	AccessToken string `json:"access_token"`
}

// GoogleUser represents payload received from google when fetching its user.
type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}
