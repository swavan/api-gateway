package access

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
)

var (
	_ persist.Adapter          = new(Adapter)
	_ persist.FilteredAdapter  = new(Adapter)
	_ persist.BatchAdapter     = new(Adapter)
	_ persist.UpdatableAdapter = new(Adapter)
)

func NewAdapter(db *sql.DB, tableName string) (*Adapter, error) {
	return NewAdapterWithContext(context.Background(), db, tableName)
}

func NewAdapterWithContext(ctx context.Context, db *sql.DB, tableName string) (*Adapter, error) {
	if ctx == nil {
		return nil, errors.New("ctx is nil")
	}

	if db == nil {
		return nil, errors.New("db is nil")
	}

	// check db connection
	err := db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	d := getDao(db, tableName)

	if !d.IsTableExist(ctx) {
		if err = d.CreateTable(ctx); err != nil {
			return nil, err
		}
	}

	adapter := Adapter{dao: d}

	return &adapter, nil
}

type Adapter struct {
	dao      dao
	filtered interface{}
}

func (Adapter) loadPolicyLine(line rule, model model.Model) error {
	return persist.LoadPolicyArray(line.Data(), model)
}

func (Adapter) genArgs(pType string, rule []string) []interface{} {
	args := make([]interface{}, maxParameterCount)
	args[0] = pType

	for idx := range rule {
		args[idx+1] = strings.TrimSpace(rule[idx])
	}

	for idx := len(rule) + 1; idx < maxParameterCount; idx++ {
		args[idx] = ""
	}

	return args
}

func (adapter *Adapter) LoadPolicy(model model.Model) error {
	lines, err := adapter.dao.SelectAll(context.Background())
	if err != nil {
		return err
	}

	adapter.filtered = nil

	for _, line := range lines {
		if err = adapter.loadPolicyLine(line, model); err != nil {
			return err
		}
	}

	return nil
}

func (adapter Adapter) SavePolicy(model model.Model) error {
	if adapter.filtered != nil {
		return errors.New("could not save filtered policies")
	}

	args := make([][]interface{}, 0, 128)

	for pType, ast := range model["p"] {
		for _, rule := range ast.Policy {
			arg := adapter.genArgs(pType, rule)
			args = append(args, arg)
		}
	}

	for pType, ast := range model["g"] {
		for _, rule := range ast.Policy {
			arg := adapter.genArgs(pType, rule)
			args = append(args, arg)
		}
	}

	return adapter.dao.DeleteAllAndInsertRows(context.Background(), args)
}

// AddPolicy  add one policy rule to the storage.
func (adapter Adapter) AddPolicy(sec string, pType string, rule []string) error {
	args := adapter.genArgs(pType, rule)
	return adapter.dao.InsertRow(context.Background(), args...)
}

// AddPolicies  add multiple policy rules to the storage.
func (adapter Adapter) AddPolicies(sec string, pType string, rules [][]string) error {
	args := make([][]interface{}, 0, len(rules))

	for _, rule := range rules {
		arg := adapter.genArgs(pType, rule)
		args = append(args, arg)
	}

	return adapter.dao.InsertRows(context.Background(), args)
}

// RemovePolicy  remove policy rules from the storage.
func (adapter Adapter) RemovePolicy(sec, pType string, rule []string) error {
	return adapter.dao.DeleteByArgs(context.Background(), pType, rule)
}

// RemoveFilteredPolicy  remove policy rules that match the filter from the storage.
func (adapter Adapter) RemoveFilteredPolicy(sec string, pType string, fieldIndex int, fieldValues ...string) error {
	whereCondition, whereArgs := adapter.dao.GenFilteredCondition(pType, fieldIndex, fieldValues...)

	return adapter.dao.DeleteByCondition(context.Background(), whereCondition, whereArgs...)
}

func (adapter Adapter) RemovePolicies(sec string, pType string, rules [][]string) (err error) {
	args := make([][]interface{}, len(rules))

	for idx, rule := range rules {
		arg := adapter.genArgs(pType, rule)
		args[idx] = arg
	}

	return adapter.dao.DeleteRows(context.Background(), args)
}

// LoadFilteredPolicy  load policy rules that match the Filter.
// filterPtr must be a pointer.
func (adapter *Adapter) LoadFilteredPolicy(model model.Model, filterPtr interface{}) error {
	if filterPtr == nil {
		return adapter.LoadPolicy(model)
	}

	filter, ok := filterPtr.(*Filter)
	if !ok {
		return errors.New("invalid filter type")
	}

	lines, err := adapter.dao.SelectByFilter(context.Background(), filter.genData())
	if err != nil {
		return err
	}

	for _, line := range lines {
		if err = adapter.loadPolicyLine(line, model); err != nil {
			return err
		}
	}

	adapter.filtered = struct{}{}

	return nil
}

// IsFiltered  returns true if the loaded policy rules has been filtered.
func (adapter Adapter) IsFiltered() bool {
	return adapter.filtered != nil
}

// UpdatePolicy update a policy rule from storage.
// This is part of the Auto-Save feature.
func (adapter Adapter) UpdatePolicy(sec, pType string, oldRule, newPolicy []string) error {
	oldArgs := adapter.genArgs(pType, oldRule)
	newArgs := adapter.genArgs(pType, newPolicy)

	return adapter.dao.UpdateRow(context.Background(), append(newArgs, oldArgs...)...)
}

// UpdatePolicies updates policy rules to storage.
func (adapter Adapter) UpdatePolicies(sec, pType string, oldRules, newRules [][]string) (err error) {
	if len(oldRules) != len(newRules) {
		return errors.New("old rules size not equal to new rules size")
	}

	args := make([][]interface{}, 0, len(oldRules)+len(newRules))

	for idx := range oldRules {
		oldArgs := adapter.genArgs(pType, oldRules[idx])
		newArgs := adapter.genArgs(pType, newRules[idx])
		args = append(args, append(newArgs, oldArgs...))
	}

	return adapter.dao.UpdateRows(context.Background(), args)
}

// UpdateFilteredPolicies deletes old rules and adds new rules.
func (adapter Adapter) UpdateFilteredPolicies(sec, pType string, newPolicies [][]string, fieldIndex int, fieldValues ...string) (oldPolicies [][]string, err error) {
	whereCondition, whereArgs := adapter.dao.GenFilteredCondition(pType, fieldIndex, fieldValues...)

	var oldRules []rule
	oldRules, err = adapter.dao.SelectByCondition(context.Background(), whereCondition, whereArgs...)
	if err != nil {
		return
	}

	args := make([][]interface{}, 0, len(newPolicies))
	for _, policy := range newPolicies {
		arg := adapter.genArgs(pType, policy)
		args = append(args, arg)
	}

	if err = adapter.dao.UpdateFilteredRows(context.Background(), whereCondition, whereArgs, args); err != nil {
		return
	}

	oldPolicies = make([][]string, 0, len(oldRules))
	for _, rule := range oldRules {
		oldPolicies = append(oldPolicies, rule.Data())
	}

	return
}
