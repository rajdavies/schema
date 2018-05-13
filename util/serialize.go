package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/clbanning/mxj"
	"github.com/naoina/toml"
	"gopkg.in/yaml.v2"
)

func MuxDecode(r io.Reader) (map[string]interface{}, error) {
	// ReadAll is necessary, since the input stream could be only
	// traversable once; we must be sure to save the data
	// into a buffer on the first pass, so that we can read it
	// *multiple* times
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	if err = json.Unmarshal(buf, &data); err == nil {
		return data, nil
	}

	data = make(map[string]interface{}) // be sure it's an empty map
	if err = yaml.Unmarshal(buf, &data); err == nil {
		return data, nil
	}

	data = make(map[string]interface{}) // be sure it's an empty map
	if err = toml.Unmarshal(buf, &data); err == nil {
		return data, nil
	}

	mv, err := mxj.NewMapXmlReader(bytes.NewReader(buf))
	if err == nil {
		return map[string]interface{}(mv), nil
	}

	return nil, fmt.Errorf("failed to recognize input data format")
}

type OutFmt string

func DemuxEncode(w io.Writer, data interface{}, outFmt OutFmt, doPretty bool) error {
	switch string(outFmt) {
	case "json":
		enc := json.NewEncoder(w)
		if doPretty {
			enc.SetIndent("", "    ")
		}
		if err := enc.Encode(&data); err != nil {
			return err
		}
	case "yaml":
		if err := yaml.NewEncoder(w).Encode(&data); err != nil {
			return err
		}
	case "toml":
		if err := toml.NewEncoder(w).Encode(&data); err != nil {
			return err
		}
	case "xml", "protobuf", "graphql":
		return fmt.Errorf("'%v' is not implemented yet", outFmt)
	default:
		return fmt.Errorf("unrecognized output format '%v'", outFmt)
	}
	return nil
}