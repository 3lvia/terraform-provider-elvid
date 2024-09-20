package elvidapiclient

type UserClientDto struct {
	Id                               int                 `json:"Id"`
	ClientId                         string              `json:"ClientId"`
	ClientName                       string              `json:"ClientName"`
	Scopes                           []string            `json:"Scopes"`
	Domains                          []string            `json:"Domains"`
	RedirectUriPaths                 []string            `json:"RedirectUriPaths"`
	PostLogoutRedirectUriPaths       []string            `json:"PostLogoutRedirectUriPaths"`
	IdPortenLoginEnabled             bool                `json:"IdPortenLoginEnabled"`
	LocalLoginEnabled                bool                `json:"LocalLoginEnabled"`
	ElviaADLoginEnabled              bool                `json:"ElviaADLoginEnabled"`
	TestUserLoginEnabled             bool                `json:"TestUserLoginEnabled"`
	RequireClientSecret              bool                `json:"RequireClientSecret"`
	AccessTokenLifetime              int                 `json:"AccessTokenLifetime"`
	AlwaysIncludeUserClaimsInIdToken bool                `json:"AlwaysIncludeUserClaimsInIdToken"`
	ClientNameLanguageKey            string              `json:"ClientNameLanguageKey"`
	AllowUseOfRefreshTokens          bool                `json:"AllowUseOfRefreshTokens"`
	OneTimeUsageForRefreshTokens     bool                `json:"OneTimeUsageForRefreshTokens"`
	RefreshTokensLifeTime            int                 `json:"RefreshTokensLifeTime"`
	ClientProperties                 []ClientPropertyDto `json:"ClientProperties"`
}

type ClientPropertyDto struct {
	Key    string   `json:"Key"`
	Values []string `json:"Values"`
}

type ClientSecretDto struct {
	Id                    int    `json:"Id"`
	Value                 string `json:"SecretValue"`
	HashedValueStartsWith string `json:"HashedValueStartsWith"`
}

type MachineClientDto struct {
	Id                   int              `json:"Id"`
	ClientId             string           `json:"ClientId"`
	ClientName           string           `json:"ClientName"`
	TestUserLoginEnabled bool             `json:"TestUserLoginEnabled"`
	IsDelegationClient   bool             `json:"IsDelegationClient"`
	AccessTokenLifeTime  int              `json:"AccessTokenLifeTime"`
	Scopes               []string         `json:"Scopes"`
	ClientClaims         []ClientClaimDto `json:"ClientClaims"`
}

type ClientClaimDto struct {
	Type   string   `json:"Type"`
	Values []string `json:"Values"`
}

type ApiScopeDto struct {
	Name                string   `json:"Name"`
	Description         string   `json:"Description"`
	UserClaims          []string `json:"UserClaims"`
	AllowMachineClients bool     `json:"AllowMachineClients"`
	AllowUserClients    bool     `json:"AllowUserClients"`
}
