package react

import (
	"runtime"
	"testing"
	"time"
)

var _ Reactor = New()

func assertCellValue(t *testing.T, c Cell, expected int, explanation string) {
	observed := c.Value()
	_, _, line, _ := runtime.Caller(1)
	if observed != expected {
		t.Fatalf("(from line %d) %s: expected %d, got %d", line, explanation, expected, observed)
	}
}

func TestSetInput(t *testing.T) {
	r := New()
	i := r.CreateInput(1)
	assertCellValue(t, i, 1, "i.Value() doesn't match initial value")
	i.SetValue(2)
	time.Sleep(100 * time.Millisecond)
	assertCellValue(t, i, 2, "i.Value() doesn't match changed value")
}

func TestBasicCompute1(t *testing.T) {
	r := New()
	i := r.CreateInput(1)
	c := r.CreateCompute1(i, func(v int) int { return v + 1 })
	assertCellValue(t, c, 2, "c.Value() isn't properly computed based on initial input cell value")
	i.SetValue(2)
	time.Sleep(10 * time.Millisecond) // @todo remove
	assertCellValue(t, c, 3, "c.Value() isn't properly computed based on changed input cell value")
}

func TestBasicCompute2(t *testing.T) {
	r := New()
	i1 := r.CreateInput(1)
	i2 := r.CreateInput(2)
	c := r.CreateCompute2(i1, i2, func(v1, v2 int) int { return v1 | v2 })
	assertCellValue(t, c, 3, "c.Value() isn't properly computed based on initial input cell values")
	i1.SetValue(4)
	assertCellValue(t, c, 6, "c.Value() isn't properly computed when first input cell value changes")
	i2.SetValue(8)
	assertCellValue(t, c, 12, "c.Value() isn't properly computed when second input cell value changes")
}

func TestCompute2Diamond(t *testing.T) {
	r := New()
	i := r.CreateInput(1)
	c1 := r.CreateCompute1(i, func(v int) int { return v + 1 })
	c2 := r.CreateCompute1(i, func(v int) int { return v - 1 })
	c3 := r.CreateCompute2(c1, c2, func(v1, v2 int) int { return v1 * v2 })
	assertCellValue(t, c3, 0, "c3.Value() isn't properly computed based on initial input cell value")
	i.SetValue(3)
	assertCellValue(t, c3, 8, "c3.Value() isn't properly computed based on changed input cell value")
}

func TestCompute1Chain(t *testing.T) {
	r := New()
	inp := r.CreateInput(1)
	var c Cell = inp
	for i := 2; i <= 8; i++ {
		digitToAdd := i
		c = r.CreateCompute1(c, func(v int) int { return v*10 + digitToAdd })
	}
	assertCellValue(t, c, 12345678, "c.Value() isn't properly computed based on initial input cell value")
	inp.SetValue(9)
	assertCellValue(t, c, 92345678, "c.Value() isn't properly computed based on changed input cell value")
}

func TestCompute2Tree(t *testing.T) {
	r := New()
	ins := make([]InputCell, 3)
	for i, v := range []int{1, 10, 100} {
		ins[i] = r.CreateInput(v)
	}

	add := func(v1, v2 int) int { return v1 + v2 }

	firstLevel := make([]ComputeCell, 2)
	for i := 0; i < 2; i++ {
		firstLevel[i] = r.CreateCompute2(ins[i], ins[i+1], add)
	}

	output := r.CreateCompute2(firstLevel[0], firstLevel[1], add)
	assertCellValue(t, output, 121, "output.Value() isn't properly computed based on initial input cell values")

	for i := 0; i < 3; i++ {
		ins[i].SetValue(ins[i].Value() * 2)
	}

	assertCellValue(t, output, 242, "output.Value() isn't properly computed based on changed input cell values")
}

func TestBasicCallback(t *testing.T) {
	r := New()
	i := r.CreateInput(1)
	c := r.CreateCompute1(i, func(v int) int { return v + 1 })
	var observed []int
	c.AddCallback(func(v int) {
		observed = append(observed, v)
	})
	if len(observed) != 0 {
		t.Fatalf("callback called before changes were made")
	}
	i.SetValue(2)
	if len(observed) != 1 {
		t.Fatalf("callback not called when changes were made")
	}
	if observed[0] != 3 {
		t.Fatalf("callback not called with proper value")
	}
}

