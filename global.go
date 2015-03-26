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

var (
  cache Config
)

func SetGlobalConfig(code string, conf Config) {
  if nil == cache {
    cache = make(Config)
  }
  if nil != conf {
    cache[code] = conf
  }
}

func HasGlobalConfig(code string) bool {
  if _, ok := cache[code]; ok {
    return true
  }
  return false
}

func SetGlobalByFile(code, filename, ftype string) (Config, error) {
  conf, err := FromFile(filename, ftype)
  if nil == err {
    SetGlobalConfig(code, conf)
  }
  return conf, err
}

func GlobalByName(code string) Config {
  if nil == cache {
    cache = make(Config)
  }
  if c, ok := cache[code]; !ok {
    c = make(Config)
    cache[code] = c.(Config)
    return c.(Config)
  } else {
    return c.(Config)
  }
}

func GlobalByFile(filename, ftype string) (Config, error) {
  return SetGlobalByFile("default", filename, ftype)
}

func Global() Config {
  return GlobalByName("default")
}

///////////////////////////////////////////////////////////////////////////////
/// Getters/Setters
///////////////////////////////////////////////////////////////////////////////

func Get(path string) (interface{}, error) {
  return Global().Get(path)
}

func GetPath(path []string) (interface{}, error) {
  return Global().GetPath(path)
}

func GetDefault(path string, def interface{}) interface{} {
  return Global().GetDefault(path, def)
}

func String(path string) string {
  return Global().String(path)
}

func StringOrDefault(path string, def string) string {
  return Global().StringOrDefault(path, def)
}

func IntOrDefault(path string, def int) int {
  return Global().IntOrDefault(path, def)
}

func Float64OrDefault(path string, def float64) float64 {
  return Global().Float64OrDefault(path, def)
}

func BoolOrDefault(path string, def bool) bool {
  return Global().BoolOrDefault(path, def)
}

// Set

func Set(path string, value interface{}) Config {
  return Global().Set(path, value)
}

func SetPath(path []string, value interface{}) Config {
  return Global().SetPath(path, value)
}
