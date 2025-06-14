// Code generated by ent, DO NOT EDIT.

package user

import (
	"stride-wars-app/ent/predicate"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
)

// ID filters vertices based on their ID field.
func ID(id uuid.UUID) predicate.User {
	return predicate.User(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uuid.UUID) predicate.User {
	return predicate.User(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uuid.UUID) predicate.User {
	return predicate.User(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uuid.UUID) predicate.User {
	return predicate.User(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uuid.UUID) predicate.User {
	return predicate.User(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uuid.UUID) predicate.User {
	return predicate.User(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uuid.UUID) predicate.User {
	return predicate.User(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uuid.UUID) predicate.User {
	return predicate.User(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uuid.UUID) predicate.User {
	return predicate.User(sql.FieldLTE(FieldID, id))
}

// ExternalUser applies equality check predicate on the "external_user" field. It's identical to ExternalUserEQ.
func ExternalUser(v uuid.UUID) predicate.User {
	return predicate.User(sql.FieldEQ(FieldExternalUser, v))
}

// Username applies equality check predicate on the "username" field. It's identical to UsernameEQ.
func Username(v string) predicate.User {
	return predicate.User(sql.FieldEQ(FieldUsername, v))
}

// ExternalUserEQ applies the EQ predicate on the "external_user" field.
func ExternalUserEQ(v uuid.UUID) predicate.User {
	return predicate.User(sql.FieldEQ(FieldExternalUser, v))
}

// ExternalUserNEQ applies the NEQ predicate on the "external_user" field.
func ExternalUserNEQ(v uuid.UUID) predicate.User {
	return predicate.User(sql.FieldNEQ(FieldExternalUser, v))
}

// ExternalUserIn applies the In predicate on the "external_user" field.
func ExternalUserIn(vs ...uuid.UUID) predicate.User {
	return predicate.User(sql.FieldIn(FieldExternalUser, vs...))
}

// ExternalUserNotIn applies the NotIn predicate on the "external_user" field.
func ExternalUserNotIn(vs ...uuid.UUID) predicate.User {
	return predicate.User(sql.FieldNotIn(FieldExternalUser, vs...))
}

// ExternalUserGT applies the GT predicate on the "external_user" field.
func ExternalUserGT(v uuid.UUID) predicate.User {
	return predicate.User(sql.FieldGT(FieldExternalUser, v))
}

// ExternalUserGTE applies the GTE predicate on the "external_user" field.
func ExternalUserGTE(v uuid.UUID) predicate.User {
	return predicate.User(sql.FieldGTE(FieldExternalUser, v))
}

// ExternalUserLT applies the LT predicate on the "external_user" field.
func ExternalUserLT(v uuid.UUID) predicate.User {
	return predicate.User(sql.FieldLT(FieldExternalUser, v))
}

// ExternalUserLTE applies the LTE predicate on the "external_user" field.
func ExternalUserLTE(v uuid.UUID) predicate.User {
	return predicate.User(sql.FieldLTE(FieldExternalUser, v))
}

// UsernameEQ applies the EQ predicate on the "username" field.
func UsernameEQ(v string) predicate.User {
	return predicate.User(sql.FieldEQ(FieldUsername, v))
}

// UsernameNEQ applies the NEQ predicate on the "username" field.
func UsernameNEQ(v string) predicate.User {
	return predicate.User(sql.FieldNEQ(FieldUsername, v))
}

// UsernameIn applies the In predicate on the "username" field.
func UsernameIn(vs ...string) predicate.User {
	return predicate.User(sql.FieldIn(FieldUsername, vs...))
}

// UsernameNotIn applies the NotIn predicate on the "username" field.
func UsernameNotIn(vs ...string) predicate.User {
	return predicate.User(sql.FieldNotIn(FieldUsername, vs...))
}

// UsernameGT applies the GT predicate on the "username" field.
func UsernameGT(v string) predicate.User {
	return predicate.User(sql.FieldGT(FieldUsername, v))
}

// UsernameGTE applies the GTE predicate on the "username" field.
func UsernameGTE(v string) predicate.User {
	return predicate.User(sql.FieldGTE(FieldUsername, v))
}

// UsernameLT applies the LT predicate on the "username" field.
func UsernameLT(v string) predicate.User {
	return predicate.User(sql.FieldLT(FieldUsername, v))
}

// UsernameLTE applies the LTE predicate on the "username" field.
func UsernameLTE(v string) predicate.User {
	return predicate.User(sql.FieldLTE(FieldUsername, v))
}

// UsernameContains applies the Contains predicate on the "username" field.
func UsernameContains(v string) predicate.User {
	return predicate.User(sql.FieldContains(FieldUsername, v))
}

// UsernameHasPrefix applies the HasPrefix predicate on the "username" field.
func UsernameHasPrefix(v string) predicate.User {
	return predicate.User(sql.FieldHasPrefix(FieldUsername, v))
}

// UsernameHasSuffix applies the HasSuffix predicate on the "username" field.
func UsernameHasSuffix(v string) predicate.User {
	return predicate.User(sql.FieldHasSuffix(FieldUsername, v))
}

// UsernameEqualFold applies the EqualFold predicate on the "username" field.
func UsernameEqualFold(v string) predicate.User {
	return predicate.User(sql.FieldEqualFold(FieldUsername, v))
}

// UsernameContainsFold applies the ContainsFold predicate on the "username" field.
func UsernameContainsFold(v string) predicate.User {
	return predicate.User(sql.FieldContainsFold(FieldUsername, v))
}

// HasActivities applies the HasEdge predicate on the "activities" edge.
func HasActivities() predicate.User {
	return predicate.User(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, ActivitiesTable, ActivitiesColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasActivitiesWith applies the HasEdge predicate on the "activities" edge with a given conditions (other predicates).
func HasActivitiesWith(preds ...predicate.Activity) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		step := newActivitiesStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasFriendship applies the HasEdge predicate on the "friendship" edge.
func HasFriendship() predicate.User {
	return predicate.User(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, FriendshipTable, FriendshipColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasFriendshipWith applies the HasEdge predicate on the "friendship" edge with a given conditions (other predicates).
func HasFriendshipWith(preds ...predicate.Friendship) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		step := newFriendshipStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasHexinfluence applies the HasEdge predicate on the "hexinfluence" edge.
func HasHexinfluence() predicate.User {
	return predicate.User(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, HexinfluenceTable, HexinfluenceColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasHexinfluenceWith applies the HasEdge predicate on the "hexinfluence" edge with a given conditions (other predicates).
func HasHexinfluenceWith(preds ...predicate.HexInfluence) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		step := newHexinfluenceStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.User) predicate.User {
	return predicate.User(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.User) predicate.User {
	return predicate.User(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.User) predicate.User {
	return predicate.User(sql.NotPredicates(p))
}
