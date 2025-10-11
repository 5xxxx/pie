package pie

// M 类型别名,简化 map 定义
type M = map[string]any

// ========== 日期表达式 ==========

// DateToString 将日期转换为字符串
func DateToString(field, format, timezone string) M {
	return M{
		"$dateToString": M{
			"date":     field,
			"format":   format,
			"timezone": timezone,
		},
	}
}

// DateFromString 将字符串转换为日期
func DateFromString(dateString any) M {
	return M{
		"$dateFromString": M{
			"dateString": dateString,
		},
	}
}

// DateToParts 将日期分解为部分
func DateToParts(field string, timezone ...string) M {
	expr := M{
		"$dateToParts": M{
			"date": field,
		},
	}
	if len(timezone) > 0 {
		expr["$dateToParts"].(M)["timezone"] = timezone[0]
	}
	return expr
}

// DateFromParts 从部分构建日期
func DateFromParts(parts M) M {
	return M{
		"$dateFromParts": parts,
	}
}

// DateAdd 添加时间到日期
func DateAdd(field string, amount any, unit string) M {
	return M{
		"$dateAdd": M{
			"startDate": field,
			"unit":      unit,
			"amount":    amount,
		},
	}
}

// DateSubtract 从日期减去时间
func DateSubtract(field string, amount any, unit string) M {
	return M{
		"$dateSubtract": M{
			"startDate": field,
			"unit":      unit,
			"amount":    amount,
		},
	}
}

// DateDiff 计算两个日期的差值
func DateDiff(startDate, endDate, unit string) M {
	return M{
		"$dateDiff": M{
			"startDate": startDate,
			"endDate":   endDate,
			"unit":      unit,
		},
	}
}

// DateTrunc 截断日期
func DateTrunc(field string, unit string, binSize ...any) M {
	expr := M{
		"$dateTrunc": M{
			"date": field,
			"unit": unit,
		},
	}
	if len(binSize) > 0 {
		expr["$dateTrunc"].(M)["binSize"] = binSize[0]
	}
	return expr
}

// Year 提取年份
func Year(field string) M {
	return M{"$year": field}
}

// Month 提取月份
func Month(field string) M {
	return M{"$month": field}
}

// Week 提取周
func Week(field string) M {
	return M{"$week": field}
}

// DayOfMonth 提取月中的天数
func DayOfMonth(field string) M {
	return M{"$dayOfMonth": field}
}

// DayOfWeek 提取周中的天数
func DayOfWeek(field string) M {
	return M{"$dayOfWeek": field}
}

// DayOfYear 提取年中的天数
func DayOfYear(field string) M {
	return M{"$dayOfYear": field}
}

// Hour 提取小时
func Hour(field string) M {
	return M{"$hour": field}
}

// Minute 提取分钟
func Minute(field string) M {
	return M{"$minute": field}
}

// Second 提取秒
func Second(field string) M {
	return M{"$second": field}
}

// Millisecond 提取毫秒
func Millisecond(field string) M {
	return M{"$millisecond": field}
}

// ISOWeek 提取 ISO 周
func ISOWeek(field string) M {
	return M{"$isoWeek": field}
}

// ISOWeekYear 提取 ISO 周年
func ISOWeekYear(field string) M {
	return M{"$isoWeekYear": field}
}

// IsoDayOfWeek 提取 ISO 周中的天数
func IsoDayOfWeek(field string) M {
	return M{"$isoDayOfWeek": field}
}

// ========== 算术表达式 ==========

// Add 加法
func Add(values ...any) M {
	return M{"$add": values}
}

// Subtract 减法
func Subtract(left, right any) M {
	return M{"$subtract": []any{left, right}}
}

// Multiply 乘法
func Multiply(values ...any) M {
	return M{"$multiply": values}
}

// Divide 除法
func Divide(left, right any) M {
	return M{"$divide": []any{left, right}}
}

// ModExpr 取模表达式
func ModExpr(left, right any) M {
	return M{"$mod": []any{left, right}}
}

