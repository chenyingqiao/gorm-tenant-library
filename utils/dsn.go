package utils

import (
	"crypto/md5"
	"errors"
	"net/url"
	"regexp"
	"strings"
)

func GetDsnQuery(dsn string, query string) (string, error) {
	values, err := url.ParseQuery(dsn)
	if err != nil {
		return "", err
	}
	queryV := values.Get(query)
	if queryV == "" {
		return "", errors.New("未找到对应值")
	}
	return queryV, nil
}

func RemoveDbholder(dsn string) string {
	reg := `&dbholder=[\w\$,]+&?`
	regObj, err := regexp.Compile(reg)
	if err != nil {
		return dsn
	}
	return regObj.ReplaceAllString(dsn, "&")
}

func Md5(str string) string {
	md5Str := md5.Sum([]byte(str))
	return string(md5Str[:])
}

func RemovePoint(dbholder string) string {
	return strings.Replace(dbholder, ".", "", 1)
}
