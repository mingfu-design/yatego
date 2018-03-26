package yatego

import (
	"strings"
)

// Player component plays list of songs
type Player struct {
	currSong int
	Base
}

// hook method
func (p *Player) initListeners() {
	//install chan.notify to get prompt eof
	p.messagesToInstall[MsgChanNotify] = InstallDef{Priority: 100}

	//on enter play song
	p.OnEnter(func(call *Call, msg *Message) *CallbackResult {
		return p.callbackResult(call)
	})

	//on chan.notify if eof play next song, or transfer
	p.Listen(MsgChanNotify, func(call *Call, msg *Message) *CallbackResult {
		msg.Processed = true
		if msg.Params["reason"] != ReasonChanNotifEOF {
			p.logger.Infof("Notify reason is not [eof], but [%s] so still waiting...", msg.Params["reason"])
			return NewCallbackResult(ResStay, "")
		}

		return p.callbackResult(call)
	})
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

// callbackResult adapts plasong out to be returned as
func (p *Player) callbackResult(call *Call) *CallbackResult {
	played := p.PlaySong(call)
	if !played {
		return p.TransferCallbackResult()
	}
	return NewCallbackResult(ResStay, "")
}

// nextSong get next song if exists
func (p *Player) nextSong(call *Call) (string, bool) {
	playlist, exists := p.Config("playlist")
	if !exists {
		p.logger.Warningf("Player [%s] has no playlist defined", p.Name())
		return "", false
	}
	songs := strings.Split(playlist.(string), ",")
	if len(songs) == 0 {
		p.logger.Warningf("Player [%s] playlist has no songs", p.Name())
		return "", false
	}
	if len(songs) <= p.currSong {
		p.logger.Debugf("Player [%s] has no more songs", p.Name())
		return "", false
	}
	song := songs[p.currSong]
	p.currSong++
	return song, true
}
