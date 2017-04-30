package errors

import (
	"errors"
	"os"
)

// RichErr satisfies the error interface. Additionally, it remembers
// individually all of the errors that have been added to the RichErr.
type RichErr struct {
	ErrSet []error
}

// Error returns the composed error string of the RichErr.
func (r RichErr) Error() string {
	s := "["
	for i, err := range r.ErrSet {
		if i != 0 {
			s = s + "; "
		}
		s = s + err.Error()
	}
	return s + "]"
}

// Compose will compose all errors together into a single rich error,
// preserving any rich context found along the way.
func Compose(errs ...error) error {
	var r RichErr
	for _, err := range errs {
		// Handle nil errors.
		if err == nil {
			continue
		}
		r.ErrSet = append(r.ErrSet, err)
	}
	if len(r.ErrSet) == 0 {
		return nil
	}
	return r
}

// Contains will check whether the base error contains the cmp error. If the
// base err is a RichErr, then it will check whether there is a match on any of
// the underlying errors.
func Contains(base, cmp error) bool {
	// Check for the easy edge cases.
	if cmp == nil || base == nil {
		return false
	}
	if base == cmp {
		return true
	}

	switch v := base.(type) {
	case RichErr:
		for _, err := range v.ErrSet {
			if Contains(err, cmp) {
				return true
			}
		}
		return false
	default:
		return false
	}
}

// Extend will extend the second error with the first error,
func Extend(err, extension error) error {
	// Check for nil edge cases. If both are nil, nil will be returned.
	if err == nil {
		return extension
	}
	if extension == nil {
		return err
	}

	var r RichErr
	// Check the original error for richness.
	switch v := err.(type) {
	case RichErr:
		r = v
	default:
		r.ErrSet = []error{v}
	}

	// Check the extension error for richness.
	switch v := extension.(type) {
	case RichErr:
		r.ErrSet = append(v.ErrSet, r.ErrSet...)
	default:
		r.ErrSet = append([]error{v}, r.ErrSet...)
	}

	// Return nil if the result has no underlying errors.
	if len(r.ErrSet) == 0 {
		return nil
	}
	return r
}

// IsOSNotExist returns true if any of the errors in the underlying composition
// return true for os.IsNotExist.
func IsOSNotExist(err error) bool {
	if err == nil {
		return false
	}

	switch v := err.(type) {
	case RichErr:
		for _, err := range v.ErrSet {
			if IsOSNotExist(err) {
				return true
			}
		}
		return false
	default:
		return os.IsNotExist(err)
	}
}

// New is a passthrough to the stdlib errors package, allowing
// NebulousLabs/errors to be a drop in replacement for the standard library
// errors.
func New(s string) error {
	return errors.New(s)
}
