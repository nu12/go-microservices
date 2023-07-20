package main

import (
	authentication "authentication/grpc"
	"context"
	"log"
)

type AuthenticationServer struct {
	authentication.UnimplementedAuthenticationServer
	Config *Config
}

func (auth *AuthenticationServer) AuthenticateWithEmailAndPassword(ctx context.Context, in *authentication.AuthRequest) (*authentication.AuthResponse, error) {
	log.Println("Processing authentication request")

	user, err := auth.Config.Repo.GetByEmail(in.Email)
	if err != nil {
		log.Println("Error processing authentication: ", err)
		return &authentication.AuthResponse{
			Success: false,
		}, err
	}

	matches, err := auth.Config.Repo.PasswordMatches(in.Password, *user)
	if err != nil {
		log.Println("Error processing authentication: ", err)
		return &authentication.AuthResponse{
			Success: false,
		}, err
	}

	return &authentication.AuthResponse{
		Success: matches,
	}, nil
}
