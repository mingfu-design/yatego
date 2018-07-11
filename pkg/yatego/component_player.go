package yatego

import (
	"os"
	"strings"
)

// Player component plays list of songs
type Player struct {
	currSong int
	Base
}

// NewPlayerComponent generates new player component
func NewPlayerComponent(base Base) *Player {
	p := &Player{
		currSong: 0,
		Base:     base,
	}
	p.Init()
	return p
}

// Init pseudo constructor
func (p *Player) Init() {
	p.logger.Debugf("Player [%s] init", p.Name())
	//install chan.notify to get prompt eof
	p.messagesToInstall[MsgChanNotify] = InstallDef{Priority: 100}

	//on enter play song
	p.OnEnter(func(call *Call, msg *Message) *CallbackResult {
		p.logger.Debugf("Player [%s] on enter", p.Name())
		playOnce := p.playOnePerEntry()
		//reset curr song in normal mode, or loop if no more songs in "playOnce" mode
		if !playOnce || (playOnce && p.currSong >= len(p.songs())) {
			p.currSong = 0
		}
		return p.callbackResult(call)
	})

	//on chan.notify if eof play next song, or transfer
	p.Listen(MsgChanNotify, func(call *Call, msg *Message) *CallbackResult {
		msg.Processed = true
		if msg.Params["reason"] != ReasonChanNotifEOF {
			p.logger.Infof("Notify reason is not [eof], but [%s] so still waiting...", msg.Params["reason"])
			return NewCallbackResult(ResStay, "")
		}

		//if just one per entrance, do not play nothing more this time
		if p.playOnePerEntry() {
			return p.TransferCallbackResult()
		}

		return p.callbackResult(call)
	})
}

// callbackResult adapts plasong out to be returned as
func (p *Player) callbackResult(call *Call) *CallbackResult {
	played := p.PlaySong(call)
	if !played {
		return p.TransferCallbackResult()
	}
	return NewCallbackResult(ResStay, "")
}

// PlaySong plays next song from the playlist or returns false
func (p *Player) PlaySong(call *Call) bool {
	song, exists := p.nextSong(call)
	if !exists {
		return false
	}
	p.PlayWave(song, call, map[string]string{})
	return true
}

// nextSong get next song if exists
func (p *Player) nextSong(call *Call) (string, bool) {
	songs := p.songs()
	p.logger.Debugf("Player [%s] has songs: %v", p.Name(), songs)
	if len(songs) == 0 {
		p.logger.Warningf("Player [%s] playlist has no songs", p.Name())
		return "", false
	}
	if len(songs) <= p.currSong {
		p.logger.Debugf("Player [%s] has no more songs", p.Name())
		return "", false
	}
	song := songs[p.currSong]

	if fi, err := os.Stat(song); os.IsNotExist(err) || !fi.Mode().IsRegular() {
		p.logger.Warningf("Player [%s] song [%s] does not exist", p.Name(), song)
		return "", false
	}

	p.currSong++

	return song, true
}

func (p *Player) songs() []string {
	playlist, ok := p.ConfigAsString("playlist")
	if !ok {
		p.logger.Warningf("Player [%s] has no playlist defined", p.Name())
		return []string{}
	}
	return strings.Split(playlist, ",")
}

func (p *Player) playOnePerEntry() bool {
	play, ok := p.ConfigAsString("play_one_per_entry")
	return ok && play == "true"
}
