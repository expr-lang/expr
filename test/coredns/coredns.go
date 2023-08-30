package coredns

import (
	"context"
	"errors"
	"net"
)

func DefaultEnv(ctx context.Context, state *Request) map[string]any {
	return map[string]any{
		"incidr": func(ipStr, cidrStr string) (bool, error) {
			ip := net.ParseIP(ipStr)
			if ip == nil {
				return false, errors.New("first argument is not an IP address")
			}
			_, cidr, err := net.ParseCIDR(cidrStr)
			if err != nil {
				return false, err
			}
			return cidr.Contains(ip), nil
		},
		"metadata": func(label string) string {
			return ""
		},
		"type":        state.Type,
		"name":        state.Name,
		"class":       state.Class,
		"proto":       state.Proto,
		"size":        state.Len,
		"client_ip":   state.IP,
		"port":        state.Port,
		"id":          func() int { return int(state.Req.Id) },
		"opcode":      func() int { return state.Req.Opcode },
		"do":          state.Do,
		"bufsize":     state.Size,
		"server_ip":   state.LocalIP,
		"server_port": state.LocalPort,
	}
}

type Request struct {
	Req  *Msg
	W    ResponseWriter
	Zone string
}

func (r *Request) NewWithQuestion(name string, typ uint16) Request {
	return Request{}
}

func (r *Request) IP() string {
	return ""
}

func (r *Request) LocalIP() string {
	return ""
}

func (r *Request) Port() string {
	return ""
}

func (r *Request) LocalPort() string {
	return ""
}

func (r *Request) RemoteAddr() string { return r.W.RemoteAddr().String() }

func (r *Request) LocalAddr() string { return r.W.LocalAddr().String() }

func (r *Request) Proto() string {
	return "udp"
}

func (r *Request) Family() int {
	return 2
}

func (r *Request) Do() bool {
	return true
}

func (r *Request) Len() int { return 0 }

func (r *Request) Size() int {
	return 0
}

func (r *Request) SizeAndDo(m *Msg) bool {
	return true
}

func (r *Request) Scrub(reply *Msg) *Msg {
	return reply
}

func (r *Request) Type() string {
	return ""
}

func (r *Request) QType() uint16 {
	return 0
}

func (r *Request) Name() string {
	return "."
}

func (r *Request) QName() string {
	return "."
}

func (r *Request) Class() string {
	return ""
}

func (r *Request) QClass() uint16 {
	return 0
}

func (r *Request) Clear() {
}

func (r *Request) Match(reply *Msg) bool {
	return true
}

type Msg struct {
	MsgHdr
	Compress bool `json:"-"`
	Question []Question
	Answer   []RR
	Ns       []RR
	Extra    []RR
}

type MsgHdr struct {
	Id                 uint16
	Response           bool
	Opcode             int
	Authoritative      bool
	Truncated          bool
	RecursionDesired   bool
	RecursionAvailable bool
	Zero               bool
	AuthenticatedData  bool
	CheckingDisabled   bool
	Rcode              int
}

type Question struct {
	Name   string `dns:"cdomain-name"`
	Qtype  uint16
	Qclass uint16
}

type RR interface {
	Header() *RR_Header
	String() string
}

type RR_Header struct {
	Name     string `dns:"cdomain-name"`
	Rrtype   uint16
	Class    uint16
	Ttl      uint32
	Rdlength uint16
}

type ResponseWriter interface {
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	WriteMsg(*Msg) error
	Write([]byte) (int, error)
	Close() error
	TsigStatus() error
	TsigTimersOnly(bool)
	Hijack()
}
