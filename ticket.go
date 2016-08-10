package zendesk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Ticket struct {
	AllowChannelback    bool          `json:"allow_channelback"`
	AssigneeId          int           `json:"assignee_id"`
	BrandId             int           `json:"brand_id"`
	CollaboratorIds     []int         `json:"collaborator_ids"`
	CreatedAt           string        `json:"created_at"`
	CustomFields        []CustomField `json:"custom_fields"`
	Description         string        `json:description"`
	DueAt               string        `json:"due_at"`
	ExternalId          string        `json:"external_id"`
	Fields              []interface{} `json:"fields"`
	FollowupIds         []int         `json:"followup_ids"`
	ForumTopicId        int           `json:"forum_topic_id"`
	GroupId             int           `json:"group_id"`
	HasIncidents        bool          `json:"has_incidents"`
	Id                  int           `json:"id"`
	OrganizationId      int           `json:"organization_id"`
	Priority            string        `json:"priority"`
	ProblemId           int           `json:"problem_id"`
	RawSubject          string        `json:"raw_subject"`
	Recipient           string        `json:"recipient"`
	RequesterId         int           `json:"requester_id"`
	SatisfactionRating  interface{}   `json:"satisfaction_rating"`
	SharingAgreementIds []int         `json:"sharing_agreement_ids"`
	Status              string        `json:"status"`
	Subject             string        `json:"subject"`
	SubmitterId         int           `json:"submitter_id"`
	Tags                []string      `json:"tags"`
	TicketFormId        int           `json:"ticket_form_id"`
	Type                string        `json:"type"`
	UpdatedAt           string        `json:"updated_at"`
	Url                 string        `json:"url"`
	Via                 Via           `json:"via"`
}

type TicketNIWrapper struct {
	Ticket interface{} `json:"ticket"`
}

type TicketsNIWrapper struct {
	Tickets interface{} `json:"tickets"`
}

type JobStatus struct {
	Id       string            `json:"id"`
	Message  string            `json:"string"`
	Progress int               `json:"progress"`
	Results  []JobStatusResult `json:"results"`
	Status   string            `json:"status"`
	Total    int               `json:"total"`
	Url      string            `json:"url"`
}

type JobStatusResult struct {
	Action  string `json:"action"`
	Errors  string `json:"errors"`
	Id      int    `json:"Id"`
	Status  string `json:"status"`
	Success bool   `json:"success"`
	Title   string `json:"title"`
}

type Audit struct {
	AuthorId  int           `json:"author_id"`
	CreatedAt string        `json:"created_at"`
	Events    []interface{} `json:"events"`
	Id        int           `json:"id"`
	Metadata  interface{}   `json:"metadata"`
	TicketId  int           `json:"ticket_id"`
	Via       Via           `json:"via"`
}

type Via struct {
	Channel string `json:"channel"`
	Source  Source `json:"source"`
}

type Source struct {
	From interface{} `json:"from"`
	Rel  interface{} `json:"rel"`
	To   interface{} `json:"to"`
}

type CustomField struct {
	Id    int         `json:"id"`
	Value interface{} `json:"value"`
}

type CreateTicketRes struct {
	Audit  Audit  `json:"audit"`
	Ticket Ticket `json:"ticket"`
}

func (a API) CreateTicket(req interface{}) (t CreateTicketRes, body []byte, res *http.Response, err error) {
	reqData, err := json.Marshal(TicketNIWrapper{Ticket: req})
	if err != nil {
		a.HandleError(err)
		return
	}

	body, res, err = a.Send(http.MethodPost, "tickets.json", reqData)
	if err != nil {
		a.HandleError(err)
		return
	}

	err = json.Unmarshal(body, &t)
	if err != nil {
		a.HandleError(err)
	}
	return
}

type CreateTicketAsyncRes struct {
	Ticket struct {
		Id int `json:"id"`
	} `json:"ticket"`
	JobStatus JobStatus `json:"job_status"`
}

func (a API) CreateTicketAsync(req interface{}) (r CreateTicketAsyncRes, body []byte, res *http.Response, err error) {
	reqData, err := json.Marshal(TicketNIWrapper{Ticket: req})
	if err != nil {
		a.HandleError(err)
		return
	}

	body, res, err = a.Send(http.MethodPost, "tickets.json?async=true", reqData)
	if err != nil {
		a.HandleError(err)
		return
	}

	err = json.Unmarshal(body, &r)
	if err != nil {
		a.HandleError(err)
	}
	return
}

func (a API) ShowTicket(id int) (t Ticket, body []byte, res *http.Response, err error) {
	body, res, err = a.Send(http.MethodGet, fmt.Sprintf("tickets/%v.json", id), nil)
	if err != nil {
		a.HandleError(err)
		return
	}

	w := TicketNIWrapper{Ticket: &t}
	err = json.Unmarshal(body, &w)
	if err != nil {
		a.HandleError(err)
	}
	return
}

func (a API) ShowTickets(ids []int) (t []Ticket, body []byte, res *http.Response, err error) {
	var idsString []string
	for _, v := range ids {
		idsString = append(idsString, strconv.Itoa(v))
	}

	body, res, err = a.Send(http.MethodGet, fmt.Sprintf("tickets/show_many.json?ids=%v", strings.Join(idsString, ",")), nil)
	if err != nil {
		a.HandleError(err)
		return
	}

	w := TicketsNIWrapper{Tickets: &t}
	err = json.Unmarshal(body, &w)
	if err != nil {
		a.HandleError(err)
	}
	return
}

func (a API) ListTickets(url, sortBy, sortOrder string) (t []Ticket, body []byte, res *http.Response, err error) {
	if len(url) == 0 {
		err = fmt.Errorf("Required argument's missed (url)")
		return
	}

	var query []string
	if len(sortBy) > 0 {
		query = append(query, "sort_by="+sortBy)
	}
	if len(sortOrder) > 0 {
		query = append(query, "sort_order="+sortOrder)
	}
	if len(query) > 0 {
		url = url + "?" + strings.Join(query, "&")
	}

	body, res, err = a.Send(http.MethodGet, url, nil)
	if err != nil {
		a.HandleError(err)
		return
	}

	w := TicketsNIWrapper{Tickets: &t}
	err = json.Unmarshal(body, &w)
	if err != nil {
		a.HandleError(err)
	}
	return
}
