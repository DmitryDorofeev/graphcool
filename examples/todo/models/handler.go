// auto generated file
package models

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/DmitryDorofeev/graphcool/errors"
	"github.com/DmitryDorofeev/graphcool/query"
	"io/ioutil"
	"net/http"
	"strconv"
)

type GQLHandler struct {
}

type Request struct {
	Query string `json:"query"`
}

func NewHandler() GQLHandler {
	return GQLHandler{}
}

type Meta struct {
	Nullable bool
	List     bool
}

type StringMeta struct {
	Meta
	Value string
}

func (s *StringMeta) Marshal() []byte {
	return []byte("\"" + s.Value + "\"")
}

type IntMeta struct {
	Meta
	Value int
}

func (s *IntMeta) Marshal() []byte {
	return []byte(strconv.Itoa(s.Value))
}

type BoolMeta struct {
	Meta
	Value bool
}

func (s *BoolMeta) Marshal() []byte {
	return []byte(strconv.FormatBool(s.Value))
}

type UserMeta struct {
	Meta
	Value User
}

func (s *UserMeta) Lookup(ctx context.Context, selections []query.Selection) ([]byte, *errors.QueryError) {
	if len(selections) == 0 {
		return nil, errors.Errorf("Objects must have selections (field User has no selections)")
	}
	res := make([][]byte, 0)
	for _, selection := range selections {
		field, ok := selection.(*query.Field)
		if !ok {
			fmt.Println("cannot cast to Field")
			continue
		}
		switch field.Name.Name {

		case "name":
			f := StringMeta{
				Value: s.Value.Name,
			}

			res = append(res, f.Marshal())

		case "friends":
			f := FriendsMeta{}

			err := f.Value.Resolve(ctx, s.Value)
			if err != nil {
				return nil, err
			}

			innerField, err := f.Lookup(ctx, field.Selections)
			if err != nil {
				return nil, err
			}

			res = append(res, innerField)

		default:
			return nil, errors.Errorf("unknown field " + field.Name.Name)
		}
	}
	return s.Marshal(ctx, selections, res)
}

func (s *UserMeta) Marshal(ctx context.Context, selections []query.Selection, fields [][]byte) ([]byte, *errors.QueryError) {
	fmt.Println("Marshal User")
	buf := bytes.Buffer{}
	buf.WriteString("{")
	for i, value := range fields {
		field, ok := selections[i].(*query.Field)
		if !ok {
			continue
		}

		if i != 0 {
			buf.WriteString(",")
		}

		buf.WriteString("\"" + field.Name.Name + "\"")
		buf.WriteString(":")
		buf.Write(value)
	}
	buf.WriteString("}")
	return buf.Bytes(), nil
}

type FriendsMeta struct {
	Meta
	Value Friends
}

func (s *FriendsMeta) Lookup(ctx context.Context, selections []query.Selection) ([]byte, *errors.QueryError) {
	if len(selections) == 0 {
		return nil, errors.Errorf("Objects must have selections (field Friends has no selections)")
	}
	res := make([][]byte, 0)
	for _, item := range s.Value {
		f := UserMeta{
			Value: item,
		}

		b, _ := f.Lookup(ctx, selections)
		res = append(res, b)
	}
	return s.Marshal(ctx, selections, res)
}

func (s *FriendsMeta) Marshal(ctx context.Context, selections []query.Selection, fields [][]byte) ([]byte, *errors.QueryError) {
	fmt.Println("Marshal Friends")
	buf := bytes.Buffer{}
	buf.WriteString("[")
	for i, value := range fields {
		if i != 0 {
			buf.WriteString(",")
		}

		buf.Write(value)
	}
	buf.WriteString("]")
	return buf.Bytes(), nil
}

type QueryMeta struct {
	Meta
	Value Query
}

