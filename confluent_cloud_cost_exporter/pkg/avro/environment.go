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
 
 type Environment struct {
	 Id string `json:"id"`
 }
 
 const EnvironmentAvroCRC64Fingerprint = "\x15\x81]\x00\xd1R\x06\x0e"
 
 func NewEnvironment() Environment {
	 r := Environment{}
	 return r
 }
 
 func DeserializeEnvironment(r io.Reader) (Environment, error) {
	 t := NewEnvironment()
	 deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	 if err != nil {
		 return t, err
	 }
 
	 err = vm.Eval(r, deser, &t)
	 return t, err
 }
 
 func DeserializeEnvironmentFromSchema(r io.Reader, schema string) (Environment, error) {
	 t := NewEnvironment()
 
	 deser, err := compiler.CompileSchemaBytes([]byte(schema), []byte(t.Schema()))
	 if err != nil {
		 return t, err
	 }
 
	 err = vm.Eval(r, deser, &t)
	 return t, err
 }
 
 func writeEnvironment(r Environment, w io.Writer) error {
	 var err error
	 err = vm.WriteString(r.Id, w)
	 if err != nil {
		 return err
	 }
	 return err
 }
 
 func (r Environment) Serialize(w io.Writer) error {
	 return writeEnvironment(r, w)
 }
 
 func (r Environment) Schema() string {
	 return "{\"fields\":[{\"name\":\"id\",\"type\":\"string\"}],\"name\":\"environment\",\"type\":\"record\"}"
 }
 
 func (r Environment) SchemaName() string {
	 return "environment"
 }
 
 func (_ Environment) SetBoolean(v bool)    { panic("Unsupported operation") }
 func (_ Environment) SetInt(v int32)       { panic("Unsupported operation") }
 func (_ Environment) SetLong(v int64)      { panic("Unsupported operation") }
 func (_ Environment) SetFloat(v float32)   { panic("Unsupported operation") }
 func (_ Environment) SetDouble(v float64)  { panic("Unsupported operation") }
 func (_ Environment) SetBytes(v []byte)    { panic("Unsupported operation") }
 func (_ Environment) SetString(v string)   { panic("Unsupported operation") }
 func (_ Environment) SetUnionElem(v int64) { panic("Unsupported operation") }
 
 func (r *Environment) Get(i int) types.Field {
	 switch i {
	 case 0:
		 w := types.String{Target: &r.Id}
 
		 return w
 
	 }
	 panic("Unknown field index")
 }
 
 func (r *Environment) SetDefault(i int) {
	 switch i {
	 }
	 panic("Unknown field index")
 }
 
 func (r *Environment) NullField(i int) {
	 switch i {
	 }
	 panic("Not a nullable field index")
 }
 
 func (_ Environment) AppendMap(key string) types.Field { panic("Unsupported operation") }
 func (_ Environment) AppendArray() types.Field         { panic("Unsupported operation") }
 func (_ Environment) HintSize(int)                     { panic("Unsupported operation") }
 func (_ Environment) Finalize()                        {}
 
 func (_ Environment) AvroCRC64Fingerprint() []byte {
	 return []byte(EnvironmentAvroCRC64Fingerprint)
 }
 
 func (r Environment) MarshalJSON() ([]byte, error) {
	 var err error
	 output := make(map[string]json.RawMessage)
	 output["id"], err = json.Marshal(r.Id)
	 if err != nil {
		 return nil, err
	 }
	 return json.Marshal(output)
 }
 
 func (r *Environment) UnmarshalJSON(data []byte) error {
	 var fields map[string]json.RawMessage
	 if err := json.Unmarshal(data, &fields); err != nil {
		 return err
	 }
 
	 var val json.RawMessage
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
	 return nil
 }
 