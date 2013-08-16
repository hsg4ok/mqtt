package mqtt

import (
    "io"
    "bytes"
    "errors"
)

func bool2byte(b bool) byte {
    if b {
        return 1
    }
    return 0
}

func byte2bool(x byte) bool {
    if (x & 1) == 1 {
        return true
    }
    return false
}

func readRemainingBytes(r *bytes.Reader) (buf []byte, err error) {
    buf, err = readBytes(r.Len(), r)
    return
}

func writeBytes(b []byte, w io.Writer) (err error) {
    n, err := w.Write(b)
    if err != nil {
        return
    }
    if len(b) != n {
        err = errors.New("Unable to write bytes")
        return
    }
    return
}

func writeByte(b byte, w io.Writer) error {
    return writeBytes([]byte{b}, w)
}

// Big-Endian style ... MSB, LSB
func writeUint16(x uint16, w io.Writer) error {
    return writeBytes([]byte{ byte(x >> 8), byte(x & 0xff) }, w)
}

func readBytes(n int, r io.Reader) (result []byte, err error) {
    result = make([]byte, n)
    actual, err := io.ReadFull(r, result)
    if err != nil {
        return
    }
    if actual < n {
        err = errors.New("Expected to read more bytes")
        return
    }
    return
}

func readByte(r io.Reader) (result byte, err error) {
    bytes_buf, err := readBytes(1, r)
    if err != nil {
        return
    }
    result = byte(bytes_buf[0])
    return
}

// Big-Endian style ... MSB, LSB
func readUint16(r io.Reader) (result uint16, err error) {
    bytes_buf, err := readBytes(2, r)
    if err != nil {
        return
    }
    result = uint16(bytes_buf[0]) << 8 | uint16(bytes_buf[1])
    return
}
