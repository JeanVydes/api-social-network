package main

var (
	SessionTokens = map[string]Session{}
)

func AssignToken(accountID string) (string, error) {
	token, err := GenerateToken(accountID)
	if err != nil {
		return "", err
	}

	SessionTokens[token] = Session{
		Token:     token,
		AccountID: accountID,
	}

	return token, nil
}

func RemoveToken(token string) {
	delete(SessionTokens, token)
}
