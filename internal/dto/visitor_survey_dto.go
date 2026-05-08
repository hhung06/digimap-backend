package dto

type SubmitVisitorSurveyRequest struct {
	PublicKey      string `json:"public_key"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phone_number"`
	VisitorType    int    `json:"visitor_type"`
	BusinessName   string `json:"business_name"`
	Interests      []int  `json:"interests"`
	OtherInterests string `json:"other_interests"`
	IsConsented    bool   `json:"is_consented"`
}
