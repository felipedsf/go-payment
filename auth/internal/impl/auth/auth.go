package auth

import (
	"context"
	"database/sql"
	"errors"
	pb "github.com/felipedsf/go-payment/auth/proto"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"os"
	"time"
)

type GrpcAuth struct {
	db *sql.DB
	pb.UnimplementedAuthServiceServer
}

func NewGrpcAuth(db *sql.DB) *GrpcAuth {
	return &GrpcAuth{
		db: db,
	}
}

func (this *GrpcAuth) GetToken(ctx context.Context, credentials *pb.Credentials) (*pb.Token, error) {
	type user struct {
		userID   string
		password string
	}
	var u user

	stmt, err := this.db.Prepare("SELECT user_id, password FROM user WHERE user_id=? AND password=?")
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = stmt.QueryRow(credentials.GetUserName(), credentials.GetUserPass()).Scan(&u.userID, &u.password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	jwtToken, err := createJWT(u.userID)
	if err != nil {
		return nil, err
	}
	return &pb.Token{Jwt: jwtToken}, nil
}

func (this *GrpcAuth) ValidateToken(ctx context.Context, token *pb.Token) (*pb.User, error) {
	key := []byte(os.Getenv("SIGNING_KEY"))

	userID, err := validateJWT(token.Jwt, key)
	if err != nil {
		return nil, err
	}
	return &pb.User{UserID: userID}, nil
}

func validateJWT(t string, signingKey []byte) (string, error) {
	type MyClains struct {
		jwt.RegisteredClaims
	}

	parsedToken, err := jwt.ParseWithClaims(t, &MyClains{}, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", status.Error(codes.Unauthenticated, "token expired")
		}
		return "", status.Error(codes.Unauthenticated, "unauthenticated")
	}

	claims, ok := parsedToken.Claims.(*MyClains)
	if !ok {
		return "", status.Error(codes.Internal, "claims type assertion failed")
	}

	return claims.RegisteredClaims.Subject, nil
}

func createJWT(id string) (string, error) {
	key := []byte(os.Getenv("SIGNING_KEY"))
	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss": "auth-service",
			"sub": id,
			"iat": now.Unix(),
			"exp": now.Add(24 * time.Hour).Unix(),
		})

	signedToken, err := token.SignedString(key)
	if err != nil {
		return "", status.Error(codes.Internal, err.Error())
	}

	return signedToken, nil
}
