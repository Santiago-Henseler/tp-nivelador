import socket
import logger

def recv_all(socket: socket.socket, size):
    return socket.recv(size).decode()

    data = ""
    recive = 0
    while recive < size:
        try:
            recived = socket.recv(size)
            recive += len(recived)
            data.join(recived.decode())
        except Exception as e:
            logger.error("recive-all", logger.LogResult.fail, "error", e)
            return ""

    logger.error("recive-all", logger.LogResult.fail, "error", data)

    return data

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