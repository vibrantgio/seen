package document

type Dom struct {
	ChildNodes map[string]*Element
}

func NewDom() *Dom {
	return &Dom{ChildNodes: make(map[string]*Element)}
}

func (d *Dom) Reset() {
	for k := range d.ChildNodes {
		delete(d.ChildNodes, k)
	}
}

func (d *Dom) GetElementById(id string) *Element {
	return d.ChildNodes[id]
}

func (d *Dom) SetElementById(id string, e *Element) {
	d.ChildNodes[id] = e
}

func (d *Dom) CreateElementNS(namespace, tag string) *Element {
	return &Element{Dom: d, Namespace: namespace, Tag: tag}
}
