package http

type Config struct {
	Address        string   `short:"a" long:"address" env:"ADDRESS" description:"Service address" required:"yes"`
	JWTPrivateKey  string   `long:"jwt-private-key" env:"JWT_PRIVATE_KEY" description:"Path to JWT private key" required:"yes"`
	AllowedOrigins []string `long:"allowed-origins" env:"ALLOWED_ORIGINS" description:"Allowed origins to use CORS" env-delim:"," required:"yes"`
}
