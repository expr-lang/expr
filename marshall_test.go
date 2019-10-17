package expr

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"testing"

	"github.com/antonmedv/expr/vm"
)

func TestProgram_Marshall(t *testing.T) {
	param := map[string]interface{}{
		"subject": map[string]interface{}{
			"age": 10,
		},
	}
	s := `subject.age == 10 && subject.age in [10,11,12] && "hello" matches "hell"`
	p, err := Compile(s, Env(param))
	if err != nil {
		t.Fatal(err)
	}
	o, err := Run(p, param)
	if err != nil {
		t.Error(err)
	}
	if v, ok := o.(bool); !ok || !v {
		t.Error("should match")
	}

	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	dec := gob.NewDecoder(&buff)
	if err := enc.Encode(p); err != nil {
		t.Fatal("encode error: ", err)
	}
	p2 := &vm.Program{}
	if err := dec.Decode(p2); err != nil {
		t.Fatal("decode error: ", err)
	}
	o, err = Run(p2, param)
	if err != nil {
		t.Error(err)
	}
	if v, ok := o.(bool); !ok || !v {
		t.Error("should match")
	}

	j, err := json.Marshal(p)
	if err != nil {
		t.Error(err)
	}
	p3 := &vm.Program{}
	err = json.Unmarshal(j, p3)
	if err != nil {
		t.Error(err)
	}
	o, err = Run(p3, param)
	if err != nil {
		t.Error(err)
	}
	if v, ok := o.(bool); !ok || !v {
		t.Error("should match")
	}
}
