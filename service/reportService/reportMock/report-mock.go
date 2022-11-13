package reportMockService

import "github.com/stretchr/testify/mock"

type ReportMock struct {
	mock.Mock
}

func (r *ReportMock) GetCashflow(userId uint, period string) (map[string]interface{}, error) {
	args := r.Called(userId, period)

	return args.Get(0).(map[string]interface{}), args.Error(1)
}
func (r *ReportMock) GetReportbyCategory(userId uint, period string, numberPeriod int) (map[string]interface{}, error) {
	args := r.Called(userId, period, numberPeriod)

	return args.Get(0).(map[string]interface{}), args.Error(1)
}
func (r *ReportMock) GetAnalyticPeriod(userId uint, period string) (map[string]interface{}, error) {
	args := r.Called(userId, period)

	return args.Get(0).(map[string]interface{}), args.Error(1)
}