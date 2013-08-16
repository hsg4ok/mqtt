package mqtt

import (
    "time"
)

// http://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html

const ProtocolVersion = byte(3)
var ProtocolName = []byte{ 0x00, 0x06, byte('M'), byte('Q'), byte('I'), byte('s'), byte('d'), byte('p') }


const IndefiniteReadTimeout = 0 * time.Second

const MQTTScheme    = "mqtt"
const MQTTSSLScheme = "mqtt-ssl"
const MQTTPort      = 1883
const MQTTSSLPort   = 8883

// C  S
//  ->  Connect
//  <-  ConnAck
//  <>  Publish
//  <>  PubAck
//  <>  PubRec
//  <>  PubRel
//  <>  PubComp
//  ->  Subscribe
//  <-  SubAck
//  ->  Unsubscribe
//  <-  UnsubAck
//  ->  PingReq*
//  <-  PingResp*
//  ->  Disconnect*
//
//    * don't require bespoke decoding
