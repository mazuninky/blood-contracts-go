package blood_contracts_go

import (
	"errors"
	"regexp"
)

type RType struct {
	mapFunc MapFunction
}

func NewType(mapFunc MapFunction) RefinementType {
	rType := RType{
		mapFunc: mapFunc,
	}

	return &rType
}

func (base *RType) IsValid(value interface{}) bool {
	_, err := base.mapFunc(value)
	if err != nil {
		return false
	}
	return true
}

func (base *RType) getMapFunction() MapFunction {
	return base.mapFunc
}

func (base *RType) Pack(value interface{}) RefinementTypeBox {
	return NewBox(base.mapFunc, value)
}

func (base *RType) And(rt RefinementType) RefinementType {
	mapFunc := func(value interface{}) (interface{}, error) {
		firstValue, firstError := base.mapFunc(value)
		if firstError != nil {
			return nil, firstError
		}

		secondValue, secondError := rt.getMapFunction()(value)
		if secondError != nil {
			return nil, secondError
		}

		return []interface{}{
			firstValue,
			secondValue,
		}, nil
	}

	return NewType(mapFunc)
}

func (base *RType) Or(rt RefinementType) RefinementType {
	mapFunc := func(value interface{}) (interface{}, error) {
		firstValue, firstError := base.mapFunc(value)
		if firstError == nil {
			return firstValue, nil
		}

		secondValue, secondError := rt.getMapFunction()(value)
		if secondError == nil {
			return secondValue, nil
		}

		// TODO Choose exception
		return nil, errors.New("can't find type for both")
	}

	return NewType(mapFunc)
}

func (base *RType) Pipe(rt RefinementType) RefinementType {
	mapFunc := func(value interface{}) (interface{}, error) {
		firstValue, firstError := base.mapFunc(value)
		if firstError != nil {
			return nil, firstError
		}

		return rt.getMapFunction()(firstValue)
	}

	return NewType(mapFunc)
}

// String

// Regex

func MustNewRegexType(regex string) RefinementType {
	reg, err := regexp.Compile(regex)
	if err != nil {
		panic(err)
	}
	regexMapFunc := createRegexMapFunc(reg)
	return NewType(regexMapFunc)
}

func NewRegexType(regex string) (RefinementType, error) {
	reg, err := regexp.Compile(regex)
	if err != nil {
		return nil, err
	}
	regexMapFunc := createRegexMapFunc(reg)
	return NewType(regexMapFunc), nil
}

// Number

func NewNumberMin(minValue interface{}) (RefinementType, error) {
	container, err := NewContainer(minValue)
	if err != nil {
		return nil, err
	}

	mapFunc := createMinMapFunc(container, true)
	return NewType(mapFunc), nil
}

func NewNumberMinExclude(minValue interface{}) (RefinementType, error) {
	container, err := NewContainer(minValue)
	if err != nil {
		return nil, err
	}

	mapFunc := createMinMapFunc(container, false)
	return NewType(mapFunc), nil
}

func NewNumberMax(maxValue interface{}) (RefinementType, error) {
	container, err := NewContainer(maxValue)
	if err != nil {
		return nil, err
	}

	mapFunc := createMaxMapFunc(container, true)
	return NewType(mapFunc), nil
}

func NewNumberMaxExclude(maxValue interface{}) (RefinementType, error) {
	container, err := NewContainer(maxValue)
	if err != nil {
		return nil, err
	}

	mapFunc := createMinMapFunc(container, false)
	return NewType(mapFunc), nil
}

func NewNumberEqual(equal interface{}) (RefinementType, error) {
	container, err := NewContainer(equal)
	if err != nil {
		return nil, err
	}

	mapFunc := createEqualMapFunc(container)
	return NewType(mapFunc), nil
}


// Struct

//func NewStructType(structType interface{}) {
//
//}

// Function

func NewFunctionType(mapFunc MapFunction) RefinementType {
	return NewType(mapFunc)
}