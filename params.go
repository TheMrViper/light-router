package router

type param struct {
	Key   string
	Value string
}

type Params []param

func (p Params) AsMap() map[string]string {
	ret := make(map[string]string)
	for k := range p {
		ret[p[k].Key] = p[k].Value
	}
	return ret
}

func (p Params) Get(key string) (string, bool) {
	for k := range p {
		if p[k].Key == key {
			return p[k].Value, true
		}
	}
	return "", false
}
