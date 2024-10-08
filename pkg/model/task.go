package model

import (
	"net"
	"strings"
	"time"

	"github.com/miekg/dns"
)

type Task struct {
	QNAME    string `json:"qname"`
	QTYPE    string `json:"qtype"`
	Resolver string `json:"resolver"`
	Timeout  int64  `json:"timeout"`
	DNS      DNS    `json:"dns"`
}

func NewTask(options ...func(*Task)) *Task {
	svr := &Task{
		DNS: DNS{},
	}
	for _, o := range options {
		o(svr)
	}
	return svr
}

func WithQNAME(qname string) func(*Task) {
	return func(t *Task) {
		t.QNAME = qname
	}
}

func WithQTYPE(qtype string) func(*Task) {
	return func(t *Task) {
		t.QTYPE = strings.ToUpper(qtype)
	}
}

func WithResolver(resolver string) func(*Task) {
	return func(t *Task) {
		// check if port is appended, append if not
		if _, _, err := net.SplitHostPort(resolver); err != nil {
			resolver = net.JoinHostPort(resolver, "53")
		}
		t.Resolver = resolver
	}
}

func WithTimeout(timeout int64) func(*Task) {
	return func(t *Task) {
		t.Timeout = timeout
	}
}

func (t *Task) Do() error {
	query := new(dns.Msg)
	query.SetQuestion(dns.Fqdn(t.QNAME), dns.StringToType[t.QTYPE])
	t.DNS.Request = NewReadableMsg(query)
	DNSClient.Timeout = time.Second * time.Duration(t.Timeout)
	answer, _, err := DNSClient.Exchange(query, t.Resolver)
	if err != nil {
		return err
	}
	t.DNS.Response = NewReadableMsg(answer)
	return nil
}
