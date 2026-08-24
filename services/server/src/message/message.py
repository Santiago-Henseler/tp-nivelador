import socket
import safe_socket
import logger #TODO
from lottery.bet import Bet

END_MESSAGE =  b'\x00'
BET_MESSAGE = b'\x01'
BATCH_MESSAGE = b'\x02'

def recive_message(socket: socket):
    type = safe_socket.recv_all(socket, 1)

    if type == BET_MESSAGE:
        return [recive_bet_message(socket)]
    elif type == END_MESSAGE:
        return None
    elif type == BATCH_MESSAGE:
        logger.info( "a", logger.LogResult.success, "messages-amount", type)
        return recibe_batch_message(socket)
    else:
        return None

def recibe_batch_message(socket: socket.socket):
    info = safe_socket.recv_all(socket, 8)
    if len(info) != 8:
        return None

    bets_size = int.from_bytes(info[0:4], byteorder='big')
    bets_lenght = int.from_bytes(info[4:8], byteorder='big')

    bets_bytes = safe_socket.recv_all(socket, bets_lenght)

    if len(bets_bytes) != bets_lenght:
        return None

    bets = [] 
    size = 0

    for i in range(bets_size):
        agency_id = int.from_bytes(bets_bytes[size:size+4], byteorder='big')
        document = int.from_bytes(bets_bytes[size+4:size+8], byteorder='big')
        number = int.from_bytes(bets_bytes[size+8:size+12], byteorder='big')
        
        first_name_len = int.from_bytes(bets_bytes[size+12:size+16], byteorder='big')
        last_name_len = int.from_bytes(bets_bytes[size+16:size+20], byteorder='big')
        birthdate_len = int.from_bytes(bets_bytes[size+20:size+24], byteorder='big')

        size += 24

        first_name = bets_bytes[size:size+first_name_len].decode()
        last_name = bets_bytes[size+first_name_len:size+first_name_len+last_name_len].decode()
        birthdate = bets_bytes[size+first_name_len+last_name_len:size+first_name_len+last_name_len+birthdate_len].decode()

        size += (first_name_len+last_name_len+birthdate_len)
        bets.append(Bet(agency_id, first_name, last_name, document, birthdate, number))

    return bets


def send_batch_message(socket: socket.socket, bets):
    bytes = BATCH_MESSAGE
    bytes += len(bets).to_bytes(4, byteorder='big')

    bet_bytes = b''

    for bet in bets:
        bet_bytes += bet_to_bytes(bet)

    bytes += len(bet_bytes).to_bytes(4, byteorder='big')
    bytes += bet_bytes

    safe_socket.send_all(socket, bytes)

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

def send_bet_message(socket: socket.socket, bet):
    bytes = BET_MESSAGE
    bytes += bet_to_bytes(bet)

    safe_socket.send_all(socket, bytes)

def end_bet_message(socket: socket.socket):
    safe_socket.send_all(socket, END_MESSAGE)

def bet_to_bytes(bet):
    bytes = b''

    bytes += bet.agency_id.to_bytes(4, byteorder='big')
    bytes += bet.document.to_bytes(4, byteorder='big')
    bytes += bet.number.to_bytes(4, byteorder='big')

    bytes += len(bet.first_name).to_bytes(4, byteorder='big')
    bytes += len(bet.last_name).to_bytes(4,  byteorder='big')
    bytes += len(bet.birthdate).to_bytes(4, byteorder='big')

    bytes += bet.first_name.encode()
    bytes += bet.last_name.encode()
    bytes += bet.birthdate.encode()

    return bytes