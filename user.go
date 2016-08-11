package zendesk

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type UserWrapper struct {
	User *User `json:"user"`
}

type UserNIWrapper struct {
	User interface{} `json:"user"`
}

type User struct {
	Active               bool                   `json:"active"`
	Alias                string                 `json:"alias"`
	ChatOnly             bool                   `json:"chat_only"`
	CreatedAt            string                 `json:"created_at"`
	CustomRoleId         int                    `json:"custom_role_id"`
	Details              string                 `json:"details"`
	Email                string                 `json:"email"`
	ExternalId           string                 `json:"external_id"`
	Id                   int                    `json:"id"`
	LastLoginAt          string                 `json:"last_login_at"`
	LocaleId             int                    `json:"locale_id"`
	Locale               string                 `json:"locale"`
	Moderator            bool                   `json:"moderator"`
	Name                 string                 `json:"name"`
	Notes                string                 `json:"notes"`
	OnlyPrivateComments  bool                   `json:"only_private_comments"`
	OrganizationId       int                    `json:"organization_id"`
	Phone                string                 `json:"phone"`
	Photo                Attachment             `json:"photo"`
	RestrictedAgent      bool                   `json:"restricted_agent"`
	Role                 string                 `json:"role"`
	SharedAgent          bool                   `json:"shared_agent"`
	SharedPhoneNumber    bool                   `json:"shared_phone_number"`
	Shared               bool                   `json:"shared"`
	Signature            string                 `json:"signature"`
	Suspended            bool                   `json:"suspended"`
	Tags                 []string               `json:"tags"`
	TicketRestriction    string                 `json:"ticket_restriction"`
	TimeZone             string                 `json:"time_zone"`
	TwoFactorAuthEnabled bool                   `json:"two_factor_auth_enabled"`
	UpdatedAt            string                 `json:"updated_at"`
	Url                  string                 `json:"url"`
	UserFields           map[string]interface{} `json:"user_fields"`
	Verified             bool                   `json:"verified"`
}

type Attachment struct {
	ContentType string        `json:"content_type"`
	ContentUrl  string        `json:"content_url"`
	FileName    string        `json:"file_name"`
	Id          int           `json:"id"`
	Inline      bool          `json:"inline"`
	Size        int           `json:"size"`
	Thumbnails  []*Attachment `json:"thumbnails,omitempty"`
}

func (a *API) CreateUser(req interface{}) (u User, body []byte, res *http.Response, err error) {
	reqData, err := json.Marshal(UserNIWrapper{User: req})
	if err != nil {
		a.HandleError(err)
		return
	}

	body, res, err = a.Send(http.MethodPost, "users.json", reqData)
	if err != nil {
		a.HandleError(err)
		return
	}

	w := UserWrapper{User: &u}
	err = json.Unmarshal(body, &w)
	if err != nil {
		a.HandleError(err)
	}
	return
}

func (a *API) ShowUser(userId int) (u User, body []byte, res *http.Response, err error) {
	body, res, err = a.Send(http.MethodGet, fmt.Sprintf("users/%v.json", userId), nil)
	if err != nil {
		a.HandleError(err)
		return
	}

	w := UserWrapper{User: &u}
	err = json.Unmarshal(body, &w)
	if err != nil {
		a.HandleError(err)
	}
	return
}

func (a *API) UpdateUser(userId int, req interface{}) (u User, body []byte, res *http.Response, err error) {
	reqData, err := json.Marshal(UserNIWrapper{User: req})
	if err != nil {
		a.HandleError(err)
		return
	}

	body, res, err = a.Send(http.MethodPut, fmt.Sprintf("users/%v.json", userId), reqData)
	if err != nil {
		a.HandleError(err)
		return
	}

	w := UserWrapper{User: &u}
	err = json.Unmarshal(body, &w)
	if err != nil {
		a.HandleError(err)
	}
	return
}

func (a *API) SetUserPassword(userId int, newPassword string) (body []byte, res *http.Response, err error) {
	r := struct {
		Password string `json:"password"`
	}{
		Password: newPassword,
	}

	reqData, err := json.Marshal(r)
	if err != nil {
		a.HandleError(err)
		return
	}

	body, res, err = a.Send(http.MethodPost, fmt.Sprintf("users/%v/password.json", userId), reqData)
	if err != nil {
		a.HandleError(err)
	}
	return
}

func (a *API) UpdateUserProfileImage(userId int, imagePath, imageLink string) (body []byte, res *http.Response, err error) {
	url := fmt.Sprintf("users/%v.json", userId)
	if len(imagePath) > 0 {
		body, res, err = a.SendFile(http.MethodPut, url, "user[photo][uploaded_data]", imagePath)
		return
	}

	if len(imageLink) > 0 {
		r := map[string]map[string]string{
			"user": map[string]string{
				"remote_photo_url": imageLink,
			},
		}

		var reqData []byte
		reqData, err = json.Marshal(r)
		if err != nil {
			a.HandleError(err)
			return
		}

		body, res, err = a.Send(http.MethodPut, url, reqData)
		return
	}

	err = fmt.Errorf("Required attributes are missed")
	return
}
