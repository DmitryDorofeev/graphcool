// hello
package models

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"github.com/DmitryDorofeev/graphcool/query"
)

type GQLHandler struct {
}

type Request struct {
	Query string `json:"query"`
}

func NewHandler() GQLHandler {
	return GQLHandler{}
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
		for _, s := range o.Selections {
			if ident, ok := s.(*query.Field); ok {
				
					
						switch(ident.Name.Name) {
							
								case "todo":
									todo := &Todo{}
									todo.Resolve(ctx)
									fmt.Println(ident.Name.Name)
							
								default:
									fmt.Println("wrong!")
						}
					
				
					
				
					
				
			}
		}
	}

	w.Write([]byte("ba"))
}

