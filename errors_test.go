package csv

import (
	"testing"
)

func TestError_Error(t *testing.T) {
	t.Parallel()

	t.Run("should return the localized error message", func(t *testing.T) {
		t.Parallel()

		err := NewError(helperLocalizer(t), ErrStructSlicePointerID, "subMessage")

		got := err.Error()
		want := "value is not a pointer to a struct slice: subMessage"

		if got != want {
			t.Errorf("Error() = %v, want %v", got, want)
		}
	})

	t.Run("should return the localized error message without subMessage", func(t *testing.T) {
		t.Parallel()

		err := NewError(helperLocalizer(t), ErrStructSlicePointerID, "")

		got := err.Error()
		want := "value is not a pointer to a struct slice"

		if got != want {
			t.Errorf("Error() = %v, want %v", got, want)
		}
	})
}

func TestError_Is(t *testing.T) {
	t.Parallel()

	t.Run("should return true if the target error is the same as the error", func(t *testing.T) {
		t.Parallel()

		err := NewError(helperLocalizer(t), ErrStructSlicePointerID, "subMessage")
		target := NewError(helperLocalizer(t), ErrStructSlicePointerID, "subMessage")

		got := err.Is(target)
		want := true

		if got != want {
			t.Errorf("Is() = %v, want %v", got, want)
		}
	})

	t.Run("should return false if the target error is not the same as the error", func(t *testing.T) {
		t.Parallel()

		err := NewError(helperLocalizer(t), ErrStructSlicePointerID, "subMessage")
		target := NewError(helperLocalizer(t), ErrEqualID, "")

		got := err.Is(target)
		want := false

		if got != want {
			t.Errorf("Is() = %v, want %v", got, want)
		}
	})
}
