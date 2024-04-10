package model

import (
	"encoding/json"
	"net"
	"time"

	"github.com/WangYihang/dns-grab/pkg/util"
	"github.com/miekg/dns"
)

var DNSClient *dns.Client

func init() {
	DNSClient = new(dns.Client)
	DNSClient.Net = "udp"
	DNSClient.Timeout = 8 * time.Second
}

type DNS struct {
	Request  *ReadableMsg `json:"request"`
	Response *ReadableMsg `json:"response"`
}

type ReadableMsg struct {
	Header    ReadableMsgHdr     `json:"header"`
	Questions []ReadableQuestion `json:"questions,omitempty"`
	Answers   []ReadableRR       `json:"answers,omitempty"`
	Ns        []ReadableRR       `json:"authority,omitempty"`
	Extra     []ReadableRR       `json:"additional,omitempty"`
}

type ReadableMsgHdr struct {
	dns.MsgHdr
}

type ReadableQuestion struct {
	dns.Question
}

type ReadableRR struct {
	dns.RR
}

func (h ReadableMsgHdr) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID     uint16 `json:"id"`
		QR     bool   `json:"qr"` // QR (Query/Response)
		OpCode string `json:"op"` // OP (Operation Code)
		AA     bool   `json:"aa"` // AA (Authoritative Answer)
		TC     bool   `json:"tc"` // TC (Truncated)
		RD     bool   `json:"rd"` // RD (Recursion Desired)
		RA     bool   `json:"ra"` // RA (Recursion Available)
		Z      bool   `json:"z"`  // Z  (Reserved)
		AD     bool   `json:"ad"` // AD (Authenticated Data)
		CD     bool   `json:"cd"` // CD (Checking Disabled)
		RCode  string `json:"rc"` // RC (Response Code)
	}{
		ID:     h.Id,
		QR:     h.Response,
		OpCode: dns.OpcodeToString[h.Opcode],
		AA:     h.Authoritative,
		TC:     h.Truncated,
		RD:     h.RecursionDesired,
		RA:     h.RecursionAvailable,
		Z:      h.Zero,
		AD:     h.AuthenticatedData,
		CD:     h.CheckingDisabled,
		RCode:  dns.RcodeToString[h.Rcode],
	})
}

func (q ReadableQuestion) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Name   string `json:"name"`
		QType  string `json:"qtype"`
		QClass string `json:"qclass"`
	}{
		Name:   q.Name,
		QType:  dns.TypeToString[q.Qtype],
		QClass: dns.ClassToString[q.Qclass],
	})
}

func (q ReadableQuestion) MarshalBSON() ([]byte, error) {
	return json.Marshal(&struct {
		Name   string `json:"name_reverse"`
		QType  string `json:"qtype"`
		QClass string `json:"qclass"`
	}{
		Name:   util.ReverseString(q.Name),
		QType:  dns.TypeToString[q.Qtype],
		QClass: dns.ClassToString[q.Qclass],
	})
}

