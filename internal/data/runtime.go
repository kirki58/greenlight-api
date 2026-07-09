package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Runtime int32

func (r Runtime) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, "\"%d mins\"", r), nil
}

var InvalidRuntimeFieldFormat error = errors.New("json: the runtime field should be in the format \"runtime-int mins\"")

func (r *Runtime) UnmarshalJSON(jsonVal []byte) error {
	unquotedStr, err := strconv.Unquote(string(jsonVal))
	if err != nil {
		return InvalidRuntimeFieldFormat
	}

	strs := strings.Split(unquotedStr, " ")
	if len(strs) != 2 || strs[1] != "mins" {
		return InvalidRuntimeFieldFormat
	}

	runtimeVal, err := strconv.ParseInt(strs[0], 10, 32)
	if err != nil {
		return InvalidRuntimeFieldFormat
	}

	*r = Runtime(runtimeVal)
	return nil
}
