// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// MeReader is a Reader for the Me structure.
type MeReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *MeReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewMeOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewMeOK creates a MeOK with default headers values
func NewMeOK() *MeOK {
	return &MeOK{}
}

/*MeOK handles this case with default header values.

ok
*/
type MeOK struct {
	Payload string
}

func (o *MeOK) Error() string {
	return fmt.Sprintf("[GET /me][%d] meOK  %+v", 200, o.Payload)
}

func (o *MeOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
