package internal

import (
    "fmt"
    "log"
    "reflect"
    "strconv"
    "strings"
    "time"

    "github.com/go-playground/locales"
    "github.com/go-playground/universal-translator"
    "gopkg.in/go-playground/validator.v9"
)


type FieldError struct {
    validator.FieldError
    locale *locale
}

func (fe *FieldError) Field() string {
    validationFieldTranslation := fe.locale.validationTranslation.FieldTranslation
    if value, ok := validationFieldTranslation[fe.FieldError.Field()]; ok {
        return value
    }
    return fe.FieldError.Field()
}

// RegisterDefaultTranslations registers a set of default translations
// for all built in tag's in validator; you may add your own as desired.
func RegisterDefaultTranslations(v *validator.Validate, locale *locale) (err error) {

    trans := locale.universalTranslator
    translationValue := locale.validationTranslation

    translations := []struct {
        tag             string
        translation     string
        override        bool
        customRegisFunc validator.RegisterTranslationsFunc
        customTransFunc validator.TranslationFunc
    }{
        {
            tag:         "required",
            translation: translationValue.Required,
            override:    false,
        },
        {
            tag: "len",
            customRegisFunc: func(ut ut.Translator) (err error) {

                if err = ut.Add("len-string", translationValue.Len.String, false); err != nil {
                    return
                }

                if err = ut.AddCardinal("len-string-character", "{0} "+translationValue.PluralRuleMap["character"].One, locales.PluralRuleOne, false); err != nil {
                    return
                }

                if err = ut.AddCardinal("len-string-character", "{0} "+translationValue.PluralRuleMap["character"].Other, locales.PluralRuleOther, false); err != nil {
                    return
                }

                if err = ut.Add("len-number", translationValue.Len.Numeric, false); err != nil {
                    return
                }

                if err = ut.Add("len-items", translationValue.Len.Array, false); err != nil {
                    return
                }
                if err = ut.AddCardinal("len-items-item", "{0} "+translationValue.PluralRuleMap["item"].One, locales.PluralRuleOne, false); err != nil {
                    return
                }

                if err = ut.AddCardinal("len-items-item", "{0} "+translationValue.PluralRuleMap["item"].Other, locales.PluralRuleOther, false); err != nil {
                    return
                }

                return

            },
            customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
                fe = &FieldError{FieldError:fe, locale: locale}

                var err error
                var t string

                var digits uint64
                var kind reflect.Kind

                if idx := strings.Index(fe.Param(), "."); idx != -1 {
                    digits = uint64(len(fe.Param()[idx+1:]))
                }

                f64, err := strconv.ParseFloat(fe.Param(), 64)
                if err != nil {
                    goto END
                }

                kind = fe.Kind()
                if kind == reflect.Ptr {
                    kind = fe.Type().Elem().Kind()
                }

                switch kind {
                case reflect.String:

                    var c string

                    c, err = ut.C("len-string-character", f64, digits, ut.FmtNumber(f64, digits))
                    if err != nil {
                        goto END
                    }

                    t, err = ut.T("len-string", fe.Field(), c)

                case reflect.Slice, reflect.Map, reflect.Array:
                    var c string

                    c, err = ut.C("len-items-item", f64, digits, ut.FmtNumber(f64, digits))
                    if err != nil {
                        goto END
                    }

                    t, err = ut.T("len-items", fe.Field(), c)

                default:
                    t, err = ut.T("len-number", fe.Field(), ut.FmtNumber(f64, digits))
                }

            END:
                if err != nil {
                    fmt.Printf("warning: error translating FieldError: %s", err)
                    return fe.(error).Error()
                }

                return t
            },
        },
        {
            tag: "min",
            customRegisFunc: func(ut ut.Translator) (err error) {

                if err = ut.Add("min-string", translationValue.Min.String, false); err != nil {
                    return
                }

                if err = ut.AddCardinal("min-string-character", "{0} "+translationValue.PluralRuleMap["character"].One, locales.PluralRuleOne, false); err != nil {
                    return
                }

                if err = ut.AddCardinal("min-string-character", "{0} "+translationValue.PluralRuleMap["character"].Other, locales.PluralRuleOther, false); err != nil {
                    return
                }

                if err = ut.Add("min-number", translationValue.Min.Numeric, false); err != nil {
                    return
                }

                if err = ut.Add("min-items", translationValue.Min.Array, false); err != nil {
                    return
                }
                if err = ut.AddCardinal("min-items-item", "{0} "+translationValue.PluralRuleMap["item"].One, locales.PluralRuleOne, false); err != nil {
                    return
                }

                if err = ut.AddCardinal("min-items-item", "{0} "+translationValue.PluralRuleMap["item"].Other, locales.PluralRuleOther, false); err != nil {
                    return
                }

                return

            },
            customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
                fe = &FieldError{FieldError:fe, locale: locale}

                var err error
                var t string

                var digits uint64
                var kind reflect.Kind

                if idx := strings.Index(fe.Param(), "."); idx != -1 {
                    digits = uint64(len(fe.Param()[idx+1:]))
                }

                f64, err := strconv.ParseFloat(fe.Param(), 64)
                if err != nil {
                    goto END
                }

                kind = fe.Kind()
                if kind == reflect.Ptr {
                    kind = fe.Type().Elem().Kind()
                }

                switch kind {
                case reflect.String:

                    var c string

                    c, err = ut.C("min-string-character", f64, digits, ut.FmtNumber(f64, digits))
                    if err != nil {
                        goto END
                    }

                    t, err = ut.T("min-string", fe.Field(), c)

                case reflect.Slice, reflect.Map, reflect.Array:
                    var c string

                    c, err = ut.C("min-items-item", f64, digits, ut.FmtNumber(f64, digits))
                    if err != nil {
                        goto END
                    }

                    t, err = ut.T("min-items", fe.Field(), c)

                default:
                    t, err = ut.T("min-number", fe.Field(), ut.FmtNumber(f64, digits))
                }

            END:
                if err != nil {
                    fmt.Printf("warning: error translating FieldError: %s", err)
                    return fe.(error).Error()
                }

                return t
            },
        },
        {
            tag: "max",
            customRegisFunc: func(ut ut.Translator) (err error) {

                if err = ut.Add("max-string", translationValue.Max.String, false); err != nil {
                    return
                }

                if err = ut.AddCardinal("max-string-character", "{0} "+translationValue.PluralRuleMap["character"].One, locales.PluralRuleOne, false); err != nil {
                    return
                }

                if err = ut.AddCardinal("max-string-character", "{0} "+translationValue.PluralRuleMap["character"].Other, locales.PluralRuleOther, false); err != nil {
                    return
                }

                if err = ut.Add("max-number", translationValue.Max.Numeric, false); err != nil {
                    return
                }

                if err = ut.Add("max-items", translationValue.Max.Array, false); err != nil {
                    return
                }
                if err = ut.AddCardinal("max-items-item", "{0} "+translationValue.PluralRuleMap["item"].One, locales.PluralRuleOne, false); err != nil {
                    return
                }

                if err = ut.AddCardinal("max-items-item", "{0} "+translationValue.PluralRuleMap["item"].Other, locales.PluralRuleOther, false); err != nil {
                    return
                }

                return

            },
            customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
                fe = &FieldError{FieldError:fe, locale: locale}

                var err error
                var t string

                var digits uint64
                var kind reflect.Kind

                if idx := strings.Index(fe.Param(), "."); idx != -1 {
                    digits = uint64(len(fe.Param()[idx+1:]))
                }

                f64, err := strconv.ParseFloat(fe.Param(), 64)
                if err != nil {
                    goto END
                }

                kind = fe.Kind()
                if kind == reflect.Ptr {
                    kind = fe.Type().Elem().Kind()
                }

                switch kind {
                case reflect.String:

                    var c string

                    c, err = ut.C("max-string-character", f64, digits, ut.FmtNumber(f64, digits))
                    if err != nil {
                        goto END
                    }

                    t, err = ut.T("max-string", fe.Field(), c)

                case reflect.Slice, reflect.Map, reflect.Array:
                    var c string

                    c, err = ut.C("max-items-item", f64, digits, ut.FmtNumber(f64, digits))
                    if err != nil {
                        goto END
                    }

                    t, err = ut.T("max-items", fe.Field(), c)

                default:
                    t, err = ut.T("max-number", fe.Field(), ut.FmtNumber(f64, digits))
                }

            END:
                if err != nil {
                    fmt.Printf("warning: error translating FieldError: %s", err)
                    return fe.(error).Error()
                }

                return t
            },
        },
        {
            tag:         "eq",
            translation: translationValue.Eq,
            override:    false,
            customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
                fe = &FieldError{FieldError:fe, locale: locale}

                t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
                if err != nil {
                    fmt.Printf("warning: error translating FieldError: %#v", fe)
                    return fe.(error).Error()
                }

                return t
            },
        },
        {
            tag:         "ne",
            translation: translationValue.Ne,
            override:    false,
            customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
                fe = &FieldError{FieldError:fe, locale: locale}

                t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
                if err != nil {
                    fmt.Printf("warning: error translating FieldError: %#v", fe)
                    return fe.(error).Error()
                }

                return t
            },
        },
        {
            tag: "lt",
            customRegisFunc: func(ut ut.Translator) (err error) {

                if err = ut.Add("lt-string", translationValue.Lt.String, false); err != nil {
                    return
                }

                if err = ut.AddCardinal("lt-string-character", "{0} "+translationValue.PluralRuleMap["character"].One, locales.PluralRuleOne, false); err != nil {
                    return
                }

                if err = ut.AddCardinal("lt-string-character", "{0} "+translationValue.PluralRuleMap["character"].Other, locales.PluralRuleOther, false); err != nil {
                    return
                }

                if err = ut.Add("lt-number", translationValue.Lt.Numeric, false); err != nil {
                    return
                }

                if err = ut.Add("lt-items", translationValue.Lt.Array, false); err != nil {
                    return
                }

                if err = ut.AddCardinal("lt-items-item", "{0} "+translationValue.PluralRuleMap["item"].One, locales.PluralRuleOne, false); err != nil {
                    return
                }

                if err = ut.AddCardinal("lt-items-item", "{0} "+translationValue.PluralRuleMap["item"].Other, locales.PluralRuleOther, false); err != nil {
                    return
                }

                if err = ut.Add("lt-datetime", translationValue.Lt.Datetime, false); err != nil {
                    return
                }

                return

            },
            customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
                fe = &FieldError{FieldError:fe, locale: locale}

                var err error
                var t string
                var f64 float64
                var digits uint64
                var kind reflect.Kind

                fn := func() (err error) {

                    if idx := strings.Index(fe.Param(), "."); idx != -1 {
                        digits = uint64(len(fe.Param()[idx+1:]))
                    }

                    f64, err = strconv.ParseFloat(fe.Param(), 64)

                    return
                }

                kind = fe.Kind()
                if kind == reflect.Ptr {
                    kind = fe.Type().Elem().Kind()
                }

                switch kind {
                case reflect.String:

                    var c string

                    err = fn()
                    if err != nil {
                        goto END
                    }

                    c, err = ut.C("lt-string-character", f64, digits, ut.FmtNumber(f64, digits))
                    if err != nil {
                        goto END
                    }

                    t, err = ut.T("lt-string", fe.Field(), c)

                case reflect.Slice, reflect.Map, reflect.Array:
                    var c string

                    err = fn()
                    if err != nil {
                        goto END
                    }

                    c, err = ut.C("lt-items-item", f64, digits, ut.FmtNumber(f64, digits))
                    if err != nil {
                        goto END
                    }

                    t, err = ut.T("lt-items", fe.Field(), c)

                case reflect.Struct:
                    if fe.Type() != reflect.TypeOf(time.Time{}) {
                        err = fmt.Errorf("tag '%s' cannot be used on a struct type", fe.Tag())
                        goto END
                    }

                    t, err = ut.T("lt-datetime", fe.Field())

                default:
                    err = fn()
                    if err != nil {
                        goto END
                    }

                    t, err = ut.T("lt-number", fe.Field(), ut.FmtNumber(f64, digits))
                }

            END:
                if err != nil {
                    fmt.Printf("warning: error translating FieldError: %s", err)
                    return fe.(error).Error()
                }

                return t
            },
        },
        {
            tag: "lte",
            customRegisFunc: func(ut ut.Translator) (err error) {

                if err = ut.Add("lte-string", translationValue.Lte.String, false); err != nil {
                    return
                }

                if err = ut.AddCardinal("lte-string-character", "{0} "+translationValue.PluralRuleMap["character"].One, locales.PluralRuleOne, false); err != nil {
                    return
                }

                if err = ut.AddCardinal("lte-string-character", "{0} "+translationValue.PluralRuleMap["character"].Other, locales.PluralRuleOther, false); err != nil {
                    return
                }

                if err = ut.Add("lte-number", translationValue.Lte.Numeric, false); err != nil {
                    return
                }

                if err = ut.Add("lte-items", translationValue.Lte.Array, false); err != nil {
                    return
                }

                if err = ut.AddCardinal("lte-items-item", "{0} "+translationValue.PluralRuleMap["item"].One, locales.PluralRuleOne, false); err != nil {
                    return
                }

                if err = ut.AddCardinal("lte-items-item", "{0} "+translationValue.PluralRuleMap["item"].Other, locales.PluralRuleOther, false); err != nil {
                    return
                }

                if err = ut.Add("lte-datetime", translationValue.Lte.Datetime, false); err != nil {
                    return
                }

                return
            },
            customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
                fe = &FieldError{FieldError:fe, locale: locale}

                var err error
                var t string
                var f64 float64
                var digits uint64
                var kind reflect.Kind

                fn := func() (err error) {

                    if idx := strings.Index(fe.Param(), "."); idx != -1 {
                        digits = uint64(len(fe.Param()[idx+1:]))
                    }

                    f64, err = strconv.ParseFloat(fe.Param(), 64)

                    return
                }

                kind = fe.Kind()
                if kind == reflect.Ptr {
                    kind = fe.Type().Elem().Kind()
                }

                switch kind {
                case reflect.String:

                    var c string

                    err = fn()
                    if err != nil {
                        goto END
                    }

                    c, err = ut.C("lte-string-character", f64, digits, ut.FmtNumber(f64, digits))
                    if err != nil {
                        goto END
                    }

                    t, err = ut.T("lte-string", fe.Field(), c)

                case reflect.Slice, reflect.Map, reflect.Array:
                    var c string

                    err = fn()
                    if err != nil {
                        goto END
                    }

                    c, err = ut.C("lte-items-item", f64, digits, ut.FmtNumber(f64, digits))
                    if err != nil {
                        goto END
                    }

                    t, err = ut.T("lte-items", fe.Field(), c)

                case reflect.Struct:
                    if fe.Type() != reflect.TypeOf(time.Time{}) {
                        err = fmt.Errorf("tag '%s' cannot be used on a struct type", fe.Tag())
                        goto END
                    }

                    t, err = ut.T("lte-datetime", fe.Field())

                default:
                    err = fn()
                    if err != nil {
                        goto END
                    }

                    t, err = ut.T("lte-number", fe.Field(), ut.FmtNumber(f64, digits))
                }

            END:
                if err != nil {
                    fmt.Printf("warning: error translating FieldError: %s", err)
                    return fe.(error).Error()
                }

                return t
            },
        },
        {
            tag: "gt",
            customRegisFunc: func(ut ut.Translator) (err error) {

                if err = ut.Add("gt-string", translationValue.Gt.String, false); err != nil {
                    return
                }

                if err = ut.AddCardinal("gt-string-character", "{0} "+translationValue.PluralRuleMap["character"].One, locales.PluralRuleOne, false); err != nil {
                    return
                }

                if err = ut.AddCardinal("gt-string-character", "{0} "+translationValue.PluralRuleMap["character"].Other, locales.PluralRuleOther, false); err != nil {
                    return
                }

                if err = ut.Add("gt-number", translationValue.Gt.Numeric, false); err != nil {
                    return
                }

                if err = ut.Add("gt-items", translationValue.Gt.Array, false); err != nil {
                    return
                }

                if err = ut.AddCardinal("gt-items-item", "{0} "+translationValue.PluralRuleMap["item"].One, locales.PluralRuleOne, false); err != nil {
                    return
                }

                if err = ut.AddCardinal("gt-items-item", "{0} "+translationValue.PluralRuleMap["item"].Other, locales.PluralRuleOther, false); err != nil {
                    return
                }

                if err = ut.Add("gt-datetime", translationValue.Gt.Datetime, false); err != nil {
                    return
                }

                return
            },
            customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
                fe = &FieldError{FieldError:fe, locale: locale}

                var err error
                var t string
                var f64 float64
                var digits uint64
                var kind reflect.Kind

                fn := func() (err error) {

                    if idx := strings.Index(fe.Param(), "."); idx != -1 {
                        digits = uint64(len(fe.Param()[idx+1:]))
                    }

                    f64, err = strconv.ParseFloat(fe.Param(), 64)

                    return
                }

                kind = fe.Kind()
                if kind == reflect.Ptr {
                    kind = fe.Type().Elem().Kind()
                }

                switch kind {
                case reflect.String:

                    var c string

                    err = fn()
                    if err != nil {
                        goto END
                    }

                    c, err = ut.C("gt-string-character", f64, digits, ut.FmtNumber(f64, digits))
                    if err != nil {
                        goto END
                    }

                    t, err = ut.T("gt-string", fe.Field(), c)

                case reflect.Slice, reflect.Map, reflect.Array:
                    var c string

                    err = fn()
                    if err != nil {
                        goto END
                    }

                    c, err = ut.C("gt-items-item", f64, digits, ut.FmtNumber(f64, digits))
                    if err != nil {
                        goto END
                    }

                    t, err = ut.T("gt-items", fe.Field(), c)

                case reflect.Struct:
                    if fe.Type() != reflect.TypeOf(time.Time{}) {
                        err = fmt.Errorf("tag '%s' cannot be used on a struct type", fe.Tag())
                        goto END
                    }

                    t, err = ut.T("gt-datetime", fe.Field())

                default:
                    err = fn()
                    if err != nil {
                        goto END
                    }

                    t, err = ut.T("gt-number", fe.Field(), ut.FmtNumber(f64, digits))
                }

            END:
                if err != nil {
                    fmt.Printf("warning: error translating FieldError: %s", err)
                    return fe.(error).Error()
                }

                return t
            },
        },
        {
            tag: "gte",
            customRegisFunc: func(ut ut.Translator) (err error) {

                if err = ut.Add("gte-string", translationValue.Gte.String, false); err != nil {
                    return
                }

                if err = ut.AddCardinal("gte-string-character", "{0} "+translationValue.PluralRuleMap["character"].One, locales.PluralRuleOne, false); err != nil {
                    return
                }

                if err = ut.AddCardinal("gte-string-character", "{0} "+translationValue.PluralRuleMap["character"].Other, locales.PluralRuleOther, false); err != nil {
                    return
                }

                if err = ut.Add("gte-number", translationValue.Gte.Numeric, false); err != nil {
                    return
                }

                if err = ut.Add("gte-items", translationValue.Gte.Array, false); err != nil {
                    return
                }

                if err = ut.AddCardinal("gte-items-item", "{0} "+translationValue.PluralRuleMap["item"].One, locales.PluralRuleOne, false); err != nil {
                    return
                }

                if err = ut.AddCardinal("gte-items-item", "{0} "+translationValue.PluralRuleMap["item"].Other, locales.PluralRuleOther, false); err != nil {
                    return
                }

                if err = ut.Add("gte-datetime", translationValue.Gte.Datetime, false); err != nil {
                    return
                }

                return
            },
            customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
                fe = &FieldError{FieldError:fe, locale: locale}

                var err error
                var t string
                var f64 float64
                var digits uint64
                var kind reflect.Kind

                fn := func() (err error) {

                    if idx := strings.Index(fe.Param(), "."); idx != -1 {
                        digits = uint64(len(fe.Param()[idx+1:]))
                    }

                    f64, err = strconv.ParseFloat(fe.Param(), 64)

                    return
                }

                kind = fe.Kind()
                if kind == reflect.Ptr {
                    kind = fe.Type().Elem().Kind()
                }

                switch kind {
                case reflect.String:

                    var c string

                    err = fn()
                    if err != nil {
                        goto END
                    }

                    c, err = ut.C("gte-string-character", f64, digits, ut.FmtNumber(f64, digits))
                    if err != nil {
                        goto END
                    }

                    t, err = ut.T("gte-string", fe.Field(), c)

                case reflect.Slice, reflect.Map, reflect.Array:
                    var c string

                    err = fn()
                    if err != nil {
                        goto END
                    }

                    c, err = ut.C("gte-items-item", f64, digits, ut.FmtNumber(f64, digits))
                    if err != nil {
                        goto END
                    }

                    t, err = ut.T("gte-items", fe.Field(), c)

                case reflect.Struct:
                    if fe.Type() != reflect.TypeOf(time.Time{}) {
                        err = fmt.Errorf("tag '%s' cannot be used on a struct type", fe.Tag())
                        goto END
                    }

                    t, err = ut.T("gte-datetime", fe.Field())

                default:
                    err = fn()
                    if err != nil {
                        goto END
                    }

                    t, err = ut.T("gte-number", fe.Field(), ut.FmtNumber(f64, digits))
                }

            END:
                if err != nil {
                    fmt.Printf("warning: error translating FieldError: %s", err)
                    return fe.(error).Error()
                }

                return t
            },
        },
        {
            tag:         "eqfield",
            translation: translationValue.Eqfield,
            override:    false,
            customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
                fe = &FieldError{FieldError:fe, locale: locale}

                t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
                if err != nil {
                    log.Printf("warning: error translating FieldError: %#v", fe)
                    return fe.(error).Error()
                }

                return t
            },
        },
        {
            tag:         "eqcsfield",
            translation: translationValue.Eqcsfield,
            override:    false,
            customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
                fe = &FieldError{FieldError:fe, locale: locale}

                t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
                if err != nil {
                    log.Printf("warning: error translating FieldError: %#v", fe)
                    return fe.(error).Error()
                }

                return t
            },
        },
        {
            tag:         "necsfield",
            translation: translationValue.Necsfield,
            override:    false,
            customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
                fe = &FieldError{FieldError:fe, locale: locale}

                t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
                if err != nil {
                    log.Printf("warning: error translating FieldError: %#v", fe)
                    return fe.(error).Error()
                }

                return t
            },
        },
        {
            tag:         "gtcsfield",
            translation: translationValue.Gtcsfield,
            override:    false,
            customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
                fe = &FieldError{FieldError:fe, locale: locale}

                t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
                if err != nil {
                    log.Printf("warning: error translating FieldError: %#v", fe)
                    return fe.(error).Error()
                }

                return t
            },
        },
        {
            tag:         "gtecsfield",
            translation: translationValue.Gtecsfield,
            override:    false,
            customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
                fe = &FieldError{FieldError:fe, locale: locale}

                t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
                if err != nil {
                    log.Printf("warning: error translating FieldError: %#v", fe)
                    return fe.(error).Error()
                }

                return t
            },
        },
        {
            tag:         "ltcsfield",
            translation: translationValue.Ltcsfield,
            override:    false,
            customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
                fe = &FieldError{FieldError:fe, locale: locale}

                t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
                if err != nil {
                    log.Printf("warning: error translating FieldError: %#v", fe)
                    return fe.(error).Error()
                }

                return t
            },
        },
        {
            tag:         "ltecsfield",
            translation: translationValue.Ltecsfield,
            override:    false,
            customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
                fe = &FieldError{FieldError:fe, locale: locale}

                t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
                if err != nil {
                    log.Printf("warning: error translating FieldError: %#v", fe)
                    return fe.(error).Error()
                }

                return t
            },
        },
        {
            tag:         "nefield",
            translation: translationValue.Nefield,
            override:    false,
            customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
                fe = &FieldError{FieldError:fe, locale: locale}

                t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
                if err != nil {
                    log.Printf("warning: error translating FieldError: %#v", fe)
                    return fe.(error).Error()
                }

                return t
            },
        },
        {
            tag:         "gtfield",
            translation: translationValue.Gtfield,
            override:    false,
            customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
                fe = &FieldError{FieldError:fe, locale: locale}

                t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
                if err != nil {
                    log.Printf("warning: error translating FieldError: %#v", fe)
                    return fe.(error).Error()
                }

                return t
            },
        },
        {
            tag:         "gtefield",
            translation: translationValue.Gtefield,
            override:    false,
            customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
                fe = &FieldError{FieldError:fe, locale: locale}

                t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
                if err != nil {
                    log.Printf("warning: error translating FieldError: %#v", fe)
                    return fe.(error).Error()
                }

                return t
            },
        },
        {
            tag:         "ltfield",
            translation: translationValue.Ltfield,
            override:    false,
            customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
                fe = &FieldError{FieldError:fe, locale: locale}

                t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
                if err != nil {
                    log.Printf("warning: error translating FieldError: %#v", fe)
                    return fe.(error).Error()
                }

                return t
            },
        },
        {
            tag:         "ltefield",
            translation: translationValue.Ltefield,
            override:    false,
            customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
                fe = &FieldError{FieldError:fe, locale: locale}

                t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
                if err != nil {
                    log.Printf("warning: error translating FieldError: %#v", fe)
                    return fe.(error).Error()
                }

                return t
            },
        },
        {
            tag:         "alpha",
            translation: translationValue.Alpha,
            override:    false,
        },
        {
            tag:         "alphanum",
            translation: translationValue.Alphanum,
            override:    false,
        },
        {
            tag:         "numeric",
            translation: translationValue.Numeric,
            override:    false,
        },
        {
            tag:         "number",
            translation: translationValue.Number,
            override:    false,
        },
        {
            tag:         "hexadecimal",
            translation: translationValue.Hexadecimal,
            override:    false,
        },
        {
            tag:         "hexcolor",
            translation: translationValue.Hexcolor,
            override:    false,
        },
        {
            tag:         "rgb",
            translation: translationValue.Rgb,
            override:    false,
        },
        {
            tag:         "rgba",
            translation: translationValue.Rgba,
            override:    false,
        },
        {
            tag:         "hsl",
            translation: translationValue.Hsl,
            override:    false,
        },
        {
            tag:         "hsla",
            translation: translationValue.Hsla,
            override:    false,
        },
        {
            tag:         "email",
            translation: translationValue.Email,
            override:    false,
        },
        {
            tag:         "url",
            translation: translationValue.Url,
            override:    false,
        },
        {
            tag:         "uri",
            translation: translationValue.Uri,
            override:    false,
        },
        {
            tag:         "base64",
            translation: translationValue.Base64,
            override:    false,
        },
        {
            tag:         "contains",
            translation: translationValue.Contains,
            override:    false,
            customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
                fe = &FieldError{FieldError:fe, locale: locale}

                t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
                if err != nil {
                    log.Printf("warning: error translating FieldError: %#v", fe)
                    return fe.(error).Error()
                }

                return t
            },
        },
        {
            tag:         "containsany",
            translation: translationValue.Containsany,
            override:    false,
            customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
                fe = &FieldError{FieldError:fe, locale: locale}

                t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
                if err != nil {
                    log.Printf("warning: error translating FieldError: %#v", fe)
                    return fe.(error).Error()
                }

                return t
            },
        },
        {
            tag:         "excludes",
            translation: translationValue.Excludes,
            override:    false,
            customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
                fe = &FieldError{FieldError:fe, locale: locale}

                t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
                if err != nil {
                    log.Printf("warning: error translating FieldError: %#v", fe)
                    return fe.(error).Error()
                }

                return t
            },
        },
        {
            tag:         "excludesall",
            translation: translationValue.Excludesall,
            override:    false,
            customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
                fe = &FieldError{FieldError:fe, locale: locale}

                t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
                if err != nil {
                    log.Printf("warning: error translating FieldError: %#v", fe)
                    return fe.(error).Error()
                }

                return t
            },
        },
        {
            tag:         "excludesrune",
            translation: translationValue.Excludesrune,
            override:    false,
            customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
                fe = &FieldError{FieldError:fe, locale: locale}

                t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
                if err != nil {
                    log.Printf("warning: error translating FieldError: %#v", fe)
                    return fe.(error).Error()
                }

                return t
            },
        },
        {
            tag:         "isbn",
            translation: translationValue.Isbn,
            override:    false,
        },
        {
            tag:         "isbn10",
            translation: translationValue.Isbn10,
            override:    false,
        },
        {
            tag:         "isbn13",
            translation: translationValue.Isbn13,
            override:    false,
        },
        {
            tag:         "uuid",
            translation: translationValue.Uuid,
            override:    false,
        },
        {
            tag:         "uuid3",
            translation: translationValue.Uuid3,
            override:    false,
        },
        {
            tag:         "uuid4",
            translation: translationValue.Uuid4,
            override:    false,
        },
        {
            tag:         "uuid5",
            translation: translationValue.Uuid5,
            override:    false,
        },
        {
            tag:         "ascii",
            translation: translationValue.Ascii,
            override:    false,
        },
        {
            tag:         "printascii",
            translation: translationValue.Printascii,
            override:    false,
        },
        {
            tag:         "multibyte",
            translation: translationValue.Multibyte,
            override:    false,
        },
        {
            tag:         "datauri",
            translation: translationValue.Datauri,
            override:    false,
        },
        {
            tag:         "latitude",
            translation: translationValue.Latitude,
            override:    false,
        },
        {
            tag:         "longitude",
            translation: translationValue.Longitude,
            override:    false,
        },
        {
            tag:         "ssn",
            translation: translationValue.Ssn,
            override:    false,
        },
        {
            tag:         "ipv4",
            translation: translationValue.Ipv4,
            override:    false,
        },
        {
            tag:         "ipv6",
            translation: translationValue.Ipv6,
            override:    false,
        },
        {
            tag:         "ip",
            translation: translationValue.Ip,
            override:    false,
        },
        {
            tag:         "cidr",
            translation: translationValue.Cidr,
            override:    false,
        },
        {
            tag:         "cidrv4",
            translation: translationValue.Cidrv4,
            override:    false,
        },
        {
            tag:         "cidrv6",
            translation: translationValue.Cidrv6,
            override:    false,
        },
        {
            tag:         "tcpAddr",
            translation: translationValue.TcpAddr,
            override:    false,
        },
        {
            tag:         "tcp4Addr",
            translation: translationValue.Tcp4Addr,
            override:    false,
        },
        {
            tag:         "tcp6Addr",
            translation: translationValue.Tcp6Addr,
            override:    false,
        },
        {
            tag:         "udpAddr",
            translation: translationValue.UdpAddr,
            override:    false,
        },
        {
            tag:         "udp4Addr",
            translation: translationValue.Udp4Addr,
            override:    false,
        },
        {
            tag:         "udp6Addr",
            translation: translationValue.Udp6Addr,
            override:    false,
        },
        {
            tag:         "ipAddr",
            translation: translationValue.IpAddr,
            override:    false,
        },
        {
            tag:         "ip4Addr",
            translation: translationValue.Ip4Addr,
            override:    false,
        },
        {
            tag:         "ip6Addr",
            translation: translationValue.Ip6Addr,
            override:    false,
        },
        {
            tag:         "unixAddr",
            translation: translationValue.UnixAddr,
            override:    false,
        },
        {
            tag:         "mac",
            translation: translationValue.Mac,
            override:    false,
        },
        {
            tag:         "unique",
            translation: translationValue.Unique,
            override:    false,
        },
        {
            tag:         "iscolor",
            translation: translationValue.Iscolor,
            override:    false,
        },
        {
            tag:         "oneof",
            translation: translationValue.Oneof,
            override:    false,
            customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
                fe = &FieldError{FieldError:fe, locale: locale}
                s, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
                if err != nil {
                    log.Printf("warning: error translating FieldError: %#v", fe)
                    return fe.(error).Error()
                }
                return s
            },
        },
    }


    var registrationFunc = func(tag string, translation string, override bool) validator.RegisterTranslationsFunc {

        return func(ut ut.Translator) (err error) {

            if err = ut.Add(tag, translation, override); err != nil {
                return
            }

            return

        }

    }

    var translateFunc = func(ut ut.Translator, fe validator.FieldError) string {
        fe = &FieldError{FieldError:fe, locale: locale}

        t, err := ut.T(fe.Tag(), fe.Field())
        if err != nil {
            log.Printf("warning: error translating FieldError: %#v", fe)
            return fe.(error).Error()
        }

        return t
    }

    for _, t := range translations {

        if t.customTransFunc != nil && t.customRegisFunc != nil {

            err = v.RegisterTranslation(t.tag, trans, t.customRegisFunc, t.customTransFunc)

        } else if t.customTransFunc != nil && t.customRegisFunc == nil {

            err = v.RegisterTranslation(t.tag, trans, registrationFunc(t.tag, t.translation, t.override), t.customTransFunc)

        } else if t.customTransFunc == nil && t.customRegisFunc != nil {

            err = v.RegisterTranslation(t.tag, trans, t.customRegisFunc, translateFunc)

        } else {
            err = v.RegisterTranslation(t.tag, trans, registrationFunc(t.tag, t.translation, t.override), translateFunc)
        }

        if err != nil {
            return
        }
    }

    locale.setValidationRegistered()

    return
}

