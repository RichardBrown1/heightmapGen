package main

import (
	"flag"
	"os"
	"testing"
)

func BenchmarkMain(b *testing.B) {
	os.Args[1] = "/home/richard/Documents/goProjects/heightmapGen/data/"
	flag.Parse()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		main()
	}
}
