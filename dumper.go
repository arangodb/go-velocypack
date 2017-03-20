//
// DISCLAIMER
//
// Copyright 2017 ArangoDB GmbH, Cologne, Germany
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Copyright holder is ArangoDB GmbH, Cologne, Germany
//
// Author Ewout Prangsma
//

package velocypack

import (
	"io"
	"strconv"
)

type Dumper struct {
	w           io.Writer
	indentation uint
}

func NewDumper(w io.Writer) *Dumper {
	return &Dumper{
		w: w,
	}
}

func (d *Dumper) Append(s Slice) error {
	w := d.w
	switch s.Type() {
	case Null:
		if _, err := w.Write([]byte("null")); err != nil {
			return WithStack(err)
		}
		return nil
	case Bool:
		if v, err := s.GetBool(); err != nil {
			return WithStack(err)
		} else if v {
			if _, err := w.Write([]byte("true")); err != nil {
				return WithStack(err)
			}
		} else {
			if _, err := w.Write([]byte("false")); err != nil {
				return WithStack(err)
			}
		}
		return nil
	case Double:
		if v, err := s.GetDouble(); err != nil {
			return WithStack(err)
		} else if err := d.appendDouble(v); err != nil {
			return WithStack(err)
		}
		return nil
	case Int, SmallInt:
		if v, err := s.GetInt(); err != nil {
			return WithStack(err)
		} else if err := d.appendInt(v); err != nil {
			return WithStack(err)
		}
		return nil
	case UInt:
		if v, err := s.GetUInt(); err != nil {
			return WithStack(err)
		} else if err := d.appendUInt(v); err != nil {
			return WithStack(err)
		}
		return nil
	case String:
		if v, err := s.GetString(); err != nil {
			return WithStack(err)
		} else if err := d.appendString(v); err != nil {
			return WithStack(err)
		}
		return nil
	case Array:
		if err := d.appendArray(s); err != nil {
			return WithStack(err)
		}
		return nil
	case Object:
		if err := d.appendObject(s); err != nil {
			return WithStack(err)
		}
		return nil
	}

	return nil
}

func (d *Dumper) appendUInt(v uint64) error {
	s := strconv.FormatUint(v, 10)
	if _, err := d.w.Write([]byte(s)); err != nil {
		return WithStack(err)
	}
	return nil
}

func (d *Dumper) appendInt(v int64) error {
	s := strconv.FormatInt(v, 10)
	if _, err := d.w.Write([]byte(s)); err != nil {
		return WithStack(err)
	}
	return nil
}

func (d *Dumper) appendDouble(v float64) error {
	s := strconv.FormatFloat(v, 'E', -1, 64)
	if _, err := d.w.Write([]byte(s)); err != nil {
		return WithStack(err)
	}
	return nil
}

func (d *Dumper) appendString(v string) error {
	s := strconv.Quote(v)
	if _, err := d.w.Write([]byte(s)); err != nil {
		return WithStack(err)
	}
	return nil
}

func (d *Dumper) appendArray(v Slice) error {
	w := d.w
	it, err := NewArrayIterator(v)
	if err != nil {
		return WithStack(err)
	}
	if _, err := w.Write([]byte{'['}); err != nil {
		return WithStack(err)
	}
	for it.IsValid() {
		if !it.IsFirst() {
			if _, err := w.Write([]byte{','}); err != nil {
				return WithStack(err)
			}
		}
		if value, err := it.Value(); err != nil {
			return WithStack(err)
		} else if err := d.Append(value); err != nil {
			return WithStack(err)
		}
		if err := it.Next(); err != nil {
			return WithStack(err)
		}
	}
	if _, err := w.Write([]byte{']'}); err != nil {
		return WithStack(err)
	}
	return nil
}

func (d *Dumper) appendObject(v Slice) error {
	w := d.w
	it, err := NewObjectIterator(v)
	if err != nil {
		return WithStack(err)
	}
	if _, err := w.Write([]byte{'{'}); err != nil {
		return WithStack(err)
	}
	for it.IsValid() {
		if !it.IsFirst() {
			if _, err := w.Write([]byte{','}); err != nil {
				return WithStack(err)
			}
		}
		if key, err := it.Key(true); err != nil {
			return WithStack(err)
		} else if err := d.Append(key); err != nil {
			return WithStack(err)
		}
		if _, err := w.Write([]byte{':'}); err != nil {
			return WithStack(err)
		}
		if value, err := it.Value(); err != nil {
			return WithStack(err)
		} else if err := d.Append(value); err != nil {
			return WithStack(err)
		}
		if err := it.Next(); err != nil {
			return WithStack(err)
		}
	}
	if _, err := w.Write([]byte{'}'}); err != nil {
		return WithStack(err)
	}
	return nil
}
