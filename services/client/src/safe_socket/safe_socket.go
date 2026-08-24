package safe_socket

import "io"

func SendAll(socket io.Writer, bytes []byte) error {

	bytes_send := 0
	for bytes_send < len(bytes){
		sended, err := socket.Write(bytes[bytes_send:])
		bytes_send += sended

		if err != nil {
			return err
		}
	}
	return nil
}

func RecvAll(socket io.Reader, size int) ([]byte, error) {
	bytes := make([]byte, size)

	recived := 0
	for recived < size {
		byte_rec, err := socket.Read(bytes[recived:])
		recived += byte_rec
		if err != nil {
			return nil, err
		}
	}

	return bytes, nil
}
