package validator

type Validator struct {
	errors []string
}

func New() *Validator {
	return &Validator{}
}

func (v *Validator) Required(value, key string) {
	if value == "" {
		v.addError(key + " is required")
	}
}

func (v *Validator) IsValid() bool {
	return len(v.errors) == 0
}

func (v *Validator) Errors() []string {
	return v.errors
}

func (v *Validator) addError(err string) {
	v.errors = append(v.errors, err)
}
