// Code generated by github.com/actgardner/gogen-avro/v10. DO NOT EDIT.
/*
 * SOURCE:
 *     cost.avsc
 */
package avro

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/actgardner/gogen-avro/v10/compiler"
	"github.com/actgardner/gogen-avro/v10/vm"
	"github.com/actgardner/gogen-avro/v10/vm/types"
)

var _ = fmt.Printf

type Resource struct {
	Display_name string `json:"display_name"`

	Id string `json:"id"`

	Environment Environment `json:"environment"`
}

const ResourceAvroCRC64Fingerprint = "вhs\x05\xd2%M"

func NewResource() Resource {
	r := Resource{}
	r.Environment = NewEnvironment()

	return r
}

func DeserializeResource(r io.Reader) (Resource, error) {
	t := NewResource()
	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func DeserializeResourceFromSchema(r io.Reader, schema string) (Resource, error) {
	t := NewResource()

	deser, err := compiler.CompileSchemaBytes([]byte(schema), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func writeResource(r Resource, w io.Writer) error {
	var err error
	err = vm.WriteString(r.Display_name, w)
	if err != nil {
		return err
	}
	err = vm.WriteString(r.Id, w)
	if err != nil {
		return err
	}
	err = writeEnvironment(r.Environment, w)
	if err != nil {
		return err
	}
	return err
}

func (r Resource) Serialize(w io.Writer) error {
	return writeResource(r, w)
}

func (r Resource) Schema() string {
	return "{\"fields\":[{\"name\":\"display_name\",\"type\":\"string\"},{\"name\":\"id\",\"type\":\"string\"},{\"name\":\"environment\",\"type\":{\"fields\":[{\"name\":\"id\",\"type\":\"string\"}],\"name\":\"environment\",\"type\":\"record\"}}],\"name\":\"resource\",\"type\":\"record\"}"
}

func (r Resource) SchemaName() string {
	return "resource"
}

func (_ Resource) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ Resource) SetInt(v int32)       { panic("Unsupported operation") }
func (_ Resource) SetLong(v int64)      { panic("Unsupported operation") }
func (_ Resource) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ Resource) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ Resource) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ Resource) SetString(v string)   { panic("Unsupported operation") }
func (_ Resource) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *Resource) Get(i int) types.Field {
	switch i {
	case 0:
		w := types.String{Target: &r.Display_name}

		return w

	case 1:
		w := types.String{Target: &r.Id}

		return w

	case 2:
		r.Environment = NewEnvironment()

		w := types.Record{Target: &r.Environment}

		return w

	}
	panic("Unknown field index")
}

func (r *Resource) SetDefault(i int) {
	switch i {
	}
	panic("Unknown field index")
}

func (r *Resource) NullField(i int) {
	switch i {
	}
	panic("Not a nullable field index")
}

func (_ Resource) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ Resource) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ Resource) HintSize(int)                     { panic("Unsupported operation") }
func (_ Resource) Finalize()                        {}

func (_ Resource) AvroCRC64Fingerprint() []byte {
	return []byte(ResourceAvroCRC64Fingerprint)
}

func (r Resource) MarshalJSON() ([]byte, error) {
	var err error
	output := make(map[string]json.RawMessage)
	output["display_name"], err = json.Marshal(r.Display_name)
	if err != nil {
		return nil, err
	}
	output["id"], err = json.Marshal(r.Id)
	if err != nil {
		return nil, err
	}
	output["environment"], err = json.Marshal(r.Environment)
	if err != nil {
		return nil, err
	}
	return json.Marshal(output)
}

func (r *Resource) UnmarshalJSON(data []byte) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	var val json.RawMessage
	val = func() json.RawMessage {
		if v, ok := fields["display_name"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.Display_name); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for display_name")
	}
	val = func() json.RawMessage {
		if v, ok := fields["id"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.Id); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for id")
	}
	val = func() json.RawMessage {
		if v, ok := fields["environment"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.Environment); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for environment")
	}
	return nil
}
