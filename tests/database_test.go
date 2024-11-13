// tests for database commands
package tests

import (
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"
)

func TestGetSet(t *testing.T) {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		t.Fatalf("Failed to connect to database server: %v", err)
	}
	defer c.Close()

	// TEST set command
	_, err = c.Do("SET", "test_key", "test_val")
	if err != nil {
		t.Fatalf("Failed to set key: %v", err)
	}

	// TEST get command
	val, err := c.Do("GET", "test_key")
	if err != nil {
		t.Fatalf("Failed to get key: %v", err)
	}

	// assert values set by SET and get by GET
	assert.Equal(t, []byte("test_val"), val, "Expected value to be 'test_val'")
}

func TestHsetHget(t *testing.T) {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		t.Fatalf("Failed to connect to database server: %v", err)
	}
	defer c.Close()

	// TEST hset command
	_, err = c.Do("HSET", "users", "uid", "u1001")
	if err != nil {
		t.Fatalf("Failed to hset key: %v", err)
	}

	// TEST hget command
	val, err := c.Do("HGET", "users", "uid")
	if err != nil {
		t.Fatalf("Failed to hget key: %v", err)
	}

	// assert values
	expected_value := []byte("u1001")

	assert.Equal(t, expected_value, val, "Expected value to be 'u1001'")
}

func TestPing(t *testing.T) {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		t.Fatalf("Failed to connect to database server: %v", err)
	}
	defer c.Close()

	// TEST ping command
	reply, err := c.Do("PING")

	expected_value := "PONG"

	assert.Equal(t, expected_value, reply, "Expected reply to be 'PONG'")
}

// sets a key, set its expire, check after the expiration time
func TestExpire(t *testing.T) {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		t.Fatalf("failed to connect to database server: %v", err)
	}
	defer c.Close()

	// set a key
	_, err = c.Do("SET", "key_1", "value_1")
	if err != nil {
		t.Fatalf("failed to set key: %v", err)
	}

	// set an expiration of 3 seconds
	_, err = c.Do("EXPIRE", "key_1", "3")
	if err != nil {
		t.Fatalf("failed to expire key: %v", err)
	}

	// sleep for 4 seconds and check for key
	time.Sleep(time.Second * 4)

	// get the key
	reply, err := c.Do("GET", "key_1")
	if err != nil {
		t.Fatalf("failed to get key: %v", err)
	}

	assert.Equal(t, nil, reply, "Expected reply to be 'nil'")
}
