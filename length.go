package mqtt

import (
    "io"
    "errors"
)

// http://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html

type Length uint32

const (
    Length1Max = Length(127)
    Length2Max = Length(16383)
    Length3Max = Length(2097151)
    Length4Max = Length(268435455)
    LengthMax = Length4Max
)

func (length Length) write(w io.Writer) (err error) {
    if length > Length4Max {
        err = errors.New("packet length too long to write")
        return
    }

    if length == 0 {
        err = writeByte(0, w)
        return
    }

    for length > 0 {
        digit := byte(length % 128)
        length = length / 128
        if length > 0 {
            digit |= 0x80
        }
        err = writeByte(digit, w)
        if err != nil {
            return
        }
    }
    return
}

func readLength(r io.Reader) (length Length, err error) {
    multiplier := 1
    length = 0
    var digit byte
    for i := 0; i < 4; i++ {
        digit, err = readByte(r)
        if err != nil {
            return
        }
        length += Length(multiplier * (int(digit) & 0x7f))
        if (digit & 0x80) == 0 {
            return // last digit < 0x80, ok
        }
        multiplier <<= 7 // *= 128 == *= 2^7 == << 7
    }
    err = errors.New("Bad length")
    return
}
