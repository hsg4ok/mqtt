package mqtt

import (
//    "fmt"
    "errors"
    "regexp"
    "os"
    "syscall"
)

type RetainedMessages struct {
    _messageCache map[Topic]*RetainedMessage
    path string
}

func (rm *RetainedMessages) loadAll(path string) error {
    // #SLASH#
}

func (rm *RetainedMessages) load(path string) error {
}

func (rm RetainedMessages) save(msg RetainedMessage) error {
}

func (rm *RetainedMessages) messageCache() map[Topic]*RetainedMessage {
    if rm._messageCache == nil {
        rm._messageCache = make(map[Topic]*RetainedMessage)
    }
    return rm._messageCache
}

func (rm *RetainedMessages) add(msg RetainedMessage) error {
    if rm == nil {
        return errors.New("retainedmessages nil")
    }
    t := msg.PublishPacket.Topic
    if t == "" {
        rm.remove(t)
        return nil
    }
    rm.messageCache()[t] = &msg
    return msg.save(rm.path)
}

func (rm *RetainedMessages) remove(t Topic) {
    if rm._messageCache != nil {
        delete(rm._messageCache, t)
    }
    // look for it on disk
    syscall.Unlink(topicToFilename(rm.path, t))
}

func (rm *RetainedMessages) loadRetainedMessageForTopic(t Topic) *RetainedMessage {
    msg := loadRetainedMessage(rm.path, t)
    if msg != nil {
        rm.messageCache()[t] = msg
    }
    return msg
}

func (rm *RetainedMessages) get(t Topic) *RetainedMessage {
    if rm == nil {
        return nil
    }
    if rm.messageCache != nil {
        msg := rm.messageCache[t]
        if msg != nil {
            return msg
        }
    }
    return rm.loadRetainedMessageForTopic(t)
}
