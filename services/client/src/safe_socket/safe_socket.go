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
	buffer := make([]byte, size)
	bytes := make([]byte, 0, size)

	var recived = 0

	for recived < size {
		byte_rec, err := socket.Read(buffer)
		if err != nil {
			return nil, err
		}
		recived += byte_rec
		bytes = append(bytes, buffer[:byte_rec]...) 
	}

	return bytes, nil
}
