// TreeFlow.jsx
import React, { useEffect } from "react";
import ReactFlow, {
  Background,
  Controls,
  useNodesState,
  useEdgesState,
} from "reactflow";
import "reactflow/dist/style.css";

let idCounter = 0;
const getId = () => `node-${idCounter++}`;

function getRandomIndex(array) {
  if (!Array.isArray(array) || array.length === 0) {
    throw new Error("Input must be a non-empty array.");
  }
  return Math.floor(Math.random() * array.length);
}


export function convertRecipeTreeToBinary(node) {
  if (!node) return null;

  let left = null;
  let right = null;

  if (node.children[0] && node.children[0].children[0]) {
    const recipe = node.children[0].children;
    left = convertRecipeTreeToBinary(recipe[0]);
    right = convertRecipeTreeToBinary(recipe[1]);
  }

  return {
    name: node.name,
    left,
    right
  };
}


function buildNodeTree(tree, depth = 0) {
  const nodeId = getId();
  const node = {
    id: nodeId,
    label: tree.name,
    children: [],
    depth,
    x: 0,
    y: depth,
  };

  if (tree.left) {
    const left = buildNodeTree(tree.left, depth + 1);
    node.children.push(left);
  }

  if (tree.right) {
    const right = buildNodeTree(tree.right, depth + 1);
    node.children.push(right);
  }

  return node;
}

function assignX(node, x = { value: 0 }) {
  node.children.forEach((child) => assignX(child, x));

  if (node.children.length === 0) {
    node.x = x.value++;
  } else {
    const mid = (node.children[0].x + node.children[node.children.length - 1].x) / 2;
    node.x = mid;
  }
}

function flatten(node, nodes = [], edges = []) {
  nodes.push({
    id: node.id,
    data: { label: node.label },
    position: { x: node.x * 200, y: node.y * 120 },
    draggable: false,
    selectable: false,
    focusable: false,
    deletable: false,
    style: {
    pointerEvents: "none", // prevent clicking/dragging
  },
  });

  node.children.forEach((child) => {
    edges.push({
      id: `${node.id}-${child.id}`,
      source: node.id,
      target: child.id,
    });
    flatten(child, nodes, edges);
  });

  return { nodes, edges };
}

const RecipeTree = ({ treeData }) => {
  const [nodes, setNodes, onNodesChange] = useNodesState([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);

  useEffect(() => {
    idCounter = 0;

    const root = buildNodeTree(treeData);

    assignX(root);

    const { nodes: treeNodes, edges: treeEdges } = flatten(root);

    setNodes(treeNodes);
    setEdges(treeEdges);
  }, [treeData]);

  return (
    <div style={{ width: "100%", height: "100vh", border: "4px solid black", borderRadius: "8px" }}>
      <ReactFlow
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        fitView
      >
        <Background />
        <Controls />
      </ReactFlow>
    </div>
  );
};

export default RecipeTree;
