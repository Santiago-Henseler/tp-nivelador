package client

import (
	"net"
	"time"
	"os"
	"strconv"
	"bufio"
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/logger"
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/message"
)

const CONNECTION_ATTEMPTS_MAX = 3
const CONNECTION_ATTEMPS_DELAY_MS = 200

type ClientConfig struct {
	ServerHost string
	ServerPort string
	AgencyId   string
	Batch int
	InputFile string
	OutputFile string
}

type Client struct {
	conn   net.Conn
	config ClientConfig
}

func NewClient(config ClientConfig) (*Client, error) {
	conn, err := connectToServer(config.ServerHost, config.ServerPort)
	if err != nil {
		logger.Warn("connect-to-server", logger.Fail)
		return nil, err
	}

	client := &Client{conn: conn, config: config}
	return client, nil
}

func connectToServer(host, port string) (net.Conn, error) {
	const action = "connect-to-server"
	var err error
	var conn net.Conn

	logger.Info(action, logger.InProgress)
	for i := range CONNECTION_ATTEMPTS_MAX {
		conn, err = net.Dial("tcp", host+":"+port)
		if err != nil {
			logger.Warn(action, logger.Fail, "attempt", i)
			time.Sleep(CONNECTION_ATTEMPS_DELAY_MS * time.Millisecond)
			continue
		}

		logger.Info(action, logger.Success)
		break
	}

	return conn, err
}

func (client *Client) Run() int {
	archivo, err := os.Open(client.config.InputFile)
	if err != nil {
		logger.Error("client-open-file", logger.Fail, "err", err)
		return 1
	}

	defer client.conn.Close()
	defer archivo.Close()

	betMesages := []message.BetMessage{}
	scanner := bufio.NewScanner(archivo)
	for scanner.Scan() {
		
		betMesage, err := message.CreateMessageBet(client.config.AgencyId, scanner.Text());
		if err != nil {
			return 1
		}
		betMesages = append(betMesages, betMesage)

		if len(betMesages) == int(client.config.Batch){
			message.SendBatchMessage(client.conn, betMesages);
			betMesages = []message.BetMessage{}
		}
	}

	if len(betMesages) != 0{
		message.SendBatchMessage(client.conn, betMesages);
	}
	message.EndBetMessages(client.conn)
	
	file, err := os.OpenFile(client.config.OutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Error("client-open-file", logger.Fail, "err", err)
		return 1
	}

	defer file.Close()

	reciving := true
	for reciving {
		betMessage, err := message.ReciveMessage(client.conn)

		if err != nil || len(betMessage) == 0{
			reciving = false
			continue
		}

		for _, bet := range betMessage {
			_, err = file.WriteString( bet.First_name + ","+bet.Last_name + "," + strconv.Itoa(bet.Document) + "," + bet.Birthdate + "," +strconv.Itoa(bet.Number) + "\n")
			if err != nil {
				logger.Error("client-write-file", logger.Fail, "err", err)
				return 1
			}
		}
	}

	logger.Info("client-send-file", logger.Success, "agency-id", client.config.AgencyId)
	
	return 0
}
