/*
	Copyright © 2014–2021 Thomas Michael Edwards. All rights reserved.
	Use of this source code is governed by a Simplified BSD License which
	can be found in the LICENSE file.
*/

package main

import (
	"encoding/json"
)

type passageJSON struct {
	Name string   `json:"name"`
	Tags []string `json:"tags,omitempty"`
	Text string   `json:"text"`
}

type passageMetadataJSON struct {
	Position string `json:"position,omitempty"` // Twine 2 (`position`) & Twine 1 (`twine-position`).
	Size     string `json:"size,omitempty"`     // Twine 2 (`size`).
}

func (p *passage) marshalMetadata() []byte {
	marshaled, err := json.Marshal(&passageMetadataJSON{
		p.metadata.position,
		p.metadata.size,
	})
	if err != nil {
		// NOTE: We should never be able to see an error here.  If we do,
		// then something truly exceptional—in a bad way—has happened, so
		// we get our panic on.
		panic(err)
	}
	return marshaled
}

func (p *passage) unmarshalMetadata(marshaled []byte) error {
	metadata := passageMetadataJSON{}
	if err := json.Unmarshal(marshaled, &metadata); err != nil {
		return err
	}
	// Drop invalid position data.
	if metadata.Position == "NaN,NaN" {
		metadata.Position = ""
	}
	p.metadata = &passageMetadata{
		position: metadata.Position,
		size:     metadata.Size,
	}
	return nil
}
