// RESP command handler and parser
package velocitydb

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	STRING = '+';
	ERROR = '-';
	INTEGER = ':';
	BULK = '$';
	ARRAY = '*';
)

type Value struct {
	typ string;
	str string;
	num int;
	bulk string;
	array []Value;
}

type Resp struct {
	reader *bufio.Reader;
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{
		reader: bufio.NewReader(rd),	
	}
}

func (r *Resp) readLine() (line []byte, n int, err error) {
	// loop and read
	for {
		b, err := r.reader.ReadByte();
		if err != nil {
			return nil, 0, err;
		}

		n+=1;
		line = append(line, b);
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break;
		}
	}

	return line[:len(line)-2], n, nil;
}

func (r* Resp) readInteger() (x int, n int, err error) {
	// read the line with reader
	line, n, err := r.readLine();
	if err != nil {
		return 0, 0, err;
	}

	// parse int
	i64, err := strconv.ParseInt(string(line), 10, 64);
	if err != nil {
		return 0, n, err;
	}

	return int(i64), n, nil;
}

func (r *Resp) readArray() (Value, error) {
	// create result value
	v := Value{};
	v.typ = "array";

	// read the length of the array
	len, _, err := r.readInteger();
	if err != nil {
		return v, err;
	}

	// foreach line, parse and read the value
	v.array = make([]Value, 0);
	for i := 0; i < len; i++ {
		val, err := r.Read();
		if err != nil {
			return v, err;
		}

		// append the parsed values to the res array
		v.array = append(v.array, val);
	}

	return v, nil;
}

func (r *Resp) readBulk() (Value, error) {
	// create result value
	v := Value{};

	v.typ = "bulk";

	len, _, err := r.readInteger();
	if err != nil {
		return v, err;
	}

	bulk := make([]byte, len);

	r.reader.Read(bulk);

	v.bulk = string(bulk);

	// read the trailing CRLF
	r.readLine();

	return v, nil;
}

// recursive RESP reader
func (r *Resp) Read() (Value, error) {
	
	// get datatype from first byte
	_type, err := r.reader.ReadByte();

	if err != nil {
		return Value{}, err;
	}
	

	// handle _types
	switch _type {
	case ARRAY:
		return r.readArray();
	case BULK:
		return r.readBulk();
	default:
		fmt.Printf("Unknown type: %v", string(_type))
		return Value{}, nil;
	}
}

