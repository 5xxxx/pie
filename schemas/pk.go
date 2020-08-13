/*
 *
 * pk.go
 * schemas
 *
 * Created by lintao on 2020/5/18 3:39 下午
 * Copyright © 2020-2020 LINTAO. All rights reserved.
 *
 */

package schemas

import (
	"bytes"
	"encoding/gob"

	"github.com/NSObjects/tugrik/utils"
)

type PK []interface{}

func NewPK(pks ...interface{}) *PK {
	p := PK(pks)
	return &p
}

func (p *PK) IsZero() bool {
	for _, k := range *p {
		if utils.IsZero(k) {
			return true
		}
	}
	return false
}

func (p *PK) ToString() (string, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(*p)
	return buf.String(), err
}

func (p *PK) FromString(content string) error {
	dec := gob.NewDecoder(bytes.NewBufferString(content))
	err := dec.Decode(p)
	return err
}
