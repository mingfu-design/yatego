package yatego

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

//Message Type constants
const (
	TypeIncoming     = "incoming"
	TypeOutgoing     = "outgoing"
	TypeDispatched   = "dispatched"
	TypeAcknowledged = "acknowledged"
	TypeAnswer       = "answer"
	TypeInstalled    = "installed"
	TypeUninstalled  = "uninstalled"
	TypeWatched      = "watched"
	TypeUnwatched    = "unwatched"
	TypeConnected    = "connected"
	TypeSetLocal     = "setlocal"
)

const letterBytes = "123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

//A Message is Yate message
type Message struct {
	Type      string
	Name      string
	RetVal    string
	ID        string
	Params    map[string]string
	Processed bool
	Time      int64
}

//NewMessageRetVal creates new outgoing message
func NewMessageRetVal(name string, retVal string, id string, msgType string) *Message {
	if id == "" {
		id = RandString(10)
	}
	if msgType == "" {
		msgType = TypeOutgoing
	}
	return &Message{
		Type:      msgType,
		Name:      name,
		RetVal:    retVal,
		ID:        id,
		Processed: false,
		Time:      time.Now().Unix(),
		Params:    make(map[string]string),
	}
}

//NewMessage creates new outgoing message without ret value
func NewMessage(name string, params map[string]string) *Message {
	m := NewMessageRetVal(name, "", "", TypeOutgoing)
	if params == nil {
		return m
	}
	for k, v := range params {
		m.Params[k] = v
	}
	return m
}

func (m *Message) encodeParams() string {
	result := ""
	if m.Params == nil {
		return result
	}
	for k, v := range m.Params {
		result += ":" + esc(k) + "=" + esc(v)
	}
	return result
}

//Encode returns a message in yate protocol format
func (m *Message) Encode() string {
	result := ""
	switch m.Type {
	//dispatch message
	case TypeOutgoing:
		result = fmt.Sprintf("%%%%>message:%s:%d:%s:%s",
			esc(m.ID), m.Time, esc(m.Name), esc(m.RetVal)) + m.encodeParams()
		//acknowledge message
	case TypeIncoming:
		result = fmt.Sprintf("%%%%<message:%s:%s:%s:%s",
			esc(m.ID), esc(bool2str(m.Processed)), esc(m.Name), esc(m.RetVal)) + m.encodeParams()
	}
	return result
}

//DecodeMessage returns a message from received yate string
func DecodeMessage(s string) (*Message, error) {
	var m *Message
	parts := strings.Split(s, ":")
	switch parts[0] {
	case "%%>message":
		m = NewMessageRetVal(unesc(parts[3]), unesc(parts[4]), unesc(parts[1]), TypeIncoming)
		ts, err := strconv.ParseInt(parts[2], 10, 64)
		if err == nil {
			m.Time = ts
		}
		m.decodeParams(parts[5:])
	case "%%<message":
		m = NewMessageRetVal(unesc(parts[3]), unesc(parts[4]), unesc(parts[1]), TypeAnswer)
		m.Processed = str2bool(parts[2])
		m.decodeParams(parts[5:])
	case "%%<install":
		m = decodeMessageNoParams(TypeInstalled, parts)
	case "%%<uninstall":
		m = decodeMessageNoParams(TypeUninstalled, parts)
	case "%%<watch":
		m = decodeMessageNoID(TypeWatched, parts)
	case "%%<unwatch":
		m = decodeMessageNoID(TypeUnwatched, parts)
	case "%%<connect":
		m = decodeMessageNoID(TypeConnected, parts)
	case "%%<setlocal":
		m = NewMessageRetVal(unesc(parts[1]), unesc(parts[2]), "", TypeSetLocal)
		m.Processed = str2bool(parts[3])
	}

	return m, nil
}

func decodeMessageNoParams(messageType string, parts []string) *Message {
	m := NewMessageRetVal(unesc(parts[2]), "", unesc(parts[1]), messageType)
	m.Processed = str2bool(parts[2])
	return m
}

func decodeMessageNoID(messageType string, parts []string) *Message {
	m := NewMessageRetVal(unesc(parts[1]), "", "", messageType)
	m.Processed = str2bool(parts[2])
	return m
}

func (m *Message) decodeParams(parts []string) {
	for _, v := range parts {
		raw := strings.SplitN(v, "=", 2)
		if len(raw) < 2 {
			continue
		}
		m.Params[raw[0]] = unesc(raw[1])
	}
}

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func esc(str string) string {
	str = str + ""
	s := ""
	n := len(str)
	i := 0
	for i < n {
		c, size := utf8.DecodeRuneInString(str[i:])
		if c < 32 || c == ':' {
			c = rune(c + 64)
			s = s + "%"
		} else if c == '%' {
			s = s + string(c)
		}
		s = s + string(c)
		i = i + size

	}
	return s
}

func unesc(str string) string {
	s := ""
	n := len(str)
	i := 0
	for i < n {
		c, size := utf8.DecodeRuneInString(str[i:])
		if c == '%' {
			i = i + size
			c, size = utf8.DecodeRuneInString(str[i:])
			if c != '%' {
				c = rune(c - 64)
			}
		}
		s = s + string(c)
		i = i + size
	}
	return s
}

func str2bool(s string) bool {
	if s == "true" {
		return true
	}
	return false
}

func bool2str(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
