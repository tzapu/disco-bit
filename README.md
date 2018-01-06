https://github.com/reactiflux/discord-irc/wiki/Creating-a-discord-bot-&-getting-a-token
Create Application
https://discordapp.com/developers/applications/me
Create a Bot User
Get Bot Token
Add Bot to Server
https://discordapp.com/oauth2/authorize?&client_id=YOUR_CLIENT_ID_HERE&scope=bot&permissions=0

mkdir /opt/config/disco-bit/
git clone https://github.com/tzapu/disco-bit
cd disco-bit
docker build -t disco-bit . 

docker run \
-d --restart always \
--name disco-bit \
-e D_TOKEN="Bot tOkEn.From.DiSc0rd" \
-v /opt/config/disco-bit/:/go/src/app/config/ \
disco-bit:latest -v
