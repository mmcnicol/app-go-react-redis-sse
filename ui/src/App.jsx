import React, { useState, useEffect } from "react";

function App() {
  const [status, setStatus] = useState({
    documents: "pending",
    labResults: "pending",
    emergencyCareSummaries: "pending",
  });
  const [errors, setErrors] = useState([]);
  const [data, setData] = useState(null);

  const userID = "123"; // Assume userID is known

  useEffect(() => {
    // Create an EventSource connection to the backend
    const eventSource = new EventSource(`http://localhost:8080/updates?userID=${userID}`);

    // Listen for updates
    eventSource.addEventListener("documents", () => {
      setStatus((prev) => ({ ...prev, documents: "fetched" }));
    });

    eventSource.addEventListener("labResults", () => {
      setStatus((prev) => ({ ...prev, labResults: "fetched" }));
    });

    eventSource.addEventListener("emergencyCareSummaries", () => {
      setStatus((prev) => ({ ...prev, emergencyCareSummaries: "fetched" }));
    });

    eventSource.addEventListener("error", (event) => {
		//console.log(event.data);
		if (typeof event.data !== 'undefined') {
			setErrors((prev) => [...prev, event.data]);
		}
    });

	eventSource.addEventListener("eagerLoading", () => {
      if (event.data == 'finished') {
		  console.log(event.data);
		  eventSource.close();
	  }
    });
	
	eventSource.onmessage = (event) => {
		if (typeof event.data !== 'undefined') {
			console.log(event.data);
		}
	}
	
    // Cleanup on unmount
    return () => {
      eventSource.close();
    };
  }, []);

  const fetchData = async (type) => {
    const response = await fetch(`http://localhost:8080/data?userID=${userID}&type=${type}`);
    const result = await response.json();
    setData(result);
  };

  return (
    <div>
      <h1>Real-Time Data Fetching</h1>
      <nav>
        <button onClick={() => fetchData("documents")} disabled={status.documents !== "fetched"}>
          Documents {status.documents === "fetched" ? "✅" : "⏳"}
        </button>
        <button onClick={() => fetchData("labResults")} disabled={status.labResults !== "fetched"}>
          Lab Results {status.labResults === "fetched" ? "✅" : "⏳"}
        </button>
        <button onClick={() => fetchData("emergencyCareSummaries")} disabled={status.emergencyCareSummaries !== "fetched"}>
          Emergency Care Summary {status.emergencyCareSummaries === "fetched" ? "✅" : "⏳"}
        </button>
      </nav>

      <h2>Data</h2>
      <pre>{JSON.stringify(data, null, 2)}</pre>

      <h2>Errors</h2>
      <ul>
        {errors.map((error, index) => (
          <li key={index}>{error}</li>
        ))}
      </ul>
    </div>
  );
}

export default App;
