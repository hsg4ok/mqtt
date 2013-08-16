package mqtt

import (
    "regexp"
    "os"
    "path"
)

type RetainedMessage struct {
    FixedHeader
    PublishPacket
}

func loadRetainedMesage(path string, t Topic) *RetainedMessage {

return nil
}

func (rm RetainedMessage) save(path string) error {
//    os.OpenFile
return nil
}

var leadingSlash *regexp.Regexp = regexp.MustCompile("^/")
const slashDir = "#SLASH#/"
var slash *regexp.Regexp = regexp.MustCompile("/")

func topicToFilename(Path string, t Topic) string {
    s := leadingSlash.ReplaceAllString(string(t), slashDir)
    s = slash.ReplaceAllString(s, string(os.PathSeparator))
    s = path.Join(Path, s)
    return s
}

func (rm RetainedMessage) filename(path string) string {
    return topicToFilename(path, rm.topic())
}


func (rm RetainedMessage) topic() Topic {
    return rm.PublishPacket.Topic
}

func (rm RetainedMessage) path() string {
return ""
}
