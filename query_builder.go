package pie

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// WhereIn checks if field value is in specified range
func (s *Session[T]) WhereIn(field string, values interface{}) *Session[T] {
	s.query.filter = append(s.query.filter, bson.E{Key: field, Value: bson.D{{Key: "$in", Value: values}}})
	return s
}

// WhereNotIn checks if field value is not in specified range
func (s *Session[T]) WhereNotIn(field string, values interface{}) *Session[T] {
	s.query.filter = append(s.query.filter, bson.E{Key: field, Value: bson.D{{Key: "$nin", Value: values}}})
	return s
}

// WhereBetween checks if field value is between min and max
func (s *Session[T]) WhereBetween(field string, min, max interface{}) *Session[T] {
	s.query.filter = append(s.query.filter, bson.E{
		Key: field,
		Value: bson.D{
			{Key: "$gte", Value: min},
			{Key: "$lte", Value: max},
		},
	})
	return s
}

// WhereNotBetween checks if field value is not between min and max
func (s *Session[T]) WhereNotBetween(field string, min, max interface{}) *Session[T] {
	s.query.filter = append(s.query.filter, bson.E{
		Key: "$or",
		Value: []bson.D{
			{{Key: field, Value: bson.D{{Key: "$lt", Value: min}}}},
			{{Key: field, Value: bson.D{{Key: "$gt", Value: max}}}},
		},
	})
	return s
}

// WhereNull checks if field value is null or does not exist
func (s *Session[T]) WhereNull(field string) *Session[T] {
	s.query.filter = append(s.query.filter, bson.E{
		Key: "$or",
		Value: []bson.D{
			{{Key: field, Value: bson.D{{Key: "$exists", Value: false}}}},
			{{Key: field, Value: nil}},
		},
	})
	return s
}

// WhereNotNull checks if field value is not null and exists
func (s *Session[T]) WhereNotNull(field string) *Session[T] {
	s.query.filter = append(s.query.filter, bson.E{
		Key: "$and",
		Value: []bson.D{
			{{Key: field, Value: bson.D{{Key: "$exists", Value: true}}}},
			{{Key: field, Value: bson.D{{Key: "$ne", Value: nil}}}},
		},
	})
	return s
}

// WhereExists 字段存在
func (s *Session[T]) WhereExists(field string) *Session[T] {
	s.query.filter = append(s.query.filter, bson.E{Key: field, Value: bson.D{{Key: "$exists", Value: true}}})
	return s
}

// WhereNotExists 字段不存在
func (s *Session[T]) WhereNotExists(field string) *Session[T] {
	s.query.filter = append(s.query.filter, bson.E{Key: field, Value: bson.D{{Key: "$exists", Value: false}}})
	return s
}

// WhereDate 日期比较查询
func (s *Session[T]) WhereDate(field string, operator string, value interface{}) *Session[T] {
	var op string
	switch operator {
	case ">", "gt":
		op = "$gt"
	case ">=", "gte":
		op = "$gte"
	case "<", "lt":
		op = "$lt"
	case "<=", "lte":
		op = "$lte"
	case "=", "==", "eq":
		op = "$eq"
	default:
		op = "$eq"
	}

	// 将字符串日期转换为 time.Time
	var dateValue time.Time
	switch v := value.(type) {
	case string:
		// 尝试解析多种日期格式
		formats := []string{
			"2006-01-02",
			"2006-01-02 15:04:05",
			time.RFC3339,
		}
		for _, format := range formats {
			if t, err := time.Parse(format, v); err == nil {
				dateValue = t
				break
			}
		}
	case time.Time:
		dateValue = v
	default:
		// 如果无法转换，直接使用原值
		s.query.filter = append(s.query.filter, bson.E{Key: field, Value: bson.D{{Key: op, Value: value}}})
		return s
	}

	s.query.filter = append(s.query.filter, bson.E{Key: field, Value: bson.D{{Key: op, Value: dateValue}}})
	return s
}

// WhereDateBetween 日期范围查询
func (s *Session[T]) WhereDateBetween(field string, start, end interface{}) *Session[T] {
	return s.WhereBetween(field, start, end)
}

// WhereMonth 按月份查询
func (s *Session[T]) WhereMonth(field string, month int) *Session[T] {
	s.query.filter = append(s.query.filter, bson.E{
		Key: "$expr",
		Value: bson.D{{
			Key: "$eq",
			Value: []interface{}{
				bson.D{{Key: "$month", Value: "$" + field}},
				month,
			},
		}},
	})
	return s
}

// WhereYear 按年份查询
func (s *Session[T]) WhereYear(field string, year int) *Session[T] {
	s.query.filter = append(s.query.filter, bson.E{
		Key: "$expr",
		Value: bson.D{{
			Key: "$eq",
			Value: []interface{}{
				bson.D{{Key: "$year", Value: "$" + field}},
				year,
			},
		}},
	})
	return s
}

// WhereRecentDays 最近N天
func (s *Session[T]) WhereRecentDays(field string, days int) *Session[T] {
	startTime := time.Now().AddDate(0, 0, -days)
	s.query.filter = append(s.query.filter, bson.E{Key: field, Value: bson.D{{Key: "$gte", Value: startTime}}})
	return s
}