func (rr ReadableRR) MarshalJSON() ([]byte, error) {
	rrType := dns.TypeToString[rr.Header().Rrtype]
	switch t := rr.RR.(type) {
	case *dns.A:
		return json.Marshal(&struct {
			RType string `json:"rtype"`
			RName string `json:"rname"`
			A     net.IP `json:"a"`
		}{
			RType: rrType,
			RName: t.Hdr.Name,
			A:     t.A,
		})
	case *dns.AAAA:
		return json.Marshal(&struct {
			RType string `json:"rtype"`
			RName string `json:"rname"`
			AAAA  net.IP `json:"aaaa"`
		}{
			RType: rrType,
			RName: t.Hdr.Name,
			AAAA:  t.AAAA,
		})
	case *dns.CNAME:
		return json.Marshal(&struct {
			RType string `json:"rtype"`
			RName string `json:"rname"`
			CNAME string `json:"cname"`
		}{
			RType: rrType,
			RName: t.Hdr.Name,
			CNAME: t.Target,
		})
	case *dns.SOA:
		return json.Marshal(&struct {
			RType   string `json:"rtype"`
			RName   string `json:"rname"`
			NS      string `json:"ns"`
			Mbox    string `json:"mbox"`
			Serial  uint32 `json:"serial"`
			Refresh uint32 `json:"refresh"`
			Retry   uint32 `json:"retry"`
			Expire  uint32 `json:"expire"`
			MinTTL  uint32 `json:"min_ttl"`
		}{
			RType:   rrType,
			RName:   t.Hdr.Name,
			NS:      t.Ns,
			Mbox:    t.Mbox,
			Serial:  t.Serial,
			Refresh: t.Refresh,
			Retry:   t.Retry,
			Expire:  t.Expire,
			MinTTL:  t.Minttl,
		})
	case *dns.NS:
		return json.Marshal(&struct {
			RType string `json:"rtype"`
			RName string `json:"rname"`
			NS    string `json:"ns"`
		}{
			RType: rrType,
			RName: t.Hdr.Name,
			NS:    t.Ns,
		})
	case *dns.TXT:
		return json.Marshal(&struct {
			RType string   `json:"rtype"`
			RName string   `json:"rname"`
			TXT   []string `json:"txt"`
		}{
			RType: rrType,
			RName: t.Hdr.Name,
			TXT:   t.Txt,
		})
	case *dns.MX:
		return json.Marshal(&struct {
			Type     string `json:"type"`
			Name     string `json:"name"`
			Priority uint16 `json:"priority"`
			Mx       string `json:"mx"`
		}{
			Type:     rrType,
			Name:     t.Hdr.Name,
			Priority: t.Preference,
			Mx:       t.Mx,
		})
	case *dns.PTR:
		return json.Marshal(&struct {
			Type string `json:"type"`
			Name string `json:"name"`
			Ptr  string `json:"ptr"`
		}{
			Type: rrType,
			Name: t.Hdr.Name,
			Ptr:  t.Ptr,
		})
	case *dns.SRV:
		return json.Marshal(&struct {
			Type     string `json:"type"`
			Name     string `json:"name"`
			Priority uint16 `json:"priority"`
			Weight   uint16 `json:"weight"`
			Port     uint16 `json:"port"`
			Target   string `json:"target"`
		}{
			Type:     rrType,
			Name:     t.Hdr.Name,
			Priority: t.Priority,
			Weight:   t.Weight,
			Port:     t.Port,
			Target:   t.Target,
		})
	case *dns.DNSKEY:
		return json.Marshal(&struct {
			Type      string `json:"type"`
			Name      string `json:"name"`
			Flags     uint16 `json:"flags"`
			Protocol  uint8  `json:"protocol"`
			Algorithm uint8  `json:"algorithm"`
			PublicKey string `json:"public_key"`
		}{
			Type:      rrType,
			Name:      t.Hdr.Name,
			Flags:     t.Flags,
			Protocol:  t.Protocol,
			Algorithm: t.Algorithm,
			PublicKey: t.PublicKey,
		})
	case *dns.RRSIG:
		return json.Marshal(&struct {
			Type       string `json:"type"`
			Name       string `json:"name"`
			Algorithm  uint8  `json:"algorithm"`
			Labels     uint8  `json:"labels"`
			OrigTtl    uint32 `json:"orig_ttl"`
			Expiration uint32 `json:"expiration"`
			Inception  uint32 `json:"inception"`
			KeyTag     uint16 `json:"key_tag"`
			SignerName string `json:"signer_name"`
			Signature  string `json:"signature"`
		}{
			Type:       rrType,
			Name:       t.Hdr.Name,
			Algorithm:  t.Algorithm,
			Labels:     t.Labels,
			OrigTtl:    t.OrigTtl,
			Expiration: t.Expiration,
			Inception:  t.Inception,
			KeyTag:     t.KeyTag,
			SignerName: t.SignerName,
			Signature:  t.Signature,
		})
	case *dns.NAPTR:
		return json.Marshal(&struct {
			Type        string `json:"type"`
			Name        string `json:"name"`
			Order       uint16 `json:"order"`
			Preference  uint16 `json:"preference"`
			Flags       string `json:"flags"`
			Regexp      string `json:"regexp"`
			Replacement string `json:"replacement"`
		}{
			Type:        rrType,
			Name:        t.Hdr.Name,
			Order:       t.Order,
			Preference:  t.Preference,
			Flags:       t.Flags,
			Regexp:      t.Regexp,
			Replacement: t.Replacement,
		})
	default:
		return json.Marshal(&struct {
			RType string `json:"rtype"`
			RData string `json:"rdata"`
		}{
			RType: rrType,
			RData: rr.String(),
		})
	}
}

func NewReadableMsg(msg *dns.Msg) *ReadableMsg {
	readable := &ReadableMsg{
		Header: ReadableMsgHdr{msg.MsgHdr},
	}

	for _, q := range msg.Question {
		readable.Questions = append(readable.Questions, ReadableQuestion{q})
	}

	for _, a := range msg.Answer {
		readable.Answers = append(readable.Answers, ReadableRR{a})
	}

	for _, ns := range msg.Ns {
		readable.Ns = append(readable.Ns, ReadableRR{ns})
	}

	for _, extra := range msg.Extra {
		readable.Extra = append(readable.Extra, ReadableRR{extra})
	}

	return readable
}
