package main

import "encoding/json"

func ToJson(j any, pretty bool) (string, error) {
	var str []byte

	var err error

	if pretty {
		str, err = json.MarshalIndent(j, "", "\t")
	} else {
		str, err = json.Marshal(j)
	}

	if err != nil {
		return "", err
	}

	return string(str), nil
}
