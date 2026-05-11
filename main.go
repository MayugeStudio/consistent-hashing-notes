package main

import (
	"fmt"
	"hash/fnv"
	"hash/crc32"
)


const (
	N = 50    // the number of the servers
	// M = 65393 // the length of the look-up-table
	M = 1_000_003
)

type Permutations [N][M]int // all(permutations[i], (e) => e == -1) == true means the server is down
type LookUpTable [M]string

// Hash functions
func hashStringToInt_1(s string) uint32 {
    h := fnv.New32a()
    h.Write([]byte(s))
    return h.Sum32()
}

func hashStringToInt_2(s string) uint32 {
	return crc32.ChecksumIEEE([]byte(s))
}


func GeneratePermutationList(names [N]string) Permutations {
	var permutations Permutations
	for i := 0; i<N; i++ {
		name := names[i]
		offset := int(hashStringToInt_1(name) % M)
		skip := int(hashStringToInt_2(name) % (M-1)) + 1
		for j := 0; j<M; j++ {
			permutations[i][j] = (offset + j*skip) % M
		}
	}
	return permutations
}

func PopulateLookUpTable(names [N]string, permutations Permutations) [M]int {
	var next [N]int
	for i := 0; i<N; i++ {
		next[i] = 0
	}

	var entry [M]int
	for i := 0; i<M; i++ {
		entry[i] = -1
	}

	n := 0

	for {
		for i := 0; i<N; i++ {
			c := permutations[i][next[i]]
			if c == -1 {
				continue
			}

			for entry[c] >= 0 {
				next[i] += 1
				c = permutations[i][next[i]]
			}
			entry[c] = i
			next[i] += 1
			n += 1
			if n == M {
				return entry
			}
		}
	}
}

func MarkServerAsDown(index int, permutations *Permutations) {
	for i := 0; i<M; i++ {
		permutations[index][i] = -1
	}
}


func main() {
	var names [N]string
	for i := 0; i<N; i++ {
		names[i] = fmt.Sprintf("SERVER%d", i)
	}

	permutations := GeneratePermutationList(names)
	
	lookUpTable := PopulateLookUpTable(names, permutations)
	MarkServerAsDown(0, &permutations) // Mark server 0 as down
	lookUpTable2 := PopulateLookUpTable(names, permutations)

	same := 0
	diff := 0

	for j := 0; j < M; j++ {
		if lookUpTable[j] == lookUpTable2[j] {
			same += 1
		} else {
			diff += 1
		}
	}

	fmt.Printf("same = %d, diff = %d, ratio = %f\n", same, diff, float64(same)/float64(same+diff))
}
