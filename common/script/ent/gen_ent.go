package main

import (
	"database/sql"
	"fmt"
	"go/format"
	"os"
	"strings"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type TableColumn struct {
	Field      string
	Type       string
	Null       string
	Key        string
	Default    sql.NullString
	Extra      sql.NullString
	Privileges sql.NullString
	Comment    string
}

type TableInfo struct {
	Name    string
	Columns []TableColumn
}

const schemaTemplate = `package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// {{.Name | ToCamel}} holds the schema definition for the {{.Name | ToCamel}} entity.
type {{.Name | ToCamel}} struct {
	ent.Schema
}

// Fields of the {{.Name | ToCamel}}.
func ({{.Name | ToCamel}}) Fields() []ent.Field {
	return []ent.Field{
{{- range .Columns}}
		{{. | GenerateField}},
{{- end}}
	}
}

// Edges of the {{.Name | ToCamel}}.
func ({{.Name | ToCamel}}) Edges() []ent.Edge {
	return []ent.Edge{}
}

func ({{.Name | ToCamel}}) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "{{.Name}}"},
	}
}
`

func main() {
	// 连接数据库
	dsn := "root:root123@tcp(localhost:3306)/db1"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// 获取所有表名
	tables, err := getTables(db)
	if err != nil {
		panic(err)
	}

	// 为每个表生成 schema
	for _, table := range tables {
		err := generateSchema(db, table)
		if err != nil {
			fmt.Printf("Error generating schema for table %s: %v\n", table, err)
		} else {
			fmt.Printf("Successfully generated schema for table %s\n", table)
		}
	}
}

func getTables(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}
	return tables, nil
}

func generateSchema(db *sql.DB, tableName string) error {
	// 获取表结构
	columns, err := getTableColumns(db, tableName)
	if err != nil {
		return err
	}

	// 准备模板数据
	tableInfo := TableInfo{
		Name:    tableName,
		Columns: columns,
	}

	// 解析模板
	funcMap := template.FuncMap{
		"ToCamel":       toCamelCase,
		"GenerateField": generateField,
		"EscapeString":  escapeString,
	}

	tmpl, err := template.New("schema").Funcs(funcMap).Parse(schemaTemplate)
	if err != nil {
		return err
	}

	// 生成代码
	var buf strings.Builder
	err = tmpl.Execute(&buf, tableInfo)
	if err != nil {
		return err
	}

	// 格式化代码
	formatted, err := format.Source([]byte(buf.String()))
	if err != nil {
		// 如果格式化失败，至少保存原始代码
		formatted = []byte(buf.String())
	}

	// 确保目录存在
	err = os.MkdirAll("./ent/schema", 0755)
	if err != nil {
		return err
	}

	// 写入文件
	filename := fmt.Sprintf("./ent/schema/%s.go", tableName)
	return os.WriteFile(filename, formatted, 0644)
}

func getTableColumns(db *sql.DB, tableName string) ([]TableColumn, error) {
	query := fmt.Sprintf("SHOW FULL COLUMNS FROM `%s`", tableName)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []TableColumn
	for rows.Next() {
		var col TableColumn
		var collation, extra, privileges sql.NullString
		err := rows.Scan(
			&col.Field, &col.Type, &collation, &col.Null, &col.Key,
			&col.Default, &extra, &privileges, &col.Comment,
		)
		if err != nil {
			return nil, err
		}
		col.Extra = extra
		col.Privileges = privileges
		columns = append(columns, col)
	}
	return columns, nil
}

func generateField(col TableColumn) string {
	// 根据 MySQL 类型生成 Ent 字段
	fieldType := mapMySQLTypeToEntType(col.Type)

	// 构建字段定义
	fieldDef := fmt.Sprintf("field.%s(\"%s\")", fieldType, col.Field)

	// 添加主键约束
	if col.Key == "PRI" {
		fieldDef += ".Immutable().Unique()"
	}

	// 添加可选性
	if col.Null == "YES" {
		fieldDef += ".Optional()"
	}

	// 添加默认值
	if col.Default.Valid && col.Default.String != "" {
		defaultValue := formatDefaultValue(col.Default.String, fieldType)
		if defaultValue != "" {
			fieldDef += fmt.Sprintf(".Default(%s)", defaultValue)
		}
	}

	// 添加注释（正确转义）
	if col.Comment != "" {
		escapedComment := escapeString(col.Comment)
		fieldDef += fmt.Sprintf(".Comment(\"%s\")", escapedComment)
	}

	// 添加 JSON 标签
	fieldDef += fmt.Sprintf(".StructTag(`json:\"%s\"`)", col.Field)

	// 添加长度限制
	if strings.Contains(col.Type, "varchar") {
		if start := strings.Index(col.Type, "("); start != -1 {
			if end := strings.Index(col.Type, ")"); end != -1 {
				length := col.Type[start+1 : end]
				if length != "" {
					fieldDef += fmt.Sprintf(".MaxLen(%s)", length)
				}
			}
		}
	}

	return fieldDef
}

func mapMySQLTypeToEntType(mysqlType string) string {
	mysqlType = strings.ToLower(mysqlType)

	switch {
	case strings.Contains(mysqlType, "bigint"):
		return "Int64"
	case strings.Contains(mysqlType, "tinyint"):
		return "Int8"
	case strings.Contains(mysqlType, "smallint"):
		return "Int16"
	case strings.Contains(mysqlType, "mediumint"):
		return "Int32"
	case strings.Contains(mysqlType, "int"):
		return "Int"
	case strings.Contains(mysqlType, "varchar") || strings.Contains(mysqlType, "char"):
		return "String"
	case strings.Contains(mysqlType, "text") || strings.Contains(mysqlType, "mediumtext") || strings.Contains(mysqlType, "longtext"):
		return "String"
	case strings.Contains(mysqlType, "datetime") || strings.Contains(mysqlType, "timestamp"):
		return "Time"
	case strings.Contains(mysqlType, "date"):
		return "Time"
	case strings.Contains(mysqlType, "decimal") || strings.Contains(mysqlType, "float") || strings.Contains(mysqlType, "double"):
		return "Float"
	case strings.Contains(mysqlType, "enum") || strings.Contains(mysqlType, "set"):
		return "String"
	default:
		return "String"
	}
}

func formatDefaultValue(value, fieldType string) string {
	// 处理 NULL 值
	if strings.ToUpper(value) == "NULL" {
		return "" // 不设置默认值
	}

	switch fieldType {
	case "String":
		// 特殊处理空字符串
		if value == "" {
			return "\"\""
		}
		escapedValue := escapeString(value)
		return fmt.Sprintf("\"%s\"", escapedValue)
	case "Int", "Int8", "Int16", "Int32", "Int64":
		if value == "" {
			return "" // 不设置默认值
		}
		return value
	case "Float":
		if value == "" {
			return "" // 不设置默认值
		}
		return value
	case "Time":
		// 时间类型的默认值特殊处理
		if value == "" || strings.ToUpper(value) == "CURRENT_TIMESTAMP" {
			return "" // 不设置默认值
		}
		escapedValue := escapeString(value)
		return fmt.Sprintf("\"%s\"", escapedValue)
	default:
		if value == "" {
			return "" // 不设置默认值
		}
		escapedValue := escapeString(value)
		return fmt.Sprintf("\"%s\"", escapedValue)
	}
}

func escapeString(s string) string {
	// 转义字符串中的特殊字符
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}

func toCamelCase(s string) string {
	parts := strings.Split(s, "_")
	caser := cases.Title(language.English)
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = caser.String(strings.ToLower(part))
		}
	}
	return strings.Join(parts, "")
}
