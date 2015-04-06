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
  "encoding/json"
  "encoding/xml"
  "fmt"
  "gopkg.in/v1/yaml"
  "io/ioutil"
  "reflect"
  "regexp"
  "strings"

  "github.com/demdxx/gocast"
)

var (
  escapeEx = regexp.MustCompile(`\{\}\*\+\(\)\[\]`)
)

type Config map[string]interface{}

func From(c interface{}) (Config, error) {
  return FromQuick(c)
}

func FromQuick(c interface{}) (Config, error) {
  if nil == c {
    return make(Config), nil
  }

  sm, err := gocast.ToSiMap(c, "field")
  if nil != err {
    return nil, err
  }

  conf := make(Config)
  for k, v := range sm {
    if nil == v {
      conf[k] = v
      continue
    }

    t := reflect.TypeOf(v)
    switch t.Kind() {
    case reflect.Map:
      conf[k], _ = FromQuick(v)
      break
    case reflect.Slice:
      conf[k], _ = FromSliceQuick(v)
      break
    default:
      conf[k] = v
      break
    }
  }
  return conf, nil
}

func New() Config {
  return make(Config)
}

func FromFile(filename, ftype string) (Config, error) {
  bytes, err := ioutil.ReadFile(filename)
  if nil != err {
    return nil, err
  }
  return FromData(bytes, ftype)
}

func FromData(data []byte, dtype string) (conf Config, err error) {
  var info interface{}
  switch strings.ToLower(dtype) {
  case "json":
    err = json.Unmarshal(data, &info)
    break
  case "xml":
    err = xml.Unmarshal(data, &info)
    break
  case "yaml":
    err = yaml.Unmarshal(data, &info)
    break
  default:
    err = ErrInvalidConfigFormat
  }
  if nil == err {
    conf, err = From(info)
  }
  return
}

///////////////////////////////////////////////////////////////////////////////
/// Getters/Setters
///////////////////////////////////////////////////////////////////////////////

func (conf Config) Get(path string) (interface{}, error) {
  return conf.GetPath(strings.Split(path, "."))
}

func (conf Config) GetPath(path []string) (interface{}, error) {
  if len(path) < 1 {
    return nil, ErrInvalidPath
  }

  curConf := conf
  for i, key := range path {
    if len(key) < 1 {
      return nil, ErrInvalidPath
    }

    isLast := i >= len(path)-1

    if "*" == key { // Response as array
      if isLast {
        return curConf, nil
      }
      response := make([]interface{}, 0)

      for _, v := range curConf {
        switch a := v.(type) {
        case ConfigArr:
          if r, err := a.GetPath(path[i:]); nil == err {
            response = append(response, r)
          }
          break
        case Config:
          if r, err := a.GetPath(path[i:]); nil == err {
            response = append(response, r)
          }
          break
        }
      }
      return response, nil
    } else {
      if it, ok := curConf[key]; !ok || nil == it {
        return nil, ErrNoValue
      } else {
        if isLast {
          return it, nil
        }

        switch a := it.(type) {
        case ConfigArr:
          return a.GetPath(path[i:])
          break
        case Config:
          curConf = a
          break
        default:
          return it, ErrNoValid
        }
      }
    }
  }
  return nil, nil
}

func (conf Config) GetDefault(path string, def interface{}) interface{} {
  val, err := conf.Get(path)
  if nil == val || nil != err {
    return def
  }
  return val
}
func (conf Config) String(path string) string {
  return conf.StringOrDefault(path, "")
}

func (conf Config) StringOrDefault(path string, def string) string {
  return gocast.ToString(conf.GetDefault(path, def))
}

func (conf Config) IntOrDefault(path string, def int) int {
  return gocast.ToInt(conf.GetDefault(path, def))
}

func (conf Config) Float64OrDefault(path string, def float64) float64 {
  return gocast.ToFloat64(conf.GetDefault(path, def))
}

