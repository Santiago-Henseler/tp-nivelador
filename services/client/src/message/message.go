package message

import (
	"net"
	"strings"
	"strconv"
	"errors"
	"encoding/binary"
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/safe_socket"
	//"github.com/7574-sistemas-distribuidos/tp-nivelador/src/logger"
)

const END_MESSAGE byte = 0b00000000
const BET_MESSAGE byte = 0b00000001
const BATCH_MESSAGE byte = 0b00000010

type BetMessage struct{
	Agency_id int
	First_name string
	Last_name string
	Document int
	Birthdate string
	Number int
}

func ReciveMessage(conn net.Conn) ([]BetMessage, error){
	msgType, err := safe_socket.RecvAll(conn, 1)

	if err != nil {
		return []BetMessage{}, err
	}
	
	if msgType[0] == BET_MESSAGE{
		return ReciveBetMessage(conn)
	}else if msgType[0] == END_MESSAGE{
		return []BetMessage{}, nil
	}else if msgType[0] == BATCH_MESSAGE{
		return ReciveBatchMessage(conn)
	}else{
		return []BetMessage{}, nil
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

func ReciveBatchMessage(conn net.Conn) ([]BetMessage, error){
	info, err := safe_socket.RecvAll(conn, 8)

	if len(info) < 8 || err != nil{
		return []BetMessage{}, errors.New("mensaje inválido")
	}

	betsSize := int32(binary.BigEndian.Uint32(info[0:4]))
	betsLenght := int32(binary.BigEndian.Uint32(info[4:8]))

	betsBytes, err := safe_socket.RecvAll(conn, int(betsLenght))
	if err != nil || len(betsBytes) < int(betsLenght) {
		return []BetMessage{}, errors.New("mensaje inválido")
	}

	bets := []BetMessage{}
	size := 0
	for i := 0; i < int(betsSize); i++ {
		agencyId := int32(binary.BigEndian.Uint32(betsBytes[size:size+4]))
		document := int32(binary.BigEndian.Uint32(betsBytes[size+4:size+8]))
		number := int32(binary.BigEndian.Uint32(betsBytes[size+8:size+12]))

		first_name_len := int32(binary.BigEndian.Uint32(betsBytes[size+12:size+16]))
		last_name_len := int32(binary.BigEndian.Uint32(betsBytes[size+16:size+20]))
		birthdate_len :=  int32(binary.BigEndian.Uint32(betsBytes[size+20:size+24]))

		size += 24

		first_name := string(betsBytes[size:first_name_len])
		last_name := string(betsBytes[size+int(first_name_len):size+int(first_name_len+last_name_len)])
		birthdate := string(betsBytes[size+int(first_name_len+last_name_len):size+int(first_name_len+last_name_len+birthdate_len)])

		size += int(first_name_len+last_name_len+birthdate_len+birthdate_len)

		bets = append(bets, BetMessage{int(agencyId), first_name, last_name, int(document), birthdate, int(number)})
	}	
	
	return bets, nil
}

func SendBatchMessage(conn net.Conn, bets []BetMessage) error{
	var bytes []byte

	bytes = append(bytes, []byte{byte(BATCH_MESSAGE)}...)

	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(len(bets)))
	bytes = append(bytes, buf...)

	var bet_bytes []byte

	for _, bet := range bets {
		bet_bytes = append(bet_bytes, betToBytes(bet)...)
	}

	buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(len(bet_bytes)))
	bytes = append(bytes, buf...)

	bytes = append(bytes, bet_bytes...)
	
	if err := safe_socket.SendAll(conn, bytes); err != nil {
		return err
	}

	return nil
}

func ReciveBetMessage(conn  net.Conn) ([]BetMessage, error){
	integers, err := safe_socket.RecvAll(conn, 24)
	if len(integers) < 24 || err != nil{
		return []BetMessage{}, errors.New("mensaje inválido")
	}

	agencyId := int32(binary.BigEndian.Uint32(integers[0:4]))
	document := int32(binary.BigEndian.Uint32(integers[4:8]))
	number := int32(binary.BigEndian.Uint32(integers[8:12]))

	first_name_len := int32(binary.BigEndian.Uint32(integers[12:16]))
	last_name_len := int32(binary.BigEndian.Uint32(integers[16:20]))
	birthdate_len :=  int32(binary.BigEndian.Uint32(integers[20:24]))

	strings, err := safe_socket.RecvAll(conn, int(first_name_len + last_name_len + birthdate_len))
	if err != nil{
		return []BetMessage{}, errors.New("mensaje inválido")
	}

	first_name := string(strings[0:first_name_len])
	last_name := string(strings[first_name_len:first_name_len+last_name_len])
	birthdate := string(strings[first_name_len+last_name_len:first_name_len+last_name_len+birthdate_len])

	return []BetMessage{{int(agencyId), first_name, last_name, int(document), birthdate, int(number)}}, nil
}

func SendBetMessage(conn net.Conn, betMesage BetMessage) error{

	var bytes []byte

	bytes = append(bytes, []byte{byte(BET_MESSAGE)}...)
	bytes = append(bytes, []byte(betToBytes(betMesage))...)

	if err := safe_socket.SendAll(conn, bytes); err != nil {
		return err
	}

	return nil
}

func EndBetMessages(conn net.Conn){
	safe_socket.SendAll(conn, []byte{byte(END_MESSAGE)})
}

func betToBytes(betMesage BetMessage) []byte{
	var bytes []byte

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

	return bytes
}