package stringhandlers

import (
	"reflect"
	"testing"
	"time"
)

func TestStringHandlers_Substr(t *testing.T) {
	input := "test"
	result := "es"

	if Substr(input, 1, 2) != result {
		t.Error("Substring are wrong")
	}
}

func TestStringHandlers_ConvertField(t *testing.T) {
	strDate := "20170923"
	strDateFormat := "20060102"

	strNumber := "120120"

	retDate, err := ConvertField("date", strDateFormat, strDate)
	if err != nil {
		t.Error("Fail to parse date value (err): " + err.Error())
	}

	if time.Date(2017, time.September, 23, 0, 0, 0, 0, time.UTC) != retDate {
		t.Error("Fail to parse date value (values)")
	}

	retInt, err := ConvertField("number", "", strNumber)
	if err != nil {
		t.Error("Fail to parse number value (err): " + err.Error())
	}

	if reflect.ValueOf(retInt).Int() != 120120 {
		t.Error("Fail to parse number value (values)")
	}

	retDefault, err := ConvertField("not_exists", "", "abc")
	if err != nil {
		t.Error("Fail default case")
	}

	if retDefault != "abc" {
		t.Error("Fail default case")
	}
}
