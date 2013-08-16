package mqtt

import (
    "io"
)

type KeepAlive uint16

func (ka KeepAlive) write(w io.Writer) error {
    return writeUint16(uint16(ka), w)
}

func readKeepAlive(r io.Reader) (ka KeepAlive, err error) {
    var tmp uint16
    tmp, err = readUint16(r)
    if err != nil {
        return
    }
    ka = KeepAlive(tmp)
    return
}