func (s *QueryMeta) Lookup(ctx context.Context, selections []query.Selection) ([]byte, *errors.QueryError) {
	if len(selections) == 0 {
		return nil, errors.Errorf("Objects must have selections (field Query has no selections)")
	}
	res := make([][]byte, 0)
	for _, selection := range selections {
		field, ok := selection.(*query.Field)
		if !ok {
			fmt.Println("cannot cast to Field")
			continue
		}
		switch field.Name.Name {

		case "todo":
			f := TaskMeta{}

			err := f.Value.Resolve(ctx, s.Value)
			if err != nil {
				return nil, err
			}

			innerField, err := f.Lookup(ctx, field.Selections)
			if err != nil {
				return nil, err
			}

			res = append(res, innerField)

		default:
			return nil, errors.Errorf("unknown field " + field.Name.Name)
		}
	}
	return s.Marshal(ctx, selections, res)
}

func (s *QueryMeta) Marshal(ctx context.Context, selections []query.Selection, fields [][]byte) ([]byte, *errors.QueryError) {
	fmt.Println("Marshal Query")
	buf := bytes.Buffer{}
	buf.WriteString("{")
	for i, value := range fields {
		field, ok := selections[i].(*query.Field)
		if !ok {
			continue
		}

		if i != 0 {
			buf.WriteString(",")
		}

		buf.WriteString("\"" + field.Name.Name + "\"")
		buf.WriteString(":")
		buf.Write(value)
	}
	buf.WriteString("}")
	return buf.Bytes(), nil
}

type TaskMeta struct {
	Meta
	Value Task
}

func (s *TaskMeta) Lookup(ctx context.Context, selections []query.Selection) ([]byte, *errors.QueryError) {
	if len(selections) == 0 {
		return nil, errors.Errorf("Objects must have selections (field Task has no selections)")
	}
	res := make([][]byte, 0)
	for _, selection := range selections {
		field, ok := selection.(*query.Field)
		if !ok {
			fmt.Println("cannot cast to Field")
			continue
		}
		switch field.Name.Name {

		case "title":
			f := StringMeta{
				Value: s.Value.Title,
			}

			res = append(res, f.Marshal())

		case "description":
			f := StringMeta{
				Value: s.Value.Description,
			}

			res = append(res, f.Marshal())

		case "done":
			f := BoolMeta{
				Value: s.Value.Done,
			}

			res = append(res, f.Marshal())

		case "user":
			f := UserMeta{}

			err := f.Value.Resolve(ctx, s.Value)
			if err != nil {
				return nil, err
			}

			innerField, err := f.Lookup(ctx, field.Selections)
			if err != nil {
				return nil, err
			}

			res = append(res, innerField)

		default:
			return nil, errors.Errorf("unknown field " + field.Name.Name)
		}
	}
	return s.Marshal(ctx, selections, res)
}

func (s *TaskMeta) Marshal(ctx context.Context, selections []query.Selection, fields [][]byte) ([]byte, *errors.QueryError) {
	fmt.Println("Marshal Task")
	buf := bytes.Buffer{}
	buf.WriteString("{")
	for i, value := range fields {
		field, ok := selections[i].(*query.Field)
		if !ok {
			continue
		}

		if i != 0 {
			buf.WriteString(",")
		}

		buf.WriteString("\"" + field.Name.Name + "\"")
		buf.WriteString(":")
		buf.Write(value)
	}
	buf.WriteString("}")
	return buf.Bytes(), nil
}

func (h GQLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	parseError := errors.Errorf("Cannot parse request")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errBytes, _ := json.Marshal(parseError)
		w.Write(errBytes)
		return
	}

	req := Request{}

	json.Unmarshal(body, &req)

	res, err := query.Parse(req.Query)

	for _, o := range res.Operations {
		switch o.Type {
		case query.Query:
			q := QueryMeta{}
			fields, err := q.Lookup(ctx, o.Selections)
			if err != nil {
				errBytes, _ := json.Marshal(err)
				w.Write(errBytes)
				return
			}

			w.Write(fields)
			return
		}
	}

	w.Write([]byte("ba"))
}
