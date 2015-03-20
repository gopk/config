//
// The MIT License (MIT)
//
// Copyright (c) 2015 Dmiptry Ponomarev <demdxx@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//

package config

import (
  "reflect"
  "unicode"
)

func prepareValueForSet(value interface{}) interface{} {
  if nil == value {
    return nil
  }

  t := reflect.TypeOf(value)
  switch t.Kind() {
  case reflect.Map:
    value, _ = From(value)
    break
  case reflect.Slice:
    value, _ = FromSlice(value)
    break
  }
  return value
}

func isArrayChain(ch string) bool {
  return "*" == ch || "+" == ch || isDigit(ch)
}

func isDigit(s string) bool {
  if len(s) < 1 {
    return false
  }
  for _, c := range s {
    if !unicode.IsDigit(c) {
      return false
    }
  }
  return true
}