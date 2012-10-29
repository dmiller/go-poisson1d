package poisson1d

import (
	"testing"
)


func TestNextValue(t *testing.T) {
	if NextValue(0,0,0) != 0 {
		t.Error("NextValue fails on 0s")
	}

	if NextValue(6.0, 2.0, 4.0) != 6.0 {
		t.Error("NextValue fails on 6,2,4")
	}
}


var rhos = []float64{0.0, 0.00, 20.0, 0.00, 0.0}
var d0 =   []float64{0.0, 0.00,  0.0, 0.00, 0.0}
var d1 =   []float64{0.0, 0.00, 10.0, 0.00, 0.0}
var d2 =   []float64{0.0, 5.00, 10.0, 5.00, 0.0}
var d3 =   []float64{2.5, 5.00, 15.0, 5.00, 2.5}
var d4 =   []float64{2.5, 8.75, 15.0, 8.75, 2.5}

func dupeF64Slice(s []float64) []float64 {
	d := make([]float64,len(s))
	copy(d,s)
	return d
}

func eqF64Slices(s1, s2[]float64) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i,v := range s1 {
		if v != s2[i] {
			return false
		}
	}
	return true
}

func eqIntSlices(s1, s2[]int) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i,v := range s1 {
		if v != s2[i] {
			return false
		}
	}
	return true
}

func checkSliceSizes(t *testing.T, n, numSlices int, result ...int) {
	if v := SliceSizes(n,numSlices); !eqIntSlices(v,result) {
		t.Errorf("SliceSizes failed on %v %v, expected %v, got %v", n, numSlices,v, result)
	}
}


func TestSliceSizes(t *testing.T) {
	checkSliceSizes(t,  32, 4,  8,  8,  8,  8)
	checkSliceSizes(t,  34, 4,  9,  9,  8,  8)
	checkSliceSizes(t, 142, 7, 21, 21, 20, 20, 20, 20, 20)
}

func TestStepSlice(t *testing.T) {
	ins := [][]float64{d0,d1,d2,d3}
	outs := [][]float64{d1,d2,d3,d4} 
	for i,in := range ins {
		dup := dupeF64Slice(in)
		StepSlice(dup,rhos,0.0,0.0)
		if ! eqF64Slices(dup,outs[i]) {
			t.Errorf("StepSlice failed on entry #%v",i)
		}
	}
}

func testSolves1(t *testing.T, rhos []float64, numSlices, numIters int, answer []float64) {
	arr := make([]float64,len(rhos))
	startSlicing(arr, rhos, numSlices, numIters)
	if !eqF64Slices(arr,answer) {
		t.Errorf("Solver failed on %v items, %v slices, %v iters", len(rhos), numSlices, numIters)
	}
}

func TestSolves(t *testing.T) {
	//testSolves1(t,rhos,1,4,d4)
	testSolves1(t,rhos,2,4,d4)
}
