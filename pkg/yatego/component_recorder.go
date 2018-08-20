package yatego

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	stPrompt    = "prompt"
	stRecord    = "rec"
	recFilePerm = 0755
)

// Recorder component records voice mail
type Recorder struct {
	Base
	status       string
	recordedFile string
}

// NewRecorderComponent generates new Recorder component
func NewRecorderComponent(base Base) *Recorder {
	r := &Recorder{
		status: stPrompt,
		Base:   base,
	}
	r.Init()
	return r
}

// Init pseudo constructor
func (r *Recorder) Init() {
	r.logger.Debugf("Recorder [%s] init", r.Name())
	r.status = stPrompt
	//install chan.notify to get prompt eof
	r.messagesToInstall[MsgChanNotify] = InstallDef{
		Priority:    100,
		FilterName:  "targetid",
		FilterValue: "{channelID}",
	}
	//install chan.disconnect to fix file if disconnected
	/*r.messagesToInstall[MsgChanDisconnected] = InstallDef{
		Priority:    100,
		FilterName:  "targetid",
		FilterValue: "{channelID}",
	}*/

	//on enter play song
	r.OnEnter(func(call *Call, msg *Message) *CallbackResult {
		if !r.PlayPrompt(call) {
			r.RecordFile(call)
		}
		return NewCallbackResult(ResStay, "")
	})

	//on chan.notify if recorder go transfer, if played start recoding nad stay
	r.Listen(MsgChanNotify, func(call *Call, msg *Message) *CallbackResult {
		r.logger.Debugf("Chan Notify event received with reason: [%s], in recorder status: [%s]", msg.Params["reason"], r.status)
		msg.Processed = true
		switch r.status {
		case stRecord:
			if msg.Params["reason"] != ReasonChanNotifMaxlen {
				r.logger.Debugf("Notify reason is not [maxlen], but [%s] so still waiting...", msg.Params["reason"])
				return NewCallbackResult(ResStay, "")
			}
			//fix file permissions
			//r.fixRecFilePerm(call)

			r.logger.Debugf("Recording done in [%s]", r.Name())
			return r.TransferCallbackResult()
		default:
			r.logger.Debugf("Going to record in [%s]", r.Name())
			r.RecordFile(call)
		}

		return NewCallbackResult(ResStay, "")
	})

	//on chan.disonnected, fix rec. file
	r.Listen(MsgChanNotify, func(call *Call, msg *Message) *CallbackResult {
		msg.Processed = true
		r.fixRecFilePerm(call)
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
	maxlen, exists := r.ConfigAsString("maxlen")
	if !exists {
		maxlen = "80000"
	}
	f := r.recordFilePath(call)
	if f == "" {
		return false
	}
	r.logger.Debugf("Recording voicemail [%s] in [%s]", f, r.Name())
	r.status = stRecord
	r.SetCallData(call, "recorded", f)
	r.Record(f, maxlen, call, map[string]string{})
	r.recordedFile = f
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
	//create file dir if does not exist
	fp := rp.Replace(f.(string))
	dir := filepath.Dir(fp)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, recFilePerm)
	}

	return fp
}

// parse config file path
func (r *Recorder) fixRecFilePerm(call *Call) {
	if r.recordedFile == "" {
		return
	}
	if _, err := os.Stat(r.recordedFile); os.IsNotExist(err) {
		return
	}
	err := os.Chmod(r.recordedFile, recFilePerm)
	if err != nil {
		r.logger.Errorf("Error chmod of file [%s]: [%v]", r.recordedFile, err)
	}
	return
}
