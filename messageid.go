package mqtt

import (
    "io"
    "errors"
)

type MessageID uint16

const InvalidMessageID = MessageID(0)

func readMessageID(r io.Reader) (m MessageID, err error) {
    x, err := readUint16(r)
    m = MessageID(x)
    return
}

func (m MessageID) IsValid() bool {
    return m != InvalidMessageID
}

func (m MessageID) write(w io.Writer) error {
    if ! m.IsValid() {
        return errors.New("Cannot write invalid mesage id (0)")
    }
    return writeUint16(uint16(m), w)
}
