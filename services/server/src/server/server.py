import socket
import logger
import message
import threading
from lottery.lottery import Lottery

OUTPUT_FILE = "output.txt"

class Server:
    def __init__(self, server_host: str, server_port: int, quorum) -> None:
        self.server_host = server_host
        self.server_port = server_port
        self.barrier = threading.Barrier(quorum)
        
    def _handle_client(self, client_socket, lottery):
        action = "handle-client"
        message_amount = 0
        try:
            logger.info(action, logger.LogResult.in_progress)
            comunication = True
            while comunication:
                client_bet = message.recive_message(client_socket)

                if client_bet == None:
                    logger.info( action, logger.LogResult.success, "messages-amount", message_amount)
                    comunication = False
                    continue
                
                lottery.store_bets(client_bet)
                message_amount += 1
        except Exception as e:
            logger.error( action, logger.LogResult.fail, "messages-amount", message_amount)
            raise e

        self.barrier.wait()

        for bet in lottery.load_bets():
            if lottery.has_won(bet):
                message.send_bet_message(client_socket,  bet)

        message.end_bet_message(client_socket)


    def run(self):
        action = "accept-connection"
        with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as server_socket:
            server_socket.bind((self.server_host, self.server_port))
            server_socket.listen()
            lottery = Lottery("output/"+OUTPUT_FILE)
            while True:
                try:
                    logger.info(action, logger.LogResult.in_progress)
                    client_socket, _ = server_socket.accept()
                except Exception as e:
                    logger.error(action, logger.LogResult.fail)
                    raise e
                logger.info(action, logger.LogResult.success)
               
                threading.Thread(target=self._handle_client, args=(client_socket, lottery), ).start()
