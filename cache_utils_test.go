package pie

import (
	"testing"
	"time"
)

func TestRandomInt(t *testing.T) {
	// 测试正常范围
	min, max := int64(10), int64(20)
	for i := 0; i < 100; i++ {
		result := randomInt(min, max)
		if result < min || result > max {
			t.Errorf("Expected result between %d and %d, got %d", min, max, result)
		}
	}

	// 测试min等于max的情况
	result := randomInt(5, 5)
	if result != 5 {
		t.Errorf("Expected result to be 5 when min equals max, got %d", result)
	}

	// 测试min大于max的情况
	result = randomInt(10, 5)
	if result != 10 {
		t.Errorf("Expected result to be 10 when min > max, got %d", result)
	}

	// 测试负数范围
	min, max = int64(-10), int64(-5)
	for i := 0; i < 100; i++ {
		result := randomInt(min, max)
		if result < min || result > max {
			t.Errorf("Expected result between %d and %d, got %d", min, max, result)
		}
	}

	// 测试跨零范围
	min, max = int64(-5), int64(5)
	for i := 0; i < 100; i++ {
		result := randomInt(min, max)
		if result < min || result > max {
			t.Errorf("Expected result between %d and %d, got %d", min, max, result)
		}
	}
}

func TestApplyJitter(t *testing.T) {
	// 测试jitter为0的情况
	ttl := 5 * time.Minute
	jitter := time.Duration(0)
	result := applyJitter(ttl, jitter)
	if result != ttl {
		t.Errorf("Expected TTL to be unchanged when jitter is 0, got %v", result)
	}

	// 测试正常jitter
	ttl = 5 * time.Minute
	jitter = 30 * time.Second

	// 运行多次测试，确保结果在合理范围内
	for i := 0; i < 100; i++ {
		result := applyJitter(ttl, jitter)
		expectedMin := ttl - jitter
		expectedMax := ttl + jitter

		if result < expectedMin || result > expectedMax {
			t.Errorf("Expected result between %v and %v, got %v", expectedMin, expectedMax, result)
		}
	}

	// 测试jitter大于TTL的情况（应该不会返回负数）
	ttl = 10 * time.Second
	jitter = 30 * time.Second

	for i := 0; i < 100; i++ {
		result := applyJitter(ttl, jitter)
		if result < 0 {
			t.Errorf("Expected result to be non-negative, got %v", result)
		}
		if result < ttl {
			t.Errorf("Expected result to be at least original TTL when jitter > TTL, got %v", result)
		}
	}

	// 测试边界情况：TTL为0
	ttl = 0
	jitter = 10 * time.Second
	result = applyJitter(ttl, jitter)
	if result < 0 {
		t.Errorf("Expected result to be non-negative when TTL is 0, got %v", result)
	}

	// 测试边界情况：TTL为负数（虽然不应该发生）
	ttl = -5 * time.Second
	jitter = 10 * time.Second
	result = applyJitter(ttl, jitter)
	if result < 0 {
		t.Errorf("Expected result to be non-negative when TTL is negative, got %v", result)
	}
}

func TestApplyJitterConsistency(t *testing.T) {
	// 测试相同输入产生不同结果（随机性）
	ttl := 5 * time.Minute
	jitter := 30 * time.Second

	results := make(map[time.Duration]int)
	for i := 0; i < 1000; i++ {
		result := applyJitter(ttl, jitter)
		results[result]++
	}

	// 验证有多个不同的结果（随机性）
	if len(results) < 10 {
		t.Errorf("Expected multiple different results due to randomness, got %d unique results", len(results))
	}

	// 验证结果分布合理（大部分结果应该在中间范围）
	expectedMin := ttl - jitter
	expectedMax := ttl + jitter

	// 检查结果是否在合理范围内分布
	for result, count := range results {
		if result < expectedMin || result > expectedMax {
			t.Errorf("Result %v is outside expected range [%v, %v]", result, expectedMin, expectedMax)
		}
		if count > 200 { // 单个值不应该出现太频繁
			t.Errorf("Result %v appears too frequently (%d times), indicating poor randomness", result, count)
		}
	}
}

func TestApplyJitterEdgeCases(t *testing.T) {
	// 测试极小值
	ttl := time.Nanosecond
	jitter := time.Nanosecond
	result := applyJitter(ttl, jitter)
	if result < 0 {
		t.Errorf("Expected non-negative result for small values, got %v", result)
	}

	// 测试极大值
	ttl = 24 * time.Hour
	jitter = time.Hour
	result = applyJitter(ttl, jitter)
	if result < 0 {
		t.Errorf("Expected non-negative result for large values, got %v", result)
	}
	if result < ttl-jitter || result > ttl+jitter {
		t.Errorf("Expected result in range [%v, %v], got %v", ttl-jitter, ttl+jitter, result)
	}

	// 测试jitter为负数（虽然不应该发生）
	ttl = 5 * time.Minute
	jitter = -30 * time.Second
	result = applyJitter(ttl, jitter)
	if result != ttl {
		t.Errorf("Expected original TTL when jitter is negative, got %v", result)
	}
}

func TestRandomIntDistribution(t *testing.T) {
	// 测试随机数分布
	min, max := int64(1), int64(10)
	counts := make(map[int64]int)

	// 生成大量随机数
	for i := 0; i < 10000; i++ {
		result := randomInt(min, max)
		counts[result]++
	}

	// 验证每个值都出现过
	for i := min; i <= max; i++ {
		if counts[i] == 0 {
			t.Errorf("Value %d never appeared in 10000 random numbers", i)
		}
	}

	// 验证分布相对均匀（每个值的出现次数不应该相差太大）
	expectedCount := int64(10000 / (max - min + 1))
	tolerance := expectedCount / 2 // 允许50%的偏差

	for i := min; i <= max; i++ {
		count := int64(counts[i])
		if count < expectedCount-tolerance || count > expectedCount+tolerance {
			t.Errorf("Value %d appears %d times, expected around %d (tolerance: %d)", i, count, expectedCount, tolerance)
		}
	}
}

func TestRandomIntSingleValue(t *testing.T) {
	// 测试单值情况
	result := randomInt(42, 42)
	if result != 42 {
		t.Errorf("Expected 42 for single value range, got %d", result)
	}

	// 测试零值
	result = randomInt(0, 0)
	if result != 0 {
		t.Errorf("Expected 0 for zero range, got %d", result)
	}
}

func TestApplyJitterZeroValues(t *testing.T) {
	// 测试零TTL
	result := applyJitter(0, 10*time.Second)
	if result < 0 {
		t.Errorf("Expected non-negative result for zero TTL, got %v", result)
	}

	// 测试零jitter
	result = applyJitter(5*time.Minute, 0)
	if result != 5*time.Minute {
		t.Errorf("Expected original TTL for zero jitter, got %v", result)
	}

	// 测试两者都为零
	result = applyJitter(0, 0)
	if result != 0 {
		t.Errorf("Expected 0 for zero TTL and zero jitter, got %v", result)
	}
}
