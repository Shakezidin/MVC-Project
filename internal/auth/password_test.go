package auth

import "testing"

func TestHashAndCheckPassword(t *testing.T) {
	password := "password123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if !CheckPassword(hash, password) {
		t.Error("CheckPassword should return true for correct password")
	}

	if CheckPassword(hash, "wrongpassword") {
		t.Error("CheckPassword should return false for incorrect password")
	}
}

func TestHashPassword_DifferentHashes(t *testing.T) {
	password := "password123"

	hash1, _ := HashPassword(password)
	hash2, _ := HashPassword(password)

	if hash1 == hash2 {
		t.Error("bcrypt should produce different hashes due to random salt")
	}
}
