package yatego

import (
	"bytes"
	"io"
	"regexp"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestDispatch(t *testing.T) {
	b := new(bytes.Buffer)
	e := engine(b, b)
	i, err := e.Dispatch(NewMessage(MsgCallAnswered, map[string]string{"fake": "test"}))
	if err != nil {
		t.Fatalf("Dispatch failed: %s", err)
	}
	if i == 0 {
		t.Fatal("Dispatch no bytes printed to out")
	}
	msg := b.String()
	t.Logf("Dispatched msg string: [%s]", msg)

	pattern := "%%>message:.*:.*:call.answered::fake=test\n"
	match, _ := regexp.MatchString(pattern, msg)
	if err != nil {
		t.Fatalf("Dispatch output match failed: %s", err)
	}
	if !match {
		t.Fatalf("Dispatch output [%s] not matched by pattern [%s]", msg, pattern)
	}
}

func TestGetEvent(t *testing.T) {
	b := new(bytes.Buffer)
	e := engine(b, b)
	b.WriteString("%%>message:0x7f4e5df469e0.672973892:1522227506:call.execute::id=sip/5:module=sip:status=incoming:address=172.28.128.1%z34084:billid=1522152787-5:answered=false:direction=incoming:callid=sip/-altWCH3K8SSDe2HGBX0sA../0ca10a05/:caller=41587000201:called=923:antiloop=19:ip_host=172.28.128.1:ip_port=34084:ip_transport=UDP:connection_id=general:connection_reliable=false:sip_uri=sip%z923@172.28.128.3;transport=UDP:sip_from=sip%z41587000201@172.28.128.3;transport=UDP:sip_to=<sip%z923@172.28.128.3;transport=UDP>:sip_callid=-altWCH3K8SSDe2HGBX0sA..:device=Z 5.2.10 rv2.8.75-mod:sip_contact=<sip%z41587000201@172.28.128.1%z34084;transport=UDP>:sip_allow=INVITE, ACK, CANCEL, BYE, NOTIFY, REFER, MESSAGE, OPTIONS, INFO, SUBSCRIBE:sip_content-type=application/sdp:sip_user-agent=Z 5.2.10 rv2.8.75-mod:sip_allow-events=presence, kpml, talk:rtp_addr=172.28.128.1:media=yes:formats=mulaw,alaw,ilbc20:transport=RTP/AVP:rtp_mapping=mulaw=0,alaw=8,ilbc20=97:rtp_rfc2833=99:rtp_port=30000:sdp_fmtp%z=minptime=20; cbr=1; maxaveragebitrate=40000; useinbandfec=1:sdp_sendrecv=:sdp_fmtp%z=0-16:sdp_fmtp%z=0-16:sdp_fmtp%z=0-16:sdp_fmtp%z=0-16:rtp_forward=possible:handlers=javascript%z15,regexroute%z100,javascript%z15,cdrbuild%z50,fileinfo%z90,subscription%z100,sip%z100,regexroute%z100,javascript%z15,gvoice%z20,queues%z45,cdrbuild%z50,yrtp%z50,lateroute%z75,dbwave%z90,dumb%z90,sip%z90,wave%z90,filetransfer%z90,tone%z90,conf%z90,iax%z90,analyzer%z90,jingle%z90,sig%z90,mgcpgw%z90,analog%z90,callgen%z100,extmodule%z100:callto=external/nodata/yatego-callflow-vars\n")
	m, err := e.GetEvent()
	if err != nil {
		t.Fatalf("GetEvent failed: %s", err)
	}
	t.Logf("GetEvent returned msg: %+v", m)
	if m.Type != TypeIncoming {
		t.Fatalf("GetEvent invalid msg type: %s", m.Type)
	}
	if m.Name != MsgCallExecute {
		t.Fatalf("GetEvent invalid msg name: %s", m.Name)
	}
	if m.ID != "0x7f4e5df469e0.672973892" {
		t.Fatalf("GetEvent invalid msg ID: %s", m.ID)
	}
	if m.Params["id"] != "sip/5" {
		t.Fatalf("GetEvent invalid ch ID: %s", m.Params["id"])
	}
	if m.Params["callto"] != "external/nodata/yatego-callflow-vars" {
		t.Fatalf("GetEvent invalid callto: %s", m.Params["callto"])
	}
}

func engine(in io.Reader, out io.Writer) *Engine {
	return &Engine{
		In:     in,
		Out:    out,
		Logger: logrus.New(),
	}
}
