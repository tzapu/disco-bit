package bot

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/tzapu/disco-bit/encryption"
	"github.com/tzapu/disco-bit/exchange"
	"github.com/tzapu/disco-bit/persistance"
	"github.com/tzapu/disco-bit/utils"

	log "github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
)

const (
	NEW_USER = iota
	GOT_PASSWORD
	GOT_KEY
	GOT_SECRET
	SHOULD_USE_PASSWORD
)

type user struct {
	Key    []byte
	Secret []byte
}

type state struct {
	next    int
	channel *discordgo.Channel
	bittrex *exchange.Bittrex
}

// Discord holds the bot
type Discord struct {
	Session   *discordgo.Session
	users     map[string]*user
	states    map[string]*state
	passwords map[string]*string
	receiver  chan utils.Message
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

	err = d.loadUsers()
	if err != nil {
		log.Error("Can't load users ", err)
	}

	for id := range d.users {
		dm, err := d.Session.UserChannelCreate(id)
		if err != nil {
			log.Error(err)
			continue
		}
		d.Session.ChannelMessageSend(dm.ID, "I have been restarted")
		d.Session.ChannelMessageSend(dm.ID, "Please send me your password so I can continue sending you notifications")
		d.states[id] = &state{
			next:    SHOULD_USE_PASSWORD,
			channel: dm,
		}
	}

	go d.receive()

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

	if c.Name != "disco-bit" && c.Type != discordgo.ChannelTypeDM {
		return
	}

	// only works on DMs
	if c.Type != discordgo.ChannelTypeDM {
		log.Debug("Channel mesage")
		s.ChannelMessageSend(m.ChannelID, "Hold you horses there boy, I only take requests in private.")
		dm, err := s.UserChannelCreate(m.Author.ID)
		utils.FatalIfError(err)
		s.ChannelMessageSend(dm.ID, "You can talk to me here")
		return
	}

	// ever seen this user before?
	userID := m.Author.ID
	u, ok := d.users[userID]
	if !ok {
		log.Debug("We don't know ", userID)
		u = &user{}
		d.users[userID] = u
		d.states[userID] = &state{
			next:    NEW_USER,
			channel: c,
		}
		e := ""
		d.passwords[userID] = &e
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
			d.passwords[userID] = &pwd
			d.states[userID].next = GOT_KEY
			s.ChannelMessageSend(m.ChannelID, "Please provide your Bittrex key")
			return
		}
	case GOT_KEY:
		{
			key := m.Content
			ek, err := encryption.Encrypt([]byte(*d.passwords[userID]), []byte(key))
			utils.FatalIfError(err)
			d.users[userID].Key = ek
			d.states[userID].next = GOT_SECRET
			s.ChannelMessageSend(m.ChannelID, "Please provide your Bittrex secret")
			return
		}
	case GOT_SECRET:
		{
			secret := m.Content
			es, err := encryption.Encrypt([]byte(*d.passwords[userID]), []byte(secret))
			utils.FatalIfError(err)
			d.users[userID].Secret = es
			err = d.saveUsers()
			utils.FatalIfError(err)

			s.ChannelMessageSend(m.ChannelID, "Your details have been saved")
			s.ChannelMessageSend(m.ChannelID, "Yoy will be asked for your password everytime the server restarts")
			s.ChannelMessageSend(m.ChannelID, "Order notification has been started")

			d.monitor(userID)

			return
		}
	case SHOULD_USE_PASSWORD:
		{
			pwd := m.Content
			d.passwords[userID] = &pwd
			d.monitor(userID)
			return
		}
	}
}

func (d *Discord) saveUsers() error {
	return persistance.Save("config/users.gob", d.users)
}

func (d *Discord) loadUsers() error {
	err := persistance.Load("config/users.gob", &d.users)
	return err
}

func (d *Discord) monitor(id string) {
	p := []byte(*d.passwords[id])
	kb, err := encryption.Decrypt(p, d.users[id].Key)
	utils.ErrorIfError(err)
	sb, err := encryption.Decrypt(p, d.users[id].Secret)
	utils.ErrorIfError(err)
	d.states[id].bittrex = exchange.NewBittrex(string(kb), string(sb), id, d.receiver)
	d.states[id].bittrex.Start()
}

func (d *Discord) receive() {
	for m := range d.receiver {
		log.Println(m)
		_, err := d.Session.ChannelMessageSend(d.states[m.ID].channel.ID, m.Text)
		utils.FatalIfError(err)
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
		Session:   s,
		users:     map[string]*user{},
		states:    map[string]*state{},
		passwords: map[string]*string{},
		receiver:  make(chan utils.Message),
	}
}
