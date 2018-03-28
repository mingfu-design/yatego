package yatego

import (
	"bytes"
	"regexp"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	t.Log("Running testrun on controller")
	in := new(bytes.Buffer)
	out := new(bytes.Buffer)
	e := engine(in, out)
	cm := &CallManager{
		calls: make(map[string]*Call),
	}

	c := &Controller{
		callManager:          cm,
		fallbackToController: false,
		singleChannelMode:    true,
		flowID:               "",
		staticComponents:     []Component{},
		callflowLoader:       nil,
		logger:               e.Logger,
		engine:               e,
	}
	//inline component
	com := NewBaseComponent("start", e, e.Logger, map[string]interface{}{})
	com.OnEnter(func(call *Call, message *Message) *CallbackResult {
		com.logger.Infoln("Inline base component entered")
		com.PlayTone("congestion", call, map[string]string{})
		return NewCallbackResult(ResStay, "")
	})
	c.AddStaticComponent(com)

	in.WriteString("%%>message:0x7f4e5df469e0.672973892:1522227506:call.execute::id=sip/5:module=sip:status=incoming:address=172.28.128.1%z34084:billid=1522152787-5:answered=false:direction=incoming:callid=sip/-altWCH3K8SSDe2HGBX0sA../0ca10a05/:caller=41587000201:called=923:antiloop=19:ip_host=172.28.128.1:ip_port=34084:ip_transport=UDP:connection_id=general:connection_reliable=false:sip_uri=sip%z923@172.28.128.3;transport=UDP:sip_from=sip%z41587000201@172.28.128.3;transport=UDP:sip_to=<sip%z923@172.28.128.3;transport=UDP>:sip_callid=-altWCH3K8SSDe2HGBX0sA..:device=Z 5.2.10 rv2.8.75-mod:sip_contact=<sip%z41587000201@172.28.128.1%z34084;transport=UDP>:sip_allow=INVITE, ACK, CANCEL, BYE, NOTIFY, REFER, MESSAGE, OPTIONS, INFO, SUBSCRIBE:sip_content-type=application/sdp:sip_user-agent=Z 5.2.10 rv2.8.75-mod:sip_allow-events=presence, kpml, talk:rtp_addr=172.28.128.1:media=yes:formats=mulaw,alaw,ilbc20:transport=RTP/AVP:rtp_mapping=mulaw=0,alaw=8,ilbc20=97:rtp_rfc2833=99:rtp_port=30000:sdp_fmtp%z=minptime=20; cbr=1; maxaveragebitrate=40000; useinbandfec=1:sdp_sendrecv=:sdp_fmtp%z=0-16:sdp_fmtp%z=0-16:sdp_fmtp%z=0-16:sdp_fmtp%z=0-16:rtp_forward=possible:handlers=javascript%z15,regexroute%z100,javascript%z15,cdrbuild%z50,fileinfo%z90,subscription%z100,sip%z100,regexroute%z100,javascript%z15,gvoice%z20,queues%z45,cdrbuild%z50,yrtp%z50,lateroute%z75,dbwave%z90,dumb%z90,sip%z90,wave%z90,filetransfer%z90,tone%z90,conf%z90,iax%z90,analyzer%z90,jingle%z90,sig%z90,mgcpgw%z90,analog%z90,callgen%z100,extmodule%z100:callto=external/nodata/yatego-callflow-vars\n")

	c.Run("")
	lines := strings.Split(out.String(), "\n")

	for i, line := range lines {
		t.Logf("output message %d: %s", i, line)
	}

	patterns := []string{
		"%%<message:.*:true:call.execute::.*targetid=yatego/.*",
		"%%>message:.*:.*:call.answered::.*:targetid=sip/.*:id=yatego/.*",
		"%%>message:.*:.*:chan.attach::.*source=tone/congestion:notify=yatego/.*:targetid=sip/.*:id=yatego/.*",
	}

	for i := 0; i < len(patterns); i++ {
		match, err := regexp.MatchString(patterns[i], lines[i])
		if err != nil {
			t.Fatalf("Dispatch output match failed: %s", err)
		}
		if !match {
			t.Fatalf("Dispatch output [%s] not matched by pattern [%s]", lines[i], patterns[i])
		}
	}

}
