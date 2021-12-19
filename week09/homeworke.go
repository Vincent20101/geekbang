package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

const (
	_packetLen    = 4
	_headerLen    = 2
	_version      = 2
	_operaction   = 4
	_seq          = 4
	_rawHeaderLen = _packetLen + _headerLen + _version + _operaction + _seq

	_packetOffset  = _packetLen
	_headerOffset  = _packetOffset + _headerLen
	_versionOffset = _headerOffset + _version
	_operOffset    = _versionOffset + _operaction
	_seqOffset     = _operOffset + _seq
)

var (
	packet, header           int
	version, oper, seq, body string
	rawHeader, rawBody       []byte
)

func BytesToInt(bys []byte) int {
	var data int
	bytebuff := bytes.NewBuffer(bys)
	binary.Read(bytebuff, binary.BigEndian, &data)
	return int(data)
}

func parse(conn net.Conn) {
	defer conn.Close()

	for {
		//元头部信息
		rawHeader = make([]byte, _rawHeaderLen)
		n, err := conn.Read(rawHeader)
		if err != nil {
			fmt.Printf("read rawHeader error:", err)
			return
		}

		if n < _rawHeaderLen {
			fmt.Println("not goim parse")
			return
		}

		//头部内容
		packet = BytesToInt(rawHeader[:_packetOffset])
		header = BytesToInt(rawHeader[_packetOffset:_headerOffset])
		version = string(rawHeader[_headerLen:_versionOffset])
		oper = string(rawHeader[_versionOffset:_operOffset])
		seq = string(rawHeader[_seqOffset:_seqOffset])

		//body
		bodyLen := packet - header
		rawBody = make([]byte, bodyLen)
		n, err = conn.Read(rawBody)
		if err != nil {
			fmt.Printf("read rawBody error:", err)
			return
		}

		if n < bodyLen {
			fmt.Println("not goim parse")
			return
		}

		fmt.Println("goim parse success")
	}
}

func main() {
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("start listen tcp with port 8080 failed :%#v", err)
		return
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("accept error :%#v", err)
		}
		go parse(conn)
	}
}
