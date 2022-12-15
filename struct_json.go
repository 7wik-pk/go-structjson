package structjson

import (
	"fmt"
	"reflect"
	"strings"
)

// type TestStructInner struct {
// 	X reflect.Type `json:"x"`
// 	Y interface{}  `json:"y"`
// 	Z *string      `json:"z"`
// }

// type TestStruct struct {
// 	A int                `json:"a"`
// 	B int64              `json:"b"`
// 	C uint32             `json:"c"`
// 	D string             `json:"d"`
// 	E bool               `json:"e"`
// 	F float64            `json:"f"`
// 	G []int              `json:"g"`
// 	H map[string]int     `json:"h"`
// 	I []string           `json:"i"`
// 	J map[string]float64 `json:"j"`
// 	K TestStructInner    `json:"k"`
// 	L []TestStructInner  `json:"l"`
// }

type structJsonSpecs struct {
	Name         string
	ReflectValue reflect.Value
	Required     bool
}

func DisplayStructJson(s interface{}) (structure string) {
	return displayStructJsonR(s, 0)
}

func displayStructJsonR(s interface{}, tabSpace int) (structure string) {

	if tabSpace < 0 {
		return ""
	}

	return generateStructJsonString(s, tabSpace)

}

func generateStructJsonString(s interface{}, tabSpace int) (structure string) {

	if tabSpace < 0 {
		return ""
	}

	structure = strings.Repeat("\t", tabSpace) + "{\n"

	jsonSpecs := getStructJsonNamesValuesAndRequiredFlags(s)

	var value reflect.Value

	for _, spec := range jsonSpecs {
		// valueValue := reflectValue.Field(i).Interface()

		value = spec.ReflectValue

		structure += strings.Repeat("\t", tabSpace+1)

		valueKind := value.Kind()
		valueType := value.Type()

		// fmt.Println(valueKind)

		switch valueKind {

		case reflect.Ptr:

			structure += fmt.Sprintf("%s : %s", spec.Name, valueType.Elem())

		case reflect.Map:

			structure += fmt.Sprintf("%s : %s", spec.Name, valueType) // TODO: needs recursive handling

		case reflect.Struct:

			structure += fmt.Sprintf("%s :\n", spec.Name)
			v := value.Addr()
			structure += displayStructJsonR(v.Interface(), tabSpace+1)

		case reflect.Slice:

			switch valueType.Elem().Kind() {

			case reflect.Struct:

				v := reflect.MakeSlice(valueType, 1, 1).Index(0).Addr()
				structure += fmt.Sprintf("%s : %s\n", spec.Name, "array with each element of the following structure:")

				// fmt.Println("adding slice element structure:\n", displayStructJsonR(v.Interface(), tabSpace+1))

				structure += displayStructJsonR(v.Interface(), tabSpace+1)

			default:
				structure += fmt.Sprintf("%s : array of %ss", spec.Name, valueType.Elem())
			}

		default:

			structure += fmt.Sprintf("%s : %s", spec.Name, valueType)

		}

		structure += "\n"
	}

	structure += strings.Repeat("\t", tabSpace) + "}\n"

	return structure
}

func getStructJsonNamesValuesAndRequiredFlags(s interface{}) (specs []structJsonSpecs) {

	reflectType := reflect.TypeOf(s).Elem()
	reflectValue := reflect.ValueOf(s).Elem()

	var name string
	// var reqd string
	var ok bool

	for i := 0; i < reflectValue.NumField(); i++ {

		name, ok = reflectType.Field(i).Tag.Lookup("json")
		if ok {

			// reqd, ok =  reflectType.Field(i).Tag.Lookup("validate")
			// if ok {
			// 	strings.Contains()
			// }

			specs = append(specs, structJsonSpecs{
				Name:         name,
				ReflectValue: reflectValue.Field(i),
			})
			// specs[].Names = append(specs.Names, name)
			// specs.Fields = append(specs.Fields, reflectValue.Field(i))
		}
	}

	return specs
}
