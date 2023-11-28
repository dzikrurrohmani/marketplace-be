package utils

import (
	"fmt"
	"log"
	"net/http"
)

type AppError struct {
	ErrorCode    string
	ErrorMessage string
	ErrorType    int
}

func (e AppError) Error() string {
	return fmt.Sprintf("type: %d, code:%s, err:%s", e.ErrorType, e.ErrorCode, e.ErrorMessage)
}

func RequiredError() error {
	return AppError{
		ErrorCode:    "X01",
		ErrorMessage: "Input cant be empty",
		ErrorType:    http.StatusBadRequest,
	}
}

func DataNotFoundError() error {
	return AppError{
		ErrorCode:    "X02",
		ErrorMessage: "Data not found",
	}
}

func WrongPasswordError() error {
	return AppError{
		ErrorCode:    "X03",
		ErrorMessage: "Wrong password",
	}
}

func UnauthorizedError() error {
	return AppError{
		ErrorCode:    "X04",
		ErrorMessage: "Unauthorized user",
		ErrorType:    http.StatusUnauthorized,
	}
}

func InsufficientBallanceError() error {
	return AppError{
		ErrorCode:    "X07",
		ErrorMessage: "Insufficient balance",
	}
}

func RaiseError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func PanicError(err error) {
	if err != nil {
		panic(err)
	}
}
