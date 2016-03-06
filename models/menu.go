// Code generated by project
// menu.go contains model for the database table [dbname=baseinfo sslmode=disable.menu]

package models

import (
	"encoding/json"
	"encoding/gob"
	"fmt"
	"strings"
	"github.com/jmoiron/sqlx"
	"database/sql"
	"time"
)
type Menu struct {
	Id int64 `json:"id"`
	Name string `json:"name"`
	RestaurantId int64 `json:"restaurant_id"`
	CorpId int64 `json:"corp_id"`
	PeriodBegin string `json:"period_begin"`
	PeriodEnd string `json:"period_end"`
	FilterBy string `json:"filter_by"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
	Creator int64 `json:"creator"`
	Operator int64 `json:"operator"`
	Removed bool `json:"removed"`
	
}

// Start of the Menu APIs.

func (obj Menu) String() string {
	if data, err := json.Marshal(obj); err != nil {
		return fmt.Sprintf("<Menu Id=%v>", obj.Id)
	} else {
		return string(data)
	}
}

func  QueryById(id int64) (*Menu,error){

 	ret:=&Menu{}
 	err:=db.Get(ret,"select * from dbname=baseinfo sslmode=disable where id=",id)
 	return ret,err
}

func InsertMenu(obj Menu) error(){

  db.Exec("insert into dbname=baseinfo sslmode=disable(id,name,restaurant_id,corp_id,period_begin,period_end,filter_by,create_time,update_time,creator,operator,removed) values(obj.Id,obj.Name,obj.RestaurantId,obj.CorpId,obj.PeriodBegin,obj.PeriodEnd,obj.FilterBy,obj.CreateTime,obj.UpdateTime,obj.Creator,obj.Operator,obj.Removed)",obj)
  return nil
}

func (obj Menu) Get(dbtx gmq.DbTx) (Menu, error) {
	filter := MenuObjs.FilterId("=", obj.Id)
	if result, err := MenuObjs.Select().Where(filter).One(dbtx); err != nil {
		return obj, err
	} else {
		return result, nil
	}
}

func (obj Menu) Insert(dbtx gmq.DbTx) (Menu, error) {
	if result, err := MenuObjs.Insert(obj).Run(dbtx); err != nil {
		return obj, err
	}else {
		if dbtx.DriverName() != "postgres" {
			if id, err := result.LastInsertId(); err != nil {
				return obj, err
			} else {
				obj.Id = id
				return obj, err
			}
		}
		return obj, nil
	}
}

func (obj Menu) Update(dbtx gmq.DbTx) (int64, error) {
	fields := []string{ "Name", "RestaurantId", "CorpId", "PeriodBegin", "PeriodEnd", "FilterBy", "CreateTime", "UpdateTime", "Creator", "Operator", "Removed" }
	filter := MenuObjs.FilterId("=", obj.Id)
	if result, err := MenuObjs.Update(obj, fields...).Where(filter).Run(dbtx); err != nil {
		return 0, err
	} else {
		return result.RowsAffected()
	}
}

func (obj Menu) Delete(dbtx gmq.DbTx) (int64, error) {
	filter := MenuObjs.FilterId("=", obj.Id)
	if result, err := MenuObjs.Delete().Where(filter).Run(dbtx); err != nil {
		return 0, err
	} else {
		return result.RowsAffected()
	}
}

// Start of the inner Query Api

type _MenuQuery struct {
	gmq.Query
}

func (q _MenuQuery) Where(f gmq.Filter) _MenuQuery {
	q.Query = q.Query.Where(f)
	return q
}

func (q _MenuQuery) OrderBy(by ...string) _MenuQuery {
	tBy := make([]string, 0, len(by))
	for _, b := range by {
		sortDir := ""
		if b[0] == '-' || b[0] == '+' {
			sortDir = string(b[0])
			b = b[1:]
		}
		if col, ok := MenuObjs.fcMap[b]; ok {
			tBy = append(tBy, sortDir+col)
		}
	}
	q.Query = q.Query.OrderBy(tBy...)
	return q
}

func (q _MenuQuery) GroupBy(by ...string) _MenuQuery {
	tBy := make([]string, 0, len(by))
	for _, b := range by {
		if col, ok := MenuObjs.fcMap[b]; ok {
			tBy = append(tBy, col)
		}
	}
	q.Query = q.Query.GroupBy(tBy...)
	return q
}

func (q _MenuQuery) Limit(offsets ...int64) _MenuQuery {
	q.Query = q.Query.Limit(offsets...)
	return q
}

func (q _MenuQuery) Page(number, size int) _MenuQuery {
	q.Query = q.Query.Page(number, size)
	return q
}

func (q _MenuQuery) Run(dbtx gmq.DbTx) (sql.Result, error) {
	return q.Query.Exec(dbtx)
}

type MenuRowVisitor func(obj Menu) bool

func (q _MenuQuery) Iterate(dbtx gmq.DbTx, functor MenuRowVisitor) error {
	return q.Query.SelectList(dbtx, func(columns []gmq.Column, rb []sql.RawBytes) bool {
		obj := MenuObjs.toMenu(columns, rb)
		return functor(obj)
	})
}

func (q _MenuQuery) One(dbtx gmq.DbTx) (Menu, error) {
	var obj Menu
	err := q.Query.SelectOne(dbtx, func(columns []gmq.Column, rb []sql.RawBytes) bool {
		obj = MenuObjs.toMenu(columns, rb)
		return true
	})
	return obj, err
}

func (q _MenuQuery) List(dbtx gmq.DbTx) ([]Menu, error) {
	result := make([]Menu, 0, 10)
	err := q.Query.SelectList(dbtx, func(columns []gmq.Column, rb []sql.RawBytes) bool {
		obj := MenuObjs.toMenu(columns, rb)
		result = append(result, obj)
		return true
	})
	return result, err
}

// Start of the model facade Apis.

type _MenuObjs struct {
	fcMap map[string]string
}

func (o _MenuObjs) Names() (schema, tbl, alias string) {
	return "dbname=baseinfo sslmode=disable", "menu", "Menu"
}

func (o _MenuObjs) Select(fields ...string) _MenuQuery {
	q := _MenuQuery{}
	if len(fields) == 0 {
		fields = []string{ "Id", "Name", "RestaurantId", "CorpId", "PeriodBegin", "PeriodEnd", "FilterBy", "CreateTime", "UpdateTime", "Creator", "Operator", "Removed" }
	}
	q.Query = gmq.Select(o, o.columns(fields...))
	return q
}

func (o _MenuObjs) Insert(obj Menu) _MenuQuery {
	q := _MenuQuery{}
	q.Query = gmq.Insert(o, o.columnsWithData(obj, "Name", "RestaurantId", "CorpId", "PeriodBegin", "PeriodEnd", "FilterBy", "Creator", "Operator", "Removed"))
	return q
}

func (o _MenuObjs) Update(obj Menu, fields ...string) _MenuQuery {
	q := _MenuQuery{}
	q.Query = gmq.Update(o, o.columnsWithData(obj, fields...))
	return q
}

func (o _MenuObjs) Delete() _MenuQuery {
	q := _MenuQuery{}
	q.Query = gmq.Delete(o)
	return q
}


///// Managed Objects Filters definition

func (o _MenuObjs) FilterId(op string, p int64, ps ...int64) gmq.Filter {
	params := make([]interface{}, 1+len(ps))
	params[0] = p
	for i := range ps {
		params[i+1] = ps[i]
	}
	return o.newFilter("id", op, params...)
}


func (o _MenuObjs) FilterName(op string, p string, ps ...string) gmq.Filter {
	params := make([]interface{}, 1+len(ps))
	params[0] = p
	for i := range ps {
		params[i+1] = ps[i]
	}
	return o.newFilter("name", op, params...)
}


func (o _MenuObjs) FilterRestaurantId(op string, p int64, ps ...int64) gmq.Filter {
	params := make([]interface{}, 1+len(ps))
	params[0] = p
	for i := range ps {
		params[i+1] = ps[i]
	}
	return o.newFilter("restaurant_id", op, params...)
}


func (o _MenuObjs) FilterCorpId(op string, p int64, ps ...int64) gmq.Filter {
	params := make([]interface{}, 1+len(ps))
	params[0] = p
	for i := range ps {
		params[i+1] = ps[i]
	}
	return o.newFilter("corp_id", op, params...)
}


func (o _MenuObjs) FilterPeriodBegin(op string, p string, ps ...string) gmq.Filter {
	params := make([]interface{}, 1+len(ps))
	params[0] = p
	for i := range ps {
		params[i+1] = ps[i]
	}
	return o.newFilter("period_begin", op, params...)
}


func (o _MenuObjs) FilterPeriodEnd(op string, p string, ps ...string) gmq.Filter {
	params := make([]interface{}, 1+len(ps))
	params[0] = p
	for i := range ps {
		params[i+1] = ps[i]
	}
	return o.newFilter("period_end", op, params...)
}


func (o _MenuObjs) FilterFilterBy(op string, p string, ps ...string) gmq.Filter {
	params := make([]interface{}, 1+len(ps))
	params[0] = p
	for i := range ps {
		params[i+1] = ps[i]
	}
	return o.newFilter("filter_by", op, params...)
}


func (o _MenuObjs) FilterCreateTime(op string, p time.Time, ps ...time.Time) gmq.Filter {
	params := make([]interface{}, 1+len(ps))
	params[0] = p
	for i := range ps {
		params[i+1] = ps[i]
	}
	return o.newFilter("create_time", op, params...)
}


func (o _MenuObjs) FilterUpdateTime(op string, p time.Time, ps ...time.Time) gmq.Filter {
	params := make([]interface{}, 1+len(ps))
	params[0] = p
	for i := range ps {
		params[i+1] = ps[i]
	}
	return o.newFilter("update_time", op, params...)
}


func (o _MenuObjs) FilterCreator(op string, p int64, ps ...int64) gmq.Filter {
	params := make([]interface{}, 1+len(ps))
	params[0] = p
	for i := range ps {
		params[i+1] = ps[i]
	}
	return o.newFilter("creator", op, params...)
}


func (o _MenuObjs) FilterOperator(op string, p int64, ps ...int64) gmq.Filter {
	params := make([]interface{}, 1+len(ps))
	params[0] = p
	for i := range ps {
		params[i+1] = ps[i]
	}
	return o.newFilter("operator", op, params...)
}


func (o _MenuObjs) FilterRemoved(op string, p bool, ps ...bool) gmq.Filter {
	params := make([]interface{}, 1+len(ps))
	params[0] = p
	for i := range ps {
		params[i+1] = ps[i]
	}
	return o.newFilter("removed", op, params...)
}



///// Managed Objects Columns definition

func (o _MenuObjs) ColumnId(p ...int64) gmq.Column {
	var value interface{}
	if len(p) > 0 {
		value = p[0]
	}
	return gmq.Column{"id", value}
}

func (o _MenuObjs) ColumnName(p ...string) gmq.Column {
	var value interface{}
	if len(p) > 0 {
		value = p[0]
	}
	return gmq.Column{"name", value}
}

func (o _MenuObjs) ColumnRestaurantId(p ...int64) gmq.Column {
	var value interface{}
	if len(p) > 0 {
		value = p[0]
	}
	return gmq.Column{"restaurant_id", value}
}

func (o _MenuObjs) ColumnCorpId(p ...int64) gmq.Column {
	var value interface{}
	if len(p) > 0 {
		value = p[0]
	}
	return gmq.Column{"corp_id", value}
}

func (o _MenuObjs) ColumnPeriodBegin(p ...string) gmq.Column {
	var value interface{}
	if len(p) > 0 {
		value = p[0]
	}
	return gmq.Column{"period_begin", value}
}

func (o _MenuObjs) ColumnPeriodEnd(p ...string) gmq.Column {
	var value interface{}
	if len(p) > 0 {
		value = p[0]
	}
	return gmq.Column{"period_end", value}
}

func (o _MenuObjs) ColumnFilterBy(p ...string) gmq.Column {
	var value interface{}
	if len(p) > 0 {
		value = p[0]
	}
	return gmq.Column{"filter_by", value}
}

func (o _MenuObjs) ColumnCreateTime(p ...time.Time) gmq.Column {
	var value interface{}
	if len(p) > 0 {
		value = p[0]
	}
	return gmq.Column{"create_time", value}
}

func (o _MenuObjs) ColumnUpdateTime(p ...time.Time) gmq.Column {
	var value interface{}
	if len(p) > 0 {
		value = p[0]
	}
	return gmq.Column{"update_time", value}
}

func (o _MenuObjs) ColumnCreator(p ...int64) gmq.Column {
	var value interface{}
	if len(p) > 0 {
		value = p[0]
	}
	return gmq.Column{"creator", value}
}

func (o _MenuObjs) ColumnOperator(p ...int64) gmq.Column {
	var value interface{}
	if len(p) > 0 {
		value = p[0]
	}
	return gmq.Column{"operator", value}
}

func (o _MenuObjs) ColumnRemoved(p ...bool) gmq.Column {
	var value interface{}
	if len(p) > 0 {
		value = p[0]
	}
	return gmq.Column{"removed", value}
}


////// Internal helper funcs

func (o _MenuObjs) newFilter(name, op string, params ...interface{}) gmq.Filter {
	if strings.ToUpper(op) == "IN" {
		return gmq.InFilter(name, params)
	}
	return gmq.UnitFilter(name, op, params[0])
}

func (o _MenuObjs) toMenu(columns []gmq.Column, rb []sql.RawBytes) Menu {
	obj := Menu{}
	if len(columns) == len(rb) {
		for i := range columns {
			switch columns[i].Name {
			case "id":
				obj.Id = gmq.AsInt64(rb[i])
			case "name":
				obj.Name = gmq.AsString(rb[i])
			case "restaurant_id":
				obj.RestaurantId = gmq.AsInt64(rb[i])
			case "corp_id":
				obj.CorpId = gmq.AsInt64(rb[i])
			case "period_begin":
				obj.PeriodBegin = gmq.AsString(rb[i])
			case "period_end":
				obj.PeriodEnd = gmq.AsString(rb[i])
			case "filter_by":
				obj.FilterBy = gmq.AsString(rb[i])
			case "create_time":
				obj.CreateTime = gmq.AsTime(rb[i])
			case "update_time":
				obj.UpdateTime = gmq.AsTime(rb[i])
			case "creator":
				obj.Creator = gmq.AsInt64(rb[i])
			case "operator":
				obj.Operator = gmq.AsInt64(rb[i])
			case "removed":
				obj.Removed = gmq.AsBool(rb[i])
			 }
		}
	}
	return obj
}

func (o _MenuObjs) columns(fields ...string) []gmq.Column {
	data := make([]gmq.Column, 0, len(fields))
	for _, f := range fields {
		switch f {
		case "Id":
			data = append(data, o.ColumnId())
		case "Name":
			data = append(data, o.ColumnName())
		case "RestaurantId":
			data = append(data, o.ColumnRestaurantId())
		case "CorpId":
			data = append(data, o.ColumnCorpId())
		case "PeriodBegin":
			data = append(data, o.ColumnPeriodBegin())
		case "PeriodEnd":
			data = append(data, o.ColumnPeriodEnd())
		case "FilterBy":
			data = append(data, o.ColumnFilterBy())
		case "CreateTime":
			data = append(data, o.ColumnCreateTime())
		case "UpdateTime":
			data = append(data, o.ColumnUpdateTime())
		case "Creator":
			data = append(data, o.ColumnCreator())
		case "Operator":
			data = append(data, o.ColumnOperator())
		case "Removed":
			data = append(data, o.ColumnRemoved())
		 }
	}
	return data
}

func (o _MenuObjs) columnsWithData(obj Menu, fields ...string) []gmq.Column {
	data := make([]gmq.Column, 0, len(fields))
	for _, f := range fields {
		switch f {
		case "Id":
			data = append(data, o.ColumnId(obj.Id))
		case "Name":
			data = append(data, o.ColumnName(obj.Name))
		case "RestaurantId":
			data = append(data, o.ColumnRestaurantId(obj.RestaurantId))
		case "CorpId":
			data = append(data, o.ColumnCorpId(obj.CorpId))
		case "PeriodBegin":
			data = append(data, o.ColumnPeriodBegin(obj.PeriodBegin))
		case "PeriodEnd":
			data = append(data, o.ColumnPeriodEnd(obj.PeriodEnd))
		case "FilterBy":
			data = append(data, o.ColumnFilterBy(obj.FilterBy))
		case "CreateTime":
			data = append(data, o.ColumnCreateTime(obj.CreateTime))
		case "UpdateTime":
			data = append(data, o.ColumnUpdateTime(obj.UpdateTime))
		case "Creator":
			data = append(data, o.ColumnCreator(obj.Creator))
		case "Operator":
			data = append(data, o.ColumnOperator(obj.Operator))
		case "Removed":
			data = append(data, o.ColumnRemoved(obj.Removed))
		 }
	}
	return data
}

var MenuObjs _MenuObjs

func init() {
	MenuObjs.fcMap = map[string]string{
		"Id": "id",
		"Name": "name",
		"RestaurantId": "restaurant_id",
		"CorpId": "corp_id",
		"PeriodBegin": "period_begin",
		"PeriodEnd": "period_end",
		"FilterBy": "filter_by",
		"CreateTime": "create_time",
		"UpdateTime": "update_time",
		"Creator": "creator",
		"Operator": "operator",
		"Removed": "removed",
		 }
	gob.Register(Menu{})
}