// Abs 绝对值
func Abs(field any) M {
	return M{"$abs": field}
}

// Ceil 向上取整
func Ceil(field any) M {
	return M{"$ceil": field}
}

// Floor 向下取整
func Floor(field any) M {
	return M{"$floor": field}
}

// Round 四舍五入
func Round(field any, place ...any) M {
	if len(place) > 0 {
		return M{"$round": []any{field, place[0]}}
	}
	return M{"$round": field}
}

// Trunc 截断
func Trunc(field any, place ...any) M {
	if len(place) > 0 {
		return M{"$trunc": []any{field, place[0]}}
	}
	return M{"$trunc": field}
}

// Sqrt 平方根
func Sqrt(field any) M {
	return M{"$sqrt": field}
}

// Pow 幂运算
func Pow(base, exponent any) M {
	return M{"$pow": []any{base, exponent}}
}

// Exp 指数
func Exp(field any) M {
	return M{"$exp": field}
}

// Ln 自然对数
func Ln(field any) M {
	return M{"$ln": field}
}

// Log 对数
func Log(number, base any) M {
	return M{"$log": []any{number, base}}
}

// Log10 以10为底的对数
func Log10(field any) M {
	return M{"$log10": field}
}

// Sin 正弦
func Sin(field any) M {
	return M{"$sin": field}
}

// Cos 余弦
func Cos(field any) M {
	return M{"$cos": field}
}

// Tan 正切
func Tan(field any) M {
	return M{"$tan": field}
}

// Asin 反正弦
func Asin(field any) M {
	return M{"$asin": field}
}

// Acos 反余弦
func Acos(field any) M {
	return M{"$acos": field}
}

// Atan 反正切
func Atan(field any) M {
	return M{"$atan": field}
}

// Atan2 四象限反正切
func Atan2(y, x any) M {
	return M{"$atan2": []any{y, x}}
}

// DegreesToRadians 度转弧度
func DegreesToRadians(field any) M {
	return M{"$degreesToRadians": field}
}

// RadiansToDegrees 弧度转度
func RadiansToDegrees(field any) M {
	return M{"$radiansToDegrees": field}
}

// ========== 字符串表达式 ==========

// Concat 连接字符串
func Concat(values ...any) M {
	return M{"$concat": values}
}

// SubStr 子字符串
func SubStr(field string, start, length int) M {
	return M{"$substr": []any{field, start, length}}
}

// SubStrBytes 按字节提取子字符串
func SubStrBytes(field string, start, length int) M {
	return M{"$substrBytes": []any{field, start, length}}
}

// SubStrCP 按代码点提取子字符串
func SubStrCP(field string, start, length int) M {
	return M{"$substrCP": []any{field, start, length}}
}

// ToUpper 转大写
func ToUpper(field string) M {
	return M{"$toUpper": field}
}

// ToLower 转小写
func ToLower(field string) M {
	return M{"$toLower": field}
}

// StrCaseCmp 字符串比较
func StrCaseCmp(left, right string) M {
	return M{"$strcasecmp": []string{left, right}}
}

// StrLenBytes 字符串字节长度
func StrLenBytes(field string) M {
	return M{"$strLenBytes": field}
}

// StrLenCP 字符串代码点长度
func StrLenCP(field string) M {
	return M{"$strLenCP": field}
}

// Split 分割字符串
func Split(field, delimiter string) M {
	return M{"$split": []string{field, delimiter}}
}

// Trim 去除首尾空白
func Trim(field string) M {
	return M{"$trim": M{"input": field}}
}

// LTrim 去除左侧空白
func LTrim(field string) M {
	return M{"$ltrim": M{"input": field}}
}

// RTrim 去除右侧空白
func RTrim(field string) M {
	return M{"$rtrim": M{"input": field}}
}

// ReplaceOne 替换第一个匹配
func ReplaceOne(field, find, replacement string) M {
	return M{"$replaceOne": M{
		"input":       field,
		"find":        find,
		"replacement": replacement,
	}}
}

