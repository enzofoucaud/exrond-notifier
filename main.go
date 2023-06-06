package main

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/webhook"
	"github.com/disgoorg/snowflake/v2"
	"github.com/enzofoucaud/exrond-notifier/config"
	"github.com/enzofoucaud/exrond-notifier/exrond"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	// LOG
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	log.Info().Msg("Starting bot")

	c, err := config.GetConfig()
	if err != nil {
		log.Err(err).Msg("error getting config")
		return
	}

	switch c.LogLevel {
	case "INFO":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "DEBUG":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	for {
		pairs, err := exrond.GetPairsQuery()
		if err != nil {
			log.Err(err).Msg("error getting pairs")
			return
		}

		for _, pair := range pairs.Data.Pairs {
			log.Debug().Msg("--------------")
			log.Debug().Msg("Checking pair " + pair.FirstToken.Name + " " + pair.SecondToken.Name)
			log.Debug().Msg("First token price: " + pair.FirstTokenPrice)
			log.Debug().Msg("Second token price: " + pair.SecondTokenPrice)
			// Check if pair is WAGMI
			for _, token := range c.Tokens {
				if token.Token == pair.SecondToken.Name {
					if token.IsBelow {
						pairFloat, _ := strconv.ParseFloat(pair.SecondTokenPrice, 64)
						if pairFloat <= token.Price {
							message := "WAGMI alert: token is below " + pair.SecondTokenPrice
							err := Discord(message, c.DiscordID, c.DiscordToken)
							if err != nil {
								log.Err(err).Msg("error sending discord notification")
							}
							log.Info().Msg("Notifier sent for " + pair.FirstToken.Name + " " + pair.SecondToken.Name)
						}
					}
					if token.IsAbove {
						pairFloat, _ := strconv.ParseFloat(pair.SecondTokenPrice, 64)
						if pairFloat >= token.Price {
							message := "WAGMI alert: token is above " + pair.SecondTokenPrice
							err := Discord(message, c.DiscordID, c.DiscordToken)
							if err != nil {
								log.Err(err).Msg("error sending discord notification")
							}
							log.Info().Msg("Notifier sent for " + pair.FirstToken.Name + " " + pair.SecondToken.Name)
						}
					}
				}
			}
		}
		log.Debug().Msg("--------------")

		time.Sleep(1 * time.Minute)
	}
}

func Discord(message, discordID, discordToken string) error {
	client := webhook.New(snowflake.MustParse(discordID), discordToken)
	defer client.Close(context.TODO())

	if _, err := client.CreateMessage(discord.NewWebhookMessageCreateBuilder().
		SetContent(message).
		Build(),
	); err != nil {
		return err
	}

	return nil
}
