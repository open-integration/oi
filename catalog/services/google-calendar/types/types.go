// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    googleCalendarEventsArgumentsSchema, err := UnmarshalGoogleCalendarEventsArgumentsSchema(bytes)
//    bytes, err = googleCalendarEventsArgumentsSchema.Marshal()
//
//    returns, err := UnmarshalReturns(bytes)
//    bytes, err = returns.Marshal()

package types

import "encoding/json"

func UnmarshalGoogleCalendarEventsArgumentsSchema(data []byte) (GoogleCalendarEventsArgumentsSchema, error) {
	var r GoogleCalendarEventsArgumentsSchema
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *GoogleCalendarEventsArgumentsSchema) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func UnmarshalReturns(data []byte) (Returns, error) {
	var r Returns
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Returns) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type GoogleCalendarEventsArgumentsSchema struct {
	CalendarID              string         `json:"CalendarID"`                        // Google calendar-id
	ICalUID                 *string        `json:"ICalUID,omitempty"`                 // Specifies event ID in the iCalendar format to be included in the response. Optional.
	MaxAttendees            *int64         `json:"MaxAttendees,omitempty"`            // The maximum number of attendees to include in the response. If there are more than the; specified number of attendees, only the participant is returned. Optional.
	MaxResults              *int64         `json:"MaxResults,omitempty"`              // Maximum number of events returned on one result page. The number of events in the; resulting page may be less than this value, or none at all, even if there are more events; matching the query. Incomplete pages can be detected by a non-empty nextPageToken field; in the response. By default the value is 250 events. The page size can never be larger; than 2500 events. Optional.
	OrderBy                 *OrderBy       `json:"OrderBy,omitempty"`                 // The order of the events returned in the result. Optional. The default is an unspecified,; stable order.
	PrivateExtendedProperty *string        `json:"PrivateExtendedProperty,omitempty"` // Extended properties constraint specified as propertyName=value. Matches only private; properties. This parameter might be repeated multiple times to return events that match; all given constraints.
	Q                       *string        `json:"Q,omitempty"`                       // Free text search terms to find events that match these terms in any field, except for; extended properties. Optional.
	ServiceAccount          ServiceAccount `json:"ServiceAccount"`                    // service-account
	SharedExtendedProperty  *string        `json:"SharedExtendedProperty,omitempty"`  // Extended properties constraint specified as propertyName=value. Matches only shared; properties. This parameter might be repeated multiple times to return events that match; all given constraints.
	ShowDeleted             *bool          `json:"ShowDeleted,omitempty"`             // Whether to include deleted events (with status equals "cancelled") in the result.; Cancelled instances of recurring events (but not the underlying recurring event) will; still be included if showDeleted and singleEvents are both False. If showDeleted and; singleEvents are both True, only single instances of deleted events (but not the; underlying recurring events) are returned. Optional. The default is False.
	ShowHiddenInvitations   *bool          `json:"ShowHiddenInvitations,omitempty"`   // Whether to include hidden invitations in the result. Optional. The default is False.
	SingleEvents            *bool          `json:"SingleEvents,omitempty"`            // Whether to expand recurring events into instances and only return single one-off events; and instances of recurring events, but not the underlying recurring events themselves.; Optional. The default is False.
	TimeMax                 *string        `json:"TimeMax,omitempty"`                 // Upper bound (exclusive) for an event's start time to filter by. Optional. The default is; not to filter by start time. Must be an RFC3339 timestamp with mandatory time zone; offset, for example, 2011-06-03T10:00:00-07:00, 2011-06-03T10:00:00Z. Milliseconds may be; provided but are ignored. If timeMin is set, timeMax must be greater than timeMin.
	TimeMin                 *string        `json:"TimeMin,omitempty"`                 // Lower bound (exclusive) for an event's end time to filter by. Optional. The default is; not to filter by end time. Must be an RFC3339 timestamp with mandatory time zone offset,; for example, 2011-06-03T10:00:00-07:00, 2011-06-03T10:00:00Z. Milliseconds may be; provided but are ignored. If timeMax is set, timeMin must be smaller than timeMax.
	TimeZone                *string        `json:"TimeZone,omitempty"`                // Time zone used in the response. Optional. The default is the time zone of the calendar.
	UpdatedMin              *string        `json:"UpdatedMin,omitempty"`              // Lower bound for an event's last modification time (as a RFC3339 timestamp) to filter by.; When specified, entries deleted since this time will always be included regardless of; showDeleted. Optional. The default is not to filter by last modification time.
}

