package service

import (
	"financial_system/internal/domain"
	"financial_system/internal/repository"
)

type SalaryProject struct {
	enterpriseRepo repository.EnterpriseRepository
	salaryRepo     repository.SalaryApplicationRepository
}

func NewSalaryProjectService(enterpriseRepo repository.EnterpriseRepository, salaryRepo repository.SalaryApplicationRepository) *SalaryProject {
	return &SalaryProject{enterpriseRepo: enterpriseRepo, salaryRepo: salaryRepo}
}

// ApplyForSalaryProject создаёт заявку на зарплатный проект (status pending). Доступно только сотрудникам предприятия.
func (s *SalaryProject) ApplyForSalaryProject(userID, enterpriseID int, amount float64) (*domain.SalaryApplication, error) {
	if amount <= 0 {
		return nil, domain.ErrInvalidAmount
	}
	ok, err := s.enterpriseRepo.IsEmployee(enterpriseID, userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, domain.ErrNotEmployee
	}
	_, err = s.enterpriseRepo.GetEnterpriseByID(enterpriseID)
	if err != nil {
		return nil, err
	}

	app := &domain.SalaryApplication{
		UserID:       userID,
		EnterpriseID: enterpriseID,
		Amount:       amount,
		Status:       domain.SalaryApplicationStatusPending,
	}
	if err := s.salaryRepo.Create(app); err != nil {
		return nil, err
	}
	return app, nil
}

// ReceiveSalary выполняет выплату по одобренной заявке на указанный счёт или вклад.
func (s *SalaryProject) ReceiveSalary(userID, applicationID int, toAccountID, toDepositID *int) error {
	app, err := s.salaryRepo.GetByID(applicationID)
	if err != nil {
		return err
	}
	if app.UserID != userID {
		return domain.ErrForbidden
	}
	hasAccount := toAccountID != nil && *toAccountID > 0
	hasDeposit := toDepositID != nil && *toDepositID > 0
	if hasAccount && hasDeposit {
		return domain.ErrInvalidTransferTarget
	}
	if !hasAccount && !hasDeposit {
		return domain.ErrInvalidTransferTarget
	}
	var accID, depID *int
	if hasAccount {
		accID = toAccountID
	} else {
		depID = toDepositID
	}
	return s.salaryRepo.PaySalary(applicationID, accID, depID)
}
