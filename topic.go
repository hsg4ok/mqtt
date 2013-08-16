package mqtt

import (
    "io"
    "errors"
    "strings"
    "regexp"
)

const (
    Wildcards      = "#+"
    InvalidDouble  = "//"
    InvalidEnd     = "/"
    SingleWildcard = '+'
    MultiWildcard  = '#'
    TopicSep       = '/'
    TopicLengthMin = 1
    TopicLengthMax = 32767
)

// http://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html
//
type Topic String

func (t Topic) IsValid() bool {
    s := string(t)
    // check length
    l := len(s)
    if l < TopicLengthMin || l > TopicLengthMax {
        return false
    }
    // does not end in / (just / is also invalid)
    if s[l-1:] == InvalidEnd {
        return false
    }
    // does not contain //
    if strings.Index(s, InvalidDouble) != -1 {
        return false
    }
    // does not contain # or +
    if strings.IndexAny(s, Wildcards) != -1 {
        return false
    }
    return true
}

// + -> [^/]+
var SingleLevelRegexp = regexp.MustCompile("\\+")
// /#$ -> .*$
var MultiLevelRegexp = regexp.MustCompile("/#$")
const NeedsEscape = ".+*?=[]{}()$\\|"

func matchToRegexp(match_string string) regexp.Regexp {
    temp := SingleLevelRegexp.ReplaceAllString(match_string, "[^/]+")
    temp = MultiLevelRegexp.ReplaceAllString(temp, ".*$")
    runes := []rune{'^'}
    for _, r := range temp {
        if strings.ContainsRune(NeedsEscape, r) {
            runes = append(runes, '\\')
        }
        runes = append(runes, r)
    }
    return *regexp.MustCompile(string(runes))
}

func (t Topic) Match(match_string string) bool {
    re := matchToRegexp(match_string)
    return re.MatchString(string(t))
}

func FilterTopics(topics []Topic, match_string string) (result []Topic) {
    re := matchToRegexp(match_string)
    result = make([]Topic, 0, len(topics))
    for _, topic := range topics {
        if re.MatchString(string(topic)) {
            result = append(result, topic)
        }
    }
    return
}

func (t Topic) write(w io.Writer) (err error) {
    if !t.IsValid() {
        err = errors.New("Invalid topic")
        return
    }
    err = (String(t)).write(w)
    return
}

func readTopic(r io.Reader) (t Topic, err error) {
    var s String
    s, err = readString(r)
    if err != nil {
        return
    }
    t = Topic(s)
    if !t.IsValid() {
        err = errors.New("Invalid topic")
        return
    }
    return
}
