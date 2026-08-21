package safe_socket

import "io"

func SendAll(socket io.Writer, bytes []byte) error {

	bytes_send := 0
	for bytes_send < len(bytes){
		sended, err := socket.Write(bytes)
		bytes_send += sended

		if err != nil {
			return err
		}
	}
	return nil
}

func RecvAll(socket io.Reader, size int) ([]byte, error) {
	//TODO short read
	buff := make([]byte, size)
	_, err := socket.Read(buff)
	if err != nil {
		return nil, err
	}
	return buff, nil
}
