package domain

// EnterpriseWithEmployees — предприятие и список user_id сотрудников.
type EnterpriseWithEmployees struct {
	Enterprise       Enterprise `json:"enterprise"`
	EmployeeUserIDs  []int      `json:"employee_user_ids"`
}
