/*
 * Copyright (C) 2017 Dgraph Labs, Inc. and Contributors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package dgo

import (
	"context"
	"encoding/json"
	"math/rand"
	"reflect"
	"sync"

	"github.com/gogo/protobuf/proto"
	"github.com/lehajam/dgo/protos/api"
	"github.com/lehajam/dgo/y"
	"github.com/mitchellh/mapstructure"
)

// Dgraph is a transaction aware client to a set of dgraph server instances.
type Dgraph struct {
	dc []api.DgraphClient

	mu      sync.Mutex
	linRead *api.LinRead
}

// NewDgraphClient creates a new Dgraph for interacting with the Dgraph store connected to in
// conns.
// The client can be backed by multiple connections (to the same server, or multiple servers in a
// cluster).
//
// A single client is thread safe for sharing with multiple go routines.
func NewDgraphClient(clients ...api.DgraphClient) *Dgraph {
	dg := &Dgraph{
		dc:      clients,
		linRead: &api.LinRead{},
	}

	return dg
}

func (d *Dgraph) mergeLinRead(src *api.LinRead) {
	d.mu.Lock()
	defer d.mu.Unlock()
	y.MergeLinReads(d.linRead, src)
}

func (d *Dgraph) getLinRead() *api.LinRead {
	d.mu.Lock()
	defer d.mu.Unlock()
	return proto.Clone(d.linRead).(*api.LinRead)
}

// By setting various fields of api.Operation, Alter can be used to do the
// following:
//
// 1. Modify the schema.
//
// 2. Drop a predicate.
//
// 3. Drop the database.
func (d *Dgraph) Alter(ctx context.Context, op *api.Operation) error {
	dc := d.anyClient()
	_, err := dc.Alter(ctx, op)
	return err
}

func (d *Dgraph) anyClient() api.DgraphClient {
	return d.dc[rand.Intn(len(d.dc))]
}

// DeleteEdges sets the edges corresponding to predicates on the node with the given uid
// for deletion.
// This helper function doesn't run the mutation on the server. It must be done by the user
// after the function returns.
func DeleteEdges(mu *api.Mutation, uid string, predicates ...string) {
	for _, predicate := range predicates {
		mu.Del = append(mu.Del, &api.NQuad{
			Subject:   uid,
			Predicate: predicate,
			// _STAR_ALL is defined as x.Star in x package.
			ObjectValue: &api.Value{&api.Value_DefaultVal{"_STAR_ALL"}},
		})
	}
}

// Unmarshal decodes the content of a dgraph response into a struct
// Dgraph returns any embedded struct in the response as a json array
// When that's the case, we use a hook to convert the json array back into a struct
func Unmarshal(data []byte, obj interface{}) error {
	var parsed map[string]interface{}
	err := json.Unmarshal(data, &parsed)
	if err != nil {
		return err
	}

	decoder, err := mapstructure.NewDecoder(
		&mapstructure.DecoderConfig{
			Result: &obj,
			DecodeHook: func(from reflect.Kind, to reflect.Kind, data interface{}) (interface{}, error) {
				if from == reflect.Slice && to == reflect.Struct {
					slice := data.([]interface{})
					if len(slice) > 0 {
						return slice[0], nil
					}
				}
				return data, nil
			},
		})
	if err != nil {
		return err
	}

	if err := decoder.Decode(parsed); err != nil {
		return err
	}

	return nil
}
