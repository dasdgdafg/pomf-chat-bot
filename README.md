chat bot for pomf.tv

### example code
`example/example.go` has a small bot that just replies to set commands. If you're looking to write your own, that's a simple place to look at as an example of how to use the API.

### using the bot
`bot/main.go` has the main functionality, which is currently only replying to set commands and the !sr command.

To use, first copy the `example/settings.json` file, and put it into the `bot` folder. Replace `put your username here` with your username. 

This program is written in Go. If you haven't used Go before, you can install it from [https://golang.org/](https://golang.org/). Once installed, you can run the bot with `go run .` from the `bot` directory.

#### commands
Guest accounts are rate limited when chatting, especially when trying to post the exact same message multiple times, so if you're going to use the generic !commands with set responses you will need to create an account for your bot, and fill in the bot's username and API key into the settings file.

Fill in the `Commands` object with the commands you want to use. Keys are the commands that the user will type, and values are the response from the bot.

#### !sr
Add `"SongRequestEnabled":true` to the settings file to enable the song request functionality. Put it at the same level as `StreamerName`/`BotName` etc.

Check your VLC settings, and make sure `Allow only one instance` is enabled (otherwise it will open a new instance of VLC for each request). You probably also want to turn off shuffle and repeat. On windows you may need to change `"vlc"` to `"vlc.exe"` in the code.

### sample settings.json file
```
{
  "StreamerName":"xxx_cool_dude_xxx_69_420",
  "BotName":"cool_robot",
  "Apikey":"69696969696969696969696969696969",
  "Commands":{
    "!discord":"https://discord.gg/0123456789",
    "!rules":"lol just shitpost",
    "!twitter":"https://twitter.com/xx_some_username_xx"
  },
  "SongRequestEnabled":true
}
```
