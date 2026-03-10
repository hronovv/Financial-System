package domain

// EnterpriseWithEmployees is an enterprise plus its employee user IDs.
type EnterpriseWithEmployees struct {
	Enterprise       Enterprise `json:"enterprise"`
	EmployeeUserIDs  []int      `json:"employee_user_ids"`
}
