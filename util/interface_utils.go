package util

import (
    "encoding/json"
    "fmt"
    "github.com/spf13/cast"
    "reflect"
)

func GetTypeOfInterface(data interface{}) string {
    if t := reflect.TypeOf(data); t.Kind() == reflect.Ptr {
        return "*" + t.Elem().Name()
    } else {
        return t.Name()
    }
}

func ConvertStructInterfaceToMap(object interface{}) map[string]interface{} {
    result := make(map[string]interface{})

    // Change object with the specific structure but interface to byte array
    // Example: object is User{ID: 1, Name: "Hello"}
    byteArray, err := json.Marshal(object)
    if err != nil {
        fmt.Printf("[INTERFACE_UTILS][ERROR] Convert interface to map fail: %v", err.Error())
        return result
    }

    // Unmarshal an old object with the specific structure
    var newObject interface{}
    err = json.Unmarshal(byteArray, &newObject)
    if err != nil {
        fmt.Printf("[INTERFACE_UTILS][ERROR] Convert interface to map fail: %v", err.Error())
        return result
    }

    // Convert the final interface to the map interface
    ok := false
    result, ok = newObject.(map[string]interface{})
    if !ok {
        fmt.Printf("[INTERFACE_UTILS][ERROR] Convert interface to map fail: %T", newObject)
        return result
    }

    return result
}

func GetInt64ValueInMapInterface(mapInterface map[string]interface{}, key string) int64 {
    if len(key) == 0 {
        return 0
    }

    valueInterface, ok := mapInterface[key]
    if !ok {
        return 0
    }

    return cast.ToInt64(valueInterface)
}

func GetStringValueInMapInterface(mapInterface map[string]interface{}, key string) string {
    if len(key) == 0 {
        return ""
    }

    valueInterface, ok := mapInterface[key]
    if !ok {
        return ""
    }

    return cast.ToString(valueInterface)
}

func GetBoolValueInMapInterface(mapInterface map[string]interface{}, key string) bool {
    if len(key) == 0 {
        return false
    }

    valueInterface, ok := mapInterface[key]
    if !ok {
        return false
    }

    return cast.ToBool(valueInterface)
}

func GetFloat64ValueInMapInterface(mapInterface map[string]interface{}, key string) float64 {
    if len(key) == 0 {
        return 0
    }

    valueInterface, ok := mapInterface[key]
    if !ok {
        return 0
    }

    return cast.ToFloat64(valueInterface)
}
