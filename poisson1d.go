package poisson1d

import (
	"fmt"
	"sync"
	)

// NextValue computes the next value at a site
func NextValue(rho, left, right float64) float64 {
	return 0.5 * (left + right + rho)
}

func StepSlice(arr, rhos []float64, lval, rval float64) {
	lastIndex := len(arr)-1
	lastVal := lval
	for i,currVal := range arr[:lastIndex] {
		arr[i] = NextValue(rhos[i],lastVal,arr[i+1])
		lastVal = currVal
	}
	arr[lastIndex] = NextValue(rhos[lastIndex], lastVal, rval)
}

func ProcessSlice(id int, arr, rhos []float64, lq, rq chan float64, niters int, lnbr, rnbr chan float64, wg *sync.WaitGroup) {
	fmt.Printf("Creating slice %v for %v iterations\n", id, niters)
	if lq == nil { fmt.Printf("Slice %v: no left queue\n", id) }
	if rq == nil { fmt.Printf("Slice %v: no right queue\n", id) }
	if lnbr == nil { fmt.Printf("Slice %v: no left neighbor\n", id) }
	if rnbr == nil { fmt.Printf("Slice %v: no right neighbor\n", id) }

	for i := 0; i<niters; i++ {

		fmt.Printf("Slice %v: Starting iteration %v\n",id,i)

		var lval, rval float64
		if lq != nil {
			fmt.Printf("Slice %v/%v: waiting on left\n",id,i)
			lval = <-lq
		}

		fmt.Printf("Slice %v/%v: received left = %v\n",id,i,lval)

		if rq != nil {
			fmt.Printf("Slice %v/%v: waiting on right\n",id,i)
			rval = <-rq
		}

		fmt.Printf("Slice %v/%v: received right = %v\n",id,i,rval)

		StepSlice(arr,rhos,lval,rval)

		if lnbr != nil {
			fmt.Printf("Slice %v/%v: sending %v to left neighbor\n",id,i,arr[0] )
			lnbr <- arr[0]
		} else {
			fmt.Printf("Slice %v/%v: sending 0.0 to own left\n",id,i)
			lq <- 0.0
		}
		if rnbr != nil {
			fmt.Printf("Slice %v/%v: sending %v to right neighbor\n",id,i,arr[len(arr)-1] )
			rnbr <- arr[len(arr)-1]
		} else {
			fmt.Printf("Slice %v/%v: sending 0.0 to own right\n",id,i)
			rq <- 0.0
		}
	}
	fmt.Printf("Slice %v: Done!\n",id)
	wg.Done()
}

func SliceSizes(n, numSlices int)  []int {
	d := n / numSlices
	r := n % numSlices
	s := make([]int,numSlices)
	for i,_ := range s {
		if i < r {
			s[i] = d+1
		} else {
			s[i] = d
		}
	}
	return s
}


func startSlicing(vals, rhos []float64, numSlices, numIters int) {
	if len(vals) != len(rhos) {
		panic("Must have same size value array as rhos")
	}
	fmt.Printf("Starting to slice, %v slices, %v iters, %v points\n",numSlices,numIters, len(rhos))
	lefts := make([]chan float64, numSlices)
	rights := make([]chan float64, numSlices)
	for i,_ := range lefts {
		lefts[i] = make(chan float64,2)
		rights[i] = make(chan float64,2)
	}
	var wg sync.WaitGroup
	sizes := SliceSizes(len(vals),numSlices)
	startIdx := 0
	for i,lc := range lefts {
		rc := rights[i]
		var lnbr, rnbr chan float64
		if i > 0 {
			lnbr = rights[i-1]
		}
		if i < len(lefts)-1 {
			rnbr = lefts[i+1]
		}
		wg.Add(1)
		fmt.Printf("Creating slice #%v [%v:%v]\n",i,startIdx,startIdx+sizes[i])
		go ProcessSlice(i,vals[startIdx:startIdx+sizes[i]], rhos, lc, rc, numIters, lnbr, rnbr, &wg) 
		lc <- 0.0
		rc <- 0.0
		startIdx += sizes[i]
	}
	wg.Wait()
	fmt.Printf("Main done!")
}


