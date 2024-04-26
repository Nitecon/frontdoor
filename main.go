package main

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
)

var (
	version = "source"
)

func setLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	target := fmt.Sprintf("https://%s%s", r.Host, r.RequestURI)
	log.Info().Msgf("Redirecting %s to %s", r.Host, target)
	http.Redirect(w, r, target, http.StatusSeeOther)
}

func startHttpServer() *http.Server {
	server := &http.Server{Addr: ":80", Handler: http.HandlerFunc(RedirectHandler)}
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("HTTP server failed")
		}
	}()
	return server
}

func startHttpsServer(host, key, cert string) *http.Server {
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   host,
		Path:   "/",
	})
	proxy.Transport = &http.Transport{
		// Enforce HTTP/1.1 for backend communication if needed
	}

	server := &http.Server{Addr: ":443", Handler: proxy}
	go func() {
		if err := server.ListenAndServeTLS(cert, key); err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("HTTPS server failed")
		}
	}()
	return server
}

func main() {
	setLogger()
	log.Info().Msgf("Starting FrontDoor (Version: %s)", version)

	app := &cli.App{
		Name:  "FrontDoor",
		Usage: "Start with ./frontdoor -key path/to/server.key -cert path/to/server.crt -backend 127.0.0.1:8080",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "key",
				Required: true,
				Aliases:  []string{"k"},
				Usage:    "Path to the certificate key file, e.g. /path/to/server.key",
			},
			&cli.StringFlag{
				Name:     "cert",
				Required: true,
				Aliases:  []string{"c"},
				Usage:    "Path to the certificate file, e.g. /path/to/server.crt",
			},
			&cli.StringFlag{
				Name:     "backend",
				Required: true,
				Aliases:  []string{"b"},
				Usage:    "Backend server address, e.g. 127.0.0.1:8080",
			},
		},
		Action: func(c *cli.Context) error {
			httpServer := startHttpServer()
			httpsServer := startHttpsServer(c.String("backend"), c.String("key"), c.String("cert"))

			stop := make(chan os.Signal, 1)
			signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

			<-stop

			// Shutdown HTTP server
			if err := httpServer.Shutdown(context.Background()); err != nil {
				log.Error().Err(err).Msg("Failed to shutdown HTTP server")
			}

			// Shutdown HTTPS server
			if err := httpsServer.Shutdown(context.Background()); err != nil {
				log.Error().Err(err).Msg("Failed to shutdown HTTPS server")
			}

			log.Info().Msg("Shutting down server")
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("Failed to start application")
	}
}
