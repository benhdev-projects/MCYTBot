package utils

import "github.com/bwmarrin/discordgo"

func GetMember(s *discordgo.Session, guild, id string) *discordgo.Member {
	member, err := s.State.Member(guild, id)
	if err != nil {
		member, _ := s.GuildMember(guild, id)
		if member != nil {
			s.State.MemberAdd(member)
		}
		return member
	}

	return member
}
