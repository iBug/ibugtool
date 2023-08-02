package zfs

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testProperty struct {
	Name         string
	Type         string
	CreationTime time.Time `zfs:"creation"`
	Skip         string    `zfs:"-"`
	Used         uint64
	Available    uint64
}

func TestBuildPropertiesArgument(t *testing.T) {
	assert.Equal(t,
		"name,type,creation,used,available",
		buildPropertiesArgument(testProperty{}))
}

func TestParseProperties(t *testing.T) {
	p := testProperty{Skip: "untouched"}
	testInput := "test-dataset\tfilesystem\t1600000000\t123456\t654321"
	assert.Nil(t, parseProperties(&p, testInput))
	assert.Equal(t, "test-dataset", p.Name)
	assert.Equal(t, "filesystem", p.Type)
	assert.Equal(t, time.Unix(1600000000, 0), p.CreationTime)
	assert.Equal(t, "untouched", p.Skip)
	assert.Equal(t, uint64(123456), p.Used)
	assert.Equal(t, uint64(654321), p.Available)
}
