package snowflake_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tackboon/snowflake"
)

func TestGenerate64BitSnowflakeID(t *testing.T) {
	count := 1000000
	generator := snowflake.NewDefaultSnowflake64Bit(1, 1744615749000)

	idMap := make(map[uint64]struct{})
	prevID := uint64(0)
	for range count {
		id, err := generator.GenerateID()
		assert.NoError(t, err, "failed to generate 64 bit snowflake id")
		assert.Less(t, prevID, id, "64 bit snowflake id not in sequence")
		idMap[id] = struct{}{}
		prevID = id
	}

	assert.Equal(t, len(idMap), count)
}

func TestDecode64BitSnowflakeID(t *testing.T) {
	generator := snowflake.NewDefaultSnowflake64Bit(1, 1744615749000)

	now := time.Now()
	id, err := generator.GenerateID()
	assert.NoError(t, err, "failed to generate 53 bit snowflake id")

	timestamp, machineID, seq := generator.DecodeID(id)
	assert.Less(t, now.UnixMilli()-timestamp, int64(1000), "ID generated time more than 1s")
	assert.Equal(t, machineID, int64(1), "invalid machine id: %d", machineID)
	assert.Equal(t, seq, int64(0), "invalid seq: %d", seq)
}
