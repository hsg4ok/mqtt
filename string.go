package mqtt

import (
    "errors"
    "io"
)

// http://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html

type String string

func (s String) write(w io.Writer) (err error) {
    string_data := []byte(s)
    bytes := len(string_data)
    if bytes > 0xffff {
        err = errors.New("String too big to encode")
        return
    }
    err = writeUint16(uint16(bytes), w)
    if err != nil {
        return
    }
    return writeBytes(string_data, w)
}

func readString(r io.Reader) (result String, err error) {
    bytes, err := readUint16(r)

    string_data, err := readBytes(int(bytes), r)
    if err != nil {
        return
    }
    result = String(string_data)
    return
}
