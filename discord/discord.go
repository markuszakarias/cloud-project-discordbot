package discord

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"projectGroup23/caching"
	"projectGroup23/utils"
	"projectGroup23/webhook"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

func InitDiscord() {
	// Initializes Discord bot with token
	token := utils.EnvVar("DC_TOKEN")
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err.Error())
	}
	if err = s.Open(); err != nil {
		log.Fatal(err.Error())
	}

	go webhook.WebhookRoutine(s)
	s.AddHandler(messageCreate)
	s.Identify.Intents = discordgo.IntentsGuildMessages

	fmt.Println("Bot is now running. Press CTRL-C to exit...")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	s.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content[:1] != "!" {
		return
	}
	var err error

	switch {
	case len(m.Content) >= 5 && m.Content[:5] == "!help":
		SendHelpMessage(s, m)
	case len(m.Content) >= 11 && m.Content[:11] == "!steamdeals":
		err = SendSteamMessage(s, m)
	case len(m.Content) >= 8 && m.Content[:8] == "!weather":
		err = SendWeatherMessage(s, m)
	case m.Content == "!mealplan":
		dur, _ := time.ParseDuration("20s")
		err = caching.CacheMeals(dur)
		SendMealplanMessage(s, m)
	case len(m.Content) >= 11 && m.Content[:11] == "!newsletter":
		err = SendNewsletterMessage(s, m)
	case len(m.Content) >= 5 && m.Content[:5] == "!todo":
		SendTodoMessage(s, m)
	case len(m.Content) >= 14 && m.Content[:14] == "!notifyweather":
		err = NotifyWeather(s, m)
	case len(m.Content) >= 5 && m.Content[:5] == "!joke":
		err = SendJokeMessage(s, m)

	default:
		s.ChannelMessageSend(m.ChannelID, "Unable to recognize command, try !help (not implemented) if you need a reminder!")
	}

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error: "+err.Error())
	}
}