package main

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestAssignmentPairsDecryptToMessages(t *testing.T) {
	want := []string{
		"Apple",
		"Orange",
		"Pineapple",
		"Cherry",
		"Grape",
		"Tangerine",
		"Watermelon",
		"Peach",
	}

	got := make([]string, 0, len(challenges))
	seen := make(map[string]bool)
	for i, challenge := range challenges {
		if seen[challenge.KeyContainer] || seen[challenge.Ciphertext] {
			t.Fatalf("challenge %d reuses a source string", i)
		}
		seen[challenge.KeyContainer] = true
		seen[challenge.Ciphertext] = true

		plaintext, err := challenge.decrypt()
		if err != nil {
			t.Fatalf("challenge %d: %v", i, err)
		}
		got = append(got, plaintext)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("decrypted messages = %q; want %q", got, want)
	}
}

func TestNewDeckUsesEveryPairOnce(t *testing.T) {
	deck := newDeck(rand.New(rand.NewSource(23)))
	if len(deck) != len(challenges) {
		t.Fatalf("deck has %d rounds; want %d", len(deck), len(challenges))
	}

	seen := make(map[int]bool)
	for _, round := range deck {
		if round.KeyPosition < 0 || round.KeyPosition > 1 {
			t.Fatalf("invalid key position %d", round.KeyPosition)
		}
		if seen[round.ChallengeIndex] {
			t.Fatalf("challenge %d appears twice", round.ChallengeIndex)
		}
		seen[round.ChallengeIndex] = true
	}
}

func TestRoundStringsPlaceKeyAtRecordedPosition(t *testing.T) {
	for challengeIndex, challenge := range challenges {
		for keyPosition := 0; keyPosition < 2; keyPosition++ {
			round := Round{ChallengeIndex: challengeIndex, KeyPosition: keyPosition}
			items := round.strings()
			if items[keyPosition] != challenge.KeyContainer {
				t.Fatalf("challenge %d: key is not at position %d", challengeIndex, keyPosition)
			}
			if items[1-keyPosition] != challenge.Ciphertext {
				t.Fatalf("challenge %d: ciphertext is not opposite the key", challengeIndex)
			}
		}
	}
}