func TestOnlyCallOnChanges(t *testing.T) {
	r := New()
	i := r.CreateInput(1)
	c := r.CreateCompute1(i, func(v int) int {
		if v > 3 {
			return v + 1
		}
		return 2
	})
	var observedCalled int
	c.AddCallback(func(int) {
		observedCalled++
	})
	i.SetValue(1)
	if observedCalled != 0 {
		t.Fatalf("observe function called even though input didn't change")
	}
	i.SetValue(2)
	if observedCalled != 0 {
		t.Fatalf("observe function called even though computed value didn't change")
	}
}

func TestCallbackAddRemove(t *testing.T) {
	r := New()
	i := r.CreateInput(1)
	c := r.CreateCompute1(i, func(v int) int { return v + 1 })
	var observed1 []int
	cb1 := c.AddCallback(func(v int) {
		observed1 = append(observed1, v)
	})
	var observed2 []int
	c.AddCallback(func(v int) {
		observed2 = append(observed2, v)
	})
	i.SetValue(2)
	if len(observed1) != 1 || observed1[0] != 3 {
		t.Fatalf("observed1 not properly called")
	}
	if len(observed2) != 1 || observed2[0] != 3 {
		t.Fatalf("observed2 not properly called")
	}
	cb1.Cancel()
	i.SetValue(3)
	if len(observed1) != 1 {
		t.Fatalf("observed1 called after removal")
	}
	if len(observed2) != 2 || observed2[1] != 4 {
		t.Fatalf("observed2 not properly called after first callback removal")
	}
}

func TestMultipleCallbackRemoval(t *testing.T) {
	r := New()
	inp := r.CreateInput(1)
	c := r.CreateCompute1(inp, func(v int) int { return v + 1 })

	numCallbacks := 5

	calls := make([]int, numCallbacks)
	cancelers := make([]Canceler, numCallbacks)
	for i := 0; i < numCallbacks; i++ {
		i := i
		cancelers[i] = c.AddCallback(func(v int) { calls[i]++ })
	}

	inp.SetValue(2)
	for i := 0; i < numCallbacks; i++ {
		if calls[i] != 1 {
			t.Fatalf("callback %d/%d should be called 1 time, was called %d times", i+1, numCallbacks, calls[i])
		}
		cancelers[i].Cancel()
	}

	inp.SetValue(3)
	for i := 0; i < numCallbacks; i++ {
		if calls[i] != 1 {
			t.Fatalf("callback %d/%d was called after it was removed", i+1, numCallbacks)
		}
	}
}

func TestRemoveIdempotence(t *testing.T) {
	r := New()
	inp := r.CreateInput(1)
	output := r.CreateCompute1(inp, func(v int) int { return v + 1 })
	timesCalled := 0
	cb1 := output.AddCallback(func(int) {})
	output.AddCallback(func(int) { timesCalled++ })
	for i := 0; i < 10; i++ {
		cb1.Cancel()
	}
	inp.SetValue(2)
	if timesCalled != 1 {
		t.Fatalf("remaining callback function was not called")
	}
}

func TestOnlyCallOnceOnMultipleDepChanges(t *testing.T) {
	r := New()
	i := r.CreateInput(1)
	c1 := r.CreateCompute1(i, func(v int) int { return v + 1 })
	c2 := r.CreateCompute1(i, func(v int) int { return v - 1 })
	c3 := r.CreateCompute1(c2, func(v int) int { return v - 1 })
	c4 := r.CreateCompute2(c1, c3, func(v1, v3 int) int { return v1 * v3 })
	changed4 := 0
	c4.AddCallback(func(int) { changed4++ })
	i.SetValue(3)
	if changed4 < 1 {
		t.Fatalf("callback function was not called")
	} else if changed4 > 1 {
		t.Fatalf("callback function was called too often")
	}
}

func TestNoCallOnDepChangesResultingInNoChange(t *testing.T) {
	r := New()
	inp := r.CreateInput(0)
	plus1 := r.CreateCompute1(inp, func(v int) int { return v + 1 })
	minus1 := r.CreateCompute1(inp, func(v int) int { return v - 1 })
	output := r.CreateCompute2(plus1, minus1, func(v1, v2 int) int { return v1 - v2 })

	timesCalled := 0
	output.AddCallback(func(int) { timesCalled++ })

	inp.SetValue(5)
	if timesCalled != 0 {
		t.Fatalf("callback function called even though computed value didn't change")
	}
}
