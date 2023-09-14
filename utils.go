package main

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type JSONButton struct {
	Label string
	Role  string
	Style string
}

type JSONCommand struct {
	Message     string `json:",omitempty"`
	Description string `json:",omitempty"`
	Buttons     [][]JSONButton
}

type Command struct {
	Message     string
	Description string
	Buttons     []discordgo.MessageComponent
}
type CommandMap map[string]Command

func convertMap(fileMap map[string]JSONCommand, roles []*discordgo.Role) CommandMap {
	result := make(CommandMap)
    overrideMsg := false
    if *message != "Pick your roles" {
        overrideMsg = true
    }
	for k, cmd := range fileMap {
		if len(cmd.Buttons) > 5 {
			log.Fatalln("Command has more than 5 action rows")
		}
		currentActionRow := []discordgo.MessageComponent{}
		for _, buttons := range cmd.Buttons {
			if len(buttons) > 5 {
				log.Fatalln("Action row has more than 5 buttons")
			}
			currentButton := []discordgo.MessageComponent{}
			for _, button := range buttons {
				var style discordgo.ButtonStyle
				switch strings.ToLower(button.Style) {
				case "blurple":
					style = discordgo.PrimaryButton
				case "grey":
					style = discordgo.SecondaryButton
				case "green":
					style = discordgo.SuccessButton
				case "red":
					style = discordgo.DangerButton
				default:
					fmt.Println(button.Style)
					panic("Unacceptable style for button")
				}
				roleIDidx := sort.Search(len(roles), func(i int) bool {
					return roles[i].Name >= button.Role
				})
				currentButton = append(currentButton, discordgo.Button{
					Label:    button.Label,
					Style:    style,
					CustomID: roles[roleIDidx].ID,
				})
			}
			currentActionRow = append(currentActionRow, discordgo.ActionsRow{
				Components: currentButton,
			})
		}
        msg := cmd.Message
        if overrideMsg || msg == "" {
            msg = *message
        }
		result[k] = Command{
			Buttons:     currentActionRow,
			Description: cmd.Description,
			Message:     msg,
		}
	}
	return result
}
