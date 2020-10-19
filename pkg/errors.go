package pkg

import "errors"

var (
	//ErrWrongFormat depicts wrong json format sent
	ErrWrongFormat = errors.New("error: Invalid format of data sent")

	//ErrInvalidMeetingID depicts wrong meeting id is sent
	ErrInvalidMeetingID = errors.New("error: Invalid meeting ID")

	//ErrParticipantEmail depicts wrong participant name is sent
	ErrParticipantEmail = errors.New("error: Invalid participant email")

	//ErrInternalServer depicts arbitary errors
	ErrInternalServer = errors.New("error: Internal Server Error")

	//ErrWrongTimestampFormat depicts wrong timestamp format
	ErrWrongTimestampFormat = errors.New("error: Invalid format of date and time sent, please check the example for correct format")
)
