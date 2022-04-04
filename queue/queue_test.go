package queue

import (
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnqueHas(t *testing.T) {
	tests := []struct {
		desc       string
		useHasFunc bool
		element    *Entry
		err        error
	}{
		{
			desc:       "has activated",
			useHasFunc: true,
			element: &Entry{
				Key:   "key1",
				Value: []byte("value1"),
			},
			err: nil,
		},
		{
			desc:       "has disactivated",
			useHasFunc: false,
			element: &Entry{
				Key:   "key1",
				Value: []byte("value1"),
			},
			err: ErrHasFuncDisactivated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			q := NewQueue(2, tt.useHasFunc)
			err := q.Enqueue(tt.element)
			require.NoError(t, err)

			has, err := q.Has(tt.element.Key)
			if tt.err != nil {
				require.EqualError(t, err, tt.err.Error())
				return
			}
			require.NoError(t, err)
			require.Equal(t, has, true)
		})
	}
}

func TestDequeueLenIsEmpty(t *testing.T) {
	q := NewQueue(3, false)

	e1 := &Entry{
		Key:   "e1",
		Value: []byte{},
	}
	err := q.Enqueue(e1)
	require.NoError(t, err)

	e2 := &Entry{
		Key:   "e2",
		Value: []byte{},
	}
	err = q.Enqueue(e2)
	require.NoError(t, err)

	e3 := &Entry{
		Key:   "e3",
		Value: []byte{},
	}
	err = q.Enqueue(e3)
	require.NoError(t, err)

	require.Equal(t, q.Len(), 3)

	e4 := &Entry{
		Key:   "e4",
		Value: []byte{},
	}
	err = q.Enqueue(e4)
	require.EqualError(t, err, ErrOverflow.Error())

	e, _ := q.Dequeue()
	require.Equal(t, e1, e)

	require.Equal(t, q.Len(), 2)

	e, _ = q.Dequeue()
	require.Equal(t, e2, e)

	require.Equal(t, q.Len(), 1)

	e, _ = q.Dequeue()
	require.Equal(t, e3, e)

	require.Equal(t, q.Len(), 0)

	_, isEmpty := q.Dequeue()
	require.Equal(t, isEmpty, true)

	require.Equal(t, q.IsEmpty(), true)

	err = q.Enqueue(e4)
	require.NoError(t, err)

	require.Equal(t, q.Len(), 1)
}

func TestEnqueuThredSafe(t *testing.T) {
	var (
		q        = NewQueue(50, false)
		wg       = &sync.WaitGroup{}
		elements = []*Entry{
			&Entry{},
			&Entry{},
			&Entry{},
			&Entry{},
			&Entry{},
			&Entry{},
			&Entry{},
			&Entry{},
			&Entry{},
			&Entry{},
		}
		pararelCount = 5
	)

	for i := 0; i < pararelCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < len(elements); j++ {
				q.Enqueue(elements[j])
			}
		}()
	}

	wg.Wait()

	if g, w := q.Len(), len(elements)*pararelCount; g != w {
		t.Errorf("got: %d, want: %d", g, w)
	}
}

func TestDequeueThredSafe(t *testing.T) {
	var (
		q = Queue{
			elements: []*Entry{
				&Entry{},
				&Entry{},
				&Entry{},
				&Entry{},
				&Entry{},
				&Entry{},
				&Entry{},
				&Entry{},
				&Entry{},
				&Entry{},
			},
		}

		wg           = &sync.WaitGroup{}
		pararelCount = 5
	)

	for i := 0; i < pararelCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				if _, isEmpty := q.Dequeue(); isEmpty {
					break
				}
			}
		}()
	}

	wg.Wait()

	if g, w := q.Len(), 0; g != w {
		t.Errorf("got: %d, want: %d", g, w)
	}
}

func BenchmarkNewQueue(b *testing.B) {
	for n := 0; n < b.N; n++ {
		es := make([]*Entry, 100)
		for i := range es {
			es[i] = &Entry{
				Key:   strconv.FormatInt(int64(n), 10) + "hoge" + strconv.FormatInt(int64(i), 10),
				Value: []byte{0x01},
			}
		}

		q := NewQueue(DefaultQueusSize, true)
		for i := range es {
			err := q.Enqueue(es[i])
			require.NoError(b, err)
		}
		// _, ok := q.Dequeue()
		// require.Equal(b, ok, false)
	}
}
