package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"flag"

	"math/rand"

	"github.com/bwmarrin/discordgo"
	"github.com/postables/dogebot/config"
	"github.com/tbruyelle/imgur"
)

var (
	configFile  = flag.String("config.file", "config.json", "config file to use")
	imgurClient *imgur.Client
)

func init() {
	flag.Parse()
}

func main() {
	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		log.Fatal(err)
	}
	if cfg.Discord.Token == "" {
		cfg.Discord.Token = os.Getenv("DISCORD_TOKEN")
	}
	imgurClient = imgur.NewClient(cfg.ImgurClientID)
	// we need to prepend Bot to allow discord
	// to assign permissions properly
	dg, err := discordgo.New("Bot " + cfg.Discord.Token)
	if err != nil {
		fmt.Println("failed to authenticate with discord")
		log.Fatal(err)
	}
	dg.AddHandler(messageCreate)
	if err := dg.Open(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("bot is now running")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// parse the message contents based off string fields
	args := strings.Fields(m.Content)
	if len(args) == 0 {
		return
	}
	// ensure the first field is a valid invocation of dpinner
	if args[0] != "!dogebot" {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if args[1] == "ping" {
		if _, err := s.ChannelMessageSend(m.ChannelID, "Pong!"); err != nil {
			fmt.Println(err)
		}
		return
	}
	// If the message is "pong" reply with "Ping!"
	if args[1] == "pong" {
		if _, err := s.ChannelMessageSend(m.ChannelID, "Ping!"); err != nil {
			fmt.Println(err)
		}
		return
	}
	if args[1] == "search" {
		imresp, _, err := imgurClient.Search(imgur.SearchOptions{All: "shiba inu", Type: "png"})
		if err != nil {
			fmt.Println(err)
			return
		}
		images := imresp.Data
		var image imgur.Image
		for i := 0; i < 1000; i++ {
			if i == 1000 {
				s.ChannelMessageSend(m.ChannelID, "failed to find picture")
				return
			}
			img := images[rand.Intn(len(images))]
			if img.Animated {
				continue
			}
			// skip albums
			if strings.TrimPrefix(img.Link, "https://imgur.com/")[0:2] == "a/" {
				continue
			}
			if strings.Contains(img.Link, "https://imgur.com/a/") {
				continue
			}
			image = img
			break
		}
		if image.Id == "" {
			s.ChannelMessageSend(m.ChannelID, "failed to find picture")
			return
		}
		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Embed: &discordgo.MessageEmbed{
				Image: &discordgo.MessageEmbedImage{
					URL: image.Link,
				},
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: image.Link,
				},
			},
		})
		fmt.Printf("%+v\n", image)
	}
}
