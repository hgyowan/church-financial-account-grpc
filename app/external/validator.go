package external

import (
	"github.com/go-playground/validator/v10"
	"github.com/hgyowan/church-financial-account-grpc/domain"
	pkgLogger "github.com/hgyowan/go-pkg-library/logger"
	"regexp"
)

type validate struct {
	validator *validator.Validate
}

func (v *validate) Validator() *validator.Validate {
	return v.validator
}

func MustNewValidator() domain.ExternalValidator {
	v := &validate{validator: validator.New()}

	if err := v.Validator().RegisterValidation("phoneNumberReg", func(fl validator.FieldLevel) bool {
		phone := fl.Field().String()
		if phone == "" {
			return true
		}
		matched, _ := regexp.MatchString(`^01[016789]-?\d{3,4}-?\d{4}$`, phone)
		return matched
	}); err != nil {
		pkgLogger.ZapLogger.Logger.Sugar().Fatalf("failed to create valid object: %v", err)
	}

	return v
}
