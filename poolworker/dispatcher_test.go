package poolworker

import (
	"testing"
	"workerpool/tasks"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	d := New(2, 10)
	require.NotNil(t, d)
}

func TestStartPool(t *testing.T) {
	d := New(2, 10)
	d.Start()
	defer d.Stop()
	require.Equal(t, len(d.workers), 2)
	for _, w := range d.workers {
		require.NotNil(t, w)
	}

	queued := d.Queue(NewTask("fibo", tasks.Fibn, 2), 0)
	require.Equal(t, queued, true)
}
