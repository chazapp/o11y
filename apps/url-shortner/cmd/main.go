package main

import (
	"context"
	"log"
	"os"

	"github.com/chazapp/apps/url-shortner/http"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:        "url-shortner",
		Usage:       "run the url-shortner microservice",
		Description: "The url-shortner. Creates short URLs for given long URLs.",
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
						Aliases: []string{"H"},
						Usage:   "Host to bind the HTTP server to",
						Value:   "0.0.0.0",
					},
					&cli.StringFlag{
						Name:    "db",
						Aliases: []string{"d"},
						Usage:   "Database connection string",
						Value:   "postgresql://admin:password@localhost:5432/url-shortner",
						Sources: cli.EnvVars("POSTGRES_URL"),
					},
				},
				Usage: "Run the auth HTTP server",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					port := cmd.Int("port")
					host := cmd.String("host")
					db := cmd.String("db")
					http.Run(ctx, port, host, db)
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
