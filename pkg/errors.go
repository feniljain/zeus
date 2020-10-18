package pkg

import "errors"

var (
	//ErrWrongFormat depicts wrong json format sent
	ErrWrongFormat = errors.New("error: Invalid format of data sent")
	//ErrInvalidMeetingID depicts wrong meeting id is sent
	ErrInvalidMeetingID = errors.New("error: Invalid meeting ID")
)
