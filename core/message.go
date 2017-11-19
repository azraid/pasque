/****************************************************************************
*
*   message.go
*
*   Written by mylee (2016-03-30)
*   Owned by mylee
*
*
*   protocol
*   [headerlen]/[totallen]/Spn/Version/Command/[header]/[body]
*   common한 protocol들을 등록
***/

package core

import (
	"encoding/json"
	"fmt"
	"reflect"
)

const (
	msgTypeConnect  byte = 'C'
	msgTypeDie      byte = 'D'
	msgTypeAccept   byte = 'A'
	msgTypePing     byte = 'P'
	msgTypeRequest  byte = 'S'
	msgTypeResponse byte = 'R'
)

const MaxBufferLength = 1 + 4 + 1024 + 5 + 65535

type msgPack struct {
	msgType byte
	header  []byte
	body    []byte
	buffer  []byte
}

type ConnHeader struct {
	Eid       string
	Federated bool `json:",string,omitempty"`
}

type ConnBody struct {
	Spn           string   `json:",,omitempty"`
	FederatedKey  string   `json:",,omitempty"`
	FederatedApis []string `json:",,omitempty"`
}

type AccptHeader struct {
	ErrCode  uint32 `json:",string,omitempty"`
	ErrText  string `json:",,omitempty"`
	ErrIssue string `json:",,omitempty"`
}

type AccptBody struct {
}

type PingHeader struct {
	Eid string
}

type DieHeader struct {
	Eid string
}

type ConnectMsg struct {
	header ConnHeader
	body   ConnBody
}

type AcceptMsg struct {
	header AccptHeader
	body   AccptBody
}

type PingMsg struct {
	header PingHeader
}

type DieMsg struct {
	header DieHeader
}

type ReqHeader struct {
	Spn      string
	Api      string
	Key      string   `json:",,omitempty"`
	TxnNo    uint64   `json:",string,omitempty"`
	ExtTxnNo uint64   `json:",string,omitempty"`
	ToEid    string   `json:",,omitempty"`
	FromEids []string `json:",,omitempty"`
}

type RequestMsg struct {
	Header ReqHeader
	Body   json.RawMessage
}

type ResHeader struct {
	TxnNo    uint64   `json:",string,omitempty"`
	ExtTxnNo uint64   `json:",string,omitempty"`
	ToEids   []string `json:",,omitempty"`
	ErrCode  uint32   `json:",string,omitempty"`
	ErrText  string   `json:",,omitempty"`
	ErrIssue string   `json:",,omitempty"`
}

type ResponseMsg struct {
	Header ResHeader
	Body   json.RawMessage
}

func (header *ResHeader) SetError(neterr NetError) {
	header.ErrCode = neterr.Code
	header.ErrText = neterr.Text
	header.ErrIssue = neterr.Issue
}

func (header ResHeader) GetError() NetError {
	return NetError{Code: header.ErrCode, Text: header.ErrText, Issue: header.ErrIssue}
}

func (out *msgPack) MsgType() byte {
	return out.msgType
}

func (out *msgPack) Header() []byte {
	return out.header
}

func (out *msgPack) Body() []byte {
	return out.body
}

func (out *msgPack) Bytes() []byte {
	if len(out.buffer) == 0 {
		out.build()
	}

	return out.buffer
}

func (out *msgPack) build() error {
	switch out.msgType {
	case msgTypeConnect:
	case msgTypeAccept:
	case msgTypePing:
	case msgTypeRequest:
	case msgTypeResponse:
	default:
		return NetError{Code: NetErrorUnknownMsgType, Text: "unknown msg type", Issue: "Infra"}
	}

	out.buffer = []byte(fmt.Sprintf("/%c%05d", out.msgType, len(out.header)))
	out.buffer = append(out.buffer, out.header...)

	if out.msgType != msgTypePing {
		out.buffer = append(out.buffer, []byte(fmt.Sprintf("%010d", len(out.body)))...)
		if len(out.body) > 0 {
			out.buffer = append(out.buffer, out.body...)
		}
	}

	if len(out.buffer) > MaxBufferLength {
		return NetError{Code: NetErrorTooLargeSize, Text: "too large size", Issue: "Infra"}
	}

	return nil
}

func (out *msgPack) Rebuild(header interface{}) (err error) {
	var msgType byte

	switch header.(type) {
	case ConnHeader:
		msgType = msgTypeConnect

	case AccptHeader:
		msgType = msgTypeAccept

	case PingHeader:
		msgType = msgTypePing

	case ReqHeader:
		msgType = msgTypeRequest

	case ResHeader:
		msgType = msgTypeResponse

	default:
		return NetError{Code: NetErrorUnknownMsgType, Text: "unknown msg type", Issue: "Infra"}
	}

	if out.msgType == 0 {
		out.msgType = msgType
	}

	if out.msgType != msgType {
		return NetError{Code: NetErrorInternal, Text: "msg type is different from original", Issue: "Infra"}
	}

	out.header, err = json.Marshal(header)
	if err != nil {
		return NetError{Code: NetErrorInternal, Text: err.Error(), Issue: "Infra"}
	}

	return out.build()
}

