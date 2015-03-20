# Config for GO projects

Simple universal config xml, json, yaml

## Config example

```yaml
app:
  name: App1
  version: 1.0.3beta
  description: "Description text"

  params:
    p1: v1
    p2: v2

  arr:
    - value1
    - value2
    - value3
```

## Get

```go
conf, _ := config.From(map)

v, err := conf.Get("app.*.p1")
fmt.Println(v) // v1
```

## Set

```go
conf, _ := config.From(map[string]interface{}{})

conf.Set("array.+", "Item 1")
conf.Set("array.+", "Item 2")
conf.Set("array.+", "Item 3")

fmt.Println(conf.JsonString())
```

```json
{"array":["Item 1", "Item 3", "Item 2"]}
```

```go
conf.Set("array.$", "New Value")
fmt.Println(conf.JsonString())
```

```json
{"array":["New Value", "New Value", "New Value"]}
```

```go
conf.Set("array.$.map", "New Value")
fmt.Println(conf.JsonString())
```

```json
{"array":["New Value", "New Value", "New Value", {"map": "Nev Value"}]}
```

# License

    The MIT License (MIT)

    Copyright (c) 2015 Dmiptry Ponomarev <demdxx@gmail.com>

    Permission is hereby granted, free of charge, to any person obtaining a copy
    of this software and associated documentation files (the "Software"), to deal
    in the Software without restriction, including without limitation the rights
    to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
    copies of the Software, and to permit persons to whom the Software is
    furnished to do so, subject to the following conditions:

    The above copyright notice and this permission notice shall be included in all
    copies or substantial portions of the Software.

    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
    AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
    OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
    SOFTWARE.

