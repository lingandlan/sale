package model

import (
	"reflect"
	"testing"
)

// allGormModels 列出所有参与 GORM AutoMigrate 的 model struct
// 新增 model 时必须在此注册
var allGormModels = []any{
	User{},
	RechargeApplication{},
	CRecharge{},
	StoreCard{},
	CardTransaction{},
	RechargeCenter{},
	RechargeOperator{},
}

// TestAllModelsHaveTableName 确保每个 GORM model 都显式定义了 TableName() 方法
// 防止 GORM 默认复数推导与实际表名不一致
func TestAllModelsHaveTableName(t *testing.T) {
	for _, m := range allGormModels {
		typ := reflect.TypeOf(m)
		name := typ.Name()

		// 检查是否有 TableName() string 方法
		method, exists := typ.MethodByName("TableName")
		if !exists {
			t.Errorf("%s: 缺少 TableName() 方法，请添加 func (%s) TableName() string", name, name)
			continue
		}

		// 验证方法签名：无参数（除 receiver），返回 string
		if method.Type.NumIn() != 1 || method.Type.NumOut() != 1 {
			t.Errorf("%s: TableName() 签名错误，应为 func (%s) TableName() string", name, name)
			continue
		}

		if method.Type.Out(0).Kind() != reflect.String {
			t.Errorf("%s: TableName() 返回值不是 string", name)
			continue
		}

		// 调用方法获取表名，验证非空
		val := reflect.New(typ)
		results := method.Func.Call([]reflect.Value{val.Elem()})
		tableName := results[0].String()
		if tableName == "" {
			t.Errorf("%s: TableName() 返回空字符串", name)
		}

		t.Logf("%s → %s ✓", name, tableName)
	}
}