// WhereLike 模糊查询（包含）
func (s *Session[T]) WhereLike(field string, pattern string) *Session[T] {
	// 去除用户输入的 % 符号，转换为正则表达式
	regexPattern := pattern
	if len(pattern) > 0 {
		if pattern[0] == '%' && pattern[len(pattern)-1] == '%' {
			// %keyword% -> keyword（包含）
			regexPattern = pattern[1 : len(pattern)-1]
		} else if pattern[0] == '%' {
			// %keyword -> keyword$（结尾）
			regexPattern = pattern[1:] + "$"
		} else if pattern[len(pattern)-1] == '%' {
			// keyword% -> ^keyword（开头）
			regexPattern = "^" + pattern[:len(pattern)-1]
		}
	}

	s.query.filter = append(s.query.filter, bson.E{
		Key: field,
		Value: bson.D{
			{Key: "$regex", Value: regexPattern},
			{Key: "$options", Value: "i"}, // 忽略大小写
		},
	})
	return s
}

// WhereStartsWith 前缀匹配
func (s *Session[T]) WhereStartsWith(field string, prefix string) *Session[T] {
	s.query.filter = append(s.query.filter, bson.E{
		Key: field,
		Value: bson.D{
			{Key: "$regex", Value: "^" + prefix},
			{Key: "$options", Value: "i"},
		},
	})
	return s
}

// WhereEndsWith 后缀匹配
func (s *Session[T]) WhereEndsWith(field string, suffix string) *Session[T] {
	s.query.filter = append(s.query.filter, bson.E{
		Key: field,
		Value: bson.D{
			{Key: "$regex", Value: suffix + "$"},
			{Key: "$options", Value: "i"},
		},
	})
	return s
}

// WhereArrayContains 数组包含指定元素
func (s *Session[T]) WhereArrayContains(field string, value interface{}) *Session[T] {
	s.query.filter = append(s.query.filter, bson.E{Key: field, Value: value})
	return s
}

// WhereArraySize 数组大小
func (s *Session[T]) WhereArraySize(field string, size int) *Session[T] {
	s.query.filter = append(s.query.filter, bson.E{Key: field, Value: bson.D{{Key: "$size", Value: size}}})
	return s
}

// WhereArrayAll 数组包含所有指定元素
func (s *Session[T]) WhereArrayAll(field string, values interface{}) *Session[T] {
	s.query.filter = append(s.query.filter, bson.E{Key: field, Value: bson.D{{Key: "$all", Value: values}}})
	return s
}

// OrWhere 添加OR条件（支持回调）
func (s *Session[T]) OrWhere(callback func(*Query) *Query) *Session[T] {
	// 创建一个新的查询用于OR条件
	orQuery := NewQuery()
	callback(orQuery)

	// 获取回调中构建的条件
	orFilter := orQuery.GetFilter()
	if len(orFilter) > 0 {
		s.query.filter = append(s.query.filter, bson.E{
			Key:   "$or",
			Value: []bson.D{orFilter},
		})
	}

	return s
}

// AndWhere 添加AND条件（支持回调）
func (s *Session[T]) AndWhere(callback func(*Query) *Query) *Session[T] {
	// 创建一个新的查询用于AND条件
	andQuery := NewQuery()
	callback(andQuery)

	// 获取回调中构建的条件
	andFilter := andQuery.GetFilter()
	if len(andFilter) > 0 {
		s.query.filter = append(s.query.filter, bson.E{
			Key:   "$and",
			Value: []bson.D{andFilter},
		})
	}

	return s
}

// Query 方法扩展（用于支持链式的Query对象）

// WhereIn 字段值在指定范围内
func (q *Query) WhereIn(field string, values interface{}) *Query {
	q.filter = append(q.filter, bson.E{Key: field, Value: bson.D{{Key: "$in", Value: values}}})
	return q
}

// WhereNotIn 字段值不在指定范围内
func (q *Query) WhereNotIn(field string, values interface{}) *Query {
	q.filter = append(q.filter, bson.E{Key: field, Value: bson.D{{Key: "$nin", Value: values}}})
	return q
}

// WhereBetween 字段值在指定范围之间
func (q *Query) WhereBetween(field string, min, max interface{}) *Query {
	q.filter = append(q.filter, bson.E{
		Key: field,
		Value: bson.D{
			{Key: "$gte", Value: min},
			{Key: "$lte", Value: max},
		},
	})
	return q
}

// WhereNull 字段值为null或不存在
func (q *Query) WhereNull(field string) *Query {
	q.filter = append(q.filter, bson.E{
		Key: "$or",
		Value: []bson.D{
			{{Key: field, Value: bson.D{{Key: "$exists", Value: false}}}},
			{{Key: field, Value: nil}},
		},
	})
	return q
}

// WhereNotNull 字段值不为null且存在
func (q *Query) WhereNotNull(field string) *Query {
	q.filter = append(q.filter, bson.E{
		Key: "$and",
		Value: []bson.D{
			{{Key: field, Value: bson.D{{Key: "$exists", Value: true}}}},
			{{Key: field, Value: bson.D{{Key: "$ne", Value: nil}}}},
		},
	})
	return q
}

// WhereExists 字段存在
func (q *Query) WhereExists(field string) *Query {
	q.filter = append(q.filter, bson.E{Key: field, Value: bson.D{{Key: "$exists", Value: true}}})
	return q
}

// WhereNotExists 字段不存在
func (q *Query) WhereNotExists(field string) *Query {
	q.filter = append(q.filter, bson.E{Key: field, Value: bson.D{{Key: "$exists", Value: false}}})
	return q
}

// Bool 辅助函数，用于创建布尔指针
func Bool(b bool) *bool {
	return &b
}

// Float64 辅助函数，用于创建float64指针
func Float64(f float64) *float64 {
	return &f
}

// Int 辅助函数，用于创建int指针
func Int(i int) *int {
	return &i
}

// String 辅助函数，用于创建string指针
func String(s string) *string {
	return &s
}
