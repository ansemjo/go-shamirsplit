package main

import (
	"github.com/alexbakker/pkcs7"
)

// Pad pads a message to a multiple of blocksize using PKCS7
var Pad = pkcs7.Pad

// Unpad removes a PKCS7 padding
var Unpad = pkcs7.Unpad
