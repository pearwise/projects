package utils

import "bufio"

func ReadLine(r *bufio.Reader) ([]byte,error) {
	data, more, err := r.ReadLine()
	if err!=nil {
		return nil, err
	}
	if !more {
		return data, nil
	}
	
	var line []byte
	for more {
		line, more, err = r.ReadLine()
		if err!=nil {
			return nil, err
		}
		data = append(data, line...)
	}
	return data, nil
}