func NewMsgPack(msgType byte, header []byte, body []byte) MsgPack {
	return &msgPack{msgType: msgType, header: header, body: body}
}

func BuildMsgPack(header interface{}, body interface{}) (MsgPack, error) {
	var out msgPack

	switch header.(type) {
	case ConnHeader:
		out.msgType = msgTypeConnect

	case AccptHeader:
		out.msgType = msgTypeAccept

	case PingHeader:
		out.msgType = msgTypePing

	case DieHeader:
		out.msgType = msgTypeDie

	case ReqHeader:
		out.msgType = msgTypeRequest

	case ResHeader:
		out.msgType = msgTypeResponse

	default:
		return nil, NetError{Code: NetErrorUnknownMsgType, Text: "unknown msg type", Issue: "Infra"}
	}

	var err error
	out.header, err = json.Marshal(header)
	if err != nil {
		return nil, NetError{Code: NetErrorInternal, Text: err.Error(), Issue: "Infra"}
	}

	if body == nil {
		out.body = []byte("{}")
	} else {
		out.body, err = json.Marshal(body)
		if err != nil {
			return nil, NetError{Code: NetErrorInternal, Text: err.Error(), Issue: "Infra"}
		}
	}

	if err := out.build(); err != nil {
		return nil, err
	}

	return &out, err
}

func ParseConnectMsg(header []byte, body []byte) *ConnectMsg {
	var msg ConnectMsg

	if err := json.Unmarshal(header, &msg.header); err != nil {
		return nil
	}

	if err := json.Unmarshal(body, &msg.body); err != nil {
		return nil
	}

	return &msg
}

func BuildConnectMsgPack(eid string, toplgy Topology) MsgPack {
	federated := false
	if len(toplgy.FederatedKey) > 0 {
		federated = true
	}

	mp, _ := BuildMsgPack(ConnHeader{Eid: eid, Federated: federated}, ConnBody{Spn: toplgy.Spn, FederatedKey: toplgy.FederatedKey, FederatedApis: toplgy.FederatedApis})

	return mp
}

func ParseAcceptMsg(header []byte, body []byte) *AcceptMsg {
	var msg AcceptMsg

	if err := json.Unmarshal(header, &msg.header); err != nil {
		return nil
	}

	if err := json.Unmarshal(body, &msg.body); err != nil {
		return nil
	}

	return &msg
}

func BuildAcceptMsgPack(ne NetError) MsgPack {
	mp, _ := BuildMsgPack(AccptHeader{ErrCode: ne.Code, ErrText: ne.Text, ErrIssue: ne.Issue}, nil)
	return mp
}

func BuildPingMsgPack(eid string) MsgPack {
	mp, _ := BuildMsgPack(PingHeader{Eid: eid}, nil)
	return mp
}

func BuildDieMsgPack(eid string) MsgPack {
	mp, _ := BuildMsgPack(DieHeader{Eid: eid}, nil)
	return mp
}

func ParseReqHeader(b []byte) *ReqHeader {
	var header ReqHeader
	if err := json.Unmarshal(b, &header); err != nil {
		return nil
	}

	return &header
}

func ParseResHeader(b []byte) *ResHeader {
	var header ResHeader
	if err := json.Unmarshal(b, &header); err != nil {
		return nil
	}

	return &header
}

func PeekFromEids(eids []string) string {
	if len(eids) == 0 {
		return ""
	}

	return eids[len(eids)-1]
}

func PopFromEids(eids []string) (string, []string, error) {
	l := len(eids)
	if l == 0 {
		return "", nil, fmt.Errorf("nothing at eids")
	}

	eid := eids[l-1]
	eids = eids[:l-1]

	return eid, eids, nil
}

func PushToEids(eid string, eids []string) []string {
	return append(eids, eid)
}

func IsValidMsg(rv reflect.Value) error {
	switch rv.Kind() {
	case reflect.Ptr:
		fallthrough
	case reflect.Interface:
		if !rv.Elem().IsValid() {
			return fmt.Errorf("nil")
		}

		return IsValidMsg(rv.Elem())

	case reflect.Struct:
		var errText string
		for i := 0; i < rv.NumField(); i += 1 {
			if len(rv.Type().Field(i).Tag.Get("required")) > 0 {
				if err := IsValidMsg(rv.Field(i)); err != nil {
					if len(errText) > 0 {
						errText += ","
					}

					errText += rv.Type().Field(i).Name
				}
			}
		}

		if len(errText) > 0 {
			return fmt.Errorf("%s", errText)
		} else {
			return nil
		}

	case reflect.String:
		fallthrough
	case reflect.Map:
		fallthrough
	case reflect.Slice:
		if rv.Len() == 0 {
			return fmt.Errorf("nil")
		}

	default:
		if rv.Interface() == 0 || rv.Interface() == nil {
			return fmt.Errorf("nil")
		}
	}

	return nil
}

func UnmarshalMsg(raw []byte, body interface{}) error {
	if err := json.Unmarshal(raw, body); err != nil {
		return err
	}

	if err := IsValidMsg(reflect.ValueOf(body)); err != nil {
		return fmt.Errorf("invalid param, %s", err.Error())
	}

	return nil
}
