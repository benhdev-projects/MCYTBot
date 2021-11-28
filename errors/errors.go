package errors

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
)

func ReturnError(err sentry.EventID, s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
		Content: fmt.Sprintf("<a:dinkdonk:907702680972058654> <a:panic:907702680938504222> Uhhh... I broke, can you send the error code '%v' to Ben#2028 (<@126429064218017802>) please c:", err),
	})
}

func ReturnAccessDenied(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "You do not have permission to use this command",
		},
	})
}
