package mqtt

import (
    "testing"
)

var valid_topics = []Topic{ "/finance", "finance", "finance/stock/ibm", "finance/stock/ibm/closingprice" }

var invalid_topics = []Topic{ "", "/", "#", "+", "//", "+/+", "/+", "///", "bad/" }


func TestTopicIsValid(t *testing.T) {
    for _, x := range valid_topics {
        if ! x.IsValid() {
            t.Error("topic should be valid", x)
        }
    }
}

func TestTopicIsValidInvalid(t *testing.T) {
    for _, x := range invalid_topics {
        if x.IsValid() {
            t.Error("topic should be invalid", x)
        }
    }
}

func BenchmarkIsValid(b *testing.B) {
    t := Topic("a")
    for i := 0; i < b.N; i++ {
        t.IsValid()
    }
}
