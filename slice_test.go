package ringslice

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTrueIndex(t *testing.T) {

	tests := []struct {
		name     string
		capacity int
		start    int
		distance int
		want     int
	}{
		{
			name:     "no wrap",
			capacity: 10,
			start:    0,
			distance: 1,
			want:     1,
		},
		{
			name:     "wrap",
			capacity: 10,
			start:    9,
			distance: 2,
			want:     1,
		},
		{
			name:     "no movement",
			capacity: 10,
			start:    9,
			distance: 0,
			want:     9,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSlice(tt.capacity, false, wipeInt)
			ind := s.trueIndex(tt.start, tt.distance)
			require.Equal(t, tt.want, ind)
		})
	}
}

func TestNext(t *testing.T) {
	tests := []struct {
		name string
		cap  int
		ind  int
		want int
	}{
		{
			name: "lowest",
			cap:  10,
			ind:  9,
			want: 0,
		},
		{
			name: "highest",
			cap:  10,
			ind:  8,
			want: 9,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NewSlice(tt.cap, false, wipeInt)
			got := n.next(tt.ind)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestDeleteCount(t *testing.T) {
	tests := []struct {
		name       string
		input      []interface{}
		count      int
		start      int
		used       int
		want       []interface{}
		wantReturn []interface{}
	}{
		{
			name:       "Delete nothing",
			start:      0,
			count:      0,
			used:       5,
			input:      []interface{}{1, 2, 3, 4, 5},
			want:       []interface{}{1, 2, 3, 4, 5},
			wantReturn: []interface{}{},
		},
		{
			name:       "Delete one",
			start:      0,
			count:      1,
			used:       5,
			input:      []interface{}{1, 2, 3, 4, 5},
			want:       []interface{}{0, 2, 3, 4, 5},
			wantReturn: []interface{}{1},
		},
		{
			name:       "Delete all",
			start:      0,
			count:      5,
			used:       5,
			input:      []interface{}{1, 2, 3, 4, 5},
			want:       []interface{}{0, 0, 0, 0, 0},
			wantReturn: []interface{}{1, 2, 3, 4, 5},
		},
		{
			name:       "Delete many wrap",
			start:      3,
			count:      4,
			used:       5,
			input:      []interface{}{1, 2, 3, 4, 5},
			want:       []interface{}{0, 0, 3, 0, 0},
			wantReturn: []interface{}{4, 5, 1, 2},
		},
		{
			name:       "Delete too many",
			start:      0,
			count:      100,
			used:       5,
			input:      []interface{}{1, 2, 3, 4, 5},
			want:       []interface{}{0, 0, 0, 0, 0},
			wantReturn: []interface{}{1, 2, 3, 4, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Slice{
				values: tt.input,
				start:  tt.start,
				cap:    len(tt.input),
				used:   tt.used,
				wipe:   wipeInt,
			}
			got := n.DeleteCount(tt.count)
			require.Equal(t, tt.want, n.values)
			require.Equal(t, tt.wantReturn, got)
		})
	}
}
func wipeInt(i int, l []interface{}) {
	l[i] = 0
}
func TestDeleteBounds(t *testing.T) {
	tests := []struct {
		name      string
		input     []interface{}
		end       int
		start     int
		wantStart int
		used      int
		want      []interface{}
	}{
		{
			name:      "Delete single",
			start:     0,
			end:       0,
			wantStart: 1,
			used:      5,
			input:     []interface{}{1, 2, 3, 4, 5},
			want:      []interface{}{0, 2, 3, 4, 5},
		},
		{
			name:      "Delete with traversal",
			start:     0,
			end:       1,
			wantStart: 2,
			used:      5,
			input:     []interface{}{1, 2, 3, 4, 5},
			want:      []interface{}{0, 0, 3, 4, 5},
		},
		{
			name:      "Delete all",
			start:     0,
			end:       4,
			wantStart: 0,
			used:      5,
			input:     []interface{}{1, 2, 3, 4, 5},
			want:      []interface{}{0, 0, 0, 0, 0},
		},
		{
			name:      "Wrap",
			start:     3,
			end:       1,
			wantStart: 2,
			used:      5,
			input:     []interface{}{1, 2, 3, 4, 5},
			want:      []interface{}{0, 0, 3, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Slice{
				values: tt.input,
				start:  tt.start,
				cap:    len(tt.input),
				wipe:   wipeInt,
				used:   tt.used,
			}
			n.DeleteBounds(tt.start, tt.end)
			require.Equal(t, tt.want, n.values)
			require.Equal(t, tt.wantStart, n.start)
		})
	}
}

func TestAppend(t *testing.T) {
	type appendExpectationStruct struct {
		val int
		err error
	}

	tests := []struct {
		name   string
		input  []interface{}
		append []appendExpectationStruct
		start  int
		want   []interface{}
	}{
		{
			name:   "Append one",
			start:  0,
			append: []appendExpectationStruct{{val: 10, err: nil}},
			input:  []interface{}{1, 2, 3, 4, 5},
			want:   []interface{}{10, 2, 3, 4, 5},
		},
		{
			name:  "Append too many",
			start: 0,
			append: []appendExpectationStruct{
				{val: 4, err: nil},
				{val: 5, err: nil},
				{val: 9, err: errors.New("Index is full cannot append")},
			},
			input: []interface{}{1, 2},
			want:  []interface{}{4, 5},
		},
		{
			name:  "Successful wrap around",
			start: 4,
			append: []appendExpectationStruct{
				{val: 10, err: nil},
				{val: 10, err: nil},
			},
			input: []interface{}{1, 2, 3, 4, 5},
			want:  []interface{}{10, 2, 3, 4, 10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Slice{
				values: tt.input,
				start:  tt.start,
				cap:    len(tt.input),
				wipe:   wipeInt,
			}
			for _, ex := range tt.append {
				err := n.Append(ex.val)
				require.Equal(t, ex.err, err)
			}
			require.Equal(t, tt.want, n.values)
		})
	}
}

func TestSlice_FindClosestBelow(t *testing.T) {
	type fields struct {
		values []interface{}
		used   int
		start  int
		end    int
		debug  bool
		cap    int
	}
	type args struct {
		want  int64
		value func(interface{}) int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "Full midpoint find",
			fields: fields{
				values: []interface{}{int64(1), int64(2), int64(3), int64(4), int64(5)},
				used:   5,
				start:  0,
				end:    4,
				debug:  false,
				cap:    5,
			},
			args: args{
				want:  3,
				value: func(i interface{}) int64 { return i.(int64) },
			},
			want: 2,
		},
		{
			name: "All equal and below",
			fields: fields{
				values: []interface{}{int64(1), int64(1), int64(1), int64(1), int64(1)},
				used:   5,
				start:  0,
				end:    4,
				debug:  false,
				cap:    5,
			},
			args: args{
				want:  3,
				value: func(i interface{}) int64 { return i.(int64) },
			},
			want: 4,
		},
		{
			name: "All equal exact",
			fields: fields{
				values: []interface{}{int64(3), int64(3), int64(3), int64(3), int64(3)},
				used:   5,
				start:  0,
				end:    4,
				debug:  false,
				cap:    5,
			},
			args: args{
				want:  3,
				value: func(i interface{}) int64 { return i.(int64) },
			},
			want: 4,
		},
		{
			name: "All equal and above",
			fields: fields{
				values: []interface{}{int64(5), int64(5), int64(5), int64(5), int64(5)},
				used:   5,
				start:  0,
				end:    4,
				debug:  false,
				cap:    5,
			},
			args: args{
				want:  3,
				value: func(i interface{}) int64 { return i.(int64) },
			},
			want: -1,
		},
		{
			name: "Equal and above midpoint",
			fields: fields{
				values: []interface{}{int64(1), int64(2), int64(3), int64(3), int64(3)},
				used:   5,
				start:  0,
				end:    4,
				debug:  false,
				cap:    5,
			},
			args: args{
				want:  3,
				value: func(i interface{}) int64 { return i.(int64) },
			},
			want: 4,
		},
		{
			name: "empty",
			fields: fields{
				used:  0,
				start: 0,
				end:   0,
			},
			args: args{
				want:  3,
				value: func(i interface{}) int64 { return i.(int64) },
			},
			want: -1,
		},
		{
			name: "single below",
			fields: fields{
				values: []interface{}{int64(1)},
				used:   1,
				start:  0,
				end:    0,
				debug:  false,
				cap:    1,
			},
			args: args{
				want:  2,
				value: func(i interface{}) int64 { return i.(int64) },
			},
			want: 0,
		},
		{
			name: "single above",
			fields: fields{
				values: []interface{}{int64(10)},
				used:   1,
				start:  0,
				end:    0,
				debug:  false,
				cap:    1,
			},
			args: args{
				want:  2,
				value: func(i interface{}) int64 { return i.(int64) },
			},
			want: -1,
		},
		{
			name: "two node boundary both above",
			fields: fields{
				values: []interface{}{int64(10), int64(11)},
				used:   2,
				start:  0,
				end:    1,
				debug:  false,
				cap:    2,
			},
			args: args{
				want:  5,
				value: func(i interface{}) int64 { return i.(int64) },
			},
			want: -1,
		},
		{
			name: "two node boundary across boundary",
			fields: fields{
				values: []interface{}{int64(3), int64(5)},
				used:   2,
				start:  0,
				end:    1,
				debug:  false,
				cap:    2,
			},
			args: args{
				want:  4,
				value: func(i interface{}) int64 { return i.(int64) },
			},
			want: 0,
		},
		{
			name: "two node boundary both below",
			fields: fields{
				values: []interface{}{int64(1), int64(2)},
				used:   2,
				start:  0,
				end:    1,
				debug:  false,
				cap:    2,
			},
			args: args{
				want:  5,
				value: func(i interface{}) int64 { return i.(int64) },
			},
			want: 1,
		},
		{
			name: "two node backwards boundary",
			fields: fields{
				values: []interface{}{int64(10), int64(1)},
				used:   2,
				start:  1,
				end:    0,
				debug:  false,
				cap:    2,
			},
			args: args{
				want:  5,
				value: func(i interface{}) int64 { return i.(int64) },
			},
			want: 1,
		},
		{
			name: "wrap around requiring multiple midpoint calculation",
			fields: fields{
				values: []interface{}{int64(1561882874), int64(1561882875), int64(1561882876), int64(0), int64(0), int64(0), int64(0), int64(0), int64(1561882872), int64(1561882873)},
				used:   5,
				start:  8,
				end:    2,
				debug:  false,
				cap:    10,
			},
			args: args{
				want:  1561882872,
				value: func(i interface{}) int64 { return i.(int64) },
			},
			want: 8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(tt.name)
			s := &Slice{
				values: tt.fields.values,
				used:   tt.fields.used,
				start:  tt.fields.start,
				cap:    tt.fields.cap,
				wipe:   wipeInt,
			}
			if got := s.FindClosestBelowOrEqual(tt.args.want, tt.args.value); got != tt.want {
				t.Errorf("Slice.FindClosestBelow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_countBetween(t *testing.T) {
	type args struct {
		start int
		end   int
		cap   int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "single",
			args: args{
				start: 0,
				end:   0,
				cap:   5,
			},
			want: 1,
		},
		{
			name: "no wrap",
			args: args{
				start: 0,
				end:   4,
				cap:   5,
			},
			want: 5,
		},
		{
			name: "wrap double",
			args: args{
				start: 4,
				end:   0,
				cap:   5,
			},
			want: 2,
		},
		{
			name: "wrap long",
			args: args{
				start: 2,
				end:   1,
				cap:   5,
			},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := countBetween(tt.args.start, tt.args.end, tt.args.cap); got != tt.want {
				t.Errorf("countBetween() = %v, want %v", got, tt.want)
			}
		})
	}
}
