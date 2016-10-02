package util

// AnyError checks each of the provided errors and returns the first one that
// is not nil. If all of them are nil, this function returns nil.
func AnyError(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}
