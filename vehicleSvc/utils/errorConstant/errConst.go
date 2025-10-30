package errorConstant

import (
	"app/utils"
	"net/http"
)

var (
	ErrVehicleNotFound = utils.CustomError{
		StatusCode: http.StatusNotFound,
		Code:       http.StatusNotFound,
	}
)
