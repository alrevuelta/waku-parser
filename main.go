package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
)

func main() {

	msgStats := NewMessageStats()
	_ = msgStats

	fmt.Println("Hello, world.")

	// if default doesnt work check "docker context list"
	cli, err := client.NewClientWithOpts(client.WithHost("unix:///Users/alrevuelta/.docker/run/docker.sock"))
	if err != nil {
		log.Fatal("could not create env client: ", err)
	}

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

	go storeSent(cli, msgStats)
	go storeReceived(cli, msgStats)

	for {
		//dirty way of keeping the main running
		time.Sleep(1 * time.Second)
	}

}

func storeReceived(cli *client.Client, msgStats *MessageStats) {
	container := "nwaku-simulator-nwaku-1"
	i, err := cli.ContainerLogs(context.Background(), container, types.ContainerLogsOptions{
		ShowStderr: true,
		ShowStdout: true,
		Timestamps: false,
		Follow:     true,
		Tail:       "40",
	})
	if err != nil {
		log.Fatal("could not create contianer logs instance: ", err)
	}
	hdr := make([]byte, 8)
	for {
		_, err := i.Read(hdr)
		if err != nil {
			log.Fatal(err)
		}
		var w io.Writer
		switch hdr[0] {
		case 1:
			w = os.Stdout
		default:
			w = os.Stderr
		}
		count := binary.BigEndian.Uint32(hdr[4:])
		dat := make([]byte, count)
		_, err = i.Read(dat)
		logLine := string(dat)
		//fmt.Fprint(w, logLine)
		_ = w
		//log.Info(logLine)

		// TODO: Write a proper regex for this this is crap
		// Example: time="2023-06-07 14:18:14" level=info msg="Published message" hash=0xb801283ad665f50d975524f7889e7c835db54c1bc5c80d09d15929698cfc6ce4 payloadSizeBytes=2000 peerId=16Uiu2HAmFuJ3MmRHbAQkuNWZBshr12CBk2ftxWPLG8sowpbfosEb pubsubTopic=/waku/2/default-waku/proto sentTime=1686147494072801380
		if strings.Contains(logLine, "waku.relay received") {

			// TODO: this is super wrong, seems like the logs contains colors so quite vulnerable to changes
			split := strings.Split(logLine[192:], " ") // very dirty wat of removing the first part of the log line (dont even know if its contatn lol)

			hash := split[2][14:]
			/*
				payloadSize, err := strconv.ParseUint(split[1][17:], 10, 64)
				if err != nil {
					panic(err)
				}*/
			peerId := split[0][11:]
			pubsubTopic := split[1][21:]
			rxTime := split[3][22:]
			log.Info("hash: ", hash)
			//log.Info("payloadSize: ", payloadSize)
			log.Info("peerId: ", peerId)
			log.Info("pubsubtopic: ", pubsubTopic)
			log.Info("sentime: ", rxTime)

			//msg := NewMessage(hash, sentTime, peerId, payloadSize)
			//msgStats.SentMessage(msg)

		}
	}
}

func storeSent(cli *client.Client, msgStats *MessageStats) {
	container := "nwaku-simulator-traffic-1"
	i, err := cli.ContainerLogs(context.Background(), container, types.ContainerLogsOptions{
		ShowStderr: true,
		ShowStdout: true,
		Timestamps: false,
		Follow:     true,
		Tail:       "40",
	})
	if err != nil {
		log.Fatal("could not create contianer logs instance: ", err)
	}
	hdr := make([]byte, 8)
	for {
		_, err := i.Read(hdr)
		if err != nil {
			log.Fatal(err)
		}
		var w io.Writer
		switch hdr[0] {
		case 1:
			w = os.Stdout
		default:
			w = os.Stderr
		}
		count := binary.BigEndian.Uint32(hdr[4:])
		dat := make([]byte, count)
		_, err = i.Read(dat)
		logLine := string(dat)
		//fmt.Fprint(w, logLine)
		_ = w
		//log.Info(logLine)

		// TODO: Write a proper regex for this this is crap
		// Example: time="2023-06-07 14:18:14" level=info msg="Published message" hash=0xb801283ad665f50d975524f7889e7c835db54c1bc5c80d09d15929698cfc6ce4 payloadSizeBytes=2000 peerId=16Uiu2HAmFuJ3MmRHbAQkuNWZBshr12CBk2ftxWPLG8sowpbfosEb pubsubTopic=/waku/2/default-waku/proto sentTime=1686147494072801380
		if strings.Contains(logLine, "Published message") {
			split := strings.Split(logLine[62:], " ") // very dirty wat of removing the time="2023-06-07 14:18:14" level=info msg="Published message" . TODO: fix this
			hash := split[0][5:]
			payloadSize, err := strconv.ParseUint(split[1][17:], 10, 64)
			if err != nil {
				panic(err)
			}
			peerId := split[2][7:]
			pubsubTopic := split[3][13:]
			sentTime := split[4][9:]
			log.Info("hash: ", hash)
			log.Info("payloadSize: ", payloadSize)
			log.Info("peerId: ", peerId)
			log.Info("pubsubtopic: ", pubsubTopic)
			log.Info("sentime: ", sentTime)

			msg := NewMessage(hash, sentTime, peerId, payloadSize)
			msgStats.SentMessage(msg)

		}
	}
}