// ReplaceAll 替换所有匹配
func ReplaceAll(field, find, replacement string) M {
	return M{"$replaceAll": M{
		"input":       field,
		"find":        find,
		"replacement": replacement,
	}}
}

// IndexOfBytes 按字节查找索引
func IndexOfBytes(field, substring string) M {
	return M{"$indexOfBytes": []string{field, substring}}
}

// IndexOfCP 按代码点查找索引
func IndexOfCP(field, substring string) M {
	return M{"$indexOfCP": []string{field, substring}}
}

// ========== 数组表达式 ==========

// ArrayElemAt 获取数组元素
func ArrayElemAt(field string, index int) M {
	return M{"$arrayElemAt": []any{field, index}}
}

// ArrayToObject 数组转对象
func ArrayToObject(field string) M {
	return M{"$arrayToObject": field}
}

// ConcatArrays 连接数组
func ConcatArrays(arrays ...any) M {
	return M{"$concatArrays": arrays}
}

// FilterArray 过滤数组
func FilterArray(field string, cond any) M {
	return M{"$filter": M{
		"input": field,
		"cond":  cond,
	}}
}

// First 第一个元素
func First(field string) M {
	return M{"$first": field}
}

// Last 最后一个元素
func Last(field string) M {
	return M{"$last": field}
}

// InArray 检查是否在数组中
func InArray(field string, array any) M {
	return M{"$in": []any{field, array}}
}

// IndexOfArray 查找数组索引
func IndexOfArray(field string, value any) M {
	return M{"$indexOfArray": []any{field, value}}
}

// IsArray 检查是否为数组
func IsArray(field string) M {
	return M{"$isArray": field}
}

// MapArray 映射数组
func MapArray(field string, as string, in any) M {
	return M{"$map": M{
		"input": field,
		"as":    as,
		"in":    in,
	}}
}

// ObjectToArray 对象转数组
func ObjectToArray(field string) M {
	return M{"$objectToArray": field}
}

// Range 生成范围数组
func Range(start, end, step int) M {
	return M{"$range": []int{start, end, step}}
}

// Reduce 归约数组
func Reduce(field string, initialValue any, in any) M {
	return M{"$reduce": M{
		"input":        field,
		"initialValue": initialValue,
		"in":           in,
	}}
}

// ReverseArray 反转数组
func ReverseArray(field string) M {
	return M{"$reverseArray": field}
}

// SizeArray 数组大小
func SizeArray(field string) M {
	return M{"$size": field}
}

// Slice 切片数组
func Slice(field string, n int, position ...int) M {
	if len(position) > 0 {
		return M{"$slice": []any{field, position[0], n}}
	}
	return M{"$slice": []any{field, n}}
}

// Zip 压缩数组
func Zip(useLongestLength bool, inputs ...any) M {
	return M{"$zip": M{
		"inputs":           inputs,
		"useLongestLength": useLongestLength,
	}}
}

// MaxN 最大N个元素
func MaxN(field string, n int) M {
	return M{"$maxN": M{
		"input": field,
		"n":     n,
	}}
}

// MinN 最小N个元素
func MinN(field string, n int) M {
	return M{"$minN": M{
		"input": field,
		"n":     n,
	}}
}

// SortArray 排序数组
func SortArray(field string, sortBy ...any) M {
	if len(sortBy) > 0 {
		return M{"$sortArray": M{
			"input":  field,
			"sortBy": sortBy[0],
		}}
	}
	return M{"$sortArray": field}
}

// ========== 条件表达式 ==========

// Cond 条件表达式
func Cond(ifExpr, thenExpr, elseExpr any) M {
	return M{"$cond": M{
		"if":   ifExpr,
		"then": thenExpr,
		"else": elseExpr,
	}}
}

// IfNull 空值处理
func IfNull(field, replacement any) M {
	return M{"$ifNull": []any{field, replacement}}
}

// Switch 多条件表达式
func Switch(branches []M, defaultCase ...any) M {
	expr := M{"$switch": M{"branches": branches}}
	if len(defaultCase) > 0 {
		expr["$switch"].(M)["default"] = defaultCase[0]
	}
	return expr
}

