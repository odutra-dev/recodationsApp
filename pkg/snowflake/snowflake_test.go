package snowflake

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSnowflakeGeneration(t *testing.T) {
	err := Init(1)
	require.NoError(t, err)

	id1 := Generate()
	id2 := Generate()

	assert.Positive(t, id1)
	assert.Positive(t, id2)
	assert.NotEqual(t, id1, id2, "Snowflake IDs generated consecutively must be distinct")

	strID := IDString()
	assert.NotEmpty(t, strID)

	parsed, err := Parse(id1)
	require.NoError(t, err)
	assert.Equal(t, id1, parsed.Int64())
}

func TestSnowflakeUniquenessConcurrent(t *testing.T) {
	err := Init(2)
	require.NoError(t, err)

	const count = 1000
	idMap := make(map[int64]bool)
	idChan := make(chan int64, count)

	for i := 0; i < count; i++ {
		go func() {
			idChan <- Generate()
		}()
	}

	for i := 0; i < count; i++ {
		id := <-idChan
		assert.False(t, idMap[id], "Duplicate Snowflake ID generated: %d", id)
		idMap[id] = true
	}
	assert.Equal(t, count, len(idMap))
}
