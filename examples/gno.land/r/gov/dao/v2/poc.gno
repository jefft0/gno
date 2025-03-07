package govdao

import (
	"errors"
	"std"

	"gno.land/p/demo/combinederr"
	"gno.land/p/demo/dao"
	"gno.land/p/demo/membstore"
	"gno.land/p/gov/executor"
)

var errNoChangesProposed = errors.New("no set changes proposed")

// NewGovDAOExecutor creates the govdao wrapped callback executor
func (_ DAO) NewGovDAOExecutor(cb func() error) dao.Executor {
	if cb == nil {
		panic(errNoChangesProposed)
	}

	return executor.NewCallbackExecutor(
		cb,
		std.CurrentRealm().PkgPath(),
	)
}

// NewMemberPropExecutor returns the GOVDAO member change executor
func (_ DAO) NewMemberPropExecutor(changesFn func() []membstore.Member) dao.Executor {
	if changesFn == nil {
		panic(errNoChangesProposed)
	}

	callback := func() error {
		errs := &combinederr.CombinedError{}
		cbMembers := changesFn()

		for _, member := range cbMembers {
			switch {
			case !members.IsMember(member.Address):
				// Addition request
				err := members.AddMember(member)

				errs.Add(err)
			case member.VotingPower == 0:
				// Remove request
				err := members.UpdateMember(member.Address, membstore.Member{
					Address:     member.Address,
					VotingPower: 0, // 0 indicated removal
				})

				errs.Add(err)
			default:
				// Update request
				err := members.UpdateMember(member.Address, member)

				errs.Add(err)
			}
		}

		// Check if there were any execution errors
		if errs.Size() == 0 {
			return nil
		}

		return errs
	}

	return GovDAO.NewGovDAOExecutor(callback)
}

func (_ DAO) NewMembStoreImplExecutor(changeFn func() membstore.MemberStore) dao.Executor {
	if changeFn == nil {
		panic(errNoChangesProposed)
	}

	callback := func() error {
		setMembStoreImpl(changeFn())

		return nil
	}

	return GovDAO.NewGovDAOExecutor(callback)
}

// setMembStoreImpl sets a new dao.MembStore implementation
func setMembStoreImpl(impl membstore.MemberStore) {
	if impl == nil {
		panic("invalid member store")
	}

	members = impl
}
