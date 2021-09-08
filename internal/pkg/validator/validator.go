package validator

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"regexp"
	"strings"
)

var validate = binding.Validator.Engine().(*validator.Validate)

type CustomValidator struct {
	Tag  string
	Func validator.Func
}

func RegisterValidator(validators []CustomValidator) error {
	for _, v := range validators {
		if err := validate.RegisterValidation(v.Tag, v.Func); err != nil {
			return err
		}
	}

	return nil
}

func Translate(errs validator.ValidationErrors, tm map[string]string) map[string]string {
	errsMap := make(map[string]string, len(errs))
	var key string

	if errs != nil {
		for _, v := range errs {
			rep := regexp.MustCompile(`\[\d\]`).ReplaceAllString(v.Namespace(), "")
			key = strings.Join(strings.Split(rep, ".")[1:], ".") + "." + v.Tag()

			if _, ok := tm[key]; ok {
				if v.Param() != "" {
					errsMap[key] = strings.Replace(tm[key], "{"+v.Tag()+"}", v.Param(), 1)
				} else {
					errsMap[key] = tm[key]
				}
			}
		}
	}

	return errsMap
}