// ========== 比较表达式 ==========

// EqExpr 等于表达式
func EqExpr(left, right any) M {
	return M{"$eq": []any{left, right}}
}

// NeExpr 不等于表达式
func NeExpr(left, right any) M {
	return M{"$ne": []any{left, right}}
}

// GtExpr 大于表达式
func GtExpr(left, right any) M {
	return M{"$gt": []any{left, right}}
}

// GteExpr 大于等于表达式
func GteExpr(left, right any) M {
	return M{"$gte": []any{left, right}}
}

// LtExpr 小于表达式
func LtExpr(left, right any) M {
	return M{"$lt": []any{left, right}}
}

// LteExpr 小于等于表达式
func LteExpr(left, right any) M {
	return M{"$lte": []any{left, right}}
}

// Cmp 比较
func Cmp(left, right any) M {
	return M{"$cmp": []any{left, right}}
}

// ========== 逻辑表达式 ==========

// AndExpr 逻辑与表达式
func AndExpr(expressions ...any) M {
	return M{"$and": expressions}
}

// OrExpr 逻辑或表达式
func OrExpr(expressions ...any) M {
	return M{"$or": expressions}
}

// NotExpr 逻辑非表达式
func NotExpr(expression any) M {
	return M{"$not": expression}
}

// ========== 类型表达式 ==========

// TypeExpr 获取类型表达式
func TypeExpr(field string) M {
	return M{"$type": field}
}

// Convert 类型转换
func Convert(field string, to string, onError, onNull any) M {
	expr := M{"$convert": M{
		"input": field,
		"to":    to,
	}}
	if onError != nil {
		expr["$convert"].(M)["onError"] = onError
	}
	if onNull != nil {
		expr["$convert"].(M)["onNull"] = onNull
	}
	return expr
}

// ToBool 转布尔
func ToBool(field string) M {
	return M{"$toBool": field}
}

// ToDate 转日期
func ToDate(field string) M {
	return M{"$toDate": field}
}

// ToDecimal 转小数
func ToDecimal(field string) M {
	return M{"$toDecimal": field}
}

// ToDouble 转双精度
func ToDouble(field string) M {
	return M{"$toDouble": field}
}

// ToInt 转整数
func ToInt(field string) M {
	return M{"$toInt": field}
}

// ToLong 转长整数
func ToLong(field string) M {
	return M{"$toLong": field}
}

// ToObjectId 转对象ID
func ToObjectId(field string) M {
	return M{"$toObjectId": field}
}

// ToString 转字符串
func ToString(field string) M {
	return M{"$toString": field}
}

// IsNumber 检查是否为数字
func IsNumber(field string) M {
	return M{"$isNumber": field}
}

// ========== 对象表达式 ==========

// MergeObjects 合并对象
func MergeObjects(objects ...any) M {
	return M{"$mergeObjects": objects}
}

// SetField 设置字段
func SetField(field string, path string, value any) M {
	return M{"$setField": M{
		"field": field,
		"input": M{"$literal": M{}},
		"path":  path,
		"value": value,
	}}
}

// GetField 获取字段
func GetField(field string, path string) M {
	return M{"$getField": M{
		"field": field,
		"input": M{"$literal": M{}},
		"path":  path,
	}}
}

// UnsetField 移除字段
func UnsetField(field string, path string) M {
	return M{"$unsetField": M{
		"field": field,
		"input": M{"$literal": M{}},
		"path":  path,
	}}
}

// ========== 聚合累加器表达式 ==========

// Sum 求和
func Sum(field string) M {
	return M{"$sum": field}
}

// Avg 平均值
func Avg(field string) M {
	return M{"$avg": field}
}

// Max 最大值
func Max(field string) M {
	return M{"$max": field}
}

// Min 最小值
func Min(field string) M {
	return M{"$min": field}
}

// StdDevPop 总体标准差
func StdDevPop(field string) M {
	return M{"$stdDevPop": field}
}

