import React, { useEffect, useState } from 'react';
import './App.css';

interface Task {
  id: string;
  status: string;
}

function App() {
  const [tasks, setTasks] = useState<Task[]>([]);

  useEffect(() => {
    fetch('/tasks')
      .then(response => response.json())
      .then(data => setTasks(data))
      .catch(error => console.error('Error fetching tasks:', error));
  }, []);

  return (
    <div className="App">
      <header className="App-header">
        <h1>Task Queue</h1>
      </header>
      <main>
        <table>
          <thead>
            <tr>
              <th>Task ID</th>
              <th>Status</th>
            </tr>
          </thead>
          <tbody>
            {tasks.map(task => (
              <tr key={task.id}>
                <td>{task.id}</td>
                <td>{task.status}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </main>
    </div>
  );
}

export default App;
