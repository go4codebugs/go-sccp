// Copyright 2019-2024 go-sccp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package sccp

import (
	"fmt"
	"io"

	"github.com/wmnsk/go-sccp/params"
)

// UDT represents a SCCP Message Unit Data(UDT).
type UDT struct {
	Type                MsgType
	ProtocolClass       *params.ProtocolClass
	Ptr1, Ptr2, Ptr3    uint8
	CalledPartyAddress  *params.PartyAddress
	CallingPartyAddress *params.PartyAddress
	DataLength          uint8
	Data                []byte
}

// NewUDT creates a new UDT.
func NewUDT(pcls int, retOnErr bool, cdpa, cgpa *params.PartyAddress, data []byte) *UDT {
	u := &UDT{
		Type: MsgTypeUDT,
		ProtocolClass: params.NewProtocolClass(
			pcls, retOnErr,
		),
		Ptr1:                3,
		CalledPartyAddress:  cdpa,
		CallingPartyAddress: cgpa,
		Data:                data,
	}
	u.Ptr2 = u.Ptr1 + uint8(cdpa.Length())
	u.Ptr3 = u.Ptr2 + uint8(cgpa.Length())
	u.SetLength()

	return u
}

// MarshalBinary returns the byte sequence generated from a UDT instance.
func (u *UDT) MarshalBinary() ([]byte, error) {
	b := make([]byte, u.MarshalLen())
	if err := u.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
// SCCP is dependent on the Pointers when serializing, which means that it might fail when invalid Pointers are set.
func (u *UDT) MarshalTo(b []byte) error {
	l := len(b)
	if l < 5 {
		return io.ErrUnexpectedEOF
	}

	b[0] = uint8(u.Type)
	b[1] = u.ProtocolClass.Value()
	b[2] = u.Ptr1
	if n := int(u.Ptr1); l < n {
		return io.ErrUnexpectedEOF
	}
	b[3] = u.Ptr2
	if n := int(u.Ptr2 + 3); l < n {
		return io.ErrUnexpectedEOF
	}
	b[4] = u.Ptr3
	if n := int(u.Ptr3 + 5); l < n {
		return io.ErrUnexpectedEOF
	}

	if err := u.CalledPartyAddress.MarshalTo(b[5:int(u.Ptr2+3)]); err != nil {
		return err
	}
	if err := u.CallingPartyAddress.MarshalTo(b[int(u.Ptr2+3):int(u.Ptr3+4)]); err != nil {
		return err
	}

	// succeed if the rest of buffer is longer than u.DataLength
	b[u.Ptr3+4] = u.DataLength
	if offset := int(u.Ptr3 + 5); len(b[offset:]) >= int(u.DataLength) {
		copy(b[offset:], u.Data)
		return nil
	}

	return io.ErrUnexpectedEOF
}

// ParseUDT decodes given byte sequence as a SCCP UDT.
func ParseUDT(b []byte) (*UDT, error) {
	u := &UDT{}
	if err := u.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	return u, nil
}

// UnmarshalBinary sets the values retrieved from byte sequence in a SCCP UDT.
func (u *UDT) UnmarshalBinary(b []byte) error {
	l := len(b)
	if l <= 5 { // where CdPA starts
		return io.ErrUnexpectedEOF
	}

	u.Type = MsgType(b[0])

	offset := 1
	pc := &params.ProtocolClass{}
	n, err := pc.Read(b[offset:])
	if err != nil {
		return err
	}
	u.ProtocolClass = pc
	offset += n

	u.Ptr1 = b[offset]
	if l < int(u.Ptr1) {
		return io.ErrUnexpectedEOF
	}
	u.Ptr2 = b[offset+1]
	if l < int(u.Ptr2+3) { // where CgPA starts
		return io.ErrUnexpectedEOF
	}
	u.Ptr3 = b[offset+2]
	if l < int(u.Ptr3+5) { // where u.Data starts
		return io.ErrUnexpectedEOF
	}

	offset += 3
	u.CalledPartyAddress, err = params.ParseCalledPartyAddress(b[offset:int(u.Ptr2+3)])
	if err != nil {
		return err
	}

	u.CallingPartyAddress, err = params.ParseCallingPartyAddress(b[int(u.Ptr2+3):int(u.Ptr3+4)])
	if err != nil {
		return err
	}

	// succeed if the rest of buffer is longer than u.DataLength
	u.DataLength = b[int(u.Ptr3+4)]
	if offset, dataLen := int(u.Ptr3+5), int(u.DataLength); l >= offset+dataLen {
		u.Data = b[offset : offset+dataLen]
		return nil
	}

	return io.ErrUnexpectedEOF
}

// MarshalLen returns the serial length.
func (u *UDT) MarshalLen() int {
	l := 6
	if param := u.CalledPartyAddress; param != nil {
		l += param.MarshalLen()
	}
	if param := u.CallingPartyAddress; param != nil {
		l += param.MarshalLen()
	}
	l += len(u.Data)

	return l
}

// SetLength sets the length in Length field.
func (u *UDT) SetLength() {
	u.DataLength = uint8(len(u.Data))
}

// String returns the UDT values in human readable format.
func (u *UDT) String() string {
	return fmt.Sprintf("%s: {ProtocolClass: %s, CalledPartyAddress: %v, CallingPartyAddress: %v, Data: %x}",
		u.Type,
		u.ProtocolClass,
		u.CalledPartyAddress,
		u.CallingPartyAddress,
		u.Data,
	)
}

// MessageType returns the Message Type in int.
func (u *UDT) MessageType() MsgType {
	return MsgTypeUDT
}

// MessageTypeName returns the Message Type in string.
func (u *UDT) MessageTypeName() string {
	return u.MessageType().String()
}

// CdGT returns the GT in CalledPartyAddress in human readable string.
func (u *UDT) CdGT() string {
	if u.CalledPartyAddress.GlobalTitle == nil {
		return ""
	}
	return u.CalledPartyAddress.Address()
}

// CgGT returns the GT in CalledPartyAddress in human readable string.
func (u *UDT) CgGT() string {
	if u.CallingPartyAddress.GlobalTitle == nil {
		return ""
	}
	return u.CallingPartyAddress.Address()
}
