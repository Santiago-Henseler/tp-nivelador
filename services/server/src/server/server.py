import socket
import logger
import message
import threading
import time
from lottery.lottery import Lottery

OUTPUT_FILE = "output.txt"

class Server:
    def __init__(self, server_host: str, server_port: int, quorum, kill) -> None:
        self.server_host = server_host
        self.server_port = server_port
        self.quorum = quorum
        self.bets = 0
        self.condVar = threading.Condition()
        self.kill = kill
        
    def _handle_client(self, client_socket, lottery):
        message_amount = 0
        client_bets = []
        try:
            logger.info("handle-client", logger.LogResult.in_progress)
            comunication = True
            while comunication:
                if self.kill.is_set():
                    client_socket.close()
                    return

                client_bet = message.recive_message(client_socket)

                if client_bet == None:
                    logger.info( "handle-client", logger.LogResult.success, "messages-amount", message_amount)
                    comunication = False
                    continue
                
                client_bets.extend(client_bet)
                message_amount += 1
        except Exception as e:
            logger.error( "handle-client", logger.LogResult.fail, "messages-amount", message_amount)
            raise e

        with self.condVar:
            self.bets += 1
            if self.bets == self.quorum:
                self.condVar.notify_all()

            while self.bets < self.quorum:
                self.condVar.wait()

        lottery.store_bets(client_bets)
        for bet in client_bets:
            if self.kill.is_set():
                client_socket.close()
                return
            
            if lottery.has_won(bet):
                message.send_bet_message(client_socket,  bet)

        message.end_bet_message(client_socket)
        client_socket.close()

    def run(self):
        with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as server_socket:
            server_socket.bind((self.server_host, self.server_port))
            server_socket.listen()
            lottery = Lottery(OUTPUT_FILE)
            while True:
                if self.kill.is_set():
                    server_socket.close()
                    return
                try:
                    logger.info("accept-connection", logger.LogResult.in_progress)
                    client_socket, _ = server_socket.accept()
                except Exception as e:
                    logger.error("accept-connection", logger.LogResult.fail)
                    raise e
                logger.info("accept-connection", logger.LogResult.success)
                threading.Thread(target=self._handle_client, args=(client_socket, lottery)).start()