package em_library

import (
	"encoding/json"
	Package_Validator "github.com/go-playground/validator/v10"
)

var Instance_Validator *Package_Validator.Validate

// https://github.com/go-playground/validator/blob/master/_examples/simple/main.go
func Init_Validator()  {
	Instance_Validator = Package_Validator.New()
	initLog.Println("[INFO]", "Successfully loaded Init_Validator.")
}

type validator struct {

}

func NewValidator() *validator {
	return &validator{}
}


func (this *validator) ValidateStruct(s interface{}) error {
	err := Instance_Validator.Struct(s)
	if err != nil {
		return err
	}
	return nil
}

func (this *validator) Validate(request interface{}, my_struct interface{}) error {
	err := this.utils_ChangeType(request, my_struct)
	if err != nil {
		return err
	}
	err = this.ValidateStruct(my_struct)
	if err != nil {
		return err
	}
	return nil
}

func (this *validator) utils_ChangeType(in interface{}, out interface{}) (error) {
	b, err := json.Marshal(in)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &out)
	if err != nil {
		return err
	}
	return nil
}