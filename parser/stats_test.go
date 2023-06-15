package parser

import (
	"fmt"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func Test_TODO2(t *testing.T) {
	msgStats := NewMessageStats([]string{"nwaku-1", "nwaku-2", "nwaku-3"})
	// Some messages are sent
	for i := 0; i < 100; i++ {
		msg := NewMessage(
			fmt.Sprintf("hash_%d", i),
			uint64(1686213729438612736),
			"traffic-1")
		msgStats.SentMessage(msg)
	}

	rx1 := NewMessage("hash_1", 1686213729438612736+1000, "nwaku-1")
	rx2 := NewMessage("hash_1", 1686213729438612736+1000, "nwaku-2")
	//rx3 := NewMessage("hash_1", 1686213729438612736+1000, "nwaku-3")

	msgStats.ReceivedMessage(rx1)
	msgStats.ReceivedMessage(rx2)
	//msgStats.ReceivedMessage(rx3)

	stats := msgStats.Stats()

	for container, stat := range stats {
		log.WithFields(log.Fields{
			"amountSent":        stat.MsgSent,
			"amountRx":          stat.MsgReceived,
			"avgDelayMicrosecs": stat.AvgDelay.Microseconds(),
			"container":         container,
		}).Info("Stats")
	}

	require.Equal(t, uint64(1), uint64(1))
}

func Test_TODO(t *testing.T) {
	require.Equal(t, uint64(1), uint64(1))

	msgStats := NewMessageStats([]string{"nwaku-1", "nwaku-2", "nwaku-3"})

	// 16Uiu2HAm8qbxAjtbcCuhP8msnGZvdK8RgHpdduDahUiCofi46AXj

	// Some messages are sent
	for i := 0; i < 100; i++ {
		msg := NewMessage(
			fmt.Sprintf("hash_%d", i),
			uint64(1686213729438612736),
			"traffic-1")
		log.Info("snd: ", msg)
		msgStats.SentMessage(msg)
	}

	// Some messages are received
	for i := 0; i < 50; i++ {
		msg := NewMessage(
			fmt.Sprintf("hash_%d", i),
			1686213729438612736+1000, // same delay for all
			"nwaku-1")
		msgStats.ReceivedMessage(msg)
	}

	for _, msg := range msgStats.msgs {
		log.Info("msg: ", msg)
	}

	log.Info(msgStats.Stats())

}
