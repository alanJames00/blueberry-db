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

// RESP: SERIALIZER
type Writer struct {
	writer io.Writer;
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer: w,
	}
}

func NewValue(typ string, str string, num int, bulk string, array []Value) *Value {
	return &Value{
		typ: typ,
		str: str,
		num: num,
		bulk: bulk,
		array: array,
	}
}

func (w* Writer) Write(v Value) error {
	var bytes = v.Marshal();

	_, err := w.writer.Write(bytes);
	if err != nil {
		return err;
	}

	return nil;
}

func (v Value) Marshal() []byte {
	// handle diff types
	switch v.typ {
	case "array":
		return v.marshalArray();
	case "bulk":
		return v.marshalBulk();
	case "string":
		return v.marshalString();
	case "null":
		return v.marshalNull();
	case "error":
		return v.marshalError();
	default:
		// return empty byte array
		return []byte{}
	}
}

// marshall simple strings
func (v Value) marshalString() []byte {
	var bytes []byte;

	bytes = append(bytes, STRING);
	bytes = append(bytes, v.str...);
	bytes = append(bytes, '\r', '\n');

	return bytes;
}

// marshall bulk
func (v Value) marshalBulk() []byte {
	var bytes []byte;
	bytes = append(bytes, BULK);
	bytes = append(bytes, strconv.Itoa(len(v.bulk))...);
	bytes = append(bytes, '\r', '\n');
	bytes = append(bytes, v.bulk...);
	bytes = append(bytes, '\r', '\n');

	return bytes;
}

// marshall array
func (v Value) marshalArray() []byte {
	len := len(v.array);

	var bytes []byte;
	bytes = append(bytes, ARRAY);
	bytes = append(bytes, strconv.Itoa(len)...);
	bytes = append(bytes, '\r', '\n');

	for i := 0; i < len; i++ {
		bytes = append(bytes, v.array[i].Marshal()...);
	}

	return bytes;
}

// marshal null
func (v Value) marshalNull() []byte {
	return []byte("$-1\r\n");
}

// marshal error
func (v Value) marshalError() []byte {
	var bytes []byte;

	bytes = append(bytes, ERROR);
	bytes = append(bytes, v.str...);
	bytes = append(bytes, '\r', '\n');

	return bytes;
}