// StdDevSamp 样本标准差
func StdDevSamp(field string) M {
	return M{"$stdDevSamp": field}
}

// First 第一个值
func FirstAccumulator(field string) M {
	return M{"$first": field}
}

// Last 最后一个值
func LastAccumulator(field string) M {
	return M{"$last": field}
}

// Push 推入数组
func Push(field string) M {
	return M{"$push": field}
}

// AddToSet 添加到集合
func AddToSet(field string) M {
	return M{"$addToSet": field}
}

// MergeObjectsAccumulator 合并对象累加器
func MergeObjectsAccumulator(field string) M {
	return M{"$mergeObjects": field}
}

// Accumulator 自定义累加器
func Accumulator(init, accumulate, merge any, initArgs, accumulateArgs, finalize, finalizeArgs any) M {
	expr := M{"$accumulator": M{
		"init":       init,
		"accumulate": accumulate,
		"merge":      merge,
	}}
	if initArgs != nil {
		expr["$accumulator"].(M)["initArgs"] = initArgs
	}
	if accumulateArgs != nil {
		expr["$accumulator"].(M)["accumulateArgs"] = accumulateArgs
	}
	if finalize != nil {
		expr["$accumulator"].(M)["finalize"] = finalize
	}
	if finalizeArgs != nil {
		expr["$accumulator"].(M)["finalizeArgs"] = finalizeArgs
	}
	return expr
}

// Bottom 底部N个
func Bottom(field string, n int) M {
	return M{"$bottom": M{
		"output": field,
		"sortBy": M{"$literal": 1},
		"n":      n,
	}}
}

// BottomN 底部N个
func BottomN(field string, n int) M {
	return M{"$bottomN": M{
		"output": field,
		"sortBy": M{"$literal": 1},
		"n":      n,
	}}
}

// Top 顶部N个
func Top(field string, n int) M {
	return M{"$top": M{
		"output": field,
		"sortBy": M{"$literal": -1},
		"n":      n,
	}}
}

// TopN 顶部N个
func TopN(field string, n int) M {
	return M{"$topN": M{
		"output": field,
		"sortBy": M{"$literal": -1},
		"n":      n,
	}}
}

// FirstN 前N个
func FirstN(field string, n int) M {
	return M{"$firstN": M{
		"input": field,
		"n":     n,
	}}
}

// LastN 后N个
func LastN(field string, n int) M {
	return M{"$lastN": M{
		"input": field,
		"n":     n,
	}}
}

// MaxN 最大N个
func MaxNAccumulator(field string, n int) M {
	return M{"$maxN": M{
		"input": field,
		"n":     n,
	}}
}

// MinNAccumulator 最小N个
func MinNAccumulator(field string, n int) M {
	return M{"$minN": M{
		"input": field,
		"n":     n,
	}}
}

// ========== 其他表达式 ==========

// Let 变量绑定
func Let(vars M, in any) M {
	return M{"$let": M{
		"vars": vars,
		"in":   in,
	}}
}

// Literal 字面值
func Literal(value any) M {
	return M{"$literal": value}
}

// ========== 常量 ==========

// Now 当前时间
func Now() string {
	return "$$NOW"
}

// Null 空值
func Null() any {
	return nil
}

// CurrentDate 当前日期
func CurrentDate() M {
	return M{"$currentDate": M{}}
}

// ROOT 根文档
func ROOT() string {
	return "$$ROOT"
}

// REMOVE 移除字段
func REMOVE() string {
	return "$$REMOVE"
}

// PRUNE 修剪字段
func PRUNE() string {
	return "$$PRUNE"
}

// KEEP 保留字段
func KEEP() string {
	return "$$KEEP"
}

// DESCEND 下降
func DESCEND() string {
	return "$$DESCEND"
}

// ========== 辅助函数 ==========

// Field 字段引用
func Field(name string) string {
	return "$" + name
}

// Var 变量引用
func Var(name string) string {
	return "$$" + name
}
