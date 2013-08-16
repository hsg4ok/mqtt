package mqtt

import (
    "bytes"
    "testing"
)

var test_strings = []string{"", "x", "Hi", "✌᧨ሧᇷ༗"}
var test_strings_writed = [][]byte {
    []byte{0, 0},
    []byte{0, 1,   120},
    []byte{0, 2,   72,  105},
    []byte{0, 15,  226, 156, 140,  225, 167, 168,  225, 136, 167,  225, 135, 183,  224, 188, 151},
}

func TestStringWrite(t *testing.T) {
    for i, x := range test_strings {
        w := bytes.NewBuffer(nil)
        err := String(x).write(w)
        if err != nil {
            t.Error("encoding error")
        }
        if bytes.Compare(w.Bytes(), test_strings_writed[i]) != 0 {
            t.Error("encoding failed")
        }
    }
}

func BenchmarkStringWrite(b *testing.B) {
    for i := 0; i < b.N; i++ {
        b.StopTimer()
        w := bytes.NewBuffer(nil)
        b.StartTimer()
        String("").write(w)
    }
}

func TestReadString(t *testing.T) {
    for i, x := range test_strings_writed {
        r := bytes.NewReader(x)
        s, err := readString(r)
        if err != nil {
            t.Error("decoding error")
        }
        if string(s) != string(test_strings[i]) {
            t.Error("decoding failed")
        }
    }
}

func BenchmarkReadString(b *testing.B) {
    for i := 0; i < b.N; i++ {
        b.StopTimer()
        r := bytes.NewReader(test_strings_writed[1])
        b.StartTimer()
        readString(r)
    }
}
