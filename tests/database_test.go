// tests for database commands
package tests

import (
	"testing"

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
