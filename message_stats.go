package main

import (
	"log"
	"sync"
)

type Message struct {
	hash             string
	time             string
	peerId           string // TODO: publisher or subscriber
	payloadSizeBytes uint64
	acks             []string // peerid that acks the messages
}

type MessageStats struct {
	// msgHash to Message
	msgs       map[string]*Message
	mutex      sync.RWMutex
	totalPeers int
}

func NewMessage(hash string, time string, peerId string, payloadSizeBytes uint64) *Message {
	return &Message{
		acks: make([]string, 0),
	}
}

func NewMessageStats() *MessageStats {
	return &MessageStats{
		msgs:       make(map[string]*Message),
		mutex:      sync.RWMutex{},
		totalPeers: 0,
	}
}

func (m *MessageStats) SentMessage(msg *Message) {
	m.mutex.Lock()
	// Store the message and init acks to empty
	m.msgs[msg.hash] = msg
	m.mutex.Unlock()
}

func (m *MessageStats) ReceivedMessage(msg *Message) {
	m.mutex.Lock()
	msg, found := m.msgs[msg.hash]
	if found {
		log.Fatal("Duplicate msg, something is wrong? TODO")
	}
	msg.acks = append(msg.acks, msg.peerId)
	m.mutex.Unlock()
}

// TODO: Some cleanups so that mem doesnt grow forever
