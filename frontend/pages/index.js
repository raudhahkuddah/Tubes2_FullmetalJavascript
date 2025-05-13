import { useState, useEffect } from 'react';
import Head from 'next/head';
import SearchForm from '../components/SearchForm';
import RecipeVisualizer from '../components/RecipeVisualizer';
import RecipeTree from '../components/RecipeTree';
import { convertRecipeTreeToBinary } from '../components/RecipeTree';
import { useEdges } from 'reactflow';

const dummy = {
  name: "Steam",
  children: [
    {
      name: "Water + Fire",
      children: [
        {
          name: "Water",
          children: []
        },
        {
          name: "Fire",
          children: []
        }
      ]
    }
  ]
};


export default function Home() {
  const [searchResults, setSearchResults] = useState(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);
  const [tree, setTree] = useState({});

  const handleSearch = async (searchParams) => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await fetch('http://localhost:8080/search', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(searchParams),
      });
      
      if (!response.ok) {
        throw new Error(`Error: ${response.statusText}`);
      }
      
      const data = await response.json();
      setSearchResults(data);
    } catch (err) {
      setError(err.message);
      console.error('Search failed:', err);
    } finally {
      setIsLoading(false);
    }
  };

  const updateTree = () => {
    try {
      const updatedTree = convertRecipeTreeToBinary(searchResults.tree);
      console.log("Updated Tree:", updatedTree);
      setTree(updatedTree);    
    }
    catch{
      console.log("data is not defined");
    }
  }

  useEffect(() => {
    updateTree();
  }, [searchResults]);

  return (
    <div className="container" style={{backgroundColor: '#f4cbb1', color: '#355070', minHeight: '100vh', margin: 0, padding: 20, width: '100%', maxWidth: '100vw',overflowX: 'hidden', boxSizing:'border-box'}} >
      <Head>
        <title>Little Alchemy 2 Recipe Search</title>
        <meta name="description" content="Search for recipes in Little Alchemy 2" />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <main>
        <SearchForm onSearch={handleSearch} isLoading={isLoading} />
        
        {error && (
          <div className="error-message" style={{ color: '#e56b6f' }}>
            <p>{error}</p>
          </div>
        )}
        
        {isLoading && (
          <div className="loading" style={{ color: '#6d597a' }}>
            <p>Searching for recipes...</p>
          </div>
        )}
        
        {Object.keys(tree).length > 0 && !isLoading &&
        <div style={{ padding: 30}}>
          <RecipeTree key={JSON.stringify(tree)} treeData={tree} />
          <h1 style={{ fontSize: '0.875rem', color: '#355070', fontWeight: 'normal',  marginTop: '1rem'}}>Algorithm = {searchResults.algorithm}, Nodes = {searchResults.visited_nodes}, Duration = {searchResults.duration_ms}ms</h1>
          {/* <button onClick={updateTree}>Change recipe</button> */}
        </div>
        }
      </main>
    </div>
  );
}