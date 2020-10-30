package onsite

const (
	OnsiteRoomJoinRequest = "ONSITE_ROOM_JOIN_REQUEST"
	OnsiteRoomJoinError   = "ONSITE_ROOM_JOIN_ERROR"
	OnsiteRoomJoinSuccess = "ONSITE_ROOM_JOIN_SUCCESS"

	OnsiteRoomLeaveRequest = "ONSITE_ROOM_LEAVE_REQUEST"
	OnsiteRoomLeaveError   = "ONSITE_ROOM_LEAVE_ERROR"
	OnsiteRoomLeaveSuccess = "ONSITE_ROOM_LEAVE_SUCCESS"

	OnsiteRoomUpdateAttendeeCounter = "ONSITE_ROOM_UPDATE_ATTENDEE_COUNTER"
)

type OnsitePayload struct {
	RoomID  string `json:"room_id,omitempty"`
	Counter *int   `json:"counter,omitempty"`
}
