package oauth2

import "context"

type UserInfo struct {
	Provider      string
	ProviderID    string
	Email         string
	EmailVerified bool
	Name          string
	AvatarURL     string
}

type Provider interface {
	Key() string
	AuthCodeURL(state string, extra map[string]string) string
	Exchange(ctx context.Context, code string) (*UserInfo, error)
}
