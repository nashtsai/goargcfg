package argcfg

import (
	"reflect"
	"os"
	"fmt"
	"strings"
	"strconv"
)

func LoadArgs(cfg interface{}) (err os.Error) {
	for _, arg := range os.Args {
		err = LoadArg(arg, cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error for \"%v\": %v\n", arg, err)
			os.Exit(1)
			return
		}
	}
	return
}

func LoadArg(arg string, cfg interface{}) (err os.Error) {
	if arg[0] != '-' {
		return
	}
	arg = arg[1:len(arg)]
	tokens := strings.Split(arg, "=", -1)
	key, val := tokens[0], tokens[1]
	keys := strings.Split(key, ".", -1)
	
	ptrValue := reflect.NewValue(cfg).(*reflect.PtrValue)
	objValue := ptrValue.Elem()
	
	err = LoadKeysVal(keys, val, objValue)
	
	return
}

func LoadKeysVal(keys []string, val string, objValue reflect.Value) (err os.Error) {
	objType := objValue.Type()	
	
	if len(keys) == 0 {
		//we're here - dump val onto obj
		switch fieldType := objType.(type) {
		case *reflect.FloatType:
			tval, err := strconv.Atof64(val)
			if err != nil {
				return
			}
			objValue.(*reflect.FloatValue).Set(float64(tval))
		case *reflect.IntType:
			tval, err := strconv.Atoi(val)
			if err != nil {
				return
			}
			objValue.(*reflect.IntValue).Set(int64(tval))
		case *reflect.UintType:
			tval, err := strconv.Atoi(val)
			if err != nil {
				return
			}
			objValue.(*reflect.UintValue).Set(uint64(tval))
		case *reflect.BoolType:
			tval, err := strconv.Atob(val)
			if err != nil {
				return
			}
			objValue.(*reflect.BoolValue).Set(tval)
		case *reflect.StringType:
			tval := val
			if err != nil {
				return
			}
			objValue.(*reflect.StringValue).Set(tval)
		}
		return
	}
	
	//otherwise obj needs to be a struct
	
	structValue, ok := objValue.(*reflect.StructValue)
	
	if !ok {
		err = os.NewError("not a struct")
		return
	}
	
	subValue := structValue.FieldByName(keys[0])
	if subValue == nil {
		return os.NewError("nil, somehow")
	}
	
	err = LoadKeysVal(keys[1:len(keys)], val, subValue)

	return
}
