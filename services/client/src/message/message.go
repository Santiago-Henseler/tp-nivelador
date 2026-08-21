package message

import (
	"net"
	"strings"
	"strconv"
	"errors"
	"encoding/binary"
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/logger"
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/safe_socket"
)

type BetMessage struct{
	Agency_id int
	First_name string
	Last_name string
	Document int
	Birthdate string
	Number int
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
		logger.Error("send-message", logger.Fail, betMesage.Agency_id)
		return err
	}

	return nil
}