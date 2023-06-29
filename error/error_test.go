package error

import (
	"fmt"
	"os"
	"testing"
)

func Benchmark_ErrorChecking(b *testing.B) {
	const runs = 5000
	_, err := os.Stat("error.go")

	b.Run("Error checking via func", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for r := 0; r < runs; r++ {
				Must(err)
			}
		}
	})

	b.Run("Error checking in-line", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for r := 0; r < runs; r++ {
				if err != nil {
					fmt.Println("Error!")
				}
			}
		}
	})
}
