package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sort"

	"github.com/bwmarrin/discordgo"
)

// Bot parameters
var (
	GuildID     = flag.String("guild", "", "Test guild ID")
	BotToken    = flag.String("token", "", "Bot access token")
	AppID       = flag.String("app", "", "Application ID")
	MappingFile = flag.String("mapping", "map.json", "Json file mapping roles with buttom labels")
	message     = flag.String("msg", "Pick your roles", "Text to be shown with the buttons")
)

var s *discordgo.Session

var buttonMap ButtonMap

func init() { flag.Parse() }

func init() {
	var err error
	file, err := os.ReadFile(*MappingFile)
	if err != nil {
		log.Fatalf("I NEED A MAP FILE: %v", err)
	}
	var mappings map[string][][]Button
	err = json.Unmarshal(file, &mappings)
	if err != nil {
		log.Fatalf("FAILED TO PARSE YOUR JSON: %v", err)
	}
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		roles, err := s.GuildRoles(*GuildID)
		if err != nil {
			log.Fatalf("Failed to fetch guild roles with error: %v", err)
		}
		sort.Slice(roles, func(i, j int) bool {
			return roles[i].Name < roles[j].Name
		})
		buttonMap = convertMap(mappings, roles)
		log.Println("Bot is up!")
	})
}

// Important note: call every command in order it's placed in the example.

func InteractionResponseEphemeral(s *discordgo.Session, i *discordgo.InteractionCreate, str string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: str,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Println("[ERROR] [InteractionResponseEphemeral]", err)
	}
}

func roleToggle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var err error
	role := i.MessageComponentData().CustomID
	member := i.Member
	if member == nil {
		InteractionResponseEphemeral(s, i, "Cant figure out your ID so... FUCK YOU")
		log.Println("[ERROR]", "Cant figure out user ID")
		return
	}

	id := i.Member.User.ID
	for _, v := range i.Member.Roles {
		if v == role {
			err = s.GuildMemberRoleRemove(i.GuildID, id, role)
			if err != nil {
				InteractionResponseEphemeral(s, i, "Failed to remove your role! Please contact an admin!")
				log.Printf("[ERROR] [Removing role {%s}] %s\n", role, err)
				return
			}
			InteractionResponseEphemeral(s, i, "Role removed")
			return
		}
	}
	err = s.GuildMemberRoleAdd(i.GuildID, id, role)
	if err != nil {
		log.Printf("[ERROR] [Adding role {%s}] %s\n", role, err)
		InteractionResponseEphemeral(s, i, "Failed to add your role! Please contact an admin!")
		return
	}
	InteractionResponseEphemeral(s, i, "Role added")
}

func sendButton(s *discordgo.Session, i *discordgo.InteractionCreate, components []discordgo.MessageComponent) {
	flags := discordgo.MessageFlagsEphemeral
	if i.Member.Permissions&discordgo.PermissionAdministrator != 0 {
		flags = 0
	}
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    *message,
			Flags:      flags,
			Components: components,
		},
	})
	if err != nil {
		log.Println("[ERROR] [sendButton]", err)
	}
}
func main() {
	var err error
	// Components are part of interactions, so we register InteractionCreate handler
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := buttonMap[i.ApplicationCommandData().Name]; ok {
				sendButton(s, i, h)
			}
		case discordgo.InteractionMessageComponent:

			roleToggle(s, i)
		}
	})
	for k := range buttonMap {
		_, err = s.ApplicationCommandCreate(*AppID, *GuildID, &discordgo.ApplicationCommand{
			Name:        k,
			Description: fmt.Sprint("Get role menu for", k),
		})
		if err != nil {
			log.Fatalf("Cannot create slash command: %v", err)
		}
	}

	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Graceful shutdown")
}
