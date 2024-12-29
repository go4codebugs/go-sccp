// Copyright 2019-2024 go-sccp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

/*
Package sccp provides encoding/decoding feature of Signalling Connection Control Part used in SS7/SIGTRAN protocol stack.

This is still an experimental project, and currently in its very early stage of development. Any part of implementations
(including exported APIs) may be changed before released as v1.0.0.
*/
package sccp

import (
	"encoding"
	"fmt"
)

// MsgType is type of SCCP message.
type MsgType uint8

// Message Type definitions.
const (
	_            MsgType = iota
	MsgTypeCR            // CR
	MsgTypeCC            // CC
	MsgTypeCREF          // CREF
	MsgTypeRLSD          // RLSD
	MsgTypeRLC           // RLC
	MsgTypeDT1           // DT1
	MsgTypeDT2           // DT2
	MsgTypeAK            // AK
	MsgTypeUDT           // UDT
	MsgTypeUDTS          // UDTS
	MsgTypeED            // ED
	MsgTypeEA            // EA
	MsgTypeRSR           // RSR
	MsgTypeRSC           // RSC
	MsgTypeERR           // ERR
	MsgTypeIT            // IT
	MsgTypeXUDT          // XUDT
	MsgTypeXUDTS         // XUDTS
	MsgTypeLUDT          // LUDT
	MsgTypeLUDTS         // LUDTS
)

// Message is an interface that defines SCCP messages.
type Message interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
	MarshalTo([]byte) error
	MarshalLen() int
	MessageType() MsgType
	MessageTypeName() string
	fmt.Stringer
}

// ParseMessage decodes the byte sequence into Message by Message Type.
// Currently this only supports UDT type of message only.
func ParseMessage(b []byte) (Message, error) {
	var m Message
	switch MsgType(b[0]) {
	/* TODO: implement!
	case CR:
	case CC:
	case CREF:
	case RLSD:
	case RLC:
	case DT1:
	case DT2:
	case AK:
	*/
	case MsgTypeUDT:
		m = &UDT{}
	/* TODO: implement!
	case UDTS:
	case ED:
	case EA:
	case RSR:
	case RSC:
	case ERR:
	case IT:
	case XUDT:
	case XUDTS:
	case LUDT:
	case LUDTS:
	*/
	default:
		return nil, UnsupportedTypeError(b[0])
	}

	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}
