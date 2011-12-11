package argcfg

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"testing"
)

type SubConfig struct {
	G   int
	Str string
}

type Config struct {
	F float64 "something"
	S SubConfig
}

func TestCFG(t *testing.T) {
	defer func() {
		e := recover()
		if e != nil {
			fmt.Printf("%v\n", e)
			for skip := 1; ; skip++ {
				_, file, line, ok := runtime.Caller(skip)
				if !ok {
					break
				}
				if file[len(file)-1] == 'c' {
					continue
				}
				fmt.Printf("%s:%d\n", file, line)
			}
		}
	}()
	var c Config
	c.F = 2
	c.S.Str = "Hello!"
	Usage(os.Stdout, &c)
	err := LoadArg("-F=.75", &c)
	if err != nil {
		t.Error(err)
	}
	err = LoadArg("-S.G=2", &c)
	if err != nil {
		t.Error(err)
	}
	v, _ := strconv.ParseFloat(".75", 64)
	if c.F != v {
		t.Fail()
	}
	if c.S.G != 2 {
		t.Fail()
	}
}
