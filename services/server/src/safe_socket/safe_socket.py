import socket
import logger

def recv_all(socket: socket.socket, size):
    bytes = b""
    recived = 0
    while recived < size:
        try:
            byte_rec = socket.recv(size - recived)
            if not byte_rec :
                return bytes

            recived += len(byte_rec)
            bytes += byte_rec
        except Exception as e:
            logger.error("recive-all", logger.LogResult.fail, "error", e)
            return b""

    return bytes

def send_all(socket: socket.socket, bytes):

    bytes_send = 0 
    while bytes_send < len(bytes):
        try:
            sended = socket.send(bytes[bytes_send:len(bytes)])
            bytes_send += sended
        except Exception as e:
            logger.error("send-all", logger.LogResult.fail, "error", e)
            return 0

    return bytes_send