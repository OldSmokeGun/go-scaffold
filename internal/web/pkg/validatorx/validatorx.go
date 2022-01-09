package validatorx

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"regexp"
	"strings"
)

type CustomValidator struct {
	Tag  string
	Func validator.Func
}

// RegisterCustomValidation 注册自定义验证器
func RegisterCustomValidation(validators []CustomValidator) error {
	validate := binding.Validator.Engine().(*validator.Validate)

	for _, v := range validators {
		if err := validate.RegisterValidation(v.Tag, v.Func); err != nil {
			return err
		}
	}

	return nil
}

// Translate 翻译 validator 校验后返回的错误
func Translate(errs validator.ValidationErrors, m map[string]string) map[string]string {
	if len(errs) == 0 {
		return nil
	}

	if len(m) == 0 {
		return nil
	}

	errsMap := make(map[string]string, len(errs))

	for _, e := range errs {
		rep := regexp.MustCompile(`\[\d\]`).ReplaceAllString(e.Namespace(), "")
		key := strings.Join(strings.Split(rep, ".")[1:], ".") + "." + e.Tag()

		if _, ok := m[key]; ok {
			if e.Param() != "" {
				errsMap[key] = strings.Replace(m[key], "{"+e.Tag()+"}", e.Param(), 1)
			} else {
				errsMap[key] = m[key]
			}
		}
	}

	return errsMap
}
