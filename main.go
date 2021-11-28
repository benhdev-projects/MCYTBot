package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"benh.codes/mcytbot/commands"
	"benh.codes/mcytbot/db"
	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
)

var guildId string
var botToken string

var s *discordgo.Session

func init() {
	godotenv.Load()
	botToken = os.Getenv("TOKEN")
	guildId = os.Getenv("GUILD")
	flag.Parse()
}

func init() {
	var err error
	s, err = discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

var (
	cmds         = commands.GetCommands()
	cmdHandlers  = commands.GetCommandHandlers()
	compHandlers = commands.GetComponentHandlers()
)

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := cmdHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}

		case discordgo.InteractionMessageComponent:
			if h, ok := compHandlers[i.MessageComponentData().CustomID]; ok {
				h(s, i)
			}
		}
	})
}

func main() {
	go func() {
		db.Init()

		err := sentry.Init(sentry.ClientOptions{
			Dsn: "https://69f57859a4fd47e6a2dafa260fc6211a@o388067.ingest.sentry.io/6056906",
		})
		if err != nil {
			log.Fatalf("sentry.Init: %s", err)
		}
		s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
			log.Println("Bot is up!")
			s.State.TrackMembers = true
			guilds, _ := s.UserGuilds(100, "", "")
			for _, g := range guilds {
				fmt.Println("Loading in information for " + g.ID)
				fullGuild, _ := s.Guild(g.ID)
				s.State.GuildAdd(fullGuild)

				var exists bool
				if db.DB.QueryRow("SELECT exists (SELECT 1 from guilds WHERE id = $1)", g.ID).Scan(&exists); !exists {
					_, err := db.DB.Exec(`INSERT INTO guilds (id) VALUES ($1)`, g.ID)
					if err != nil {
						fmt.Printf("Error inserting guild into database: %v\n", err)
					}
				}
			}
		})
		err = s.Open()
		if err != nil {
			sentry.CaptureException(err)
			log.Fatalf("Cannot open the session: %v", err)
		}

		for _, v := range cmds {
			_, err := s.ApplicationCommandCreate(s.State.User.ID, guildId, v)
			if err != nil {
				sentry.CaptureException(err)
				log.Panicf("Cannot create '%v' command: %v", v.Name, err)
			}
		}

	}()

	shutdown := make(chan int)

	//create a notification channel to shutdown
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("Shutting down...")
		shutdown <- 1
	}()

	<-shutdown

	s.Close()
	db.DB.Close()

	cmds, _ := s.ApplicationCommands(s.State.User.ID, guildId)
	for _, a := range cmds {
		err := s.ApplicationCommandDelete(s.State.User.ID, guildId, a.ID)
		if err != nil {
			sentry.CaptureException(err)
			log.Printf("Cannot delete '%v' command: %v", a.Name, err)
		}
	}
}
