package ai

import "context"

type FakeClient struct {
	Next Classification
	Err  error
}

func (f FakeClient) Classify(ctx context.Context, title, description, affectedService string) (Classification, error) {
	if f.Err != nil {
		return Classification{}, f.Err
	}
	return f.Next, nil
}
