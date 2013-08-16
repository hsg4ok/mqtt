package mqtt

import (
    "io"
    "net"
    "crypto/sha256"
    "crypto/rand"
    "encoding/hex"
    "errors"
    "time"
    "unicode/utf8"
)

const (
    ClientIDMinRunes = 1
    ClientIDMaxRunes = 23
)

type ClientID String

func (c ClientID) write(w io.Writer) error {
    if ! c.IsValid() {
        return errors.New("Invalid ClientID")
    }
    return String(c).write(w)
}

func readClientID(r io.Reader) (result ClientID, err error) {
    s, err := readString(r)
    if err != nil {
        return
    }
    result = ClientID(s)
    if ! result.IsValid() {
        err = errors.New("Invalid ClientID")
        return
    }
    return
}

func (c ClientID) IsValid() bool {
    l := utf8.RuneCountInString(string(c))
    return (l >= ClientIDMinRunes) && (l <= ClientIDMaxRunes)
}

func generateClientID() ClientID {
    hash := sha256.New()

    // mac address
    ifs, _ := net.Interfaces()
    if len(ifs) > 0 {
        hash.Write(ifs[0].HardwareAddr)
    }

    // time
    t := time.Now().Unix()
    hash.Write([]byte{byte(t>>56), byte(t>>48), byte(t>>40), byte(t>>32), byte(t>>24), byte(t>>16), byte(t>>8), byte(t>>0)})

    // 8 bytes of random data
    buf := make([]byte, 8)
    _, err := rand.Read(buf)
    if err == nil {
        hash.Write(buf)
    }

    return ClientID(hex.EncodeToString(hash.Sum(nil))[:ClientIDMaxRunes])
}
