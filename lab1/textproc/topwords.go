// Find the top K most common words in a text document.
// Input path: location of the document, K top words
// Output: Slice of top K words
// For this excercise, word is defined as characters separated by a whitespace

// Note: You should use `checkError` to handle potential errors.

package textproc

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
)

func topWords(path string, K int) []WordCount {
	// Erwin
	dat, err := os.Open(path)
	checkError(err)

	scanner := bufio.NewScanner(dat)
	scanner.Split(bufio.ScanWords)

	// Jessica
	// Map
	m := make(map[string]int)

	// Passing words into map
	for scanner.Scan() {
		//fmt.Println(scanner.Text())
		m[scanner.Text()]++
	}

	// Cole
	// Slice
	WordCounts := make([]WordCount, 0, len(m))

	// Adds words and their counts to slice
	for key, val := range m {
		WordCounts = append(WordCounts, WordCount{Word: key, Count: val})
	}

	// Sort WordCount
	sortWordCounts(WordCounts)

	// Will
	// Removes elements past index K
	// Can also do: return WordCounts[:K] to return the first K elements

	/*for len(WordCounts) > K {
		WordCounts = WordCounts[:len(WordCounts)-1]
	}

	// Prints Top K words for verification
	for i := 0; i < K; i++ {
		fmt.Println(WordCounts[i].Word, ":", WordCounts[i].Count)
	}
	*/

	return WordCounts[:K]
}

//--------------- DO NOT MODIFY----------------!

// A struct that represents how many times a word is observed in a document
type WordCount struct {
	Word  string
	Count int
}

// Method to convert struct to string format

func (wc WordCount) String() string {
	return fmt.Sprintf("%v: %v", wc.Word, wc.Count)
}

// Helper function to sort a list of word counts in place.
// This sorts by the count in decreasing order, breaking ties using the word.

func sortWordCounts(wordCounts []WordCount) {
	sort.Slice(wordCounts, func(i, j int) bool {
		wc1 := wordCounts[i]
		wc2 := wordCounts[j]
		if wc1.Count == wc2.Count {
			return wc1.Word < wc2.Word
		}
		return wc1.Count > wc2.Count
	})
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
