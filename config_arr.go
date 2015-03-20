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
  "strconv"
  "strings"

  "github.com/demdxx/gocast"
)

type ConfigArr []interface{}

func FromSlice(c interface{}) (ConfigArr, error) {
  if nil == c {
    return make(ConfigArr, 0), nil
  }

  conf := make(ConfigArr, 0)
  for _, v := range gocast.ToInterfaceSlice(c) {
    t := reflect.TypeOf(v)
    switch t.Kind() {
    case reflect.Map:
      nc, _ := From(v)
      conf = append(conf, nc)
      break
    case reflect.Slice:
      nc, _ := FromSlice(v)
      conf = append(conf, nc)
      break
    default:
      conf = append(conf, v)
      break
    }
  }
  return conf, nil
}

///////////////////////////////////////////////////////////////////////////////
/// Getters/Setters
///////////////////////////////////////////////////////////////////////////////

func (conf ConfigArr) Get(path string) (interface{}, error) {
  return conf.GetPath(strings.Split(path, "."))
}

func (conf ConfigArr) GetPath(path []string) (interface{}, error) {
  if len(path) < 1 {
    return nil, ErrInvalidPath
  }

  key := path[0]
  path = path[1:]

  if "$" == key || "+" == key || "*" == key {
    if len(path) < 1 {
      return conf, nil
    } else {
      // For each item
      result := make([]interface{}, 0)
      for _, it := range conf {
        switch a := it.(type) {
        case Config:
          it, _ := a.GetPath(path)
          result = append(result, it)
          break
        case ConfigArr:
          it, _ := a.GetPath(path)
          result = append(result, it)
          break
        default:
          return a, ErrNoValid
          break
        }
      }
      return result, nil
    }
  } else if isDigit(key) { // If digit index
    index, _ := strconv.Atoi(key)
    if index < len(conf) {
      it := conf[index]
      if len(path) == 1 {
        return it, nil
      }

      switch a := it.(type) {
      case Config:
        return a.GetPath(path)
        break
      case ConfigArr:
        return a.GetPath(path)
        break
      default:
        return a, ErrNoValid
        break
      }
    }
    return conf, ErrNoValid
  }
  return nil, ErrInvalidPath
}

/// Set

func (conf ConfigArr) Set(path string, value interface{}) ConfigArr {
  return conf.SetPath(strings.Split(path, "."), value)
}

func (conf ConfigArr) SetPath(path []string, value interface{}) ConfigArr {
  if len(path) < 1 {
    return nil
  }

  key := path[0]
  path = path[1:]

  if "+" == key {
    if len(path) == 0 {
      conf = append(conf, prepareValueForSet(value))
    } else {
      if isArrayChain(path[0]) {
        conf = append(conf, make(ConfigArr, 1).SetPath(path, value))
      } else if "*" == path[0] {
        // I don't know what me do!
      } else {
        conf = append(conf, make(Config).SetPath(path, value))
      }
    }
  } else if "$" == key || "*" == key {
    for i, it := range conf {
      if len(path) == 0 {
        conf[i] = prepareValueForSet(value)
      } else {
        switch a := it.(type) {
        case Config:
          if isArrayChain(path[0]) {
            conf[i] = append(conf, make(ConfigArr, 1).SetPath(path, value))
          } else {
            a.SetPath(path, value)
          }
          break
        case ConfigArr:
          if isArrayChain(path[0]) {
            a.SetPath(path, value)
          } else {
            conf[i] = append(conf, make(Config).SetPath(path, value)) // Replace value
          }
          break
        default:
          if isArrayChain(path[0]) {
            conf[i] = append(conf, make(ConfigArr, 1).SetPath(path, value))
          } else if "*" == path[0] {
            // I don't know what me do!
          } else {
            conf[i] = append(conf, make(Config).SetPath(path, value))
          }
          break
        }
      }
    }
  } else if isDigit(key) { // If digit index
    index, _ := strconv.Atoi(key)
    if len(conf) > index {
      conf[index] = prepareValueForSet(value)
    }
  }
  return conf
}

///////////////////////////////////////////////////////////////////////////////
/// Convertion
///////////////////////////////////////////////////////////////////////////////
