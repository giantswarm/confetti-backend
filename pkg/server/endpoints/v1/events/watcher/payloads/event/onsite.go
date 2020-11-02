package event

const (
	// OnsiteRoomJoinRequest is sent by the client
	// when they try to join a onsite room.
	OnsiteRoomJoinRequest = "ONSITE_ROOM_JOIN_REQUEST"
	// OnsiteRoomJoinError is sent by the server
	// when there was a problem with the
	// client's room join request.
	OnsiteRoomJoinError = "ONSITE_ROOM_JOIN_ERROR"
	// OnsiteRoomJoinSuccess is sent by the server
	// as a confirmation that a client was able
	// to join a room successfully.
	OnsiteRoomJoinSuccess = "ONSITE_ROOM_JOIN_SUCCESS"

	// OnsiteRoomJoinRequest is sent by the client
	// when they try to leave a onsite room.
	OnsiteRoomLeaveRequest = "ONSITE_ROOM_LEAVE_REQUEST"
	// OnsiteRoomLeaveError is sent by the server
	// when there was a problem with the
	// client's room leave request.
	OnsiteRoomLeaveError = "ONSITE_ROOM_LEAVE_ERROR"
	// OnsiteRoomLeaveSuccess is sent by the server
	// as a confirmation that a client was able
	// to leave a room successfully.
	OnsiteRoomLeaveSuccess = "ONSITE_ROOM_LEAVE_SUCCESS"

	// OnsiteRoomUpdateAttendeeCounter is broadcasted
	// by the server to all clients, whenever a client
	// joins or leaves a onsite room.
	//
	// Also, this message is sent to each client,
	// when a new connection is established.
	OnsiteRoomUpdateAttendeeCounter = "ONSITE_ROOM_UPDATE_ATTENDEE_COUNTER"
)

// OnsitePayload contains all the onsite event-specific
// keys in a message payload.
type OnsitePayload struct {
	// RoomID represents a onsite event's room.
	RoomID string `json:"room_id,omitempty"`
	// AttendeeCounter represents the number
	// of attendees in a onsite's room.
	AttendeeCounter *int `json:"attendee_counter,omitempty"`
}

type ConfigurationOnsiteRoom struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	ConferenceURL string `json:"conference_url"`
}
