package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {

	if err := runA("input.txt"); err != nil {
		fmt.Fprintf(os.Stderr, "puzzle errored: %v\n", err)
		os.Exit(1)
	}
	if err := runB("input.txt"); err != nil {
		fmt.Fprintf(os.Stderr, "puzzle errored: %v\n", err)
		os.Exit(1)
	}
}

type block struct {
	blkid int
	size  int
	moved bool
	nodes []*block
}

type disk []*block

func (d disk) String() string {
	buf := bytes.Buffer{}
	for _, blk := range d {
		buf.WriteString(blk.String())
	}
	return buf.String()
}

func (d disk) Flatten() disk {
	flat := disk{}
	for _, blk := range d {
		flat = append(flat, blk.flatten()...)
	}
	return flat
}

func (d disk) onlySpaceAtEnd() bool {
	spaceFound := false
	for _, blk := range d {
		if blk.blkid < 0 {
			spaceFound = true
			continue
		}
		if spaceFound {
			return false
		}
	}
	return true
}

func (d disk) Checksum() int {
	i, sum := 0, 0
	for _, blk := range d {
		if blk.blkid < 0 {
			i += blk.size
			continue
		}
		for range blk.size {
			sum += i * blk.blkid
			i++
		}
	}
	return sum
}

func (d disk) Defrag() disk {
	defrag := d
	itr := 0
	for {
		itr++
		defrag = defrag.Flatten()

		if defrag.onlySpaceAtEnd() {
			return defrag
		}

		// find first free space
		var free *block
		freeI := 0
		for i := range defrag {
			if defrag[i].blkid < 0 {
				free = defrag[i]
				freeI = i
				break
			}
		}
		if free == nil {
			// nothing free, stop
			return defrag
		}

		// find non-free from the back
		var nonFree *block
		for i := len(defrag) - 1; i >= freeI; i-- {
			if defrag[i].blkid >= 0 {
				nonFree = defrag[i]
				break
			}
		}
		if nonFree == nil {
			return defrag
		}

		// fmt.Println("filling free space at pos", freeI)

		if nonFree.size == free.size {
			free.blkid = nonFree.blkid
			nonFree.blkid = -1
			continue
		}

		if free.size > nonFree.size {
			free.nodes = []*block{
				{blkid: nonFree.blkid, size: nonFree.size},
				{blkid: -1, size: free.size - nonFree.size},
			}
			nonFree.blkid = -1
			continue
		}

		if free.size < nonFree.size {
			free.blkid = nonFree.blkid
			nonFree.nodes = []*block{
				{blkid: nonFree.blkid, size: nonFree.size - free.size},
				{blkid: -1, size: free.size},
			}
			continue
		}
	}
}

func (d disk) FitFiles() disk {
	defrag := d.Flatten()
	itr := 0
	nonFreeIdx := len(defrag)
	for {
		itr++
		defrag = defrag.Flatten()

		// find non-free from the back
		var nonFree *block
		for i := nonFreeIdx - 1; i >= 0; i-- {
			if defrag[i].blkid >= 0 && !defrag[i].moved {
				nonFreeIdx = i
				nonFree = defrag[i]
				break
			}
		}
		if nonFree == nil {
			return defrag
		}
		// fmt.Println("fitting file", nonFree.blkid, itr)

		// mark file as checked
		nonFree.moved = true

		// find first free space that the file fits in
		for i := 0; i < nonFreeIdx; i++ {
			// storage block
			if defrag[i].blkid >= 0 {
				continue
			}
			// not enought space
			if defrag[i].size < nonFree.size {
				continue
			}

			// exact fit
			if nonFree.size == defrag[i].size {
				defrag[i].blkid = nonFree.blkid
				defrag[i].moved = true
				nonFree.blkid = -1
				break
			}

			// more space
			if defrag[i].size > nonFree.size {
				defrag[i].nodes = []*block{
					{blkid: nonFree.blkid, size: nonFree.size, moved: true},
					{blkid: -1, size: defrag[i].size - nonFree.size},
				}
				nonFree.blkid = -1
				break
			}
			// else leave file where it is
		}
	}
}

func (b *block) flatten() []*block {
	if len(b.nodes) == 0 {
		return []*block{b}
	}
	blks := []*block{}
	for _, blk := range b.nodes {
		blks = append(blks, blk.flatten()...)
	}
	return blks
}

func (b *block) String() string {
	// super block
	if len(b.nodes) > 0 {
		buf := bytes.Buffer{}
		for _, n := range b.nodes {
			buf.WriteString(n.String())
		}
		return buf.String()
	}
	// free block
	if b.blkid < 0 {
		return strings.Repeat(".", b.size)
	}

	// storage block
	return strings.Repeat(strconv.Itoa(b.blkid), b.size)
}

func runA(file string) error {
	layout, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	d := disk{}
	id := 0
	for i, b := range layout {
		blkId := -1
		// its a block
		if i%2 == 0 {
			blkId = id
			id++
		}
		d = append(d, &block{blkid: blkId, size: int(b - '0')})
	}

	d = d.Defrag()
	fmt.Println(d.Checksum())
	return nil
}

func runB(file string) error {
	layout, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	d := disk{}
	id := 0
	for i, b := range layout {
		blkId := -1
		// its a block
		if i%2 == 0 {
			blkId = id
			id++
		}
		d = append(d, &block{blkid: blkId, size: int(b - '0')})
	}

	d = d.FitFiles()
	fmt.Println(d.Checksum())
	return nil
}
