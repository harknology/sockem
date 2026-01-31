package config

func ClientAuthAllowed() bool {
	return CLIENT_KEY != "" && len(ALLOWED_HOSTS) > 0
}
