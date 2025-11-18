package main

import (
	"context"
	"log"
	"os"

	"github.com/chazapp/o11y/apps/auth/http"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:        "auth",
		Usage:       "run the auth microservice",
		Description: "The auth microservice. Provides JWT authentication and user management.",
		Commands: []*cli.Command{
			{
				Name: "run",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "port",
						Aliases: []string{"p"},
						Usage:   "Port to bind the HTTP server to",
						Value:   8080,
					},
					&cli.StringFlag{
						Name:    "host",
						Aliases: []string{"h"},
						Usage:   "Host to bind the HTTP server to",
						Value:   "0.0.0.0",
					},
					&cli.StringFlag{
						Name:    "db",
						Aliases: []string{"d"},
						Usage:   "Database connection string",
						Value:   "postgresql://admin:password@localhost:5432/auth",
						Sources: cli.EnvVars("POSTGRES_URL"),
					},
					&cli.StringFlag{
						Name:     "jwt-priv",
						Aliases:  []string{"k"},
						Usage:    "Path to the JWT private key file",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "jwt-pub",
						Aliases:  []string{"K"},
						Usage:    "Path to the JWT public key file",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "domain",
						Aliases: []string{"D"},
						Usage:   "Domain for authentication cookies",
						Value:   "localhost",
					},
				},
				Usage: "Run the auth HTTP server",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					port := cmd.Int("port")
					host := cmd.String("host")
					db := cmd.String("db")
					jwtPrivateKeyPath := cmd.String("jwt-priv")
					jwtPublicKeyPath := cmd.String("jwt-pub")
					domain := cmd.String("domain")

					srv := http.AuthServer{
						Port:              port,
						Host:              host,
						DbConn:            db,
						JwtPrivateKeyPath: jwtPrivateKeyPath,
						JwtPublicKeyPath:  jwtPublicKeyPath,
						Domain:            domain,
					}

					srv.Run(ctx)
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
