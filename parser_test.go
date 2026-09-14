package halvector

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/koykov/byteconv"
)

type stage struct {
	hal    string
	expect string
	err    error
}

func loadTestDS() (stages []stage, err error) {
	files, err := filepath.Glob("testdata/*.hal.txt")
	if err != nil {
		return nil, err
	}

	fileExists := func(path string) bool {
		_, err := os.Stat(path)
		return err == nil
	}
	fileLoad := func(path string) (contents []byte, err error) {
		if contents, err = os.ReadFile(path); err != nil {
			return
		}
		if len(contents) > 0 && contents[len(contents)-1] == '\n' {
			contents = contents[:len(contents)-1]
		}
		return
	}

	stages = make([]stage, 0, len(files))
	for i := 0; i < len(files); i++ {
		var rawHAL, rawFmt, rawErr []byte
		if rawHAL, err = fileLoad(files[i]); err != nil {
			break
		}
		if ffmt := strings.Replace(files[i], "hal.txt", "fmt.json", 1); fileExists(ffmt) {
			if rawFmt, err = fileLoad(ffmt); err != nil {
				break
			}
		}
		if ferr := strings.Replace(files[i], "hal.txt", "err.txt", 1); fileExists(ferr) {
			if rawErr, err = fileLoad(ferr); err != nil {
				break
			}
		}
		stg := stage{
			hal:    string(rawHAL),
			expect: string(rawFmt),
			err:    errors.New(string(rawErr)),
		}
		stages = append(stages, stg)
	}
	return
}

func TestParser(t *testing.T) {
	stages, err := loadTestDS()
	if err != nil {
		t.Fatal(err)
	}
	for _, stg := range stages {
		t.Run(stg.hal, func(t *testing.T) {
			var buf bytes.Buffer
			vec := Acquire()
			if err := vec.ParseString(stg.hal); err != nil {
				if stg.err != nil {
					if stg.err.Error() != err.Error() {
						t.Error(err)
					}
				}
				return
			}
			_ = vec.Sort().Beautify(&buf)
			if stg.expect != buf.String() {
				t.Errorf("expect: %s\ngot: %s", stg.expect, buf.String())
			}
		})
	}
}

func BenchmarkParser(b *testing.B) {
	stages, err := loadTestDS()
	if err != nil {
		b.Fatal(err)
	}
	for _, stg := range stages {
		b.Run(stg.hal, func(b *testing.B) {
			b.ReportAllocs()
			var buf bytes.Buffer
			for i := 0; i < b.N; i++ {
				buf.Reset()
				vec := Acquire()
				if err := vec.ParseString(stg.hal); err != nil {
					if stg.err != nil {
						if stg.err.Error() != err.Error() {
							b.Error(err)
						}
					}
					return
				}
				_ = vec.Sort().Beautify(&buf)
				exp := buf.Bytes()
				if !bytes.Equal(byteconv.S2B(stg.expect), exp) {
					b.Errorf("expect: %s\ngot: %s", stg.expect, buf.String())
				}
				Release(vec)
			}
		})
	}
}

func BenchmarkLimit(b *testing.B) {
	longHAL, err := os.ReadFile("testdata/very-very-long.hal.txt")
	if err != nil {
		b.Fatal(err)
	}
	fn := func(b *testing.B, limit int) {
		b.ReportAllocs()
		vec := NewVector()
		for i := 0; i < b.N; i++ {
			_ = vec.SetLimit(limit).Parse(longHAL)
			vec.Reset()
		}
	}
	b.Run("no limit", func(b *testing.B) { fn(b, 0) })
	b.Run("limit 50", func(b *testing.B) { fn(b, 50) })
	b.Run("limit 10", func(b *testing.B) { fn(b, 10) })
	b.Run("limit 5", func(b *testing.B) { fn(b, 5) })
	b.Run("limit 3", func(b *testing.B) { fn(b, 3) })
}
