import { useState } from 'react';
import Head from 'next/head';
import SearchForm from '../components/SearchForm';
import RecipeVisualizer from '../components/RecipeVisualizer';

export default function Home() {
  const [searchResults, setSearchResults] = useState(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);

  const handleSearch = async (searchParams) => {
    setIsLoading(true);
    setError(null);
    
    try {
      const response = await fetch('http://localhost:8080/api/search', {
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

  return (
    <div className="container">
      <Head>
        <title>Little Alchemy 2 Recipe Search</title>
        <meta name="description" content="Search for recipes in Little Alchemy 2" />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <main>
        <SearchForm onSearch={handleSearch} isLoading={isLoading} />
        
        {error && (
          <div className="error-message">
            <p>{error}</p>
          </div>
        )}
        
        {isLoading && (
          <div className="loading">
            <p>Searching for recipes...</p>
          </div>
        )}
        
        {searchResults && !isLoading && (
          <RecipeVisualizer 
            recipes={searchResults.recipes}
            nodesVisited={searchResults.nodesVisited}
            searchTime={searchResults.searchTimeMs}
          />
        )}
      </main>
    </div>
  );
}