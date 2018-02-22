// auto generated file
package models

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

type FriendsMeta struct {
	Meta
	Value Friends
}

func (s *FriendsMeta) Lookup(ctx context.Context, selections []query.Selection) ([]byte, error) {
	fmt.Println("Looking up Friends")
	res := make([][]byte, 0)
	for _, selection := range selections {
		field, ok := selection.(*query.Field)
		if !ok {
			fmt.Println("cannot cast to Field")
			continue
		}
		switch field.Name.Name {

		default:
			return nil, fmt.Errorf("unknown field " + field.Name.Name)
		}
	}
	return s.Marshal(ctx, selections, res)
}

func (s *FriendsMeta) Marshal(ctx context.Context, selections []query.Selection, fields [][]byte) ([]byte, error) {
	fmt.Println("Marshal Friends")
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

type QueryMeta struct {
	Meta
	Value Query
}

func (s *QueryMeta) Lookup(ctx context.Context, selections []query.Selection) ([]byte, error) {
	fmt.Println("Looking up Query")
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

			err := f.Value.Resolve(ctx)
			if err != nil {
				return nil, err
			}

			innerField, err := f.Lookup(ctx, field.Selections)
			if err != nil {
				return nil, err
			}

			res = append(res, innerField)

		default:
			return nil, fmt.Errorf("unknown field " + field.Name.Name)
		}
	}
	return s.Marshal(ctx, selections, res)
}

func (s *QueryMeta) Marshal(ctx context.Context, selections []query.Selection, fields [][]byte) ([]byte, error) {
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

func (s *TaskMeta) Lookup(ctx context.Context, selections []query.Selection) ([]byte, error) {
	fmt.Println("Looking up Task")
	res := make([][]byte, 0)
	for _, selection := range selections {
		field, ok := selection.(*query.Field)
		if !ok {
			fmt.Println("cannot cast to Field")
			continue
		}
		switch field.Name.Name {

		case "title":
			fmt.Println("touch title")
			f := StringMeta{
				Value: s.Value.Title,
			}

			res = append(res, f.Marshal())

		case "description":
			fmt.Println("touch description")
			f := StringMeta{
				Value: s.Value.Description,
			}

			res = append(res, f.Marshal())

		case "done":
			fmt.Println("touch done")
			f := BoolMeta{
				Value: s.Value.Done,
			}

			res = append(res, f.Marshal())

		case "user":
			f := UserMeta{}

			err := f.Value.Resolve(ctx)
			if err != nil {
				return nil, err
			}

			innerField, err := f.Lookup(ctx, field.Selections)
			if err != nil {
				return nil, err
			}

			res = append(res, innerField)

		default:
			return nil, fmt.Errorf("unknown field " + field.Name.Name)
		}
	}
	return s.Marshal(ctx, selections, res)
}

func (s *TaskMeta) Marshal(ctx context.Context, selections []query.Selection, fields [][]byte) ([]byte, error) {
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

type UserMeta struct {
	Meta
	Value User
}

func (s *UserMeta) Lookup(ctx context.Context, selections []query.Selection) ([]byte, error) {
	fmt.Println("Looking up User")
	res := make([][]byte, 0)
	for _, selection := range selections {
		field, ok := selection.(*query.Field)
		if !ok {
			fmt.Println("cannot cast to Field")
			continue
		}
		switch field.Name.Name {

		case "name":
			fmt.Println("touch name")
			f := StringMeta{
				Value: s.Value.Name,
			}

			res = append(res, f.Marshal())

		case "friends":
			f := FriendsMeta{}

			err := f.Value.Resolve(ctx)
			if err != nil {
				return nil, err
			}

			innerField, err := f.Lookup(ctx, field.Selections)
			if err != nil {
				return nil, err
			}

			res = append(res, innerField)

		default:
			return nil, fmt.Errorf("unknown field " + field.Name.Name)
		}
	}
	return s.Marshal(ctx, selections, res)
}

func (s *UserMeta) Marshal(ctx context.Context, selections []query.Selection, fields [][]byte) ([]byte, error) {
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

func (h GQLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("error"))
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
				w.Write([]byte("pzdc: " + err.Error()))
				return
			}

			w.Write(fields)
			return
		}
	}

	w.Write([]byte("ba"))
}