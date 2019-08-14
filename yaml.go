package n

// // FromYaml return a numerable from the given M
// func FromYaml(yml string) (result *OldNumerable, err error) {
// 	data := map[string]interface{}{}
// 	if err = yaml.Unmarshal([]byte(yml), &data); err == nil {
// 		result = Q(data)
// 	}
// 	return
// }

// // FromYamlTmplFile loads a yaml file from disk and processes any templating
// // provided by the tmpl package returning an unmarshaled yaml block numerable.
// func FromYamlTmplFile(filepath string, vars map[string]string) *OldNumerable {
// 	if data, err := ioutil.ReadFile(filepath); err == nil {
// 		if tpl, err := tmpl.New(string(data), "{{", "}}"); err == nil {
// 			if result, err := tpl.Process(vars); err == nil {
// 				m := map[string]interface{}{}
// 				if err := yaml.Unmarshal([]byte(result), &m); err == nil {
// 					return Q(m)
// 				}
// 			}
// 		}
// 	}
// 	return Nil()
// }

// // M gets data by key which can be dot delimited
// // returns nil numerable on errors or keys not found
// func (q *Numerable) M(key string) (result *Numerable) {
// 	keys := A(key).Split(".")
// 	if key, ok := keys.TakeFirst(); ok {
// 		switch x := q.v.Interface().(type) {
// 		case map[string]interface{}:
// 			if !A(key).ContainsAny(":", "[", "]") {
// 				if v, ok := x[key]; ok {
// 					result = Q(v)
// 				}
// 			}
// 		case []interface{}:
// 			k, v := A(key).TrimPrefix("[").TrimSuffix("]").Split(":").YamlPair()
// 			if v == nil {
// 				if i, err := strconv.Atoi(k); err == nil {
// 					result = q.At(i)
// 				} else {
// 					panic(errors.New("Failed to convert index to an int"))
// 				}
// 			} else {
// 				for i := range x {
// 					if m, ok := x[i].(map[string]interface{}); ok {
// 						if entry, ok := m[k]; ok {
// 							if v == entry {
// 								result = Q(m)
// 								break
// 							}
// 						}
// 					}
// 				}
// 			}
// 		}
// 		if keys.Len() != 0 && result != nil && result.Any() {
// 			result = result.M(keys.Join(".").A())
// 		}
// 	}
// 	if result == nil {
// 		result = Nil()
// 	}
// 	return
// }

// // YamlMerge merges the given yaml into the existing yaml
// // with the given yaml having a higher precedence than the existing
// func (q *Numerable) YamlMerge(data interface{}) (result *Numerable, err error) {
// 	result = q
// 	switch x := q.v.Interface().(type) {
// 	case map[string]interface{}:
// 		switch y := data.(type) {
// 		case map[string]interface{}:
// 			q.Copy(MergeMap(x, y))
// 		default:
// 			err = errors.Errorf("Invalid merge type")
// 		}
// 	case []interface{}:
// 		switch y := data.(type) {
// 		case []interface{}:
// 			for i := range y {
// 				if i < len(x) {
// 					x[i] = y[i]
// 				} else {
// 					x = append(x, y[i])
// 					q.Copy(x)
// 				}
// 			}
// 		default:
// 			err = errors.Errorf("Invalid merge type")
// 		}
// 	}
// 	return
// }

// // YamlSet sets data by key which can be dot delimited
// func (q *Numerable) YamlSet(key string, data interface{}) (result *Numerable, err error) {
// 	keys := A(key).Split(".")
// 	if key, ok := keys.TakeFirst(); ok {
// 		switch x := q.v.Interface().(type) {
// 		case map[string]interface{}:

// 			// Current target is a map key
// 			if !A(key).ContainsAny(":", "[", "]") {

// 				// No more keys so we've reached our destination
// 				if !keys.Any() {
// 					x[key] = data
// 				} else {
// 					var v interface{}
// 					if v, ok = x[key]; !ok {
// 						// Doesn't exist so create
// 						if keys.First().Contains("[") {
// 							x[key] = []interface{}{}
// 						} else {
// 							x[key] = map[string]interface{}{}
// 						}
// 						v = x[key]
// 					}
// 					if result, err = Q(v).YamlSet(keys.Join(".").A(), data); err == nil {
// 						x[key] = result.O()
// 					}
// 				}
// 			}
// 		case []interface{}:
// 			insert := false
// 			var k string
// 			var v interface{}
// 			if A(key).Contains("[[") {
// 				insert = true
// 				k, v = A(key).TrimPrefix("[[").TrimSuffix("]]").Split(":").YamlPair()
// 			} else {
// 				k, v = A(key).TrimPrefix("[").TrimSuffix("]").Split(":").YamlPair()
// 			}
// 			if v == nil {
// 				var i int
// 				if i, err = strconv.Atoi(k); err == nil {

// 					// No more keys so we've reached our destination
// 					if !keys.Any() {
// 						if i < q.Len() && !insert {
// 							q.Set(i, data)
// 						} else {
// 							q.Insert(i, data)
// 						}
// 					} else {
// 						if i >= q.Len() {
// 							if keys.First().Contains("[") {
// 								q.Append([]interface{}{})
// 							} else {
// 								q.Append(map[string]interface{}{})
// 							}
// 						}
// 						result, err = q.At(i).YamlSet(keys.Join(".").A(), data)
// 					}
// 				} else {
// 					return
// 				}
// 			} else {
// 				for i := range x {
// 					if m, ok := x[i].(map[string]interface{}); ok {
// 						if entry, ok := m[k]; ok {
// 							if v == entry {
// 								if !keys.Any() {
// 									if insert {
// 										q.Insert(i, data)
// 									} else {
// 										q.Set(i, data)
// 									}
// 								} else {
// 									result, err = q.At(i).YamlSet(keys.Join(".").A(), data)
// 								}
// 								break
// 							}
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// 	result = q
// 	return
// }

// // Insert/set data in the unmarshalled yaml
// func (q *Numerable) yamlSet() (err error) {
// 	return
// }
