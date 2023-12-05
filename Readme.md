# Reaction role discord bot

## Using

Download the binary for your system and run with the disered arguments

``` no
-token
Your discord bot token from the discord dev console
-app
Your discord app ID from the discord dev console
-guild
The guild id for the server
-mapping
A json file mapping a button label to a role name
```

## map.json file

### Limitations

- You can only have a maximum of 5 rows per command and 5 buttons per row (so a total of 25 roles per command) (Discord limitation)
- Command names should follow this `/^[-_\p{L}\p{N}\p{sc=Deva}\p{sc=Thai}]{1,32}$/gu` JS style regex (Discord limitation)
  - You can test with `"your-command-name".match(/^[-_\p{L}\p{N}\p{sc=Deva}\p{sc=Thai}]{1,32}$/gu)` in a JS runtime such as your browsers console

### Available Styles

- ![#5865f2](images/red.png) red
- ![#5865f2](images/blurple.png) blurple
- ![#5865f2](images/green.png) green
- ![#5865f2](images/grey.png) grey

### Example

Inside your map.json file you'll use this schema:

```json
{
    "first-command-name": {
        "Description": "Description for your first command",
        "Message": "Message that will be sent with every invocation of this command",
        "Buttons": [
            [
                {
                    "Role":"role name in your server",
                    "Label": "Text written in this button",
                    "Style": "red"
                },
                {
                    "Role":"Another role",
                    "Label": "Different label",
                    "Style": "grey"
                }
            ],
            [
                {
                    "Role":"A role in a different row!",
                    "Label": "Label in a different row!",
                    "Style": "blurple"
                }
            ]
        ]
    },
    "second-command": {
        "Decription": "Description of your second-command",
        "Message": "A massage to be sent with the second-command",
        "Buttons": [
            [
                {
                    "Role":"Different command",
                    "Label": "Role in a diffferent command",
                    "Style": "red"
                }
            ],
            [
                {
                    "Role":"d",
                    "Label": "Role deez nuts",
                    "Style": "blurple"
                }
            ]
        ]
    }
}
```
