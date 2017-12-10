package bot

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/tzapu/disco-bit/utils"

	log "github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
)

const (
	NEW_USER = iota
	GOT_PASSWORD
	GOT_KEY
	GOT_SECRET
)

type user struct {
	key    string
	secret string
}

type state struct {
	password string
	next     int
}

// Discord holds the bot
type Discord struct {
	Session *discordgo.Session
	users   map[string]*user
	states  map[string]*state
}

func (d *Discord) Start() (err error) {
	// Verify the Token is valid and grab user information
	d.Session.State.User, err = d.Session.User("@me")
	if err != nil {
		log.Fatal("error fetching user information, %s\n", err)
	}
	log.Debug("User info ", d.Session.State.User)

	// Open a websocket connection to Discord
	err = d.Session.Open()
	if err != nil {
		log.Printf("error opening connection to Discord, %s\n", err)
		os.Exit(1)
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	d.Session.AddHandler(d.messageCreate)

	// Wait for a CTRL-C
	log.Printf(`Now running. Press CTRL-C to exit.`)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Clean up
	d.Session.Close()

	return nil
}

func (d *Discord) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	//spew.Dump(s, m)
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	c, err := d.Session.State.Channel(m.ChannelID)
	utils.FatalIfError(err)

	// only works on DMs
	if c.Type != discordgo.ChannelTypeDM {
		log.Debug("Channel mesage")
		s.ChannelMessageSend(m.ChannelID, "Hold you horses there boy, I only take requests in private.")
		dm, err := s.UserChannelCreate(m.Author.ID)
		utils.FatalIfError(err)
		s.ChannelMessageSend(dm.ID, "You can talk to me here")
		//		s.ChannelMessageSend(dm.ID, "I need a password. This will be used to encrypt your key/secret. It will not be save anywhere so please remember it.")
		//		d.states[m.Author.String()].next = GOT_PASSWORD
		return
	}

	// ever seen this user before?
	userID := m.Author.String()
	u, ok := d.users[userID]
	if !ok {
		log.Debug("We don't know ", userID)
		u = &user{}
		d.users[userID] = u
		d.states[userID] = &state{
			password: "",
			next:     NEW_USER,
		}
	}
	log.Println(u)

	switch d.states[userID].next {
	case NEW_USER:
		{
			d.states[userID].next = GOT_PASSWORD
			s.ChannelMessageSend(m.ChannelID, "I need a password. This will be used to encrypt your key/secret. It will not be save anywhere so please remember it.")
			return
		}
	case GOT_PASSWORD:
		{
			pwd := m.Content
			d.states[userID].password = pwd
			d.states[userID].next = GOT_KEY
			s.ChannelMessageSend(m.ChannelID, "Please provide your Bittrex key")
			return
		}
	case GOT_KEY:
		{
			key := m.Content
			d.users[userID].key = key
			d.states[userID].next = GOT_SECRET
			s.ChannelMessageSend(m.ChannelID, "Please provide your Bittrex secret")
			return
		}
	case GOT_SECRET:
		{
			secret := m.Content
			d.users[userID].secret = secret
			s.ChannelMessageSend(m.ChannelID, "Your details have been saved")
			s.ChannelMessageSend(m.ChannelID, "Yoy will be asked for your password everytime the server restarts")
			s.ChannelMessageSend(m.ChannelID, "Order notification has been started")
			return
		}

	}
}

// NewDiscord returns a new Discord bot
func NewDiscord(token string) *Discord {
	var s, _ = discordgo.New()

	s.Token = token

	// Verify a Token was provided
	if s.Token == "" {
		log.Fatal("You must provide a Discord authentication token.")
	}

	return &Discord{
		Session: s,
		users:   map[string]*user{},
		states:  map[string]*state{},
	}
}
