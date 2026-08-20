package main

import (
	"crypto/aes"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"strings"
)

const aesKeyLength = 16

// Challenge pairs one key-container string with the ciphertext it decrypts.
// KeyOffset is zero-based and points to the hidden 16-character AES-128 key.
type Challenge struct {
	KeyContainer string
	Ciphertext   string
	KeyOffset    int
}

// The strings come from Project I: AES Challenge in the Fall 2023 BS203
// Linear Algebra assignment. Each source string appears in exactly one pair.
var challenges = []Challenge{
	{KeyContainer: "81DA76C3C0C24F1E6F52D43909D7B09B", Ciphertext: "7425F291E2A123680787B92717383865", KeyOffset: 1},
	{KeyContainer: "137DEEA96A8BCBCD53BF0DBB426FF277", Ciphertext: "9767FABBD4EF9CB42281DB3E213C6739", KeyOffset: 6},
	{KeyContainer: "0C650471980296619760FD0A87194520", Ciphertext: "36E3D7C84F674A3703C03DCC46FA512C", KeyOffset: 1},
	{KeyContainer: "FD6CB42234261285888AF4E7CA7C6818", Ciphertext: "86E7FD9E6052FA335FD360D1BBC9EE3C", KeyOffset: 6},
	{KeyContainer: "EB4BA86D490188F8F214C45517C7A16B", Ciphertext: "15FDA2B1EED3B3502E15C674FEC1C750", KeyOffset: 15},
	{KeyContainer: "9F7151E17940B8D06027319E7F6D4F5B", Ciphertext: "3AAD1352AEDE2A33730B4DDEFC5A8889", KeyOffset: 14},
	{KeyContainer: "1FF6524F3B976F7D21AF8A12FA412472", Ciphertext: "023AE08E8DCE9F2A123CE6AC669FB203", KeyOffset: 14},
	{KeyContainer: "9DAD17C5BA473E79B23CD92EF0049FA6", Ciphertext: "04D59DB6A65796EDEECB069BA689A7C4", KeyOffset: 8},
}

func (c Challenge) actualKey() (string, error) {
	if c.KeyOffset < 0 || c.KeyOffset+aesKeyLength > len(c.KeyContainer) {
		return "", errors.New("hidden key is outside the key-container string")
	}
	return c.KeyContainer[c.KeyOffset : c.KeyOffset+aesKeyLength], nil
}

func (c Challenge) decrypt() (string, error) {
	key, err := c.actualKey()
	if err != nil {
		return "", err
	}

	ciphertext, err := hex.DecodeString(c.Ciphertext)
	if err != nil {
		return "", fmt.Errorf("decode ciphertext: %w", err)
	}
	if len(ciphertext) != aes.BlockSize {
		return "", fmt.Errorf("ciphertext has %d bytes; want %d", len(ciphertext), aes.BlockSize)
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", fmt.Errorf("create AES cipher: %w", err)
	}

	plaintext := make([]byte, aes.BlockSize)
	block.Decrypt(plaintext, ciphertext)
	return strings.TrimRight(string(plaintext), " \x00"), nil
}

type Round struct {
	ChallengeIndex int
	KeyPosition    int
}

func newDeck(random *rand.Rand) []Round {
	permutation := random.Perm(len(challenges))
	deck := make([]Round, len(permutation))
	for i, challengeIndex := range permutation {
		deck[i] = Round{
			ChallengeIndex: challengeIndex,
			KeyPosition:    random.Intn(2),
		}
	}
	return deck
}

func (r Round) strings() [2]string {
	challenge := challenges[r.ChallengeIndex]
	if r.KeyPosition == 0 {
		return [2]string{challenge.KeyContainer, challenge.Ciphertext}
	}
	return [2]string{challenge.Ciphertext, challenge.KeyContainer}
}
