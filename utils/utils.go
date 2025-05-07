// Copyright 2019-2024 go-sccp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

// Package utils provides some utilities which might be useful specifically for GTP(or other telco protocols).
package utils

import (
	"encoding/hex"
)

// BCDEncode encodes a string into BCD-encoded bytes.
func BCDEncode(s string) ([]byte, error) {
	return StrToSwappedBytes(s, "f")
}

// MustBCDEncode is the same as BCDEncode but panics if any error occurs.
// Use this function only when you are sure that the input string is valid.
func MustBCDEncode(s string) []byte {
	b, err := BCDEncode(s)
	if err != nil {
		panic(err)
	}
	return b
}

// BCDDecode decodes BCD-encoded bytes into a string.
func BCDDecode(isOdd bool, b []byte) string {
	return SwappedBytesToStr(b, isOdd)
}

// StrToSwappedBytes returns swapped bits from a byte.
// It is used for some values where some values are represented in swapped format.
//
// The second parameter is the hex character(0-f) to fill the last digit when
// handling a odd number. "f" is used In most cases.
func StrToSwappedBytes(s, filler string) ([]byte, error) {
	var raw []byte
	var err error
	if len(s)%2 == 0 {
		raw, err = hex.DecodeString(s)
	} else {
		raw, err = hex.DecodeString(s + filler)
	}
	if err != nil {
		return nil, err
	}

	return swap(raw), nil
}

// SwappedBytesToStr decodes raw swapped bytes into string.
// It is used for some values where some values are represented in swapped format.
//
// The second parameter is to decide whether to cut the last digit or not.
func SwappedBytesToStr(raw []byte, cutLastDigit bool) string {
	s := hex.EncodeToString(swap(raw))
	if cutLastDigit {
		s = s[:len(s)-1]
	}

	return s
}

func swap(raw []byte) []byte {
	swapped := make([]byte, len(raw))
	for n := range raw {
		t := ((raw[n] >> 4) & 0xf) + ((raw[n] << 4) & 0xf0)
		swapped[n] = t
	}
	return swapped
}

// Uint24To32 converts 24bits-length []byte value into the uint32 with 8bits of zeros as prefix.
// This function is used for the fields with 3 octets.
func Uint24To32(b []byte) uint32 {
	if len(b) != 3 {
		return 0
	}
	return uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2])
}

// Uint32To24 converts the uint32 value into 24bits-length []byte. The values in 25-32 bit are cut off.
// This function is used for the fields with 3 octets.
func Uint32To24(n uint32) []byte {
	return []byte{uint8(n >> 16), uint8(n >> 8), uint8(n)}
}

// Uint40To64 converts 40bits-length []byte value into the uint64 with 8bits of zeros as prefix.
// This function is used for the fields with 3 octets.
func Uint40To64(b []byte) uint64 {
	if len(b) != 5 {
		return 0
	}
	return uint64(b[0])<<32 | uint64(b[1])<<24 | uint64(b[2])<<16 | uint64(b[3])<<8 | uint64(b[4])
}

// Uint64To40 converts the uint64 value into 40bits-length []byte. The values in 25-64 bit are cut off.
// This function is used for the fields with 3 octets.
func Uint64To40(n uint64) []byte {
	return []byte{uint8(n >> 32), uint8(n >> 24), uint8(n >> 16), uint8(n >> 8), uint8(n)}
}

// EncodePLMN encodes MCC and MNC as BCD-encoded bytes.
func EncodePLMN(mcc, mnc string) ([]byte, error) {
	c, err := StrToSwappedBytes(mcc, "f")
	if err != nil {
		return nil, err
	}
	n, err := StrToSwappedBytes(mnc, "f")
	if err != nil {
		return nil, err
	}

	// 2-digit
	b := make([]byte, 3)
	if len(mnc) == 2 {
		b = append(c, n...)
		return b, nil
	}

	// 3-digit
	b[0] = c[0]
	b[1] = (c[1] & 0x0f) | (n[1] << 4 & 0xf0)
	b[2] = n[0]

	return b, nil
}

// DecodePLMN decodes BCD-encoded bytes into MCC and MNC.
func DecodePLMN(b []byte) (mcc, mnc string, err error) {
	raw := hex.EncodeToString(b)
	mcc = string(raw[1]) + string(raw[0]) + string(raw[3])
	mnc = string(raw[5]) + string(raw[4])
	if string(raw[2]) != "f" {
		mnc += string(raw[2])
	}

	return
}
