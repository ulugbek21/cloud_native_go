package retry

import "context"

type Effector func(context.Context) (string, error)