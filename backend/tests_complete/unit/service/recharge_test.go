package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRechargeService_CalculatePoints_Tier1(t *testing.T) {
	svc := &RechargeService{}

	points, basePoints, rebatePoints := svc.CalculatePoints(5000, 0)
	assert.Equal(t, 5000, points)
	assert.Equal(t, 5000, basePoints)
	assert.Equal(t, 0, rebatePoints)
}

func TestRechargeService_CalculatePoints_Tier2(t *testing.T) {
	svc := &RechargeService{}

	points, basePoints, rebatePoints := svc.CalculatePoints(30000, 0)
	assert.Equal(t, 30300, points)
	assert.Equal(t, 30000, basePoints)
	assert.Equal(t, 300, rebatePoints)
}

func TestRechargeService_CalculatePoints_Tier3(t *testing.T) {
	svc := &RechargeService{}

	points, basePoints, rebatePoints := svc.CalculatePoints(80000, 0)
	assert.Equal(t, 80800, points)
	assert.Equal(t, 80000, basePoints)
	assert.Equal(t, 800, rebatePoints)
}

func TestRechargeService_CalculatePoints_Tier4(t *testing.T) {
	svc := &RechargeService{}

	points, basePoints, rebatePoints := svc.CalculatePoints(150000, 0)
	assert.Equal(t, 153000, points)
	assert.Equal(t, 150000, basePoints)
	assert.Equal(t, 3000, rebatePoints)
}

func TestRechargeService_CalculatePoints_WithLastMonthConsumption(t *testing.T) {
	svc := &RechargeService{}

	points, basePoints, rebatePoints := svc.CalculatePoints(30000, 100000)
	assert.Equal(t, 30600, points)
	assert.Equal(t, 30000, basePoints)
	assert.Equal(t, 600, rebatePoints)
}

func TestRechargeService_CalculatePoints_EdgeCases(t *testing.T) {
	svc := &RechargeService{}

	t.Run("zero amount", func(t *testing.T) {
		points, basePoints, rebatePoints := svc.CalculatePoints(0, 0)
		assert.Equal(t, 0, points)
		assert.Equal(t, 0, basePoints)
		assert.Equal(t, 0, rebatePoints)
	})

	t.Run("exact boundary 10000", func(t *testing.T) {
		points, basePoints, rebatePoints := svc.CalculatePoints(10000, 0)
		assert.Equal(t, 10100, points)
		assert.Equal(t, 10000, basePoints)
		assert.Equal(t, 100, rebatePoints)
	})
}

func TestRechargeService_GenerateCardNo(t *testing.T) {
	svc := &RechargeService{}

	cardNo := svc.generateCardNo()
	assert.NotEmpty(t, cardNo)
	assert.Contains(t, cardNo, "TJ")
	assert.GreaterOrEqual(t, len(cardNo), 9)
}
