package parser

import (
	"fmt"
	"net"
	"strings"

	"github.com/triasbrata/golibs/pkg/eventDriven/internals/model"
	"github.com/triasbrata/golibs/pkg/eventDriven/internals/types"
	"github.com/vmihailenco/msgpack/v5"
)

func readUDPMessage(con *net.UDPConn, msgMaxLength int64) (*net.UDPAddr, error, []byte) {
	msg := make([]byte, msgMaxLength)
	_, remoteAddr, err := con.ReadFromUDP(msg)
	trimMsg := make([]byte, 0)
	for _, b := range msg {
		if b != 0 {
			trimMsg = append(trimMsg, b)
		}
	}
	return remoteAddr, err, trimMsg
}

func ListenNewMessage(in types.ReaderUDP, eInvoke types.EventInvoke) {
	for {
		if in.IsConClose() {
			return
		}
		remoteAddr, err, msg := readUDPMessage(in.GetCon(), in.GetMaxLengthMessage())
		if err != nil {
			if strings.Contains(err.Error(), "closed network connection") && in.IsConClose() {
				return
			}
			fmt.Printf("Some error  %v", err)
			return
		}
		payload := &model.StdDto{}
		err = msgpack.Unmarshal(msg, payload)
		if err != nil {
			fmt.Printf("Some error when parse %v", err)
		}
		eInvoke(payload, remoteAddr)
	}
}
