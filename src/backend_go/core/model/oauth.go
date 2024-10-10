package model

const (
	ProviderCustom    OAuthProvider = "custom"
	ProviderMicrosoft OAuthProvider = "microsoft"
	ProviderGoogle    OAuthProvider = "google"
)

var ConnectorAuthProvider = map[SourceType]OAuthProvider{
	SourceTypeOneDrive:    ProviderMicrosoft,
	SourceTypeMsTeams:     ProviderMicrosoft,
	SourceTypeGoogleDrive: ProviderGoogle,
}

// OAuthProvider represents enum for oauth providers
type OAuthProvider string
