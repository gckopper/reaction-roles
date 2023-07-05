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
-msg
The message that will come with the buttons
```

## map.json file

Inside your map.json file you'll use this schema:

```json
{
    "Label1": "roleName1",
    "Label2": "roleName2",
}
```
