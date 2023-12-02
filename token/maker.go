package token

import "time"

//Maker is an interface for maneging tokens
type Maker interface {
	//CreateToken cretaes a new token for a pacefic username and duration
	CreateToken(username string, duration time.Duration) (string, error)

	//Verify checks if the token is valid or not
	VerifyToken(token string) (*Payload, error)
}
