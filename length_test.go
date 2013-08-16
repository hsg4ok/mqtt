package mqtt

import (
    "bytes"
    "testing"
)

var test_lengths = []Length{0, 1, 126, 127, 128, 16382, 16383, 16384, 2097150, 2097151, 2097152, 268435454, 268435455}

var test_lengths_writed = [][]byte{
    []byte{ 0x00 }, // 0
    []byte{ 0x01 }, // 1
    []byte{ 0x7E }, // 126
    []byte{ 0x7F }, // 127
    []byte{ 0x80, 0x01 }, // 128
    []byte{ 0xFE, 0x7F }, // 16382
    []byte{ 0xFF, 0x7F }, // 16383
    []byte{ 0x80, 0x80, 0x01 }, // 16384
    []byte{ 0xFE, 0xFF, 0x7F }, // 2097150
    []byte{ 0xFF, 0xFF, 0x7F }, // 2097151
    []byte{ 0x80, 0x80, 0x80, 0x01 }, // 2097152
    []byte{ 0xFE, 0xFF, 0xFF, 0x7F }, // 268435454
    []byte{ 0xFF, 0xFF, 0xFF, 0x7F }, // 268435455
}

func TestLengthWrite(t *testing.T) {
    for i, l := range test_lengths {
        w := bytes.NewBuffer(nil)
        err := l.write(w)
        if err != nil {
            t.Error("encoding error")
        }
        expected := test_lengths_writed[i]
        actual := w.Bytes()
        if bytes.Compare(expected, actual) != 0 {
            t.Error("encoding wrong expected =", expected, " actual=", actual)
        }
    }
}

func BenchmarkLengthWrite(b *testing.B) {
    l := test_lengths[0]
    for i := 0; i < b.N; i++ {
        b.StopTimer()
        w := bytes.NewBuffer(nil)
        b.StartTimer()
        l.write(w)
    }
}

func TestReadLength(t *testing.T) {
    for i, writed := range test_lengths_writed {
        r := bytes.NewBuffer(writed)
        actual, err := readLength(r)
        if err != nil {
            t.Error("decoding error")
        }
        if actual != test_lengths[i] {
            t.Error("decoding wrong")
        }
    }
}

func BenchmarkReadLength(b *testing.B) {
    for i := 0; i < b.N; i++ {
        b.StopTimer()
        r := bytes.NewBuffer(test_lengths_writed[1])
        b.StartTimer()
        readLength(r)
    }
}
