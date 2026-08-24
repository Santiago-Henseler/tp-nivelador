package message

import (
	"net"
	"strings"
	"strconv"
	"errors"
	"encoding/binary"
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/safe_socket"
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/logger"
)

const END_MESSAGE byte = 0b00000000
const BET_MESSAGE byte = 0b00000001
const BACH_MESSAGE byte = 0b00000010

type BetMessage struct{
	Agency_id int
	First_name string
	Last_name string
	Document int
	Birthdate string
	Number int
}

func ReciveMessage(conn net.Conn) (BetMessage, error){
	msgType, err := safe_socket.RecvAll(conn, 1)

	if err != nil {
		return BetMessage{}, err
	}
	if msgType[0] == BET_MESSAGE{
		return ReciveBetMessage(conn)
	}else if msgType[0] == END_MESSAGE{
		return BetMessage{}, nil
	}else if msgType[0] == BACH_MESSAGE{
		return BetMessage{}, nil
	}else{
		return BetMessage{}, nil
	}
}

func CreateMessageBet(agencyIdStr string, message string) (BetMessage, error){

	info := strings.Split(message, ",")

	if len(info) != 5 {
		return BetMessage{}, errors.New("mensaje inválido")
	}

	agencyId, err := strconv.Atoi(agencyIdStr)
	if err != nil {
		return BetMessage{}, err
	}
	document, err := strconv.Atoi(info[2])
	if err != nil {
		return BetMessage{}, err
	}
	number, err := strconv.Atoi(info[4])
	if err != nil {
		return BetMessage{}, err
	}
	
	betMesage := BetMessage{agencyId, info[0], info[1], document, info[3], number}

	return betMesage, nil
}

func SendBetMessage(conn net.Conn, betMesage BetMessage) error{

	var bytes []byte

	bytes = append(bytes, []byte{byte(BET_MESSAGE)}...)

	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(betMesage.Agency_id))
	bytes = append(bytes, buf...)

	buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(betMesage.Document))
	bytes = append(bytes, buf...)

	buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(betMesage.Number))
	bytes = append(bytes, buf...)
	
	buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(len(betMesage.First_name)))
	bytes = append(bytes, buf...)

	buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(len(betMesage.Last_name)))
	bytes = append(bytes, buf...)

	buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(len(betMesage.Birthdate)))
	bytes = append(bytes, buf...)

	bytes = append(bytes, []byte(betMesage.First_name)...)
	bytes = append(bytes, []byte(betMesage.Last_name)...)
	bytes = append(bytes, []byte(betMesage.Birthdate)...)

	if err := safe_socket.SendAll(conn, bytes); err != nil {
		return err
	}

	return nil
}

func ReciveBetMessage(conn  net.Conn) (BetMessage, error){
	integers, err := safe_socket.RecvAll(conn, 24)
	if len(integers) < 24 || err != nil{
		return BetMessage{}, errors.New("mensaje inválido")
	}

	agencyId := int32(binary.BigEndian.Uint32(integers[0:4]))
	document := int32(binary.BigEndian.Uint32(integers[4:8]))
	number := int32(binary.BigEndian.Uint32(integers[8:12]))

	first_name_len := int32(binary.BigEndian.Uint32(integers[12:16]))
	last_name_len := int32(binary.BigEndian.Uint32(integers[16:20]))
	birthdate_len :=  int32(binary.BigEndian.Uint32(integers[20:24]))

	strings, err := safe_socket.RecvAll(conn, int(first_name_len + last_name_len + birthdate_len))
	if err != nil{
		return BetMessage{}, errors.New("mensaje inválido")
	}

	first_name := string(strings[0:first_name_len])
	last_name := string(strings[first_name_len:first_name_len+last_name_len])
	birthdate := string(strings[first_name_len+last_name_len+birthdate_len:])

	return BetMessage{int(agencyId), first_name, last_name, int(document), birthdate, int(number)}, nil
}

func EndBetMessages(conn net.Conn){
	safe_socket.SendAll(conn, []byte{byte(END_MESSAGE)})
}