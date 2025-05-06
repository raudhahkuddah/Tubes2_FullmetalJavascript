import { useState } from 'react';

export default function SearchForm({ onSearch, isLoading }) {
  const [element, setElement] = useState('');
  const [algorithm, setAlgorithm] = useState('BFS');
  const [multipleRecipes, setMultipleRecipes] = useState(false);
  const [maxRecipes, setMaxRecipes] = useState(5);

  const handleSubmit = (e) => {
    e.preventDefault();
    
    if (!element.trim()) {
      alert('Please enter an element name');
      return;
    }
    
    onSearch({
      element,
      algorithm,
      multipleRecipes,
      maxRecipes: multipleRecipes ? maxRecipes : 1
    });
  };

  return (
    <form onSubmit={handleSubmit} className="search-form">
      <h2>Little Alchemy 2 Recipe Finder</h2>
      
      <div className="form-group">
        <label htmlFor="element">Element to search:</label>
        <input
          id="element"
          type="text"
          value={element}
          onChange={(e) => setElement(e.target.value)}
          placeholder="Enter element name (e.g., 'glass', 'human')"
          required
        />
      </div>
      
      <div className="form-group">
        <label htmlFor="algorithm">Search Algorithm:</label>
        <select
          id="algorithm"
          value={algorithm}
          onChange={(e) => setAlgorithm(e.target.value)}
        >
          <option value="BFS">BFS (Breadth-First Search)</option>
          <option value="DFS">DFS (Depth-First Search)</option>
          <option value="Bidirectional">Bidirectional Search</option>
        </select>
      </div>
      
      <div className="form-group checkbox-group">
        <input
          id="multipleRecipes"
          type="checkbox"
          checked={multipleRecipes}
          onChange={(e) => setMultipleRecipes(e.target.checked)}
        />
        <label htmlFor="multipleRecipes">Find multiple recipes</label>
      </div>
      
      {multipleRecipes && (
        <div className="form-group">
          <label htmlFor="maxRecipes">Maximum recipes to find:</label>
          <input
            id="maxRecipes"
            type="number"
            min="1"
            max="50"
            value={maxRecipes}
            onChange={(e) => setMaxRecipes(parseInt(e.target.value))}
          />
        </div>
      )}
      
      <button type="submit" className="submit-btn" disabled={isLoading}>
        {isLoading ? 'Searching...' : 'Search Recipes'}
      </button>
    </form>
  );
}