package models

import (
	"testing"
)

func BenchmarkHashPassword(b *testing.B) {
	user := &User{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user.HashPassword("mypassword123")
	}
}

func BenchmarkVerifyPassword_Correct(b *testing.B) {
	user := &User{}
	user.HashPassword("mypassword123")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user.VerifyPassword("mypassword123")
	}
}

func BenchmarkVerifyPassword_Wrong(b *testing.B) {
	user := &User{}
	user.HashPassword("mypassword123")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user.VerifyPassword("wrongpassword")
	}
}
