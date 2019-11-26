package tradewars

//Can be owned by a specific session
//Callsign for persistant ownership, sessionId for validating requests

type OwnershipComponent struct {
	callsign  string
	sessionId int64
}
