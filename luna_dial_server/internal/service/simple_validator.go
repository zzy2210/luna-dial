package service

import (
    "errors"
    "fmt"
    "reflect"
    "strconv"
    "strings"
    "time"
)

// SimpleValidator 实现 Echo 的 Validator 接口，支持常用规则：
// - required
// - omitempty
// - oneof=a b c (字符串)
// - min= / max= （整数）
// - dive 与 oneof 联用以校验切片元素
type SimpleValidator struct{}

func NewSimpleValidator() *SimpleValidator { return &SimpleValidator{} }

func (v *SimpleValidator) Validate(i interface{}) error {
    val := reflect.ValueOf(i)
    if val.Kind() == reflect.Ptr {
        val = val.Elem()
    }
    if val.Kind() != reflect.Struct {
        return nil
    }
    t := val.Type()
    for idx := 0; idx < t.NumField(); idx++ {
        field := t.Field(idx)
        tag := field.Tag.Get("validate")
        if tag == "" {
            continue
        }
        if err := v.validateField(field, val.Field(idx), tag); err != nil {
            return err
        }
    }
    return nil
}

func (v *SimpleValidator) validateField(sf reflect.StructField, fv reflect.Value, tag string) error {
    tokens := strings.Split(tag, ",")
    // dive: 针对切片/数组元素应用后续规则
    for i, tok := range tokens {
        if tok == "dive" {
            if fv.Kind() == reflect.Ptr {
                if fv.IsNil() {
                    return nil
                }
                fv = fv.Elem()
            }
            if fv.Kind() != reflect.Slice && fv.Kind() != reflect.Array {
                return fmt.Errorf("field %s: dive on non-slice", sf.Name)
            }
            elemRules := strings.Join(tokens[i+1:], ",")
            for j := 0; j < fv.Len(); j++ {
                if err := v.applyRules(sf, fv.Index(j), elemRules); err != nil {
                    return err
                }
            }
            return nil
        }
    }
    return v.applyRules(sf, fv, tag)
}

func (v *SimpleValidator) applyRules(sf reflect.StructField, fv reflect.Value, tag string) error {
    tokens := strings.Split(tag, ",")
    omitEmpty := false
    for _, tok := range tokens {
        if tok == "omitempty" {
            omitEmpty = true
            break
        }
    }

    // 处理 required
    for _, tok := range tokens {
        if tok == "required" {
            if isZeroValue(fv) {
                return fmt.Errorf("field %s is required", jsonName(sf))
            }
            break
        }
    }

    // omitempty: 空值直接跳过其它校验
    if omitEmpty && isZeroValue(fv) {
        return nil
    }

    // 其它规则
    for _, tok := range tokens {
        if tok == "required" || tok == "omitempty" || tok == "dive" {
            continue
        }
        if strings.HasPrefix(tok, "oneof=") {
            opts := strings.TrimPrefix(tok, "oneof=")
            allowed := strings.Fields(opts)
            var s string
            vv := deref(fv)
            if vv.Kind() != reflect.String {
                // 仅支持字符串
                return fmt.Errorf("field %s must be one of %v", jsonName(sf), allowed)
            }
            s = vv.String()
            if !contains(allowed, s) {
                return fmt.Errorf("field %s must be one of %v", jsonName(sf), allowed)
            }
            continue
        }
        if strings.HasPrefix(tok, "min=") || strings.HasPrefix(tok, "max=") {
            vv := deref(fv)
            if !isIntKind(vv.Kind()) {
                return errors.New("min/max only supported on integers")
            }
            val := vv.Int()
            if strings.HasPrefix(tok, "min=") {
                minStr := strings.TrimPrefix(tok, "min=")
                m, _ := strconv.ParseInt(minStr, 10, 64)
                if val < m {
                    return fmt.Errorf("field %s must be >= %s", jsonName(sf), minStr)
                }
            } else {
                maxStr := strings.TrimPrefix(tok, "max=")
                m, _ := strconv.ParseInt(maxStr, 10, 64)
                if val > m {
                    return fmt.Errorf("field %s must be <= %s", jsonName(sf), maxStr)
                }
            }
            continue
        }
    }
    return nil
}

func deref(v reflect.Value) reflect.Value {
    for v.IsValid() && v.Kind() == reflect.Ptr {
        if v.IsNil() {
            return v
        }
        v = v.Elem()
    }
    return v
}

func isZeroValue(v reflect.Value) bool {
    v = deref(v)
    if !v.IsValid() {
        return true
    }
    switch v.Kind() {
    case reflect.String:
        return v.Len() == 0
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        return v.Int() == 0
    case reflect.Slice, reflect.Array, reflect.Map:
        return v.Len() == 0
    case reflect.Struct:
        if v.Type() == reflect.TypeOf(time.Time{}) {
            t := v.Interface().(time.Time)
            return t.IsZero()
        }
        // 其它结构体不默认 required
        return false
    case reflect.Bool:
        return !v.Bool()
    default:
        return !v.IsValid()
    }
}

func contains(list []string, s string) bool {
    for _, v := range list {
        if v == s {
            return true
        }
    }
    return false
}

func isIntKind(k reflect.Kind) bool {
    switch k {
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        return true
    default:
        return false
    }
}

func jsonName(sf reflect.StructField) string {
    if j := sf.Tag.Get("json"); j != "" {
        if idx := strings.Index(j, ","); idx > 0 {
            return j[:idx]
        }
        return j
    }
    return sf.Name
}

