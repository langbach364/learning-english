package main

import ()

type Word struct {
	Data string `json:"data"`
}

type AnswerData struct {
    Detail interface{}`json:"detail"`
	Structure string `json:"structure"`
}

type ErrorResponse struct {
    Status  int    `json:"status"`
    Message string `json:"message"`
    Code    string `json:"code"`
}

type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type LoginResponse struct {
    Token string `json:"token"`
}