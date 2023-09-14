package main

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type ButtonMap map[string][]discordgo.MessageComponent

type Button struct {
	Label string
	Role  string
	Style string
}

func convertMap(fileMap map[string][][]Button, roles []*discordgo.Role) ButtonMap {
	result := make(ButtonMap)
	for k, actionsRow := range fileMap {
		if len(actionsRow) > 5 {
			log.Fatalln("Command has more than 5 action rows")
		}
		currentActionRow := []discordgo.MessageComponent{}
		for _, buttons := range actionsRow {
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
		result[k] = currentActionRow
	}
	return result
}
