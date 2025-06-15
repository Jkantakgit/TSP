import { useState, useRef } from 'react'
import "./App.css"

type City = {
  name: string;
  x: number; // Store as percentage of container width (0-100)
  y: number; // Store as percentage of container height (0-100)
};

// Main App component
function App() {

	// State to manage cities, route, error messages, and loading state
	const [cities, setCities] = useState<City[]>([]);
	const [route, setRoute] = useState<City[] | null>(null);
	const [error, setError] = useState<string | null>(null);
	const [loading, setLoading] = useState(false);
	const containerRef = useRef<HTMLDivElement>(null);

  	// Handle click event to add a new city
  	const handleClick = (e: React.MouseEvent<HTMLDivElement, MouseEvent>) => {
		const rect = e.currentTarget.getBoundingClientRect();
		// Convert to percentages (0-100)
		const x = ((e.clientX - rect.left) / rect.width) * 100;
		const y = ((e.clientY - rect.top) / rect.height) * 100;

		// Create a new city object
		const newCity: City = {
			name: `City ${cities.length + 1}`,
			x: x,
			y: y,
		};

		// Add the new city to the list
		setCities([...cities, newCity]);
		setRoute(null);
	};

  	// Function to call the backend API to solve the TSP problem
  	const solveTSP = async () => {
		setLoading(true);
		setError(null);
		setRoute(null);

		// Send the cities to the backend for optimization
		try {
			const res = await fetch("http://localhost:8080/optimize", {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({ cities }),
			});

			// Check if the response is ok, otherwise throw an error
			if (!res.ok) {
				throw new Error(await res.text());
			}
			const data = await res.json();
			setRoute(data.route);
		} catch (e: any) {
			setError(e.message || "Unknown error");
		} finally {
			setLoading(false);
		}
	};

  	// Render the component
  	return (
    	<div className="container">
        <h1>Travelling Salesman</h1>

        <div ref={containerRef} onClick={handleClick} className="tsp-wrapper">

        {/* Route lines */}
        {route && (
            <svg width="100%" height="100%">
            {route.slice(0, -1).map((city, i) => {
                const next = route[i + 1];
                return (
                <line
                    key={i}
                    x1={`${city.x}%`}
                    y1={`${city.y}%`}
                    x2={`${next.x}%`}
                    y2={`${next.y}%`}
                />
                )
            })}
            </svg>
        )}

        {/* Cities */}
      	{cities.map((city, idx) => (
        <div
          	key={city.name}
          	title={`${city.name} (${city.x.toFixed(1)}%, ${city.y.toFixed(1)}%)`}
          	className="city"
          	style={{
            	top: `${city.y}%`,
            	left: `${city.x}%`
          	}}
        >
          	{idx + 1}
        </div>
      	))}
    </div>

    <button onClick={solveTSP} disabled={cities.length < 2 || loading} className="tsp-button">
        {loading ? "Solving…" : "Solve TSP"}
    </button>

    <button onClick={() => { setCities([]); setRoute(null); } } className="reset-button">
    Reset
    </button>

    {error && <p className="error-message">{error}</p>}

    {route && (
        <div className="route-info">
        <h2>Route order</h2>
        <ol>
            {route.map((city) => (
                <li key={city.name}>
                {city.name} — (<strong>{city.x.toFixed(1)}</strong>, <strong>{city.y.toFixed(1)}</strong>)
                </li>
        ))}
        </ol>
  		</div>
	)}

    </div>
  	);
}

export default App;
