package cron

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCommand_AddFunc(t *testing.T) {
	s := &Service{}
	ok, err := s.Init(nil, &Config{WorkDir: "tests"})
	assert.NoError(t, err)
	assert.True(t, ok)
	go func() {
		assert.NoError(t, s.Serve())
	}()
	cmd := Command{}
	err = cmd.AddFunc("* * * * * *", "testCron.php")
	assert.NoError(t, err)
	assert.Len(t, cmd.runningTasks(), 1)
}
