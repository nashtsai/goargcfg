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
	
	v := reflect.NewValue(cfg)
	if v.Kind() != reflect.Ptr {
		return os.NewError(fmt.Sprintf("%v is not a pointer", cfg))
	}
	objValue := v.Elem()
	
	//ptrValue := reflect.NewValue(cfg).(*reflect.PtrValue)
	//objValue := ptrValue.Elem()
	
	err = LoadKeysVal(keys, val, objValue)
	
	return
}

func LoadKeysVal(keys []string, val string, objValue reflect.Value) (err os.Error) {
	//objType := objValue.Type()	
	
	if len(keys) == 0 {
	
		switch objValue.Kind() {
		case reflect.Float32, reflect.Float64:
			tval, err := strconv.Atof64(val)
			if err != nil {
				return
			}
			objValue.SetFloat(tval)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			tval, err := strconv.Atoi64(val)
			if err != nil {
				return
			}
			objValue.SetInt(tval)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			tval, err := strconv.Atoui64(val)
			if err != nil {
				return
			}
			objValue.SetUint(tval)
		case reflect.Bool:
			tval, err := strconv.Atob(val)
			if err != nil {
				return
			}
			objValue.SetBool(tval)
		case reflect.String:
			objValue.SetString(val)
		}
		/*
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
		*/
		return
	}
	
	//otherwise obj needs to be a struct
	
	if objValue.Kind() != reflect.Struct {
		err = os.NewError("not a struct")
		return
	}
	
	
	subValue := objValue.FieldByName(keys[0])
	if !subValue.IsValid() {
		return os.NewError("nil, somehow")
	}
	
	err = LoadKeysVal(keys[1:len(keys)], val, subValue)

	return
}
