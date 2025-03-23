package types


type CreateMessage struct{
  Name string `json:"name"`
  Email string `json:"email"`
  Phone string `json:"phone"`
  Subject string `json:"subject"`
  Message string `json:"message"`
}
