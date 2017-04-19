// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lcio

import (
	"bytes"
	"fmt"
	"strings"

	"go-hep.org/x/hep/sio"
)

type CalorimeterHits struct {
	Flags  Flags
	Params Params
	Hits   []CalorimeterHit
}

func (hits CalorimeterHits) String() string {
	o := new(bytes.Buffer)
	fmt.Fprintf(o, "%[1]s print out of CalorimeterHit collection %[1]s\n\n", strings.Repeat("-", 15))
	fmt.Fprintf(o, "  flag:  0x%x\n%v", hits.Flags, hits.Params)
	fmt.Fprintf(o, "  -> LCIO::RCHBIT_LONG   : %v\n", hits.Flags.Test(RChBitLong))
	fmt.Fprintf(o, "     LCIO::RCHBIT_BARREL : %v\n", hits.Flags.Test(RChBitBarrel))
	fmt.Fprintf(o, "     LCIO::RCHBIT_ID1    : %v\n", hits.Flags.Test(RChBitID1))
	fmt.Fprintf(o, "     LCIO::RCHBIT_TIME   : %v\n", hits.Flags.Test(RChBitTime))
	fmt.Fprintf(o, "     LCIO::RCHBIT_NO_PTR : %v\n", hits.Flags.Test(RChBitNoPtr))
	fmt.Fprintf(o, "     LCIO::RCHBIT_ENERGY_ERROR  : %v\n", hits.Flags.Test(RChBitEnergyError))

	// FIXME(sbinet): CellIDDecoder

	fmt.Fprintf(o, "\n")

	head := " [   id   ] |cellId0 |cellId1 |  energy  |energyerr |        position (x,y,z)           \n"
	tail := "------------|--------|--------|----------|----------|-----------------------------------\n"
	fmt.Fprintf(o, head)
	fmt.Fprintf(o, tail)
	for _, hit := range hits.Hits {
		fmt.Fprintf(o, " [%08d] |%08d|%08d|%+.3e|%+.3e|", 0, hit.CellID0, hit.CellID1, hit.Energy, hit.EnergyErr)
		if hits.Flags.Test(ChBitLong) {
			fmt.Fprintf(o, "+%.3e, %+.3e, %+.3e", hit.Pos[0], hit.Pos[1], hit.Pos[2])
		} else {
			fmt.Fprintf(o, "    no position available         ")
		}
		// FIXME(sbinet): CellIDDecoder
		fmt.Fprintf(o, "\n        id-fields: --- unknown/default ----   ")
		fmt.Fprintf(o, "\n")
	}
	fmt.Fprintf(o, tail)
	return string(o.Bytes())
}

func (*CalorimeterHits) VersionSio() uint32 {
	return Version
}

func (hits *CalorimeterHits) MarshalSio(w sio.Writer) error {
	panic("not implemented")
}

func (hits *CalorimeterHits) UnmarshalSio(r sio.Reader) error {
	dec := sio.NewDecoder(r)
	dec.Decode(&hits.Flags)
	dec.Decode(&hits.Params)
	var n int32
	dec.Decode(&n)
	hits.Hits = make([]CalorimeterHit, int(n))
	for i := range hits.Hits {
		hit := &hits.Hits[i]
		dec.Decode(&hit.CellID0)
		if r.VersionSio() == 8 || hits.Flags.Test(RChBitID1) {
			dec.Decode(&hit.CellID1)
		}
		dec.Decode(&hit.Energy)
		if r.VersionSio() > 1009 && hits.Flags.Test(RChBitEnergyError) {
			dec.Decode(&hit.EnergyErr)
		}
		if r.VersionSio() > 1002 && hits.Flags.Test(RChBitTime) {
			dec.Decode(&hit.Time)
		}
		if hits.Flags.Test(RChBitLong) {
			dec.Decode(&hit.Pos)
		}
		if r.VersionSio() > 1002 {
			dec.Decode(&hit.Type)
			dec.Pointer(&hit.Raw)
		}
		if r.VersionSio() > 1002 {
			// the logic of the pointer bit has been inverted in v1.3
			if !hits.Flags.Test(RChBitNoPtr) {
				dec.Tag(hit)
			}
		} else {
			if hits.Flags.Test(RChBitNoPtr) {
				dec.Tag(hit)
			}
		}
	}
	return dec.Err()
}

type CalorimeterHit struct {
	CellID0   int32
	CellID1   int32
	Energy    float32
	EnergyErr float32
	Time      float32
	Pos       [3]float32
	Type      int32
	Raw       *RawCalorimeterHit
}

var _ sio.Codec = (*CalorimeterHits)(nil)