// service-account
type ServiceAccount struct {
	AuthProviderX509CERTURL string `json:"auth_provider_x509_cert_url"`
	AuthURI                 string `json:"auth_uri"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	ClientX509CERTURL       string `json:"client_x509_cert_url"`
	PrivateKey              string `json:"private_key"`
	PrivateKeyID            string `json:"private_key_id"`
	ProjectID               string `json:"project_id"`
	TokenURI                string `json:"token_uri"`
	Type                    string `json:"type"`
}

type Returns struct {
	Events []Event `json:"events,omitempty"`
}

type Event struct {
	AnyoneCanAddSelf        *bool               `json:"anyoneCanAddSelf,omitempty"`        // Whether anyone can invite themselves to the event (currently works for Google+ events; only). Optional. The default is False.
	Attachments             []EventAttachment   `json:"attachments,omitempty"`             // File attachments for the event. Currently only Google Drive attachments are supported.; In order to modify attachments the supportsAttachments request parameter should be set to; true.; There can be at most 25 attachments per event,
	Attendees               []EventAttendee     `json:"attendees,omitempty"`               // The attendees of the event. See the Events with attendees guide for more information on; scheduling events with other calendar users.
	AttendeesOmitted        *bool               `json:"attendeesOmitted,omitempty"`        // Whether attendees may have been omitted from the event's representation. When retrieving; an event, this may be due to a restriction specified by the maxAttendee query parameter.; When updating an event, this can be used to only update the participant's response.; Optional. The default is False.
	ColorID                 *string             `json:"colorId,omitempty"`                 // The color of the event. This is an ID referring to an entry in the event section of the; colors definition (see the  colors endpoint). Optional.
	ConferenceData          *ConferenceData     `json:"conferenceData,omitempty"`          // The conference-related information, such as details of a Google Meet conference. To; create new conference details use the createRequest field. To persist your changes,; remember to set the conferenceDataVersion request parameter to 1 for all event; modification requests.
	Created                 *string             `json:"created,omitempty"`                 // Creation time of the event (as a RFC3339 timestamp). Read-only.
	Creator                 *Creator            `json:"creator,omitempty"`                 // The creator of the event. Read-only.
	Description             *string             `json:"description,omitempty"`             // Description of the event. Can contain HTML. Optional.
	End                     *EventDateTime      `json:"end,omitempty"`                     // The (exclusive) end time of the event. For a recurring event, this is the end time of the; first instance.
	EndTimeUnspecified      *bool               `json:"endTimeUnspecified,omitempty"`      // Whether the end time is actually unspecified. An end time is still provided for; compatibility reasons, even if this attribute is set to True. The default is False.
	Etag                    *string             `json:"etag,omitempty"`                    // ETag of the resource.
	ExtendedProperties      *ExtendedProperties `json:"extendedProperties,omitempty"`      // Extended properties of the event.
	Gadget                  *Gadget             `json:"gadget,omitempty"`                  // A gadget that extends this event.
	GuestsCanInviteOthers   *bool               `json:"guestsCanInviteOthers,omitempty"`   // Whether attendees other than the organizer can invite others to the event. Optional. The; default is True.
	GuestsCanModify         *bool               `json:"guestsCanModify,omitempty"`         // Whether attendees other than the organizer can modify the event. Optional. The default is; False.
	GuestsCanSeeOtherGuests *bool               `json:"guestsCanSeeOtherGuests,omitempty"` // Whether attendees other than the organizer can see who the event's attendees are.; Optional. The default is True.
	HangoutLink             *string             `json:"hangoutLink,omitempty"`             // An absolute link to the Google+ hangout associated with this event. Read-only.
	HTMLLink                *string             `json:"htmlLink,omitempty"`                // An absolute link to this event in the Google Calendar Web UI. Read-only.
	ICalUID                 *string             `json:"iCalUID,omitempty"`                 // Event unique identifier as defined in RFC5545. It is used to uniquely identify events; accross calendaring systems and must be supplied when importing events via the import; method.; Note that the icalUID and the id are not identical and only one of them should be; supplied at event creation time. One difference in their semantics is that in recurring; events, all occurrences of one event have different ids while they all share the same; icalUIDs.
	ID                      *string             `json:"id,omitempty"`                      // Opaque identifier of the event. When creating new single or recurring events, you can; specify their IDs. Provided IDs must follow these rules:; - characters allowed in the ID are those used in base32hex encoding, i.e. lowercase; letters a-v and digits 0-9, see section 3.1.2 in RFC2938; - the length of the ID must be between 5 and 1024 characters; - the ID must be unique per calendar  Due to the globally distributed nature of the; system, we cannot guarantee that ID collisions will be detected at event creation time.; To minimize the risk of collisions we recommend using an established UUID algorithm such; as one described in RFC4122.; If you do not specify an ID, it will be automatically generated by the server.; Note that the icalUID and the id are not identical and only one of them should be; supplied at event creation time. One difference in their semantics is that in recurring; events, all occurrences of one event have different ids while they all share the same; icalUIDs.
	Kind                    *string             `json:"kind,omitempty"`                    // Type of the resource ("calendar#event").
	Location                *string             `json:"location,omitempty"`                // Geographic location of the event as free-form text. Optional.
	Locked                  *bool               `json:"locked,omitempty"`                  // Whether this is a locked event copy where no changes can be made to the main event fields; "summary", "description", "location", "start", "end" or "recurrence". The default is; False. Read-Only.
	Organizer               *Organizer          `json:"organizer,omitempty"`               // The organizer of the event. If the organizer is also an attendee, this is indicated with; a separate entry in attendees with the organizer field set to True. To change the; organizer, use the move operation. Read-only, except when importing an event.
	OriginalStartTime       *EventDateTime      `json:"originalStartTime,omitempty"`       // For an instance of a recurring event, this is the time at which this event would start; according to the recurrence data in the recurring event identified by recurringEventId.; It uniquely identifies the instance within the recurring event series even if the; instance was moved to a different time. Immutable.
	PrivateCopy             *bool               `json:"privateCopy,omitempty"`             // If set to True, Event propagation is disabled. Note that it is not the same thing as; Private event properties. Optional. Immutable. The default is False.
	Recurrence              []string            `json:"recurrence,omitempty"`              // List of RRULE, EXRULE, RDATE and EXDATE lines for a recurring event, as specified in; RFC5545. Note that DTSTART and DTEND lines are not allowed in this field; event start and; end times are specified in the start and end fields. This field is omitted for single; events or instances of recurring events.
	RecurringEventID        *string             `json:"recurringEventId,omitempty"`        // For an instance of a recurring event, this is the id of the recurring event to which this; instance belongs. Immutable.
	Reminders               *Reminders          `json:"reminders,omitempty"`               // Information about the event's reminders for the authenticated user.
	Sequence                *int64              `json:"sequence,omitempty"`                // Sequence number as per iCalendar.
	Source                  *Source             `json:"source,omitempty"`                  // Source from which the event was created. For example, a web page, an email message or any; document identifiable by an URL with HTTP or HTTPS scheme. Can only be seen or modified; by the creator of the event.
	Start                   *EventDateTime      `json:"start,omitempty"`                   // The (inclusive) start time of the event. For a recurring event, this is the start time of; the first instance.
	Status                  *string             `json:"status,omitempty"`                  // Status of the event. Optional. Possible values are:; - "confirmed" - The event is confirmed. This is the default status.; - "tentative" - The event is tentatively confirmed.; - "cancelled" - The event is cancelled (deleted). The list method returns cancelled; events only on incremental sync (when syncToken or updatedMin are specified) or if the; showDeleted flag is set to true. The get method always returns them.; A cancelled status represents two different states depending on the event type:; - Cancelled exceptions of an uncancelled recurring event indicate that this instance; should no longer be presented to the user. Clients should store these events for the; lifetime of the parent recurring event.; Cancelled exceptions are only guaranteed to have values for the id, recurringEventId and; originalStartTime fields populated. The other fields might be empty.; - All other cancelled events represent deleted events. Clients should remove their; locally synced copies. Such cancelled events will eventually disappear, so do not rely on; them being available indefinitely.; Deleted events are only guaranteed to have the id field populated.   On the organizer's; calendar, cancelled events continue to expose event details (summary, location, etc.) so; that they can be restored (undeleted). Similarly, the events to which the user was; invited and that they manually removed continue to provide details. However, incremental; sync requests with showDeleted set to false will not return these details.; If an event changes its organizer (for example via the move operation) and the original; organizer is not on the attendee list, it will leave behind a cancelled event where only; the id field is guaranteed to be populated.
	Summary                 *string             `json:"summary,omitempty"`                 // Title of the event.
	Transparency            *string             `json:"transparency,omitempty"`            // Whether the event blocks time on the calendar. Optional. Possible values are:; - "opaque" - Default value. The event does block time on the calendar. This is equivalent; to setting Show me as to Busy in the Calendar UI.; - "transparent" - The event does not block time on the calendar. This is equivalent to; setting Show me as to Available in the Calendar UI.
	Updated                 *string             `json:"updated,omitempty"`                 // Last modification time of the event (as a RFC3339 timestamp). Read-only.
	Visibility              *string             `json:"visibility,omitempty"`              // Visibility of the event. Optional. Possible values are:; - "default" - Uses the default visibility for events on the calendar. This is the default; value.; - "public" - The event is public and event details are visible to all readers of the; calendar.; - "private" - The event is private and only event attendees may view event details.; - "confidential" - The event is private. This value is provided for compatibility reasons.
}

type EventAttachment struct {
	FileID   *string `json:"fileId,omitempty"`   // ID of the attached file. Read-only.; For Google Drive files, this is the ID of the corresponding Files resource entry in the; Drive API.
	FileURL  *string `json:"fileUrl,omitempty"`  // URL link to the attachment.; For adding Google Drive file attachments use the same format as in alternateLink property; of the Files resource in the Drive API.; Required when adding an attachment.
	IconLink *string `json:"iconLink,omitempty"` // URL link to the attachment's icon. Read-only.
	MIMEType *string `json:"mimeType,omitempty"` // Internet media type (MIME type) of the attachment.
	Title    *string `json:"title,omitempty"`    // Attachment title.
}

type EventAttendee struct {
	AdditionalGuests *int64  `json:"additionalGuests,omitempty"` // Number of additional guests. Optional. The default is 0.
	Comment          *string `json:"comment,omitempty"`          // The attendee's response comment. Optional.
	DisplayName      *string `json:"displayName,omitempty"`      // The attendee's name, if available. Optional.
	Email            *string `json:"email,omitempty"`            // The attendee's email address, if available. This field must be present when adding an; attendee. It must be a valid email address as per RFC5322.; Required when adding an attendee.
	ID               *string `json:"id,omitempty"`               // The attendee's Profile ID, if available. It corresponds to the id field in the People; collection of the Google+ API
	Optional         *bool   `json:"optional,omitempty"`         // Whether this is an optional attendee. Optional. The default is False.
	Organizer        *bool   `json:"organizer,omitempty"`        // Whether the attendee is the organizer of the event. Read-only. The default is False.
	Resource         *bool   `json:"resource,omitempty"`         // Whether the attendee is a resource. Can only be set when the attendee is added to the; event for the first time. Subsequent modifications are ignored. Optional. The default is; False.
	ResponseStatus   *string `json:"responseStatus,omitempty"`   // The attendee's response status. Possible values are:; - "needsAction" - The attendee has not responded to the invitation.; - "declined" - The attendee has declined the invitation.; - "tentative" - The attendee has tentatively accepted the invitation.; - "accepted" - The attendee has accepted the invitation.
	Self             *bool   `json:"self,omitempty"`             // Whether this entry represents the calendar on which this copy of the event appears.; Read-only. The default is False.
}

// The conference-related information, such as details of a Google Meet conference. To
// create new conference details use the createRequest field. To persist your changes,
// remember to set the conferenceDataVersion request parameter to 1 for all event
// modification requests.
type ConferenceData struct {
	ConferenceID       *string                  `json:"conferenceId,omitempty"`       // The ID of the conference.; Can be used by developers to keep track of conferences, should not be displayed to users.; Values for solution types:; - "eventHangout": unset.; - "eventNamedHangout": the name of the Hangout.; - "hangoutsMeet": the 10-letter meeting code, for example "aaa-bbbb-ccc".; - "addOn": defined by 3P conference provider.  Optional.
	ConferenceSolution *ConferenceSolution      `json:"conferenceSolution,omitempty"` // The conference solution, such as Hangouts or Google Meet.; Unset for a conference with a failed create request.; Either conferenceSolution and at least one entryPoint, or createRequest is required.
	CreateRequest      *CreateConferenceRequest `json:"createRequest,omitempty"`      // A request to generate a new conference and attach it to the event. The data is generated; asynchronously. To see whether the data is present check the status field.; Either conferenceSolution and at least one entryPoint, or createRequest is required.
	EntryPoints        []EntryPoint             `json:"entryPoints,omitempty"`        // Information about individual conference entry points, such as URLs or phone numbers.; All of them must belong to the same conference.; Either conferenceSolution and at least one entryPoint, or createRequest is required.
	Notes              *string                  `json:"notes,omitempty"`              // Additional notes (such as instructions from the domain administrator, legal notices) to; display to the user. Can contain HTML. The maximum length is 2048 characters. Optional.
	Parameters         *ConferenceParameters    `json:"parameters,omitempty"`         // Additional properties related to a conference. An example would be a solution-specific; setting for enabling video streaming.
	Signature          *string                  `json:"signature,omitempty"`          // The signature of the conference data.; Generated on server side. Must be preserved while copying the conference data between; events, otherwise the conference data will not be copied.; Unset for a conference with a failed create request.; Optional for a conference with a pending create request.
}

// The conference solution, such as Hangouts or Google Meet.
// Unset for a conference with a failed create request.
// Either conferenceSolution and at least one entryPoint, or createRequest is required.
type ConferenceSolution struct {
	IconURI *string                `json:"iconUri,omitempty"` // The user-visible icon for this solution.
	Key     *ConferenceSolutionKey `json:"key,omitempty"`     // The key which can uniquely identify the conference solution for this event.
	Name    *string                `json:"name,omitempty"`    // The user-visible name of this solution. Not localized.
}

// The key which can uniquely identify the conference solution for this event.
//
// The conference solution, such as Hangouts or Google Meet.
type ConferenceSolutionKey struct {
	Type *string `json:"type,omitempty"` // The conference solution type.; If a client encounters an unfamiliar or empty type, it should still be able to display; the entry points. However, it should disallow modifications.; The possible values are:; - "eventHangout" for Hangouts for consumers (http://hangouts.google.com); - "eventNamedHangout" for classic Hangouts for G Suite users (http://hangouts.google.com); - "hangoutsMeet" for Google Meet (http://meet.google.com); - "addOn" for 3P conference providers
}

// A request to generate a new conference and attach it to the event. The data is generated
// asynchronously. To see whether the data is present check the status field.
// Either conferenceSolution and at least one entryPoint, or createRequest is required.
type CreateConferenceRequest struct {
	ConferenceSolutionKey *ConferenceSolutionKey   `json:"conferenceSolutionKey,omitempty"` // The conference solution, such as Hangouts or Google Meet.
	RequestID             *string                  `json:"requestId,omitempty"`             // The client-generated unique ID for this request.; Clients should regenerate this ID for every new request. If an ID provided is the same as; for the previous request, the request is ignored.
	Status                *ConferenceRequestStatus `json:"status,omitempty"`                // The status of the conference create request.
}

// The status of the conference create request.
type ConferenceRequestStatus struct {
	StatusCode *string `json:"statusCode,omitempty"` // The current status of the conference create request. Read-only.; The possible values are:; - "pending": the conference create request is still being processed.; - "success": the conference create request succeeded, the entry points are populated.; - "failure": the conference create request failed, there are no entry points.
}

type EntryPoint struct {
	AccessCode         *string  `json:"accessCode,omitempty"`         // The access code to access the conference. The maximum length is 128 characters.; When creating new conference data, populate only the subset of {meetingCode, accessCode,; passcode, password, pin} fields that match the terminology that the conference provider; uses. Only the populated fields should be displayed.; Optional.
	EntryPointFeatures []string `json:"entryPointFeatures,omitempty"` // Features of the entry point, such as being toll or toll-free. One entry point can have; multiple features. However, toll and toll-free cannot be both set on the same entry point.
	EntryPointType     *string  `json:"entryPointType,omitempty"`     // The type of the conference entry point.; Possible values are:; - "video" - joining a conference over HTTP. A conference can have zero or one video entry; point.; - "phone" - joining a conference by dialing a phone number. A conference can have zero or; more phone entry points.; - "sip" - joining a conference over SIP. A conference can have zero or one sip entry; point.; - "more" - further conference joining instructions, for example additional phone numbers.; A conference can have zero or one more entry point. A conference with only a more entry; point is not a valid conference.
	Label              *string  `json:"label,omitempty"`              // The label for the URI. Visible to end users. Not localized. The maximum length is 512; characters.; Examples:; - for video: meet.google.com/aaa-bbbb-ccc; - for phone: +1 123 268 2601; - for sip: 12345678@altostrat.com; - for more: should not be filled; Optional.
	MeetingCode        *string  `json:"meetingCode,omitempty"`        // The meeting code to access the conference. The maximum length is 128 characters.; When creating new conference data, populate only the subset of {meetingCode, accessCode,; passcode, password, pin} fields that match the terminology that the conference provider; uses. Only the populated fields should be displayed.; Optional.
	Passcode           *string  `json:"passcode,omitempty"`           // The passcode to access the conference. The maximum length is 128 characters.; When creating new conference data, populate only the subset of {meetingCode, accessCode,; passcode, password, pin} fields that match the terminology that the conference provider; uses. Only the populated fields should be displayed.
	Password           *string  `json:"password,omitempty"`           // The password to access the conference. The maximum length is 128 characters.; When creating new conference data, populate only the subset of {meetingCode, accessCode,; passcode, password, pin} fields that match the terminology that the conference provider; uses. Only the populated fields should be displayed.; Optional.
	Pin                *string  `json:"pin,omitempty"`                // The PIN to access the conference. The maximum length is 128 characters.; When creating new conference data, populate only the subset of {meetingCode, accessCode,; passcode, password, pin} fields that match the terminology that the conference provider; uses. Only the populated fields should be displayed.; Optional.
	RegionCode         *string  `json:"regionCode,omitempty"`         // The CLDR/ISO 3166 region code for the country associated with this phone access. Example:; "SE" for Sweden.; Calendar backend will populate this field only for EntryPointType.PHONE.
	URI                *string  `json:"uri,omitempty"`                // The URI of the entry point. The maximum length is 1300 characters.; Format:; - for video, http: or https: schema is required.; - for phone, tel: schema is required. The URI should include the entire dial sequence; (e.g., tel:+12345678900,,,123456789;1234).; - for sip, sip: schema is required, e.g., sip:12345678@myprovider.com.; - for more, http: or https: schema is required.
}

// Additional properties related to a conference. An example would be a solution-specific
// setting for enabling video streaming.
type ConferenceParameters struct {
	AddOnParameters *ConferenceParametersAddOnParameters `json:"addOnParameters,omitempty"` // Additional add-on specific data.
}

// Additional add-on specific data.
type ConferenceParametersAddOnParameters struct {
	Parameters map[string]string `json:"parameters,omitempty"`
}

// The creator of the event. Read-only.
type Creator struct {
	DisplayName *string `json:"displayName,omitempty"` // The creator's name, if available.
	Email       *string `json:"email,omitempty"`       // The creator's email address, if available.
	ID          *string `json:"id,omitempty"`          // The creator's Profile ID, if available. It corresponds to the id field in the People; collection of the Google+ API
	Self        *bool   `json:"self,omitempty"`        // Whether the creator corresponds to the calendar on which this copy of the event appears.; Read-only. The default is False.
}

// The (exclusive) end time of the event. For a recurring event, this is the end time of the
// first instance.
//
// For an instance of a recurring event, this is the time at which this event would start
// according to the recurrence data in the recurring event identified by recurringEventId.
// It uniquely identifies the instance within the recurring event series even if the
// instance was moved to a different time. Immutable.
//
// The (inclusive) start time of the event. For a recurring event, this is the start time of
// the first instance.
type EventDateTime struct {
	Date     *string `json:"date,omitempty"`     // The date, in the format "yyyy-mm-dd", if this is an all-day event.
	DateTime *string `json:"dateTime,omitempty"` // The time, as a combined date-time value (formatted according to RFC3339). A time zone; offset is required unless a time zone is explicitly specified in timeZone.
	TimeZone *string `json:"timeZone,omitempty"` // The time zone in which the time is specified. (Formatted as an IANA Time Zone Database; name, e.g. "Europe/Zurich".) For recurring events this field is required and specifies; the time zone in which the recurrence is expanded. For single events this field is; optional and indicates a custom time zone for the event start/end.
}

// Extended properties of the event.
type ExtendedProperties struct {
	Private map[string]string `json:"private,omitempty"` // Properties that are private to the copy of the event that appears on this calendar.
	Shared  map[string]string `json:"shared,omitempty"`  // Properties that are shared between copies of the event on other attendees' calendars.
}

// A gadget that extends this event.
type Gadget struct {
	Display     *string           `json:"display,omitempty"`     // The gadget's display mode. Optional. Possible values are:; - "icon" - The gadget displays next to the event's title in the calendar view.; - "chip" - The gadget displays when the event is clicked.
	Height      *int64            `json:"height,omitempty"`      // The gadget's height in pixels. The height must be an integer greater than 0. Optional.
	IconLink    *string           `json:"iconLink,omitempty"`    // The gadget's icon URL. The URL scheme must be HTTPS.
	Link        *string           `json:"link,omitempty"`        // The gadget's URL. The URL scheme must be HTTPS.
	Preferences map[string]string `json:"preferences,omitempty"` // Preferences.
	Title       *string           `json:"title,omitempty"`       // The gadget's title.
	Type        *string           `json:"type,omitempty"`        // The gadget's type.
	Width       *int64            `json:"width,omitempty"`       // The gadget's width in pixels. The width must be an integer greater than 0. Optional.
}

// The organizer of the event. If the organizer is also an attendee, this is indicated with
// a separate entry in attendees with the organizer field set to True. To change the
// organizer, use the move operation. Read-only, except when importing an event.
type Organizer struct {
	DisplayName *string `json:"displayName,omitempty"` // The organizer's name, if available.
	Email       *string `json:"email,omitempty"`       // The organizer's email address, if available. It must be a valid email address as per; RFC5322.
	ID          *string `json:"id,omitempty"`          // The organizer's Profile ID, if available. It corresponds to the id field in the People; collection of the Google+ API
	Self        *bool   `json:"self,omitempty"`        // Whether the organizer corresponds to the calendar on which this copy of the event; appears. Read-only. The default is False.
}

// Information about the event's reminders for the authenticated user.
type Reminders struct {
	Overrides  []EventReminder `json:"overrides,omitempty"`  // If the event doesn't use the default reminders, this lists the reminders specific to the; event, or, if not set, indicates that no reminders are set for this event. The maximum; number of override reminders is 5.
	UseDefault *bool           `json:"useDefault,omitempty"` // Whether the default reminders of the calendar apply to the event.
}

type EventReminder struct {
	Method  *string `json:"method,omitempty"`  // The method used by this reminder. Possible values are:; - "email" - Reminders are sent via email.; - "popup" - Reminders are sent via a UI popup.; Required when adding a reminder.
	Minutes *int64  `json:"minutes,omitempty"` // Number of minutes before the start of the event when the reminder should trigger. Valid; values are between 0 and 40320 (4 weeks in minutes).; Required when adding a reminder.
}

// Source from which the event was created. For example, a web page, an email message or any
// document identifiable by an URL with HTTP or HTTPS scheme. Can only be seen or modified
// by the creator of the event.
type Source struct {
	Title *string `json:"title,omitempty"` // Title of the source; for example a title of a web page or an email subject.
	URL   *string `json:"url,omitempty"`   // URL of the source pointing to a resource. The URL scheme must be HTTP or HTTPS.
}

// The order of the events returned in the result. Optional. The default is an unspecified,
// stable order.
type OrderBy string

const (
	StartTime OrderBy = "startTime"
	Updated   OrderBy = "updated"
)
