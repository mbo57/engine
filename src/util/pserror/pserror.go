package pserror

import "errors"

var NotFoundIndex = errors.New("index not found")
var NotFoundDoc = errors.New("doc not found")

func IsNotFoundIndex(err error) bool {
	return errors.Is(err, NotFoundIndex)
}

func IsNotFoundDoc(err error) bool {
	return errors.Is(err, NotFoundDoc)
}
