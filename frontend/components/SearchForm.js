import { useEffect, useState } from 'react';

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
      maxRecipes
    });
  };

  useEffect(() => {
    if (!multipleRecipes) setMaxRecipes(1);
    else setMaxRecipes(2);
  }, [multipleRecipes]);

  return (
    <form onSubmit={handleSubmit} className="search-form">
      <h2 style={{ color: '#355070' }}>Little Alchemy 2 Recipe Finder</h2>
      
      <div className="form-group">
        <label htmlFor="element"style={{ color: '#6d597a' }}>Element to search:</label>
        <input
          id="element"
          type="text"
          value={element}
          onChange={(e) => setElement(e.target.value)}
          placeholder="Enter element name (e.g., 'glass', 'human')"
          required
           style={{
            backgroundColor: '#fdf0ed',
            border: '1px solid #b56576',
            borderRadius: '4px',
            padding: '0.5rem',
            color: '#355070',
            width: '100%',
          }}
        />
      </div>
      
      <div className="form-group">
        <label htmlFor="algorithm"style={{ color: '#6d597a' }}>Search Algorithm:</label>
        <select
          id="algorithm"
          value={algorithm}
          onChange={(e) => setAlgorithm(e.target.value)}
          style={{
            backgroundColor: '#fdf0ed',
            border: '1px solid #b56576',
            borderRadius: '4px',
            padding: '0.5rem',
            color: '#355070',
            width: '100%',
          }}
        >
          <option value="BFS">BFS (Breadth-First Search)</option>
          <option value="DFS">DFS (Depth-First Search)</option>
        </select>
      </div>
      
      <div className="form-group checkbox-group"style={{ color: '#6d597a' }}>
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
          <label htmlFor="maxRecipes"style={{ color: '#6d597a' }}>Maximum recipes to find:</label>
          <input
            id="maxRecipes"
            type="number"
            min="1"
            max="50"
            value={maxRecipes}
            onChange={(e) => setMaxRecipes(parseInt(e.target.value))}
            style={{
              backgroundColor: '#fdf0ed',
              border: '1px solid #b56576',
              borderRadius: '4px',
              padding: '0.5rem',
              color: '#355070',
              width: '100%',
            }}
          />
        </div>
      )}
      
      <button type="submit" className="submit-btn" disabled={isLoading}style={{backgroundColor: '#b56576', color: 'white', padding: '0.5rem 1rem', border: 'none', borderRadius: '4px',marginTop: '1rem', cursor: 'pointer', opacity: isLoading ? 0.6 : 1,}}>
        {isLoading ? 'Searching...' : 'Search Recipes'}
      </button>
    </form>
  );
}