func (conf Config) BoolOrDefault(path string, def bool) bool {
  return gocast.ToBool(conf.GetDefault(path, def))
}

/// Set

func (conf Config) Set(path string, value interface{}) Config {
  return conf.SetPath(strings.Split(path, "."), value)
}

func (conf Config) SetPath(fullpath []string, value interface{}) Config {
  if len(fullpath) < 1 {
    return conf
  }

  key := fullpath[len(fullpath)-1]
  path := fullpath[:len(fullpath)-1]
  curConf := conf

  if len(path) > 0 {
    if isArrayChain(fullpath[0]) {
      return conf // Invalid path
    }

    for i, key := range path {
      isArray := false
      if i < len(path)-1 {
        isArray = isArrayChain(path[i+1])
      }

      if it, ok := curConf[key]; !ok || nil == it { // We have path key
        if isArray {
          newConf := make(ConfigArr, 0)
          curConf[key] = newConf.SetPath(fullpath[i:], value)
          return conf
        } else {
          newConf := make(Config)
          curConf[key] = newConf
          curConf = newConf
        }
        continue
      } else { // It`s something eles
        switch a := it.(type) {
        case ConfigArr:
          if isArray {
            a.SetPath(fullpath[i:], value)
            return conf
          } else {
            newConf := make(Config)
            curConf[key] = newConf
            curConf = newConf
          }
          break
        case Config:
          if isArray {
            newConf := make(ConfigArr, 0)
            curConf[key] = newConf.SetPath(fullpath[i:], value)
            return conf
          } else {
            curConf = a
          }
          break
        default: // Replace any other value
          newConf := make(Config)
          curConf[key] = newConf
          curConf = newConf
        }
      }
    }
  }

  if isArrayChain(key) {
    // Invalid! I cant set as for array
  } else {
    curConf[key] = prepareValueForSet(value)
  }
  return conf
}

func (conf Config) Update(conf2 Config) Config {
  if nil != conf2 {
    for k, v := range conf2 {
      conf[k] = v
    }
  }
  return conf
}

func (conf Config) UpdateByPath(path string) Config {
  conf2, _ := conf.Get(path)
  if nil != conf2 {
    switch conf2.(type) {
    case Config:
      return conf.Update(conf2.(Config))
      break
    }
  }
  return conf
}

///////////////////////////////////////////////////////////////////////////////
/// Processing
///////////////////////////////////////////////////////////////////////////////

func (conf Config) Prepare(escLeft, escRight string) {
  escLeft = escapeEx.ReplaceAllString(escLeft, `\$1`)
  escRight = escapeEx.ReplaceAllString(escRight, `\$1`)
  r := regexp.MustCompile(fmt.Sprintf(`%s([\w\.\d+*$]+)%s`, escLeft, escRight))
  prepare(conf, conf, r)
}

///////////////////////////////////////////////////////////////////////////////
/// Helpers
///////////////////////////////////////////////////////////////////////////////

func prepare(root, conf Config, r *regexp.Regexp) {
  if nil == conf || len(conf) < 1 {
    return
  }

  for k, v := range conf {
    if V := prepareItem(root, v, r); nil != V {
      conf[k] = V
    }
  }
}

func prepareArray(root Config, a ConfigArr, r *regexp.Regexp) {
  if nil == a || len(a) < 1 {
    return
  }
  for i, it := range a {
    if V := prepareItem(root, it, r); nil != V {
      a[i] = V
    }
  }
}

func prepareItem(root Config, it interface{}, r *regexp.Regexp) interface{} {
  switch s := it.(type) {
  case string:
    return r.ReplaceAllStringFunc(s, func(s string) string {
      return fmt.Sprintf("%v", root.GetDefault(s[2:len(s)-2], ""))
    })
    break
  case Config:
    prepare(root, s, r)
    break
  case ConfigArr:
    prepareArray(root, s, r)
    break
  }
  return nil
}
