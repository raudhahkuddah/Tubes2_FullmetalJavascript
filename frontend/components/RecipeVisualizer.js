import { useMemo } from 'react';
import Tree from 'react-d3-tree';

export default function RecipeVisualizer({ recipes, nodesVisited, searchTime }) {
  // Convert recipe steps to tree structure for visualization
  const treeData = useMemo(() => {
    if (!recipes || recipes.length === 0) return null;
    
    // Process the first recipe (can add UI to switch between multiple recipes)
    const currentRecipe = recipes[0];
    
    // Function to build tree nodes from recipe steps
    const buildTreeFromSteps = (steps) => {
      if (!steps || steps.length <= 1) return null;
      
      const element = steps[0];
      steps = steps.slice(1); // Remove the first element (target)
      
      // For leaf nodes (basic elements)
      if (steps.length === 0) {
        return { name: element };
      }
      
      // Group steps into pairs for combinations
      const children = [];
      for (let i = 0; i < steps.length; i += 2) {
        if (i + 1 < steps.length) {
          children.push({
            name: `${steps[i]} + ${steps[i+1]}`,
            children: [
              { name: steps[i] },
              { name: steps[i+1] }
            ]
          });
        } else {
          children.push({ name: steps[i] });
        }
      }
      
      return {
        name: element,
        children
      };
    };
    
    return buildTreeFromSteps(currentRecipe.steps);
  }, [recipes]);

  // Custom node rendering
  const renderForeignObjectNode = ({
    nodeDatum,
    toggleNode,
    foreignObjectProps
  }) => (
    <g>
      <circle r={15} fill="#8bc34a" />
      <foreignObject {...foreignObjectProps}>
        <div className="tree-node">
          <h3>{nodeDatum.name}</h3>
          {nodeDatum.children && (
            <button onClick={toggleNode}>
              {nodeDatum.__rd3t.collapsed ? 'Expand' : 'Collapse'}
            </button>
          )}
        </div>
      </foreignObject>
    </g>
  );

  const nodeSize = { x: 200, y: 100 };
  const foreignObjectProps = {
    width: nodeSize.x,
    height: nodeSize.y,
    x: -100,
    y: -50
  };

  if (!treeData) {
    return <div className="recipe-visualizer empty">No recipe data available</div>;
  }

  return (
    <div className="recipe-visualizer">
      <div className="stats-panel">
        <h3>Search Results</h3>
        <div className="stats">
          <p><strong>Recipes found:</strong> {recipes.length}</p>
          <p><strong>Nodes visited:</strong> {nodesVisited}</p>
          <p><strong>Search time:</strong> {searchTime.toFixed(2)} ms</p>
        </div>
      </div>
      
      <div className="tree-container">
        <Tree
          data={treeData}
          orientation="vertical"
          pathFunc="step"
          nodeSize={nodeSize}
          separation={{ siblings: 2, nonSiblings: 2 }}
          renderCustomNodeElement={(rd3tProps) =>
            renderForeignObjectNode({ ...rd3tProps, foreignObjectProps })
          }
          translate={{ x: 300, y: 50 }}
        />
      </div>
    </div>
  );
}