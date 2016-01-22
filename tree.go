package router

import "strings"

type tree struct {
	path    string
	handler iHandler

	name      string
	isParam   bool
	isCatche  bool
	haveParam bool

	nextIndex   int
	paramsIndex int
	paramsCount int

	parent *tree
	child  []*tree
}

func newTree() *tree {
	return &tree{
		child: make([]*tree, 0, 1),
	}
}

func (t *tree) getParamIndex() int {
	if t.parent != nil {
		if t.isParam || t.isCatche {
			return t.parent.getParamIndex() + 1
		}
		return t.parent.getParamIndex() + 0
	}
	return 0
}

func (t *tree) newTree(path string, paramsCount int) *tree {
	newTree := &tree{}
	newTree.child = make([]*tree, 0, 1)

	newTree.parent = t
	newTree.path = path

	if t.haveParam {
		for k := range t.child {
			if t.child[k].isCatche || t.child[k].isParam && t.child[k].paramsCount < paramsCount {
				t.child[k].paramsCount = paramsCount
			}
			paramsCount = t.child[k].paramsCount
		}
	}

	if len(path) > 1 {
		if path[1] == '$' {
			t.haveParam = true
			newTree.isParam = true
			newTree.name = path[2:]
			newTree.paramsIndex = t.getParamIndex()
			newTree.paramsCount = paramsCount
		}

		if path[1] == '@' {
			t.haveParam = true
			newTree.isCatche = true
			newTree.name = path[2:]
			newTree.paramsIndex = t.getParamIndex()
			newTree.paramsCount = paramsCount
		}

	}

	newTree.nextIndex = len(t.child) + 1
	t.child = append(t.child, newTree)

	return newTree
}

func buildPath(path string) (splited []string, count int) {
	count = 0
	splited = strings.Split(path, "/")[1:]

	for k := range splited {
		if len(splited[k]) > 0 && (splited[k][0] == '@' || splited[k][0] == '$') {
			count++
		}

		splited[k] = "/" + splited[k]
	}
	return
}

func (t *tree) AddRoute(path string, handler iHandler) {
	splited, paramsCount := buildPath(path)
	t.addRoute(splited, handler, paramsCount)
}

func (t *tree) addRoute(path []string, h iHandler, paramsCount int) {

	if len(path) == 0 {
		t.handler = h
		return
	}

	for k := range t.child {
		child := t.child[k]
		rel := child.isParam || child.isCatche

		if !rel && child.path == path[0] {

			child.addRoute(path[1:], h, paramsCount)
			return
		}

		if rel && child.paramsCount < paramsCount {
			child.paramsCount = paramsCount
		}

		if rel && child.path == path[0] {

			child.addRoute(path[1:], h, paramsCount)
			return
		}
	}

	child := t.newTree(path[0], paramsCount)
	child.addRoute(path[1:], h, paramsCount)
}

func (t *tree) FindRoute(path string) (handler iHandler, params Params) {
start:
	i := 0

	if t.isCatche || t.isParam && (len(t.child) == 0) && len(path) > 0 {
	parent:
		i = t.nextIndex
		path = "/" + params[t.paramsIndex].Value + path

		t = t.parent

		if i > len(t.child)-1 {
			goto parent
		}
	}
	for ; i < len(t.child); i++ {
		if len(path) == 0 {
			handler = t.handler
			return
		}

		if t.child[i].isParam {

			end := 1
			t = t.child[i]
			for end < len(path) && path[end] != '/' {
				end++
			}

			if params == nil {
				params = make(Params, t.paramsCount)
			}
			params[t.paramsIndex].Key = t.name
			params[t.paramsIndex].Value = path[1:end]

			path = path[end:]
			goto start
		}

		if t.child[i].isCatche {
			t = t.child[i]

			if params == nil {
				params = make(Params, t.paramsCount)
			}
			params[t.paramsIndex].Key = t.name
			params[t.paramsIndex].Value = path

			handler = t.handler
			return
		}

		if path == t.child[i].path {
			handler = t.child[i].handler
			return
		} else if len(path) > len(t.child[i].path) && len(t.child[i].path) > 1 {
			if path[:len(t.child[i].path)] == t.child[i].path {
				t = t.child[i]
				path = path[len(t.path):]

				goto start
			}
		}
	}
	if len(path) == 0 {
		handler = t.handler
	}
	return
}
