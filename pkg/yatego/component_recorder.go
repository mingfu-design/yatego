package yatego

import (
	"strings"
	"time"
)

const (
	stPrompt = "prompt"
	stRecord = "rec"
)

// Recorder component records voice mail
type Recorder struct {
	Base
	status string
}

// Init pseudo constructor
func (r *Recorder) Init() {
	r.Base.Init()
	r.logger.Debugf("Recorder [%s] init", r.Name())
	r.status = stPrompt
	//install chan.notify to get prompt eof
	r.messagesToInstall[MsgChanNotify] = InstallDef{Priority: 100}

	//on enter play song
	r.OnEnter(func(call *Call, msg *Message) *CallbackResult {
		if !r.PlayPrompt(call) {
			r.RecordFile(call)
		}
		return NewCallbackResult(ResStay, "")
	})

	//on chan.notify if recorder go transfer, if played start recoding nad stay
	r.Listen(MsgChanNotify, func(call *Call, msg *Message) *CallbackResult {
		msg.Processed = true
		switch r.status {
		case stRecord:
			if msg.Params["reason"] != ReasonChanNotifMaxlen {
				r.logger.Debugf("Notify reason is not [maxlen], but [%s] so still waiting...", msg.Params["reason"])
				return NewCallbackResult(ResStay, "")
			}
			r.logger.Debugf("Recording done in [%s]", r.Name())
			return r.TransferCallbackResult()
		default:
			r.logger.Debugf("Going to record in [%s]", r.Name())
			r.RecordFile(call)
		}

		return NewCallbackResult(ResStay, "")
	})
}

// PlayPrompt plays prompt if defined
func (r *Recorder) PlayPrompt(call *Call) bool {
	prompt, exists := r.Config("prompt")
	if !exists {
		r.logger.Debugf("Recorder [%s] has no prompt defined", r.Name())
		return false
	}
	r.status = stPrompt
	r.PlayWave(prompt.(string), call, map[string]string{})
	return true
}

// RecordFile record a voicemail
func (r *Recorder) RecordFile(call *Call) bool {
	maxlen, exists := r.Config("maxlen")
	if !exists {
		maxlen = "80000"
	} else {
	}
	f := r.recordFilePath(call)
	if f == "" {
		return false
	}
	r.logger.Debugf("Recording voicemail [%s] in [%s]", f, r.Name())
	r.status = stRecord
	r.SetCallData(call, "recorded", f)
	r.Record(f, maxlen.(string), call, map[string]string{})
	return true
}

// parse config file path
func (r *Recorder) recordFilePath(call *Call) string {
	f, exists := r.Config("file")
	if !exists {
		r.logger.Errorf("File path not defined for recorder [%s]", r.Name())
		return ""
	}
	rp := strings.NewReplacer(
		"{caller}", call.Caller,
		"{called}", call.Called,
		"{billingId}", call.BillingID,
		"{time}", time.Now().Format("2006-01-02T15:04:05Z"),
	)
	return rp.Replace(f.(string))
}
