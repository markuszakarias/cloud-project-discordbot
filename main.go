package main

import (
	"fmt"
	"os"
	"os/signal"
	"projectGroup23/firebase"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	var err error

	firebase.InitFirebase()

	token := "ODM2OTgzNjUyMjUxMzM2Nzc1.YIl7xQ.cuxQXG5lW9Sqmylm6rx4INNiLpc"

	var s, _ = discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	s.AddHandler(messageCreate)
	s.Identify.Intents = discordgo.IntentsGuildMessages

	if err = s.Open(); err != nil {
		panic(err)
	}

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
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "drit!")
	}

	if m.Content == "jokes" {
		jokes := firebase.GetAllJokes()
		for _, a := range jokes {
			s.ChannelMessageSend(m.ChannelID, a)
		}
	}

	if strings.HasPrefix(m.Content, "!createjoke ") {
		joke := m.Content[12:]
		err := firebase.CreateJoke(s.State.User.ID, joke)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
		} else {
			s.ChannelMessageSend(m.ChannelID, "Joke created")
		}
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}

	if m.Content == "minbrukerid" {
		fmt.Println("asd")
		s.ChannelMessageSend(m.ChannelID, s.State.User.ID)
	}

	fmt.Println(s)

}
