package mqtt

import (
    "io"
)

type WildCardTopic String

func (t WildCardTopic) write(w io.Writer) error {
    return (String(t)).write(w)
}

func readWildCardTopic(r io.Reader) (t WildCardTopic, err error) {
    var s String
    s, err = readString(r)
    if err != nil {
        return
    }
    t = WildCardTopic(s)
    return
}
