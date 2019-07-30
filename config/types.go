package config

// Config is used to handle configuration of dpinner
type Config struct {
	Discord       `json:"discord"`
	ImgurClientID string `json:"imgur_client_id"`
}

// Discord is used to configure discord both authentication
type Discord struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Token        string `json:"token"`
}
