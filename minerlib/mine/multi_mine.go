// TODO swap map for sync.map
package mine

import (
	// "fmt"
	"sync"
)

type nonceHashPair struct {
	nonce uint32
	hash string
}

// signal channels
var doneTesting chan struct{}
var doneHashing chan struct{}
var finished chan struct{}

// pipeline channels
var newNonces chan uint32
var hashes chan nonceHashPair
var final chan uint32

// waitgroups
var generators sync.WaitGroup
var hashers sync.WaitGroup
var testers sync.WaitGroup

var numConcurrent = 6
var UsedNonces sync.Map

func MultiMine(data []byte, difficulty uint8) (result uint32, err error) {
	doneTesting = make(chan struct{})
	doneHashing = make(chan struct{})
	finished = make(chan struct{})
	newNonces = make(chan uint32, numConcurrent * 1000)
	hashes = make(chan nonceHashPair, numConcurrent * 1000)
	final = make(chan uint32, 1)
	UsedNonces = sync.Map{}
	// spawn
	for i:= 0; i < numConcurrent; i++ {
		go GenerateValues()
		generators.Add(1)
		go HashValues(data)
		hashers.Add(1)
		go TestValues(difficulty)
		testers.Add(1)
	}
	// wait
	result = <- final
	// fmt.Printf("MAIN: Got Result: %d\n", result)
	// clean up
	testers.Wait()
	close(doneTesting)
	// fmt.Println("MAIN: Done Testing")
	hashers.Wait()
	close(doneHashing)
	// fmt.Println("MAIN: Done Hashing")
	generators.Wait()
	// fmt.Println("MAIN: Done Generating")
	// exit
	return result, nil
}




func GenerateValues() {
	for {
		if finishedHashing() {
			break
		}
		n, err := concurrentNewNonce()
		if err != nil {
			break
		}
		// fmt.Printf("Generated %d\n", n)
		newNonces <- n
	}
	generators.Done()
}

func HashValues(data []byte) {
	for {
		if finishedTesting() {
			break
		}
		nonce := <- newNonces
		hash := MD5Hash(data, nonce)
		// fmt.Printf("Hashing %d into %s\n", nonce, hash)
		hashes <- nonceHashPair{nonce: nonce, hash: hash}
	}
	hashers.Done()
}

func TestValues(difficulty uint8) {
	for {
		if finalFinish() {
			break
		}
		pair := <- hashes
		// fmt.Printf("Checking %v\n", pair)
		if Valid(pair.hash, difficulty) {
			// fmt.Printf("Found One!\n")
			if !finalFinish() {
				final <- pair.nonce
				// fmt.Printf("Sent %d\n", pair.nonce)
				// fmt.Printf("Done!\n")
				close(finished)
			}
			break
		}
	}
	testers.Done()
}

// Helpers

func concurrentNewNonce() (n uint32, err error) {
	for  {
		n = randomNonce()
		used, ok := UsedNonces.Load(n)
		if !ok || used == nil {
			UsedNonces.Store(n, true)
			return n, nil
		}
	}
}

func finalFinish() bool {
	select {
	case <- finished:
		return true
	default:
		return false
	}
}

func finishedTesting() bool {
	select {
	case <- doneTesting:
		return true
	default:
		return false
	}
}

func finishedHashing() bool {
	select {
	case <- doneHashing:
		return true
	default:
		return false
	}
}
