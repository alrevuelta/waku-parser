package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
	"time"

	"wakuparser/parser"

	"github.com/acarl005/stripansi"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
)

func main() {

	nwakuInstances := make([]string, 0)

	nwakuReeplicas := 5
	for i := 0; i < nwakuReeplicas; i++ {
		nwakuInstances = append(nwakuInstances, "nwaku-simulator-nwaku-"+strconv.Itoa(i+1))
	}

	msgStats := parser.NewMessageStats(nwakuInstances)

	// if default doesnt work check "docker context list"
	cli, err := client.NewClientWithOpts(client.WithHost("unix:///Users/alrevuelta/.docker/run/docker.sock"))
	if err != nil {
		log.Fatal("could not create env client: ", err)
	}

	// List and ensure the ones we are interested exist
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
		//Size    bool
		All: true,
		//Latest  bool
		//Since   string
		//Before  string
		//Limit   int
		//Filters filters.Args
	})

	//log.Info("containers: ", containers)
	_ = containers

	// start first sent messages.
	go storeSent(cli, msgStats)
	time.Sleep(1 * time.Second) // dirty way to avoid detecting mesg received before sent. may need some extra chekds
	go storeReceived(cli, msgStats)
	go runEvery(msgStats, 10)

	for {
		//dirty way of keeping the main running
		time.Sleep(1 * time.Second)
	}

}

