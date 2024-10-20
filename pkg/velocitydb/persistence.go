// uses Append Only File(AOF) for persistence
package velocitydb

import (
	"bufio"
	"os"
	"sync"
	"time"
)

type Aof struct {
	file *os.File;
	rd *bufio.Reader;
	mu sync.Mutex;
}


// create new bufio
func NewAof(path string) (*Aof, error) {
	// open file
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666);
	if err != nil {
		return nil, err;
	}

	aof := &Aof{
		file: f,
		rd: bufio.NewReader(f),
	}

	// goroutine to sync aof to disk
	go func ()  {
		
		for {
			aof.mu.Lock();

			aof.file.Sync();

			aof.mu.Unlock();

			time.Sleep(time.Second);
		}
	}();

	return aof, nil;
}

// AOF close file when server shutdown
func (aof *Aof) Close() error {
	aof.mu.Lock();
	defer aof.mu.Unlock();

	return aof.file.Close();
}

// AOF write to file
func (aof *Aof) Write(value Value) error {
	aof.mu.Lock();
	defer aof.mu.Unlock();
	
	// write commands after marshal to aof file
	_, err := aof.file.Write(value.Marshal())
	if err != nil {
		return err;
	}

	return nil;
}
