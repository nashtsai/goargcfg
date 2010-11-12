package cfg

import (
	"testing"
	"fmt"
	"runtime"
)

type SubConfig struct {
	G int
}

type Config struct {
	F float
	S SubConfig
}

func TestCFG(t *testing.T) {
	defer func() {
		e := recover()
		if e != nil {
			fmt.Printf("%v\n", e)
			for skip:=1; ; skip++ {
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
	err := LoadArg("-F=.75", &c)
	err = LoadArg("-S.G=2", &c)
	fmt.Printf("err = %v\n", err)
	println(c.F)
	println(c.S.G)
}