# Little Alchemy 2 Recipe Search

This project is a web application for searching recipes of elements in the game Little Alchemy 2. It utilizes both BFS and DFS strategies for searching recipes and is built with a frontend in Next.js and a backend in Golang.

## Project Structure

The project is organized into two main directories: `frontend` and `backend`.

### Frontend

- **components/**: Contains React components for the application.
  - `SearchForm.js`: A form for users to input elements and select search algorithms.
  - `RecipeVisualizer.js`: Visualizes the found recipes in a tree structure.

- **pages/**: Contains the pages of the application.
  - `_app.js`: Custom App component for Next.js.
  - `index.js`: Main page that integrates the search form and recipe visualizer.

- **styles/**: Contains global CSS styles for the application.
  - `global.css`: Styles for body, headings, buttons, and tree visualization.

- **package.json**: Configuration file for npm, listing dependencies and scripts.

- **next.config.js**: Configuration settings for the Next.js application.

### Backend

- **cmd/server/**: Contains the entry point for the backend application.
  - `main.go`: Sets up the HTTP server and routes for API requests.

- **pkg/scraper/**: Contains the implementation for scraping data.
  - `scraper.go`: Scrapes element data and recipes from the Fandom Little Alchemy 2 website.

- **pkg/search/**: Implements search algorithms.
  - `bfs.go`: BFS algorithm for searching recipes.
  - `dfs.go`: DFS algorithm for searching recipes.
  - `bidirectional.go`: Bidirectional search algorithm.

- **pkg/api/**: Contains HTTP handlers for the API.
  - `handlers.go`: Endpoints for searching recipes and returning results.

- **go.mod**: Go module file defining dependencies and versions.

- **go.sum**: Checksums for module dependencies.

### Setup Instructions

1. **Clone the repository**:
   ```
   git clone <repository-url>
   cd little-alchemy-2-recipe-search
   ```

2. **Frontend Setup**:
   - Navigate to the `frontend` directory:
     ```
     cd frontend
     ```
   - Install dependencies:
     ```
     npm install
     ```
   - Start the development server:
     ```
     npm run dev
     ```

3. **Backend Setup**:
   - Navigate to the `backend` directory:
     ```
     cd ../backend
     ```
   - Install Go dependencies:
     ```
     go mod tidy
     ```
   - Run the server:
     ```
     go run cmd/server/main.go
     ```

### Features

- Search for recipes using BFS, DFS, or Bidirectional algorithms.
- Option to find multiple recipes with a maximum limit.
- Visualize recipes as a tree structure.
- Display search time and number of nodes visited.

### Acknowledgments

This project utilizes the Goquery library for web scraping and the Treebeard library for visualizing the recipe tree.