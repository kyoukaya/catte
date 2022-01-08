package main

import (
	"context"
	"os"
	"testing"

	"gotest.tools/assert"

	"github.com/bwmarrin/discordgo"
	"github.com/golang/mock/gomock"
	"github.com/kyoukaya/catte/internal/fflogs"
	"github.com/kyoukaya/catte/internal/xivdata"
	"github.com/kyoukaya/catte/mocks"
)

var (
	fflogsClientID = os.Getenv(envPrefix + "FFLOGSID")
	fflogsToken    = os.Getenv(envPrefix + "FFLOGSTOKEN")
)

func TestF(t *testing.T) {
	// TODO: fix up tests with mocks
	input := ""
	fflogsClient, err := fflogs.NewClient(fflogsClientID, fflogsToken)
	assert.NilError(t, err)
	ds, err := xivdata.NewDataSource()
	assert.NilError(t, err)
	app := &App{
		xivds:  ds,
		fflogs: fflogsClient,
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	session := mocks.NewMockDiscordSession(ctrl)
	app.DamageInHandler(context.Background(), session, &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{Data: discordgo.ApplicationCommandInteractionData{
			Options: []*discordgo.ApplicationCommandInteractionDataOption{{Value: input}},
		}},
	})
}

func TestF2(t *testing.T) {
	// TODO: fix up tests with mocks
	input := ""
	fflogsClient, err := fflogs.NewClient(fflogsClientID, fflogsToken)
	assert.NilError(t, err)
	ds, err := xivdata.NewDataSource()
	assert.NilError(t, err)
	app := &App{
		xivds:  ds,
		fflogs: fflogsClient,
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	session := mocks.NewMockDiscordSession(ctrl)
	app.MitigHandler(context.Background(), session, &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{Data: discordgo.ApplicationCommandInteractionData{
			Options: []*discordgo.ApplicationCommandInteractionDataOption{{Value: input}},
		}},
	})
}
