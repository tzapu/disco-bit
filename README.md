https://github.com/reactiflux/discord-irc/wiki/Creating-a-discord-bot-&-getting-a-token
Create Application
https://discordapp.com/developers/applications/me
Create a Bot User
Get Bot Token
Add Bot to Server
https://discordapp.com/oauth2/authorize?&client_id=YOUR_CLIENT_ID_HERE&scope=bot&permissions=0


```
docker run -d --restart always --no-cache=true --name disco-bit -e D_TOKEN="Bot tOkEn.From.DiSc0rd" -v /opt/config/disco-bit/:/go/src/app/config/ tzapu/disco-bit:latest -v
```