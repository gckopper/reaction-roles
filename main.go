package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

// Bot parameters
var (
	GuildID     = flag.String("guild", "", "Test guild ID")
	BotToken    = flag.String("token", "", "Bot access token")
	AppID       = flag.String("app", "", "Application ID")
	MappingFile = flag.String("mapping", "map.json", "Json file mapping roles with buttom labels")
)

var s *discordgo.Session

var mappings map[string]string

func init() { flag.Parse() }

func init() {
	var err error
	file, err := os.ReadFile(*MappingFile)
	if err != nil {
		log.Fatalf("I NEED A MAP FILE: %v", err)
	}
	err = json.Unmarshal(file, &mappings)
	if err != nil {
		log.Fatalf("FAILED TO PARSE YOUR JSON: %v", err)
	}
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

// Important note: call every command in order it's placed in the example.

func InteractionResponseFailure(s *discordgo.Session, i *discordgo.InteractionCreate, str string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: str,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		panic(err)
	}
}

func roleToggle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var err error
	input := i.MessageComponentData().CustomID
	roleName, exist := mappings[input]
	if !exist {
		InteractionResponseFailure(s, i, "This button is not mapped to anything!")
		return
	}
	guildRoles, err := s.GuildRoles(i.GuildID)
	if err != nil {
		InteractionResponseFailure(s, i, "Failed to fetch guild roles")
		return
	}
	var role string
	for _, v := range guildRoles {
		if v.Name == roleName {
			role = v.ID
		}
	}
	if role == "" {
		InteractionResponseFailure(s, i, "Role does not exist")
		return
	}
	member := i.Member
	if member == nil {
		InteractionResponseFailure(s, i, "Cant figure out your ID so... FUCK YOU")
		return
	}

	id := i.Member.User.ID
	for _, v := range i.Member.Roles {
		if v == role {
			s.GuildMemberRoleRemove(i.GuildID, id, role)
			InteractionResponseFailure(s, i, "Role removed")
			return
		}
	}
	err = s.GuildMemberRoleAdd(i.GuildID, id, role)
	if err != nil {
		panic(err)
	}
	InteractionResponseFailure(s, i, "Role added")
}

var (
	commandsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"roles": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			components := make([]discordgo.MessageComponent, ((len(mappings)-1)/5)+1)
			index := 0
			var current *[]discordgo.MessageComponent
			for k, v := range mappings {
				if index%5 == 0 {
					current = &[]discordgo.MessageComponent{}
				}
				// fmt.Println("current ", index, ": ", current)
				*current = append(*current, discordgo.Button{
					Label: k,
					Style: discordgo.SuccessButton,
					// Disabled allows bot to disable some buttons for users.
					Disabled: false,
					CustomID: v,
				})
				index += 1
				if index%5 == 0 {
					components[(index-1)/5] = discordgo.ActionsRow{
						Components: *current,
					}
				}
			}
			if index%5 != 0 {
				components[((index - 1) / 5)] = discordgo.ActionsRow{
					Components: *current,
				}
			}
			flags := discordgo.MessageFlagsEphemeral
			if i.Member.Permissions == discordgo.PermissionAdministrator {
				flags = 0
			}
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content:    "Pick your roles",
					Flags:      flags,
					Components: components,
				},
			})
			if err != nil {
				panic(err)
			}
		},
	}
)

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is up!")
	})
	// Components are part of interactions, so we register InteractionCreate handler
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := commandsHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case discordgo.InteractionMessageComponent:

			roleToggle(s, i)
		}
	})
	_, err := s.ApplicationCommandCreate(*AppID, *GuildID, &discordgo.ApplicationCommand{
		Name:        "roles",
		Description: "Get roles menu",
	})

	if err != nil {
		log.Fatalf("Cannot create slash command: %v", err)
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
