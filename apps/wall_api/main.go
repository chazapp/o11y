package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "Run the API server",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "dbUser",
						EnvVars: []string{"PG_USER"},
						Usage:   "PostgreSQL User",
					},
					&cli.StringFlag{
						Name:    "dbPassword",
						EnvVars: []string{"PG_PASSWORD"},
						Usage:   "PostgreSQL Password",
					},
					&cli.StringFlag{
						Name:    "dbHost",
						Value:   "localhost",
						EnvVars: []string{"PG_HOST"},
						Usage:   "PostgreSQL Host",
					},
					&cli.StringFlag{
						Name:    "dbName",
						Value:   "mydb",
						EnvVars: []string{"PG_DBNAME"},
						Usage:   "PostgreSQL Database Name",
					},
					&cli.IntFlag{
						Name:    "port",
						Value:   8080,
						EnvVars: []string{"PG_DBNAME"},
						Usage:   "Port of the API Server",
					},
					&cli.StringSliceFlag{
						Name:    "allowedOrigins",
						EnvVars: []string{"ALLOWED_ORIGINS"},
						Usage:   "CORS Allowed Origins",
					},
				},
				Action: func(c *cli.Context) error {
					// Log using zerolog
					log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

					// Read PostgreSQL configuration from command-line flags or environment variables
					dbUser := c.String("dbUser")
					dbPassword := c.String("dbPassword")
					dbHost := c.String("dbHost")
					dbName := c.String("dbName")
					port := c.Int("port")
					allowedOrigins := c.StringSlice("allowedOrigins")
					return API(dbUser, dbPassword, dbHost, dbName, port, allowedOrigins)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("Failed to run the application")
	}
}
