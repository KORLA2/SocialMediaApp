package auth

import "github.com/golang-jwt/jwt/v5"

type JWTAuthenticator struct {
	secret   string
	audience string
	issuer   string
}

func NewJWTAuthenticator(secret, audience, issuer string) *JWTAuthenticator {
	return &JWTAuthenticator{
		secret,
		audience,
		issuer,
	}

}

func (j *JWTAuthenticator) GenerateToken(claims jwt.Claims) (string,error) {

   token:= jwt.NewWithClaims(jwt.SigningMethodHS256,claims)

       tokenString,err:=token.SignedString([]byte(j.secret))

	   if err!=nil{
		return "",err;
	   }
	   return tokenString,nil;
}

func (j *JWTAuthenticator) validateToken() bool {

	return false
}
