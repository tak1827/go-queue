package queue

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnqueHas(t *testing.T) {
	tests := []struct {
		desc       string
		useHasFunc bool
		element    Entry
		err        error
	}{
		{
			desc:       "has activated",
			useHasFunc: true,
			element: Entry{
				Key:   "key1",
				Value: "value1",
			},
			err: nil,
		},
		{
			desc:       "has disactivated",
			useHasFunc: false,
			element: Entry{
				Key:   "key1",
				Value: "value1",
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

	elm1 := "string"
	err := q.Enqueue(elm1)
	require.NoError(t, err)

	elm2 := 123
	err = q.Enqueue(elm2)
	require.NoError(t, err)

	elm3 := Entry{"key3", true}
	err = q.Enqueue(elm3)
	require.NoError(t, err)

	require.Equal(t, q.Len(), 3)

	elm4 := false
	err = q.Enqueue(elm4)
	require.EqualError(t, err, ErrOverflow.Error())

	elm, _ := q.Dequeue()
	elm = elm.(string)
	require.Equal(t, elm, elm1)

	require.Equal(t, q.Len(), 2)

	elm, _ = q.Dequeue()
	elm = elm.(int)
	require.Equal(t, elm, elm2)

	require.Equal(t, q.Len(), 1)

	elm, _ = q.Dequeue()
	elm = elm.(Entry)
	require.Equal(t, elm, elm3)

	require.Equal(t, q.Len(), 0)

	_, isEmpty := q.Dequeue()
	require.Equal(t, isEmpty, true)

	require.Equal(t, q.IsEmpty(), true)

	err = q.Enqueue(elm4)
	require.NoError(t, err)

	require.Equal(t, q.Len(), 1)
}

func TestEnqueuThredSafe(t *testing.T) {
	var (
		q        = NewQueue(50, false)
		wg       = &sync.WaitGroup{}
		elements = []interface{}{
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
			elements: []interface{}{
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
