package jwt_test

import (
	"testing"
	"urlshortener/pkg/jwt"
)

func TestJWTCreate(t *testing.T) {

	const email = "a@a.ru"

	jwtService := jwt.NewJWT("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ")
	token, err := jwtService.Create(jwt.JWTData{
		Email: email,
	})

	if err != nil {
		t.Fatal(err)
	}

	isValid, data := jwtService.Parse(token)
	if !isValid {
		t.Fatalf("Token in invalid")
	}

	if data.Email != email {
		t.Fatalf("Email %s not eq %s", data.Email, email)
	}
}
