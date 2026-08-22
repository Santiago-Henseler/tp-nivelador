package client

import (
	"net"
	"time"
	"os"
	"bufio"
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/logger"
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/message"
)

const CONNECTION_ATTEMPTS_MAX = 3
const CONNECTION_ATTEMPS_DELAY_MS = 200

const FILE_NAME = "input/input-0.csv"

type ClientConfig struct {
	ServerHost string
	ServerPort string
	AgencyId   string
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

func (client *Client) Run() error {
	archivo, err := os.Open(FILE_NAME)
	if err != nil {
		logger.Error("client-open-file", logger.Fail, "err", err)
		return nil
	}

	defer client.conn.Close()
	defer archivo.Close()

	scanner := bufio.NewScanner(archivo)
	for scanner.Scan() {
	
		betMesage, err := message.CreateMessageBet(client.config.AgencyId, scanner.Text());
		if err != nil {
			return err
		}

		message.SendBetMessage(client.conn, betMesage);
	}
	
	a := false
	for a {
		betMesage, err := message.ReciveBetMessage(client.conn)

		if err != nil{
			a = false
			continue
		}
		logger.Info("client-champion", logger.Success, "agency-id", betMesage.First_name)
	}

	logger.Info("client-send-file", logger.Success, "agency-id", client.config.AgencyId)
	
	return nil
}
