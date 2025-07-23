package git

// MockExecutor implements Executor for testing purposes.
type MockExecutor struct {
	CurrentBranchFunc  func() (string, error)
	LogFunc            func(branch, format string) (string, error)
	LsTreeFunc         func(branch, path string) ([]string, error)
	RemoteBranchesFunc func() ([]string, error)
}

func (m *MockExecutor) CurrentBranch() (string, error) {
	if m.CurrentBranchFunc != nil {
		return m.CurrentBranchFunc()
	}
	return "", nil
}
func (m *MockExecutor) Checkout(branch string) error     { return nil }
func (m *MockExecutor) Branches() ([]string, error)      { return nil, nil }
func (m *MockExecutor) Add(paths ...string) error        { return nil }
func (m *MockExecutor) Commit(message string) error      { return nil }
func (m *MockExecutor) Push(branch string) error         { return nil }
func (m *MockExecutor) StatusPorcelain() (string, error) { return "", nil }
func (m *MockExecutor) LastCommit() (string, error)      { return "", nil }
func (m *MockExecutor) IsRepository() bool               { return true }
func (m *MockExecutor) Log(branch, format string) (string, error) {
	if m.LogFunc != nil {
		return m.LogFunc(branch, format)
	}
	return "", nil
}
func (m *MockExecutor) LsTree(branch, path string) ([]string, error) {
	if m.LsTreeFunc != nil {
		return m.LsTreeFunc(branch, path)
	}
	return nil, nil
}
func (m *MockExecutor) RemoteBranches() ([]string, error) {
	if m.RemoteBranchesFunc != nil {
		return m.RemoteBranchesFunc()
	}
	return nil, nil
}
