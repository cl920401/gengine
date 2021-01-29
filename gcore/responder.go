package gcore

import (
	"fmt"
	"gengine/format"
	"github.com/gin-gonic/gin"
	"reflect"
	"sync"
)

var responderList []Responder
var once_resp_list sync.Once

func get_responder_list() []Responder {
	once_resp_list.Do(func() {
		responderList = []Responder{(StringResponder)(nil),
			(JsonResponder)(nil),
			(ViewResponder)(nil),
			(SqlResponder)(nil),
			(SqlQueryResponder)(nil),
		}
	})
	return responderList
}

type Responder interface {
	RespondTo() gin.HandlerFunc
}

func Convert(handler interface{}) gin.HandlerFunc {
	h_ref := reflect.ValueOf(handler)
	for _, r := range get_responder_list() {
		r_ref := reflect.TypeOf(r)
		if h_ref.Type().ConvertibleTo(r_ref) {
			return h_ref.Convert(r_ref).Interface().(Responder).RespondTo()
		}
	}
	return nil
}

type StringResponder func(*Context) string

func (this StringResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.String(200, getFilters().Filter(this, &Context{Context: context}).(string))
	}
}

type Json interface{}
type JsonResponder func(*Context) Json

func (this JsonResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.JSON(200, getFilters().Filter(this, &Context{Context: context}))
	}
}

type SqlQueryResponder func(*Context) Query

func (this SqlQueryResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		getQuery := getFilters().Filter(this, &Context{Context: context}).(Query)
		ret, err := format.QueryForMapsByInterface(getQuery)
		if err != nil {
			panic(err)
		}
		context.JSON(200, ret)
	}
}

type SqlResponder func(*Context) SimpleQuery

func (this SqlResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		getSql := getFilters().Filter(this, &Context{Context: context}).(SimpleQuery)
		ret, err := format.QueryForMaps(string(getSql), nil, []interface{}{}...)
		if err != nil {
			panic(err)
		}
		context.JSON(200, ret)
	}
}

// Deprecated: 暂时不提供View的解析
type View string
type ViewResponder func(*Context) View

func (this ViewResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.HTML(200, string(this(&Context{Context: context}))+".html", context.Keys)
	}
}

type Query interface {
	Sql() string
	Args() []interface{}
	Mapping() map[string]string
	First() bool
	Key() string
	Get() interface{}
}

type SimpleQueryWithArgs struct {
	sql        string
	args       []interface{}
	mapping    map[string]string
	fetchFirst bool
	datakey    string //data:{datakey:xxxx}   add by shenyi
}

func NewSimpleQueryWithArgs(sql string, args []interface{}) *SimpleQueryWithArgs {
	return &SimpleQueryWithArgs{sql: sql, args: args}
}
func NewSimpleQueryWithMapping(sql string, mapping map[string]string) *SimpleQueryWithArgs {
	return &SimpleQueryWithArgs{sql: sql, mapping: mapping}
}
func NewSimpleQueryWithFetchFirst(sql string) *SimpleQueryWithArgs {
	return &SimpleQueryWithArgs{sql: sql, fetchFirst: true}
}
func NewSimpleQueryWithKey(sql string, key string) *SimpleQueryWithArgs {
	return &SimpleQueryWithArgs{sql: sql, datakey: key}
}
func (this *SimpleQueryWithArgs) Sql() string {
	return this.sql
}
func (this *SimpleQueryWithArgs) Mapping() map[string]string {
	return this.mapping
}
func (this *SimpleQueryWithArgs) Args() []interface{} {
	return this.args
}
func (this *SimpleQueryWithArgs) First() bool {
	return this.fetchFirst
}
func (this *SimpleQueryWithArgs) Key() string {
	return this.datakey
}
func (this *SimpleQueryWithArgs) WithMapping(mapping map[string]string) *SimpleQueryWithArgs {
	this.mapping = mapping
	return this
}
func (this *SimpleQueryWithArgs) WithFirst() *SimpleQueryWithArgs {
	this.fetchFirst = true
	return this
}
func (this *SimpleQueryWithArgs) WithKey(key string) *SimpleQueryWithArgs {
	this.datakey = key
	return this
}
func (this *SimpleQueryWithArgs) Get() interface{} {
	ret, err := format.QueryForMapsByInterface(this)
	if err != nil {
		fmt.Println("query get error:", err)
		return nil
	}
	return ret
}

type SimpleQuery string

func (this SimpleQuery) WithArgs(args ...interface{}) *SimpleQueryWithArgs {
	return NewSimpleQueryWithArgs(string(this), args)
}
func (this SimpleQuery) WithMapping(mapping map[string]string) *SimpleQueryWithArgs {
	return NewSimpleQueryWithMapping(string(this), mapping)
}
func (this SimpleQuery) WithFirst() *SimpleQueryWithArgs {
	return NewSimpleQueryWithFetchFirst(string(this))
}
func (this SimpleQuery) WithKey(key string) *SimpleQueryWithArgs {
	return NewSimpleQueryWithKey(string(this), key)
}
func (this SimpleQuery) First() bool {
	return false
}
func (this SimpleQuery) Sql() string {
	return string(this)
}
func (this SimpleQuery) Key() string {
	return ""
}
func (this SimpleQuery) Args() []interface{} {
	return []interface{}{}
}
func (this SimpleQuery) Mapping() map[string]string {
	return map[string]string{}
}
func (this SimpleQuery) Get() interface{} {
	return NewSimpleQueryWithArgs(string(this), nil).Get()
}
