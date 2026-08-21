import socket
import safe_socket
import logger
from lottery.bet import Bet

def recive_bet_message(socket: socket.socket):

    integers = safe_socket.recv_all(socket, 24)

    if len(integers) < 24 :
        return None

    agency_id = int.from_bytes(integers[0:4], byteorder='big')
    document = int.from_bytes(integers[4:8], byteorder='big')
    number = int.from_bytes(integers[8:12], byteorder='big')

    first_name_len = int.from_bytes(integers[12:16], byteorder='big')
    last_name_len = int.from_bytes(integers[16:20], byteorder='big')
    birthdate_len = int.from_bytes(integers[20:24], byteorder='big')

    strings = safe_socket.recv_all(socket, first_name_len + last_name_len + birthdate_len)

    first_name = strings[0:first_name_len].decode()
    last_name = strings[first_name_len:first_name_len+last_name_len].decode()
    birthdate = strings[first_name_len+last_name_len:].decode()


    return Bet(agency_id, first_name, last_name, document, birthdate, number)


def send_bet_message(socket: socket.socket):
    