package repository_test

import "testing"

//TODO to properly test this we probably need to set up local postgres container????

func TestUserRepository_FindByUUID(t *testing.T) {
	t.Parallel()

	t.Run("correctly find by uuid", func(t *testing.T) {
	})
}

func TestUserRepository_FindByUsername(t *testing.T) {
	t.Parallel()

	t.Run("correctly find by username", func(t *testing.T) {

	})
}

func TestUserRepository_CreateUser(t *testing.T) {
	t.Parallel()

	t.Run("correctly create user", func(t *testing.T) {

	})
}

func TestUserRepository_UpdateUsername(t *testing.T) {
	t.Parallel()

	t.Run("correctly update user", func(t *testing.T) {

	})
}
