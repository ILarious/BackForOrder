package errors

import stderrors "errors"

var ErrEmptyUsername = stderrors.New("order: username is required")
