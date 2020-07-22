package shodan

//BaseURL is the ... base url
const BaseURL = "https://api.shodan.io"

//Client struct has an apiKey
type Client struct {
	apiKey string
}

//New creates a new Shodan client given an API key
func New(apiKey string) *Client {
	return &Client{apiKey: apiKey}
}
