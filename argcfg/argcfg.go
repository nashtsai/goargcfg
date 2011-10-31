package argcfg

import (
	"reflect"
	"os"
	"io"
	"fmt"
	"strings"
	"strconv"
)

func Usage(out io.Writer, cfg interface{}) {
	obj := reflect.ValueOf(cfg)
	if obj.Kind() == reflect.Ptr {
		usageAux(out, "", obj.Elem())
	}
}

func usageAux(out io.Writer, prefix string, obj reflect.Value) {

	if obj.Kind() == reflect.Struct {
		typ := obj.Type()
		for i := 0; i < obj.NumField(); i++ {
			ftyp := typ.Field(i)

			if ftyp.PkgPath != "" {
				continue
			}

			fobj := obj.Field(i)
			name := ftyp.Name
			if prefix != "" {
				name = fmt.Sprintf("%s.%s", prefix, ftyp.Name)
			}

			if fobj.Kind() == reflect.Struct {
				usageAux(out, name, fobj)
			} else {
				valstr := fmt.Sprintf("%v", fobj.Interface())
				if fobj.Kind() == reflect.String {
					valstr = fmt.Sprintf("%q", fobj.Interface())
				}
				fmt.Fprintf(out, "  -%s=%s (%T)", name, valstr, fobj.Interface())
				if ftyp.Tag != "" {
					fmt.Fprintf(out, ": %s", ftyp.Tag)
				}
				fmt.Fprintf(out, "\n")
			}
		}
	}
}

func LoadArgs(cfg interface{}) (ok bool, err os.Error) {
	for _, arg := range os.Args {
		if arg == "-?" || arg == "-help" || arg == "--help" {
			Usage(os.Stderr, cfg)
			ok = false
			return
		}
		err = LoadArg(arg, cfg)
		if err != nil {
			//fmt.Fprintf(os.Stderr, "Error for \"%v\": %v\n", arg, err)
			//os.Exit(1)
			return
		}
	}
	ok = true
	return
}

func LoadArg(arg string, cfg interface{}) (err os.Error) {
	if arg[0] != '-' {
		return
	}
	arg = arg[1:len(arg)]
	tokens := strings.Split(arg, "=")
	key, val := tokens[0], tokens[1]
	keys := strings.Split(key, ".")

	v := reflect.ValueOf(cfg)
	if v.Kind() != reflect.Ptr {
		return os.NewError(fmt.Sprintf("%v is not a pointer", cfg))
	}

	err = LoadKeysVal(keys, val, v.Elem())

	return
}

func LoadKeysVal(keys []string, val string, objValue reflect.Value) (err os.Error) {
	if len(keys) == 0 {
		//we're here - dump val onto obj

		if !objValue.CanSet() {
			err = os.NewError("Attempting to set an unexported field")
		}

		switch objValue.Kind() {
		case reflect.Float32, reflect.Float64:
			var tval float64
			tval, err = strconv.Atof64(val)
			if err != nil {
				return
			}
			objValue.SetFloat(tval)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			var tval int64
			tval, err = strconv.Atoi64(val)
			if err != nil {
				return
			}
			objValue.SetInt(tval)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			var tval uint64
			tval, err = strconv.Atoui64(val)
			if err != nil {
				return
			}
			objValue.SetUint(tval)
		case reflect.Bool:
			var tval bool
			tval, err = strconv.Atob(val)
			if err != nil {
				return
			}
			objValue.SetBool(tval)
		case reflect.String:
			objValue.SetString(val)
		}
		return
	}

	//otherwise obj needs to be a struct
	if objValue.Kind() != reflect.Struct {
		err = os.NewError("not a struct")
		return
	}

	subValue := objValue.FieldByName(keys[0])
	if !subValue.IsValid() {
		return os.NewError("Invalid field")
	}

	err = LoadKeysVal(keys[1:len(keys)], val, subValue)

	return
}
