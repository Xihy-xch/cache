package local_cache

type Node struct {
	pre  *Node
	next *Node
	key  string
	val  item
}

func NewNode(key string, val item) *Node {
	return &Node{
		pre:  nil,
		next: nil,
		key:  key,
		val:  val,
	}
}

type NodeList struct {
	head *Node //最近使用
	end  *Node //最久未使用
}

func NewNodeList() *NodeList {
	head := NewNode("head", item{})
	end := NewNode("end", item{})

	head.next = end
	end.pre = head
	return &NodeList{
		head: head,
		end:  end,
	}
}

func (n *NodeList) front() *Node {
	if n.isEmpty() {
		return nil
	}

	return n.head.next
}

func (n *NodeList) back() *Node {
	if n.isEmpty() {
		return nil
	}

	return n.end.pre
}

func (n *NodeList) isEmpty() bool {
	return n.head.next == n.end
}

func (n *NodeList) pushFront(node *Node) {
	node.next = n.head.next
	node.next.pre = node
	n.head.next = node
	node.pre = n.head
}

func (n *NodeList) popBack() {
	node := n.end.pre
	if node == n.head {
		return
	}

	node.pre.next = n.end
	n.end.pre = node.pre
}

func (n *NodeList) moveToFront(node *Node) {
	node.next.pre = node.pre
	node.pre.next = node.next

	node.pre = n.head
	node.next = n.head.next
	node.next.pre = node
	n.head.next = node
}

func (n *NodeList) moveToBack(node *Node) {
	node.next.pre = node.pre
	node.pre.next = node.next

	node.next = n.end
	node.pre = n.end.pre
	node.pre.next = node
	n.end.pre = node
}

func (n *NodeList) delete(node *Node) {
	node.next.pre = node.pre
	node.pre.next = node.next
}
