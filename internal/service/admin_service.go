package service

import (
	"context"
	"encoding/json"
	"errors"

	"financial_system/internal/domain"
	"financial_system/internal/repository"
)

type Admin struct {
	logRepo repository.ActionLogRepository
}

func NewAdminService(logRepo repository.ActionLogRepository) *Admin {
	return &Admin{logRepo: logRepo}
}

// GetAllLogs возвращает все записи action_logs в порядке убывания даты.
func (a *Admin) GetAllLogs() ([]domain.ActionLog, error) {
	return a.logRepo.GetAll()
}

// UndoAction выполняет логический откат действия, описанного в action_logs.
func (a *Admin) UndoAction(logID int, deps *repository.Repositories) error {
	ctx := context.Background()

	// Получаем запись под блокировкой.
	entry, err := a.logRepo.GetByIDForUpdate(ctx, logID)
	if err != nil {
		return err
	}
	if entry.IsUndone {
		return domain.ErrApplicationAlreadyPaid // переиспользуем как "уже отменено" (можно завести отдельную ошибку)
	}

	// userID, от имени которого выполнялось исходное действие (если есть).
	actorUserID := 0
	if entry.UserID != nil {
		actorUserID = *entry.UserID
	}

	var details map[string]any
	if len(entry.Details) > 0 {
		if err := json.Unmarshal(entry.Details, &details); err != nil {
			return err
		}
	}

	switch entry.Action {
	// CLIENT
	case "client_open_account":
		id, ok := details["account_id"].(float64)
		if !ok {
			return errors.New("account_id not found in log details")
		}
		accountID := int(id)
		acc, err := deps.Account.GetAccountByID(accountID)
		if err != nil {
			return err
		}
		if acc.Balance != 0 || acc.IsBlocked {
			return errors.New("cannot undo: account is blocked or balance is not zero")
		}
		if err := deps.Account.SetAccountBlocked(accountID, true); err != nil {
			return err
		}

	case "client_close_account":
		id, ok := details["account_id"].(float64)
		if !ok {
			return errors.New("account_id not found in log details")
		}
		accountID := int(id)
		acc, err := deps.Account.GetAccountByID(accountID)
		if err != nil {
			return err
		}
		if !acc.IsBlocked {
			return errors.New("cannot undo: account is already unblocked")
		}
		if err := deps.Account.SetAccountBlocked(accountID, false); err != nil {
			return err
		}

	case "client_transfer_from_account":
		fromID := int(details["from_account_id"].(float64))
		amount := details["amount"].(float64)
		var toAccountID, toDepositID *int
		if v, ok := details["to_account_id"].(float64); ok {
			id := int(v)
			toAccountID = &id
		}
		if v, ok := details["to_deposit_id"].(float64); ok {
			id := int(v)
			toDepositID = &id
		}
		if toAccountID != nil {
			// undo: с получателя обратно на отправителя
			if err := deps.Account.TransferAccountToAccount(actorUserID, *toAccountID, fromID, amount); err != nil {
				return err
			}
		} else if toDepositID != nil {
			// undo: с вклада обратно на счёт
			if err := deps.Deposit.TransferDepositToAccount(actorUserID, *toDepositID, fromID, amount); err != nil {
				return err
			}
		} else {
			return errors.New("invalid transfer details")
		}

	case "client_open_deposit":
		id := int(details["deposit_id"].(float64))
		dep, err := deps.Deposit.GetDepositByID(id)
		if err != nil {
			return err
		}
		if dep.Balance != 0 || dep.IsBlocked {
			return errors.New("cannot undo: deposit is blocked or balance is not zero")
		}
		if err := deps.Deposit.SetDepositBlocked(id, true); err != nil {
			return err
		}

	case "client_close_deposit":
		id := int(details["deposit_id"].(float64))
		dep, err := deps.Deposit.GetDepositByID(id)
		if err != nil {
			return err
		}
		if !dep.IsBlocked {
			return errors.New("cannot undo: deposit is already unblocked")
		}
		if err := deps.Deposit.SetDepositBlocked(id, false); err != nil {
			return err
		}

	case "client_transfer_from_deposit":
		fromID := int(details["from_deposit_id"].(float64))
		amount := details["amount"].(float64)
		var toAccountID, toDepositID *int
		if v, ok := details["to_account_id"].(float64); ok {
			id := int(v)
			toAccountID = &id
		}
		if v, ok := details["to_deposit_id"].(float64); ok {
			id := int(v)
			toDepositID = &id
		}
		if toAccountID != nil {
			if err := deps.Account.TransferAccountToDeposit(actorUserID, *toAccountID, fromID, amount); err != nil {
				return err
			}
		} else if toDepositID != nil {
			if err := deps.Deposit.TransferDepositToDeposit(actorUserID, *toDepositID, fromID, amount); err != nil {
				return err
			}
		} else {
			return errors.New("invalid transfer details")
		}

	case "client_accumulate_deposit":
		fromAccountID := int(details["from_account_id"].(float64))
		depositID := int(details["deposit_id"].(float64))
		amount := details["amount"].(float64)
		// undo: вклад -> счёт
		if err := deps.Deposit.TransferDepositToAccount(actorUserID, depositID, fromAccountID, amount); err != nil {
			return err
		}

	case "client_salary_application_create":
		appID := int(details["application_id"].(float64))
		app, err := deps.SalaryApplication.GetByID(appID)
		if err != nil {
			return err
		}
		if app.Status != domain.SalaryApplicationStatusPending || app.PaidAt != nil {
			return errors.New("cannot undo: application is not pending or already paid")
		}
		// мягкий undo: пометить rejected
		if err := deps.SalaryApplication.UpdateStatus(appID, domain.SalaryApplicationStatusRejected); err != nil {
			return err
		}

	case "client_salary_receive":
		appID := int(details["application_id"].(float64))
		var toAccID, toDepID *int
		if v, ok := details["to_account_id"].(float64); ok {
			id := int(v)
			toAccID = &id
		}
		if v, ok := details["to_deposit_id"].(float64); ok {
			id := int(v)
			toDepID = &id
		}
		if err := deps.SalaryApplication.UndoPaySalary(appID, toAccID, toDepID); err != nil {
			return err
		}

	// MANAGER
	case "manager_approve_user":
		userID := int(details["target_user_id"].(float64))
		user, err := deps.User.GetUserByID(userID)
		if err != nil {
			return err
		}
		if !user.IsActive {
			return errors.New("cannot undo: user is already inactive")
		}
		if err := deps.User.SetUserActive(userID, false); err != nil {
			return err
		}

	case "manager_block_account":
		accountID := int(details["account_id"].(float64))
		return deps.Account.SetAccountBlocked(accountID, false)

	case "manager_unblock_account":
		accountID := int(details["account_id"].(float64))
		return deps.Account.SetAccountBlocked(accountID, true)

	case "manager_block_deposit":
		depositID := int(details["deposit_id"].(float64))
		return deps.Deposit.SetDepositBlocked(depositID, false)

	case "manager_unblock_deposit":
		depositID := int(details["deposit_id"].(float64))
		return deps.Deposit.SetDepositBlocked(depositID, true)

	case "manager_add_employee":
		enterpriseID := int(details["enterprise_id"].(float64))
		userID := int(details["user_id"].(float64))
		return deps.Enterprise.RemoveEmployee(enterpriseID, userID)

	case "manager_remove_employee":
		enterpriseID := int(details["enterprise_id"].(float64))
		userID := int(details["user_id"].(float64))
		return deps.Enterprise.AddEmployee(enterpriseID, userID)

	case "manager_approve_salary_application":
		appID := int(details["application_id"].(float64))
		app, err := deps.SalaryApplication.GetByID(appID)
		if err != nil {
			return err
		}
		if app.Status != domain.SalaryApplicationStatusApproved || app.PaidAt != nil {
			return errors.New("cannot undo: application is not in approved state or already paid")
		}
		if err := deps.SalaryApplication.UpdateStatus(appID, domain.SalaryApplicationStatusPending); err != nil {
			return err
		}

	default:
		return errors.New("undo is not supported for this action type")
	}

	if err := a.logRepo.MarkUndone(ctx, logID); err != nil {
		return err
	}

	return nil
}