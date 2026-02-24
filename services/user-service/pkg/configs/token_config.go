package configs

type Token struct {
	AccessToken string `json:"access_token"`
	AccessUuid  string `json:"access_uuid"`
	AtExpires   int64  `json:"at_expires"`

	RefreshToken     string `json:"refresh_token"`
	RefreshUuid      string `json:"refresh_uuid"`
	RefreshAtExpires int64  `json:"rat_expires"`
}

func (t *Token) GetRefreshTokenUUid() string {
	return t.RefreshUuid
}

func (t *Token) SetTokenRefreshAtExpires(exp int64) {
	t.RefreshAtExpires = exp
}

func (t *Token) GetRefreshUUid() string {
	return t.RefreshUuid
}

func (t *Token) GetRefreshAtExpires() int64 {
	return t.RefreshAtExpires
}
