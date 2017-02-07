// +build ignore

package main

import (
	"fmt"

	"github.com/go-hep/croot"
)

const ARRAYSZ = 10

type Event struct {
	I32 int32
	I64 int64
	U32 uint32
	U64 uint64
	F32 float32
	F64 float64

	ArrayI32 [ARRAYSZ]int32
	ArrayI64 [ARRAYSZ]int64
	ArrayU32 [ARRAYSZ]uint32
	ArrayU64 [ARRAYSZ]uint64
	ArrayF32 [ARRAYSZ]float32
	ArrayF64 [ARRAYSZ]float64
}

func main() {
	const fname = "test-small.root"
	const evtmax = 100
	const splitlevel = 32
	const bufsiz = 32000
	const compress = 1
	const netopt = 0

	f, err := croot.OpenFile(fname, "recreate", "small event file", compress, netopt)
	if err != nil {
		panic(err.Error())
	}

	// create a tree
	tree := croot.NewTree("tree", "tree", splitlevel)

	e := Event{}

	_, err = tree.Branch("evt", &e, bufsiz, 0)
	if err != nil {
		panic(err.Error())
	}

	// fill some events with random numbers
	for iev := int64(0); iev != evtmax; iev++ {
		if iev%1000 == 0 {
			fmt.Printf(":: processing event %d...\n", iev)
		}

		e.I32 = int32(iev)
		e.I64 = int64(iev)
		e.U32 = uint32(iev)
		e.U64 = uint64(iev)
		e.F32 = float32(iev)
		e.F64 = float64(iev)

		for ii := 0; ii < ARRAYSZ; ii++ {
			e.ArrayI32[ii] = int32(iev)
			e.ArrayI64[ii] = int64(iev)
			e.ArrayU32[ii] = uint32(iev)
			e.ArrayU64[ii] = uint64(iev)
			e.ArrayF32[ii] = float32(iev)
			e.ArrayF64[ii] = float64(iev)
		}

		_, err = tree.Fill()
		if err != nil {
			panic(err.Error())
		}
	}
	f.Write("", 0, 0)
	f.Close("")

}

// EOF