func runEvery(msgStats *parser.MessageStats, tickerTimeInSeconds int64) {
	ticker := time.NewTicker(time.Duration(tickerTimeInSeconds) * time.Second)
	quit := make(chan struct{})
	for {
		select {
		case <-ticker.C:
			start := time.Now().UnixNano() / int64(time.Millisecond)

			// do stuff
			stats := msgStats.Stats()
			for container, stat := range stats {
				log.WithFields(log.Fields{
					"amountSent":        stat.MsgSent,
					"amountRx":          stat.MsgReceived,
					"avgDelayMicrosecs": stat.AvgDelay.Microseconds(),
					"container":         container,
				}).Info("Stats")
			}

			log.Info("total msg pending: ", msgStats.TotalMessages())

			end := time.Now().UnixNano() / int64(time.Millisecond)

			diff := end - start
			fmt.Println("Duration(ms):", diff)
			// do something if message rate is greater than what can be handled
			if diff > tickerTimeInSeconds*1000 {
				fmt.Println("Warning: took more than 1 second")
			}
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func storeReceived(cli *client.Client, msgStats *parser.MessageStats) {
	//container := "nwaku-simulator-nwaku-1"
	for _, container := range msgStats.Containers() {
		//if i == 2 { // TODO: simulate that a node losses
		//	continue
		//}
		go func(container string) {
			i, err := cli.ContainerLogs(context.Background(), container, types.ContainerLogsOptions{
				ShowStderr: true,
				ShowStdout: true,
				Timestamps: false,
				Follow:     true,
				Tail:       "0",
			})
			if err != nil {
				log.Fatal("could not create contianer logs instance: ", err)
			}
			hdr := make([]byte, 8)
			for {

				//continue
				// TODO:; perhaps add some delay to ensure we dont read the rx before we read the sent
				//time.Sleep(1 * time.Second)
				_, err := i.Read(hdr)
				if err != nil {
					log.Fatal(err)
				}

				count := binary.BigEndian.Uint32(hdr[4:])
				dat := make([]byte, count)
				_, err = i.Read(dat)
				if err != nil {
					log.Fatal(err)
				}
				logLine := string(dat)

				// TODO: Write a proper regex for this this is crap
				// Example: time="2023-06-07 14:18:14" level=info msg="Published message" hash=0xb801283ad665f50d975524f7889e7c835db54c1bc5c80d09d15929698cfc6ce4 payloadSizeBytes=2000 peerId=16Uiu2HAmFuJ3MmRHbAQkuNWZBshr12CBk2ftxWPLG8sowpbfosEb pubsubTopic=/waku/2/default-waku/proto sentTime=1686147494072801380
				logLine = stripansi.Strip(logLine) // remove colors
				if strings.Contains(logLine, "waku.relay received") {

					// simulate msg loss. TODO: remove
					//if rand.Intn(2) == 0 {
					//	//log.Info("simulating mesage loss: ", container, "line; ", logLine)
					//	continue // simulate loosing messages
					//}

					logLine = logLine[125:]

					// TODO: this is super wrong, seems like the logs contains colors so quite vulnerable to changes
					// very dirty wat of removing the first part of the log line (dont even know if its contatn lol)
					split := strings.Split(logLine, " ")

					hash := split[2][5:]
					rxTime := strings.TrimSpace(split[3][13:])

					rxTimeUint64, err := strconv.ParseUint(rxTime, 10, 64)
					if err != nil {
						log.Fatal("could not parse rxTime: ", rxTime, " ", err)
					}

					msg := parser.NewMessage(hash, rxTimeUint64, container)

					if !msgStats.WasMsgPublished(msg) {
						time.Sleep(2 * time.Second)
						log.Info("waiting")
					}

					/*
						log.WithFields(log.Fields{
							"hash":      hash,
							"rxTime":    rxTime,
							"container": container,
						}).Info("Detected received msg")*/

					msgStats.ReceivedMessage(msg)
				}
			}
		}(container)
	}
}

func storeSent(cli *client.Client, msgStats *parser.MessageStats) {
	container := "nwaku-simulator-traffic-1"
	i, err := cli.ContainerLogs(context.Background(), container, types.ContainerLogsOptions{
		ShowStderr: true,
		ShowStdout: true,
		Timestamps: false,
		Follow:     true,
		Tail:       "0",
	})
	if err != nil {
		log.Fatal("could not create contianer logs instance: ", err)
	}
	hdr := make([]byte, 8)
	for {
		// TODO: hoist this out to a common function
		_, err := i.Read(hdr)
		if err != nil {
			log.Fatal(err)
		}
		count := binary.BigEndian.Uint32(hdr[4:])
		dat := make([]byte, count)
		_, err = i.Read(dat)
		if err != nil {
			log.Fatal(err)
		}
		logLine := string(dat)

		// TODO: Write a proper regex for this this is crap
		// Example: time="2023-06-07 14:18:14" level=info msg="Published message" hash=0xb801283ad665f50d975524f7889e7c835db54c1bc5c80d09d15929698cfc6ce4 payloadSizeBytes=2000 peerId=16Uiu2HAmFuJ3MmRHbAQkuNWZBshr12CBk2ftxWPLG8sowpbfosEb pubsubTopic=/waku/2/default-waku/proto sentTime=1686147494072801380
		logLine = stripansi.Strip(logLine) // remove colors
		if strings.Contains(logLine, "Published message") {
			// fucking dirty. unnecesary duplicated
			/*
				timestamp := strings.Split(logLine, " ")[0] + " " + strings.Split(logLine, " ")[1] // quick shit
				timestamp = timestamp[6:]
				timestamp = timestamp[:len(timestamp)-1]
				log.Info("timestamo: ", timestamp)*/

			split := strings.Split(logLine[62:], " ")
			hash := split[0][5:]
			sentTime := strings.TrimSpace(split[3][9:])

			/*
				log.WithFields(log.Fields{
					"hash":      hash,
					"sentTime":  sentTime,
					"container": container,
				}).Info("Detected published msg")*/

			sentTimeUint64, err := strconv.ParseUint(sentTime, 10, 64)
			if err != nil {
				log.Fatal("error parsing sentTime: ", sentTime, " ", err)
			}

			msg := parser.NewMessage(hash, sentTimeUint64, container)
			msgStats.SentMessage(msg)

		}
	}
}
